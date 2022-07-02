package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

// 用于初始化项目的基本配置
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config") // 设置配置文件名称
	// 设置配置文件相对路径， viper 允许多个配置路径，可以不断调用 AddConfigPath()
	//vp.AddConfigPath("configs/")

	// 添加可变更配置文件路径
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml") // 设置配置文件类型

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

// 新增热更新的监听和变更处理
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		// 如果配置文件发生了改变就重新读取配置项
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
