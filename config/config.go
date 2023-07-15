package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

var Configs *Config

func GetConfig(path string) error {
	Configs = new(Config)
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, Configs)
	if err != nil {
		return err
	}
	if len(Configs.Account) == 0 {
		panic("配置文件未配置，启动失败")
	}
	return nil
}

type Config struct {
	Account []*Account `yaml:"account" json:"account"`
}

type Account struct {
	Phone       string  `yaml:"phone" json:"phone"`             // 手机号
	UserId      string  `yaml:"userId" json:"userId"`           // 茅台用户id
	Token       string  `yaml:"token" json:"token"`             // 登录token,有效期一个月
	Lat         float64 `yaml:"lat" json:"lat"`                 // 纬度
	Lng         float64 `yaml:"lng" json:"lng"`                 // 经度
	Item        string  `yaml:"item" json:"item"`               // 预约code，预约哪种酒
	Province    string  `yaml:"province" json:"province"`       // 省份
	City        string  `yaml:"city" json:"city"`               // 城市
	ReserveType int     `yaml:"reserveType" json:"reserveType"` // 预约方式 1-预约本市出货量最大的门店 2-预约你的位置附近门店
}
