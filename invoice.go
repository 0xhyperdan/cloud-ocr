package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var invoiceMethod = OcrMethodData{"GeneralFastOCR", "2018-11-19"}

type InvoiceResultData struct {
	// 开票信息
	CompanyName    string `json:"companyName"`    // 公司名称
	CompanyAddress string `json:"companyAddress"` // 公司地址
	BankName       string `json:"bankName"`       // 开户行
	BankAccount    string `json:"bankAccount"`    // 银行汇款账号
	TaxCode        string `json:"taxCode"`        // 纳税人识别号
	ZipCode        string `json:"zipCode"`        // 邮政编码
	Telephone      string `json:"telephone"`      // 电话
	DepotAddress   string `json:"depotAddress"`   // 仓库地址
}

func (generalData GeneralData) toResultData() InvoiceResultData {
	var isCompanyAddr bool
	var invoiceResultData InvoiceResultData
	for _, v := range generalData.Response.TextDetections {
		item := strings.ReplaceAll(v.DetectedText, "：", ":")
		if strings.Contains(item, "电子承兑信息") {
			break
		}
		if strings.Contains(item, ":") {
			vale := strings.Split(item, ":")[1]
			if strings.Contains(item, "公司名称") {
				isCompanyAddr = true
				invoiceResultData.CompanyName = vale
			} else if strings.Contains(item, "公司地址") {
				invoiceResultData.CompanyAddress = vale
			} else if strings.Contains(item, "开户银行") {
				invoiceResultData.BankName = vale
			} else if strings.Contains(item, "银行汇款账号") {
				invoiceResultData.BankAccount = vale
			} else if strings.Contains(item, "纳税人识别号") {
				invoiceResultData.TaxCode = vale
			} else if strings.Contains(item, "邮编") {
				invoiceResultData.ZipCode = vale
			} else if strings.Contains(item, "电话") {
				invoiceResultData.Telephone = vale
			} else if strings.Contains(item, "仓库地址") {
				isCompanyAddr = false
				invoiceResultData.DepotAddress = vale
			}
		} else if strings.Contains(item, "号") ||
			strings.Contains(item, "单元") ||
			strings.Contains(item, "仓库") ||
			strings.Contains(item, "公司") ||
			strings.Contains(item, "栋") ||
			strings.Contains(item, "幢") ||
			strings.Contains(item, "楼") {
			if isCompanyAddr {
				invoiceResultData.CompanyAddress += item
			} else {
				invoiceResultData.DepotAddress += item
			}
		}
	}
	return invoiceResultData
}

func invoice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := GeneralData{}
		if resp, err := request(invoiceMethod, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.error() {
			ResponseData{10002, res.Response.Error.Message, res}.response(w)
		} else {
			ResponseData{0, "success", res.toResultData()}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
