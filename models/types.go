package models

import (
	"gorm.io/gorm"
	"imaotai/db"
)

func Init() {
	db.Gormdb.AutoMigrate(&MtVersion{}, &Session{}, &SessionItem{}, &ShopItem{}, &ShopBean{})
}

type MtVersion struct {
	gorm.Model
	Version string `json:"version"` // app版本号
}

type Session struct {
	gorm.Model
	SessionId int `json:"session_id"` // 预约产品的归属id
}

// 产品信息
type SessionItem struct {
	gorm.Model
	SessionID int    `json:"session_id"`
	Content   string `json:"content"`
	ItemCode  string `json:"itemCode"`
	JumpURL   string `json:"jumpUrl"`
	Picture   string `json:"picture"`
	PictureV2 string `json:"pictureV2"`
	Title     string `json:"title"`
}

// 门店投放产品和数量
type ShopItem struct {
	gorm.Model
	ShopID              string `json:"shop_id"`
	Count               int    `json:"count"`
	MaxReserveCount     int    `json:"maxReserveCount"`
	DefaultReserveCount int    `json:"defaultReserveCount"`
	ItemID              string `json:"itemId"`
	Inventory           int    `json:"inventory"`
	OwnerName           string `json:"ownerName"`
	Province            string `json:"province"`
}

// 门店信息
type ShopBean struct {
	gorm.Model
	Address       string  `json:"address"`
	City          int     `json:"city"`
	CityName      string  `json:"cityName"`
	District      int     `json:"district"`
	DistrictName  string  `json:"districtName"`
	FullAddress   string  `json:"fullAddress"`
	Lat           float64 `json:"lat"`
	Layaway       bool    `json:"layaway"`
	Lng           float64 `json:"lng"`
	Name          string  `json:"name"`
	OpenEndTime   string  `json:"openEndTime"`
	OpenStartTime string  `json:"openStartTime"`
	Province      int     `json:"province"`
	ProvinceName  string  `json:"provinceName"`
	ShopID        string  `json:"shopId"`
	TenantName    string  `json:"tenantName"`
	Distance      float64 `json:"distance" gorm:"-"` //计算，不入库
}
