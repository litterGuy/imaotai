package service

import (
	"errors"
	"fmt"
	"imaotai/config"
	"imaotai/db"
	"imaotai/models"
	"math"
	"sort"
	"strconv"
	"strings"
)

// GetShopId 获取预约门店id
func GetShopId(itemCode string, account config.Account) (string, error) {
	//查询所在省市的投放产品和数量
	var shopitems []*models.ShopItem
	err := db.Gormdb.Where("province", account.Province).Where("item_id", itemCode).Find(&shopitems).Error
	if err != nil {
		return "", err
	}
	if len(shopitems) == 0 {
		return "", errors.New("查询投放产品信息和数量失败")
	}
	//获取门店列表
	var shops []*models.ShopBean
	err = db.Gormdb.Find(&shops).Error
	if err != nil {
		return "", err
	}
	if len(shops) == 0 {
		return "", errors.New("查询门店失败")
	}
	shopId := ""
	if account.ReserveType == 1 {
		//预约本市出货量最大的门店
		shopId = getMaxInventoryShopId(shopitems, shops, account.City)
		if len(shopId) == 0 {
			shopId = getMinDistanceShopId(shops, account.Province, account.Lat, account.Lng)
		}
	} else if account.ReserveType == 2 {
		shoplist := getPartyShopList(shopitems, shops)
		// 预约本省距离最近的门店
		shopId = getMinDistanceShopId(shoplist, account.Province, account.Lat, account.Lng)
	}
	return shopId, nil
}

func getPartyShopList(list1 []*models.ShopItem, list2 []*models.ShopBean) []*models.ShopBean {
	rt := make([]*models.ShopBean, 0)

	tmpmap := make(map[string]*models.ShopBean)
	for _, bean := range list2 {
		tmpmap[bean.ShopID] = bean
	}
	for _, item := range list1 {
		if t, ok := tmpmap[item.ShopID]; ok {
			rt = append(rt, t)
		}
	}
	return rt
}

// 预约本市出货量最大的门店
func getMaxInventoryShopId(list1 []*models.ShopItem, list2 []*models.ShopBean, city string) string {
	// 按照出货量大小排序
	sort.Slice(list1, func(i, j int) bool {
		return list1[i].Inventory > list1[j].Inventory
	})

	cityShopIDList := make([]string, 0)
	for _, iShop := range list2 {
		if strings.Contains(iShop.CityName, city) {
			cityShopIDList = append(cityShopIDList, iShop.ShopID)
		}
	}
	for _, i := range list1 {
		if contains(cityShopIDList, i.ShopID) {
			return i.ShopID
		}
	}
	return ""
}

// 预约本省距离最近的门店
func getMinDistanceShopId(list2 []*models.ShopBean, province string, lat, lng float64) string {
	iShopList := make([]*models.ShopBean, 0)
	for _, iShop := range list2 {
		if strings.Contains(iShop.ProvinceName, province) {
			iShopList = append(iShopList, iShop)
		}
	}

	myPoint := MapPoint{
		Latitude:  parseLocation(fmt.Sprintf("%v", lat)),
		Longitude: parseLocation(fmt.Sprintf("%v", lng)),
	}

	for i := range iShopList {
		point := MapPoint{
			Latitude:  parseLocation(fmt.Sprintf("%v", iShopList[i].Lat)),
			Longitude: parseLocation(fmt.Sprintf("%v", iShopList[i].Lng)),
		}
		distance := getDistance(myPoint, point)
		iShopList[i].Distance = distance
	}

	sort.Slice(iShopList, func(i, j int) bool {
		return iShopList[i].Distance < iShopList[j].Distance
	})

	return iShopList[0].ShopID
}

func parseLocation(str string) float64 {
	data, _ := strconv.ParseFloat(str, 64)
	return (data * math.Pi) / 180 //将角度换算为弧度
}

func getDistance(point1, point2 MapPoint) float64 {
	lat1 := point1.Latitude
	lat2 := point2.Latitude
	lonDifference := point1.Longitude - point2.Longitude
	//计算两点之间距离   6378137.0 取自WGS84标准参考椭球中的地球长半径(单位:m)
	distance := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin((lat1-lat2)/2), 2)+
		math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(lonDifference/2), 2))) * 6378137.0
	return distance
}

type MapPoint struct {
	Latitude  float64
	Longitude float64
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
