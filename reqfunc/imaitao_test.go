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
