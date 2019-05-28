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
	var invoiceResultData InvoiceResultData
	var text string
	for _, v := range generalData.Response.TextDetections {
		if strings.Contains(v.DetectedText, "电子承兑信息") {
			break
		}
		text += v.DetectedText
	}
	fmt.Printf("开票信息识别结果: %s\n", text)
	text = strings.ReplaceAll(text, "：", ":")
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "（", "(")
	text = strings.ReplaceAll(text, "）", ")")
	text = strings.Replace(text, "开票信息", "", 1)
	text = strings.Replace(text, "公司名称:", "", -1)
	text = strings.Replace(text, "公司地址:", ",", -1)
	text = strings.Replace(text, "开户银行:", ",", -1)
	text = strings.Replace(text, "银行汇款账号:", ",", -1)
	text = strings.Replace(text, "纳税人识别号:", ",", -1)
	text = strings.Replace(text, "邮编:", ",", -1)
	text = strings.Replace(text, "电话:", ",", -1)
	text = strings.Replace(text, "仓库地址:", ",", -1)

	infos := strings.Split(text, ",")
	fmt.Printf("开票信息分割结果: %v\n", infos)
	for i, v := range infos {
		switch i {
		case 0: // 公司名
			invoiceResultData.CompanyName = v
			break
		case 1: // 公司地址
			invoiceResultData.CompanyAddress = v
			break
		case 2: // 开户行
			invoiceResultData.BankName = v
			break
		case 3: // 汇款账号
			invoiceResultData.BankAccount = v
			break
		case 4: // 纳税人识别号
			invoiceResultData.TaxCode = v
			break
		case 5: // 邮编
			invoiceResultData.ZipCode = v
			break
		case 6: // 电话
			invoiceResultData.Telephone = v
			break
		case 7: // 仓库地址
			invoiceResultData.DepotAddress = v
			break
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
