package reqfunc

import (
	"encoding/json"
	"testing"
)

// 获取茅台app版本号
func TestGetMTVersion(t *testing.T) {
	mtversion, err := GetMTVersion()
	if err != nil {
		t.Error(err)
	}
	println(mtversion)
}

// 发送短信验证码
func TestSendCode(t *testing.T) {
	mobile := "18610847758"
	flag, err := SendCode(mobile)
	if err != nil {
		t.Error(err)
	}
	println(flag)
}

// 登录获取用户信息
func TestLogin(t *testing.T) {
	mobile := "18610847758"
	code := "473433"
	token, err := Login(mobile, code)
	if err != err {
		t.Error(err)
	}
	body, err := json.Marshal(token)
	if err != nil {
		t.Error(err)
	}
	println(string(body))
}

// 获取当天的预约商品列表
func TestGetCurrentSessionId(t *testing.T) {
	sessionData, err := GetCurrentSessionId()
	if err != nil {
		t.Error(err)
	}
	body, err := json.Marshal(sessionData)
	if err != nil {
		t.Error(err)
	}
	println(string(body))
}

// 获取门店列表
func TestGetShopList(t *testing.T) {
	shops, err := GetShopList()
	if err != nil {
		t.Error(err)
	}
	body, err := json.Marshal(shops)
	if err != nil {
		t.Error(err)
	}
	println(string(body))
}

// 省市的投放产品和数量
func TestGetShopsByProvince(t *testing.T) {
	shops, err := GetShopsByProvince("山东省", "10213", "686")
	if err != nil {
		t.Error(err)
	}
	body, err := json.Marshal(shops)
	if err != nil {
		t.Error(err)
	}
	println(string(body))
}

// 获取要预约的门店id
func TestGetShopId(t *testing.T) {
	shopId, err := GetShopId(2, "10213", "山东省", "临沂市", "686", "118.257331", "35.002533")
	if err != nil {
		t.Error(err)
	}
	println(shopId)
}

// 根据地址查询坐标
func TestGetLocationByAddress(t *testing.T) {
	geocodes, err := GetLocationByAddress("前崔庄社区")
	if err != nil {
		t.Error(err)
	}
	body, err := json.Marshal(geocodes)
	if err != nil {
		t.Error(err)
	}
	println(string(body))
}

// 预约
func TestReservation(t *testing.T) {
	itemcode := "10213"
	shopId := "137373002002"
	sessionid := "686"
	mtversion := "1.4.3"
	req := UserInfo{
		UserId: "",
		Lat:    "118.257331",
		Lng:    "35.002533",
		Token:  "",
	}
	rt, err := Reservation(req, itemcode, shopId, sessionid, mtversion)
	if err != nil {
		t.Error(err)
	}
	println(rt)
}
