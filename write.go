package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var WriteUrl = "https://recognition.image.myqcloud.com/ocr/handwriting"

type WriteData struct {
	Code    int
	Message string
	Data    struct {
		SessionId string `json:"sessionId"`
		Items     [] struct {
			ItemString string `json:"itemString"`
			ItemCoord  struct {
				X      int `json:"x"`
				Y      int `json:"y"`
				Width  int `json:"width"`
				Height int `json:"height"`
			}
			ItemConf float64 `json:"itemConf"`
			Words    []struct {
				Character  string  `json:"character"`
				Confidence float64 `json:"confidence"`
			} `json:"words"`
		} `json:"items"`
	}
}

func (writeData WriteData) success() bool {
	return writeData.Code == 0
}

func write(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := WriteData{}
		if resp, err := requestOld(WriteUrl, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.success() {
			ResponseData{res.Code, res.Message, res.Data}.response(w)
		} else {
			ResponseData{Code: 10002, Msg: res.Message}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		response := ResponseData{Code: 10000, Msg: msg}
		response.response(w)
	}
}
