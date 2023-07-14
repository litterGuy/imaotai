package reqfunc

import (
	"bytes"
	"encoding/json"
	"imaotai/common/errorx"
	"io/ioutil"
	"net/http"
)

// Reservation 预约
func Reservation(user UserInfo, itemCode, shopId string, sessionId string, mtversion string) (string, error) {
	itemArray := make([]map[string]interface{}, 0)
	info := map[string]interface{}{
		"count":  1,
		"itemId": itemCode,
	}
	itemArray = append(itemArray, info)

	params := map[string]interface{}{
		"itemInfoList": itemArray,
		"sessionId":    sessionId,
		"userId":       user.UserId,
		"shopId":       shopId,
	}
	actParam, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	params["actParam"], err = AesEncrypt(string(actParam))
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", errorx.NewDefaultError("Failed to marshal JSON data", err)
	}

	url := "https://app.moutai519.com.cn/xhr/front/mall/reservation/add"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errorx.NewDefaultError("Failed to create HTTP request", err)
	}

	req.Header.Set("MT-Lat", user.Lat)
	req.Header.Set("MT-K", "1675213490331")
	req.Header.Set("MT-Lng", user.Lng)
	req.Header.Set("Host", "app.moutai519.com.cn")
	req.Header.Set("MT-User-Tag", "0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("MT-Network-Type", "WIFI")
	req.Header.Set("MT-Token", user.Token)
	req.Header.Set("MT-Team-ID", "")
	req.Header.Set("MT-Info", "028e7f96f6369cafe1d105579c5b9377")
	req.Header.Set("MT-Device-ID", "2F2075D0-B66C-4287-A903-DBFF6358342A")
	req.Header.Set("MT-Bundle-ID", "com.moutai.mall")
	req.Header.Set("Accept-Language", "en-CN;q=1, zh-Hans-CN;q=0.9")
	req.Header.Set("MT-Request-ID", "167560018873318465")
	req.Header.Set("MT-APP-Version", mtversion)
	req.Header.Set("User-Agent", "iOS;16.3;Apple;?unrecognized?")
	req.Header.Set("MT-R", "clips_OlU6TmFRag5rCXwbNAQ/Tz1SKlN8THcecBp/HGhHdw==")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("userId", user.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errorx.NewDefaultError("HTTP request failed", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errorx.NewDefaultError("Failed to read response body", err)
	}
	return string(respBody), nil
}

type UserInfo struct {
	UserId string `json:"userId"`
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
	Token  string `json:"token"`
}
