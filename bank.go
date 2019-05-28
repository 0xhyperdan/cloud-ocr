package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var BankUrl = "https://recognition.image.myqcloud.com/ocr/bankcard"

type BankData struct {
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

func (bankData BankData) success() bool {
	return bankData.Code == 0
}

type BankResultData struct {
	// 银行卡信息
	BankCardCode string `json:"bankCardCode"` // 卡号
	BankCardType string `json:"bankCardType"` // 卡类型
	BankCardName string `json:"bankCardName"` // 卡名字
	BankInfo     string `json:"bankInfo"`     // 银行信息
	ValidityDate string `json:"validityDate"` // 有效期
}

func (bankData BankData) toResultData() BankResultData {
	var bankResultData BankResultData
	for _, v := range bankData.Data.Items {
		fmt.Printf("item => %v\n", v.Item)
		fmt.Printf("item_value => %v\n", v.ItemString)
		switch v.Item {
		case "卡号":
			bankResultData.BankCardCode = v.ItemString
			break
		case "卡类型":
			bankResultData.BankCardType = v.ItemString
			break
		case "卡名字":
			bankResultData.BankCardName = v.ItemString
			break
		case "银行信息":
			bankResultData.BankInfo = v.ItemString
			break
		case "有效期":
			bankResultData.ValidityDate = v.ItemString
			break
		}
	}
	return bankResultData
}

func bank(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := BankData{}
		if resp, err := requestOld(BankUrl, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.success() {
			ResponseData{res.Code, res.Message, res.toResultData()}.response(w)
		} else {
			ResponseData{Code: 10002, Msg: res.Message}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
