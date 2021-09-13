package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	Port                 int
	PushToReadyQueueSize int
}

var AppSetting = &App{}

type Redis struct {
	Host         string
	DB           int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
}

var RedisSetting = &Redis{}
var cfg *ini.File

func Setup(confPath string) {
	var err error
	cfg, err = ini.Load(confPath)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse '%s': %v", confPath, err)
	}

	mapTo("app", AppSetting)
	mapTo("redis", RedisSetting)

	RedisSetting.DialTimeout = RedisSetting.DialTimeout * time.Second
	RedisSetting.ReadTimeout = RedisSetting.ReadTimeout * time.Second
	RedisSetting.WriteTimeout = RedisSetting.WriteTimeout * time.Second

	if AppSetting.PushToReadyQueueSize == 0 {
		AppSetting.PushToReadyQueueSize = 500
	}
}

// mapTo map section
func mapTo(section string, v interface{}) {
	if _, err := cfg.GetSection(section); err != nil {
		log.Fatalf("cfg must config with section: %s", section)
	}
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
