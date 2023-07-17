package msg

import (
	"bytes"
	"encoding/json"
	"imaotai/common/errorx"
	"imaotai/config"
	"io/ioutil"
	"net/http"
)

func SendPushPlus(content string) {
	pushplus := PushPlusReq{
		Token: config.Configs.PushPlus.Token,
		Topic: config.Configs.PushPlus.Topic,
		Title: "imaotai",
	}
	pushplus.Content = content
	_ = pushPlus(pushplus)
}

func pushPlus(req PushPlusReq) error {
	url := "http://www.pushplus.plus/send/"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return errorx.NewDefaultError("Failed to marshal JSON data", err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errorx.NewDefaultError("Failed to create HTTP request", err)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return errorx.NewDefaultError("HTTP request failed", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorx.NewDefaultError("Failed to read response body", err)
	}

	rt := new(PushPlusResp)
	err = json.Unmarshal(respBody, &rt)
	if err != nil {
		return errorx.NewDefaultError("Failed to unmarshal response body", err)
	}
	if rt.Code == 200 {
		return nil
	} else {
		return errorx.NewDefault(rt.Msg)
	}
}

type PushPlusReq struct {
	Token   string `json:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Topic   string `json:"topic"`
}

type PushPlusResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
