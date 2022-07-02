package main

import (
	"context"
	"demo/ch02/global"
	"demo/ch02/internal/model"
	"demo/ch02/internal/routers"
	"demo/ch02/pkg/logger"
	"demo/ch02/pkg/setting"
	"demo/ch02/pkg/tracer"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port         string
	runMode      string
	config       string
	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

// Go 中的执行顺序: 全局变量初始化 =>init() => main()
// 在 main() 之前自动执行，进行初始化操作
func init() {
	// 添加命令行参数处理逻辑
	setupFlag()
	// 获取初始化配置
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	// 数据库初始化
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	// 日志初始化
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	// 链路追踪初始化
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言项目实战学习
// @termsOfService https://github.com/go-programming-tour-book
func main() {
	// 添加版本信息
	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
		return
	}

	// 使用映射好的配置设置 gin 的运行模式: debug
	gin.SetMode(global.ServerSetting.RunMode)
	// 不再使用默认路由而使用项目下自定义的路由
	// router := gin.Default()
	router := routers.NewRouter()
	// 自定义 http.Server
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort, // 设置监听端口
		Handler:        router,                              // 设置处理程序
		ReadTimeout:    global.ServerSetting.ReadTimeout,    // 允许读取最大时间
		WriteTimeout:   global.ServerSetting.WriteTimeout,   // 允许写入最大时间
		MaxHeaderBytes: 1 << 20,                             // 请求头最大字节数
	}
	//// 调用 ListenAndServe() 监听
	//if err := s.ListenAndServe(); err != nil {
	//	log.Fatalf("监听失败：%v", err)
	//}
	// 从此处开始修改 使项目支持优雅重启和停止
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal)
	// 接受 syscall.SIGINT 和 syscall.SIGTERM 信号 两个都是终止信号
	/*
		signal.Notify()
		通知使包信号将传入信号中继到 quit。
		如果没有提供信号，所有传入的信号将被中继到 quit。
		否则，只有提供的信号会。
	*/
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shut down server...")

	// 最大时间控制，用于通知该服务端它有 5s 的时间来处理原有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// cancel 类: type CancelFunc func()
	defer cancel()
	/*
		Shutdown 优雅地关闭服务器而不中断任何活动连接。
		关闭首先关闭所有打开的侦听器，然后关闭所有空闲连接，
		然后无限期地等待连接返回空闲状态，然后关闭。
		如果提供的上下文在关闭完成之前过期，则 Shutdown 返回上下文的错误，
		否则返回关闭服务器的底层侦听器返回的任何错误。
	*/
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func setupSetting() error {
	//settings, err := setting.NewSetting()
	// 如果存在则覆盖原有的文件配置
	settings, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = settings.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = settings.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = settings.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = settings.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	err = settings.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.AppSetting.DefaultContextTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	// global.ServerSetting.ReadTimeout *=1000，将秒转换成毫秒
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	// 如果存在则覆盖原有的文件配置
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

func setupDBEngine() error {
	var err error
	// 初始化数据库连接信息，注意此处不是 := 而是 =，使用前者会导致在其他包中调用该变量时值为 nil
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	// 使用了 lumberjack 作为日志库的 io.Writer
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		// 设置生成日志文件存储相对位置与文件名
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,  // 设置最大占用空间
		MaxAge:    10,   // 设置日志文件最大生存周期
		LocalTime: true, // 设置日志文件名的时间格式为本地时间
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupTracer() error {
	// 6831 端口 以 compact 协议接受 jaeger .thrift 数据
	jaegerTracer, _, err := tracer.NewJaegerTrace("blog_service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	// 添加版本信息
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
