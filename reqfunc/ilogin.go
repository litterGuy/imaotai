package reqfunc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"imaotai/common/errorx"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GetMTVersion 获取app版本号
func GetMTVersion() (string, error) {
	url := "https://apps.apple.com/cn/app/i%E8%8C%85%E5%8F%B0/id1600482450"
	response, err := http.Get(url)
	if err != nil {
		return "", errorx.NewDefault(err.Error())
	}
	defer response.Body.Close()

	htmlContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errorx.NewDefault(err.Error())
	}
	pattern := regexp.MustCompile("new__latest__version\">(.*?)</p>")
	matches := pattern.FindStringSubmatch(string(htmlContent))
	if len(matches) > 1 {
		mtVersion := matches[1]
		mtVersion = strings.ReplaceAll(mtVersion, "版本 ", "")
		return mtVersion, nil
	}
	return "", errorx.NewDefault("获取版本号错误")
}

// SendCode 发送短信验证码
func SendCode(mobile string) (bool, error) {
	mtversion, err := GetMTVersion()
	if err != nil {
		return false, errorx.NewDefault(err.Error())
	}
	data := map[string]interface{}{
		"mobile":    mobile,
		"md5":       signature(mobile),
		"timestamp": fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, errorx.NewDefaultError("Failed to marshal JSON data", err)
	}

	url := "https://app.moutai519.com.cn/xhr/front/user/register/vcode"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, errorx.NewDefaultError("Failed to create HTTP request", err)
	}

	req.Header.Set("MT-Lat", "28.499562")
	req.Header.Set("MT-K", "1675213490331")
	req.Header.Set("MT-Lng", "102.182324")
	req.Header.Set("Host", "app.moutai519.com.cn")
	req.Header.Set("MT-User-Tag", "0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("MT-Network-Type", "WIFI")
	req.Header.Set("MT-Team-ID", "")
	req.Header.Set("MT-Info", "028e7f96f6369cafe1d105579c5b9377")
	req.Header.Set("MT-Device-ID", "2F2075D0-B66C-4287-A903-DBFF6358342A")
	req.Header.Set("MT-Bundle-ID", "com.moutai.mall")
	req.Header.Set("Accept-Language", "en-CN;q=1, zh-Hans-CN;q=0.9")
	req.Header.Set("MT-Request-ID", "167560018873318465")
	req.Header.Set("MT-APP-Version", mtversion)
	req.Header.Set("User-Agent", "iOS;16.3;Apple;?unrecognized?")
	req.Header.Set("MT-R", "clips_OlU6TmFRag5rCXwbNAQ/Tz1SKlN8THcecBp/HGhHdw==")
	req.Header.Set("Content-Length", "93")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errorx.NewDefaultError("HTTP request failed", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errorx.NewDefaultError("Failed to read response body", err)
	}

	rt := new(ImaotaiResp)
	err = json.Unmarshal(respBody, &rt)
	if err != nil {
		return false, errorx.NewDefaultError("Failed to unmarshal response body", err)
	}

	if rt.Code == 2000 {
		return true, nil
	} else {
		return false, errorx.NewDefault("发送验证码错误:" + rt.Message)
	}
}

// Login 登录获取用户信息
func Login(mobile string, code string) (*LoginUser, error) {
	mtversion, err := GetMTVersion()
	if err != nil {
		return nil, errorx.NewDefault(err.Error())
	}
	data := map[string]string{
		"mobile":         mobile,
		"vCode":          code,
		"ydToken":        "",
		"ydLogId":        "",
		"md5":            signature(mobile + code + "" + ""),
		"timestamp":      fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)),
		"MT-APP-Version": mtversion,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errorx.NewDefaultError("Failed to marshal JSON data", err)
	}

	url := "https://app.moutai519.com.cn/xhr/front/user/register/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errorx.NewDefaultError("Failed to create HTTP request", err)
	}

	req.Header.Set("MT-Lat", "28.499562")
	req.Header.Set("MT-K", "1675213490331")
	req.Header.Set("MT-Lng", "102.182324")
	req.Header.Set("Host", "app.moutai519.com.cn")
	req.Header.Set("MT-User-Tag", "0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("MT-Network-Type", "WIFI")
	req.Header.Set("MT-Team-ID", "")
	req.Header.Set("MT-Info", "028e7f96f6369cafe1d105579c5b9377")
	req.Header.Set("MT-Device-ID", "2F2075D0-B66C-4287-A903-DBFF6358342A")
	req.Header.Set("MT-Bundle-ID", "com.moutai.mall")
	req.Header.Set("Accept-Language", "en-CN;q=1, zh-Hans-CN;q=0.9")
	req.Header.Set("MT-Request-ID", "167560018873318465")
	req.Header.Set("MT-APP-Version", mtversion)
	req.Header.Set("User-Agent", "iOS;16.3;Apple;?unrecognized?")
	req.Header.Set("MT-R", "clips_OlU6TmFRag5rCXwbNAQ/Tz1SKlN8THcecBp/HGhHdw==")
	req.Header.Set("Content-Length", "93")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errorx.NewDefaultError("HTTP request failed", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errorx.NewDefaultError("Failed to read response body", err)
	}

	rt := new(LoginResp)
	err = json.Unmarshal(respBody, &rt)
	if err != nil {
		return nil, errorx.NewDefaultError("Failed to unmarshal response body", err)
	}

	if rt.Code == 2000 {
		return rt.Data, nil
	} else {
		return nil, errorx.NewDefault("登录请求-失败:" + rt.Message)
	}
}

const salt = "2af72f100c356273d46284f6fd1dfc08"

/*
*
获取验证码的md5签名，密钥+手机号+时间
登录的md5签名：密钥+mobile+vCode+ydLogId+ydToken
*/
func signature(content string) string {
	text := salt + content + fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	hash := md5.Sum([]byte(text))
	md5Hash := hex.EncodeToString(hash[:])
	return md5Hash
}
