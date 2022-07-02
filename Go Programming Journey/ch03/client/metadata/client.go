package main

import (
	"ch03/global"
	"ch03/internal/middleware"
	"ch03/pkg/tracer"
	pb "ch03/proto"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

type Auth struct {
	Appkey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.Appkey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

func main() {
	auth := Auth{
		Appkey:    "admin",
		AppSecret: "go-learning",
	}
	ctx := context.Background()
	// 调用 grpc.WithPerRPCCredentials() 进行注册
	// grpc.WithPerRPCCredentials() 设置凭据并在每个出站 RPC 上放置身份验证状态
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&auth),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				middleware.UnaryContextTimeout(),
				grpc_retry.UnaryClientInterceptor(
					grpc_retry.WithMax(2),
					grpc_retry.WithCodes(
						codes.Unknown,
						codes.Internal,
						codes.DeadlineExceeded,
					),
				),
				middleware.ClientTracing(),
			)),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				middleware.StreamContextTimeout())),
	}
	clientConn, _ := GetClientConn(ctx, "tag-service", opts)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Golang"})
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, serviceName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	config := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second * 60,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	r := &naming.GRPCResolver{Client: cli}
	target := fmt.Sprintf("/etcdv3://go-learning-test/grpc/%s", serviceName)

	// grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
	// 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
	//opts = append(opts, grp c.WithInsecure())
	// insecure.NewCredentials 返回一个禁用传输安全的凭据
	opts = append(opts, grpc.WithInsecure(), grpc.WithBalancer(grpc.RoundRobin(r)), grpc.WithBlock())

	/*
		grpc.DialContext() 创建到给定目标的客户端连接。
		默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
		要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
	*/
	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	var err error
	jaegerTracer, _, err := tracer.NewJaegerTracer("article-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
