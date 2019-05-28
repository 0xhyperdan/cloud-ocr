package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var LicenseUrl = "https://recognition.image.myqcloud.com/ocr/bizlicense"

type LicenseData struct {
	Code    int
	Message string
	Data    struct {
		SessionId string
		Items     [] struct {
			Item       string
			ItemString string
			ItemCoord  struct {
				X      int
				Y      int
				Width  int
				Height int
			}
			ItemConf float64
		}
	}
}

func (licenseData LicenseData) success() bool {
	return licenseData.Code == 0
}

type LicenseResultData struct {
	RegisterId      string `json:"registerId"`      // 注册号
	Corporation     string `json:"corporation"`     // 法人
	Address         string `json:"address"`         // 注册地址
	CompanyName     string `json:"companyName"`     // 公司名称
	OperatingPeriod string `json:"operatingPeriod"` // 营业期限
	BusinessScope   string `json:"businessScope"`   // 经营范围
}

func (licenseData LicenseData) toResultData() LicenseResultData {
	var licenseResultData LicenseResultData
	for _, v := range licenseData.Data.Items {
		switch v.Item {
		case "注册号":
			licenseResultData.RegisterId = v.ItemString
		case "公司名称":
			licenseResultData.CompanyName = v.ItemString
		case "地址":
			licenseResultData.Address = v.ItemString
		case "法定代表人":
			licenseResultData.Corporation = v.ItemString
		case "营业期限":
			licenseResultData.OperatingPeriod = v.ItemString
		case "经营范围":
			licenseResultData.BusinessScope = v.ItemString
		}
	}
	return licenseResultData
}

func license(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := LicenseData{}
		if resp, err := requestOld(LicenseUrl, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.success() {
			ResponseData{res.Code, "success", res.toResultData()}.response(w)
		} else {
			ResponseData{Code: res.Code, Msg: res.Message}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
