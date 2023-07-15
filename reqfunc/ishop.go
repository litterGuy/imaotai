package reqfunc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetCurrentSessionId 获取预约商品列表
func GetCurrentSessionId() (*SessionData, error) {
	dayTime := getCurrentDayTime()

	res, err := http.Get(fmt.Sprintf("https://static.moutai519.com.cn/mt-backend/xhr/front/mall/index/session/get/%d", dayTime))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(SessionResp)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data.Code != 2000 {
		return nil, errors.New(data.Message)
	}

	return &data.Data, nil
}

// GetShopsByProvince 省市的投放产品和数量
func GetShopsByProvince(province, itemcode, sessionid string) ([]*ShopItemBean, error) {
	dayTime := getCurrentDayTime()

	url := fmt.Sprintf("https://static.moutai519.com.cn/mt-backend/xhr/front/mall/shop/list/slim/v3/%s/%s/%s/%d", sessionid, province, itemcode, dayTime)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(ShopResp)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data.Code != 2000 {
		return nil, errors.New(data.Message)
	}

	rt := make([]*ShopItemBean, 0)
	for _, shop := range data.Data.Shops {
		for _, item := range shop.Items {
			t := new(ShopItemBean)
			t.ShopRespItems = item
			t.ShopID = shop.ShopID
			rt = append(rt, t)
		}
	}

	return rt, nil
}

// GetShopList 获取门店列表
func GetShopList() (map[string]ShopBean, error) {
	// 获取资源列表
	res, err := http.Get("https://static.moutai519.com.cn/mt-backend/xhr/front/mall/resource/get")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resourceData := new(ShopResourceResp)
	err = json.Unmarshal(body, &resourceData)
	if err != nil {
		return nil, err
	}
	if resourceData.Code != 2000 {
		return nil, errors.New(resourceData.Message)
	}
	// 从资源列表中获取门店列表
	shops, err := getShopList(resourceData.Data.MtshopsPc.Url)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func getShopList(url string) (map[string]ShopBean, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rt := make(map[string]ShopBean)
	err = json.Unmarshal(body, &rt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

/**
@Deprecated
shopType 	1：预约本市出货量最大的门店，2：预约你的位置附近门店
itemCode   	预约项目code
province 	省份，例如：河北省，北京市
city     	市：例如石家庄市
*/

// GetShopId 获取要预约的门店id
func GetShopId(shopType int, itemCode, province, city string, sessionid string, lat, lng string) (string, error) {
	//查询所在省市的投放产品和数量
	shops, err := GetShopsByProvince(province, itemCode, sessionid)
	if err != nil {
		return "", err
	}
	//获取门店列表
	ishops, err := GetShopList()
	if err != nil {
		return "", err
	}
	//获取今日的门店信息列表
	list := make([]ShopBean, 0)
	for _, shop := range shops {
		if t, ok := ishops[shop.ShopID]; ok {
			list = append(list, t)
		}
	}

	shopId := ""
	if shopType == 1 {
		//预约本市出货量最大的门店
		shopId = getMaxInventoryShopId(shops, list, city)
		if len(shopId) == 0 {
			shopId = getMinDistanceShopId(list, province, lat, lng)
		}
	} else if shopType == 2 {
		// 预约本省距离最近的门店
		shopId = getMinDistanceShopId(list, province, lat, lng)
	}
	return shopId, nil
}

// 预约本市出货量最大的门店
func getMaxInventoryShopId(list1 []*ShopItemBean, list2 []ShopBean, city string) string {
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
func getMinDistanceShopId(list2 []ShopBean, province, lat, lng string) string {
	iShopList := make([]ShopBean, 0)
	for _, iShop := range list2 {
		if strings.Contains(iShop.ProvinceName, province) {
			iShopList = append(iShopList, iShop)
		}
	}

	myPoint := MapPoint{
		Latitude:  parseLocation(lat),
		Longitude: parseLocation(lng),
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

func getCurrentDayTime() int64 {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 处理错误
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	dayTime := startOfDay.UnixNano() / 1e6
	return dayTime
}
