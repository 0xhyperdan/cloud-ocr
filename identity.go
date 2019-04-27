package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IdentityData struct {
	Response struct {
		RequestId    string `json:"requestId"`    // 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		Name         string `json:"name"`         // 证件姓名
		Sex          string `json:"sex"`          // 性别
		Nation       string `json:"nation"`       // 民族
		Birth        string `json:"birth"`        // 出生日期
		Address      string `json:"address"`      // 地址
		IdNum        string `json:"idNum"`        // 身份证号
		ValidDate    string `json:"validDate"`    // 证件的有效期
		Authority    string `json:"authority"`    // 发证机关
		AdvancedInfo string `json:"advancedInfo"` // 扩展信息，根据请求的可选字段返回对应内容，不请求则不返回。目前支持的扩展字段为： IdCard身份证照片，请求CropIdCard时返回； Portrait人像照片，请求CropPortrait时返回； WarnInfos告警信息（Code告警码，Msg告警信息），识别出翻拍件或复印件时返回。
		Error        struct {
			Code    string
			Message string
		}
	}
}

func (identity IdentityData) success() bool {
	return identity.Response.Error.Code == ""
}

var IDCardOCR = OcrMethodData{"IDCardOCR", "2018-11-19"}

func identity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := IdentityData{}
		if resp, err := request(IDCardOCR, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.success() {
			ResponseData{0, "success", res.Response}.response(w)
		} else {
			ResponseData{Code: 10002, Msg: res.Response.Error.Message}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
