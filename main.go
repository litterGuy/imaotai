package main

import (
	"encoding/json"
	"flag"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"os"
)

func main() {
	var configpath string
	var dbpath string
	flag.StringVar(&configpath, "path", "config.yml", "配置文件路径")
	flag.StringVar(&dbpath, "db", "imaotai.db", "数据库路径")

	config, err := getConfig(configpath)
	if err != nil {
		panic(err)
	}
	s, _ := json.Marshal(config)
	println(string(s))

	// 数据库
	err = Init(dbpath)
	if err != nil {
		panic(err)
	}

	err = gormdb.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	err = gormdb.Create(&User{Mobile: config.Account[0].Phone}).Error
	if err != nil {
		panic(err)
	}
	var user User
	err = gormdb.Find(&user).Error
	if err != nil {
		panic(err)
	}
	a, _ := json.Marshal(user)
	println(string(a))

	err = gormdb.Model(&User{}).Where("ID", user.ID).Update("UserName", "常理").Error
	if err != nil {
		panic(err)
	}

	err = gormdb.Delete(&user).Error
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	UserID       int64  `json:"userId"`
	UserName     string `json:"userName"`
	Mobile       string `json:"mobile"`
	VerifyStatus int    `json:"verifyStatus"`
	IDCode       string `json:"idCode"`
	IDType       int    `json:"idType"`
	Token        string `json:"token"`
	UserTag      int    `json:"userTag"`
	Cookie       string `json:"cookie"`
	Did          string `json:"did"`

	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	FormattedAddress string  `json:"formatted_address"`
	Country          string  `json:"country"`
	Province         string  `json:"province"`
	City             string  `json:"city"`
}

func getConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type Config struct {
	Account []*Account `yaml:"account" json:"account"`
}

type Account struct {
	Phone string `yaml:"phone" json:"phone"`
}
