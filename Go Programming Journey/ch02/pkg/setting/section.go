package setting

import "time"

// 声明配置属性的结构体
// 服务配置结构体
type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 应用配置结构体
type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	UploadSavePath        string
	UploadServerUrl       string
	UploadImageMaxSize    int
	UploadImageAllowExts  []string
	DefaultContextTimeout time.Duration
}

// 数据库配置结构体
type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

// JWT 配置结构体
type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

// Email 结构体
type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

var sections = make(map[string]interface{})

// 读取相应配置的配置方法
func (s *Setting) ReadSection(k string, v interface{}) error {
	// 将配置文件 按照 父节点读取到相应的struct中
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	// 针对重载应用配置项，新增处理方法
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// 重新读取配置
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
