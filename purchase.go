package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/**
采购授权委托书
*/
var purchaseMethod = OcrMethodData{"GeneralFastOCR", "2018-11-19"}

type PurchaseResultData struct {
	CompanyName   string `json:"companyName"`   // 公司名
	EmployeeName  string `json:"employeeName"`  // 被授权人
	EmployeeId    string `json:"employeeId"`    // 被授权人身份证
	//EmployeePhone string `json:"employeePhone"` // 被授权人电话
	TimeLimit     string `json:"timeLimit"`     // 期限
}
type PurchaseData struct {
	Response struct {
		Language       string `json:"language"`
		RequestId      string `json:"requestId"`
		TextDetections []struct {
			DetectedText string `json:"detectedText"` // 识别出的文本行内容
			Confidence   int    `json:"confidence"`   // 置信度 0 ~100
			Polygon      []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"polygon"`                        // 文本行坐标，以四个顶点坐标表示 注意：此字段可能返回 null，表示取不到有效值。
			AdvancedInfo string `json:"advancedInfo"` //此字段为扩展字段。 GeneralBasicOcr接口返回段落信息Parag ，包含ParagNo。
		} `json:"textDetections"`
		Error struct {
			Code    string
			Message string
		}
	}
}

func (purchase PurchaseData) error() bool {
	return purchase.Response.Error.Code != ""
}
func (purchase PurchaseData) convert() PurchaseResultData {
	var purchaseResult PurchaseResultData
	text := ""
	for _, d := range purchase.Response.TextDetections {
		text += d.DetectedText
	}
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "，", ",")
	text = strings.ReplaceAll(text, "：", ":")
	text = strings.ReplaceAll(text, "（", "(")
	text = strings.ReplaceAll(text, "）", ")")

	text = strings.Replace(text, "采购授权委托书", "", -1)
	text = strings.Replace(text, ":兹授权", ",", -1)
	text = strings.Replace(text, "(先生/女士)为我单位采购代表,按药品相关法律法规要求,在许可范围内,负责与贵公司的采购事宜,请给予支持与合作。", "",
		-1)
	text = strings.Replace(text, "被授权人身份证号码:", ",", 1)
	//text = strings.Replace(text, "被授权人联系电话:", ",", -1)
	text = strings.Replace(text, "有效期限:", ",", -1)
	text = strings.Replace(text, "附:被授权人员身份证复印件", ",", -1)

	fmt.Printf("采购授权委托书 ocr result: %s\n", text)
	splits := strings.Split(text, ",")
	fmt.Printf("采购授权委托书 split result: %v\n", splits)
	for k, v := range splits {
		switch k {
		case 0: // 公司名
			purchaseResult.CompanyName = v
			break
		case 1: // 被授权采购人
			purchaseResult.EmployeeName = v
			break
		case 2: // 被授权采购人身份证
			purchaseResult.EmployeeId = v
			break
		case 3: // 期限
			purchaseResult.TimeLimit = v
			break
		}
	}
	return purchaseResult
}
func purchase(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := PurchaseData{}
		if resp, err := request(purchaseMethod, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.error() {
			ResponseData{Code: 10002, Msg: res.Response.Error.Message}.response(w)
		} else {
			ResponseData{Code: 0, Msg: "success", Data: res.convert()}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
