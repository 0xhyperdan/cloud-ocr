package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/**
收货授权委托书
*/
var consigneeMethod = OcrMethodData{"GeneralFastOCR", "2018-11-19"}

type ConsigneeResultData struct {
	CompanyName   string `json:"companyName"`   // 公司名
	EmployeeName  string `json:"employeeName"`  // 被授权人
	EmployeeId    string `json:"employeeId"`    // 被授权人身份证
	EmployeePhone string `json:"employeePhone"` // 被授权人电话
	TimeLimit     string `json:"timeLimit"`     // 期限
	Address       string `json:"address"`       // 收货地址
}
type ConsigneeData struct {
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

func (consignee ConsigneeData) error() bool {
	return consignee.Response.Error.Code != ""
}
func (consignee ConsigneeData) convert() ConsigneeResultData {
	var consigneeResult ConsigneeResultData
	text := ""
	for _, d := range consignee.Response.TextDetections {
		text += d.DetectedText
	}
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "，", ",")
	text = strings.ReplaceAll(text, "：", ":")
	text = strings.ReplaceAll(text, "（", "(")
	text = strings.ReplaceAll(text, "）", ")")

	text = strings.Replace(text, "收货授权委托书", "", -1)
	text = strings.Replace(text, ":兹委托我公司员工:", ",", -1)
	text = strings.Replace(text, "(先生/女士)", "", -1)
	text = strings.Replace(text, "(身份证号码:", "", -1)
	text = strings.Replace(text, "联系电话:", "", -1)
	text = strings.Replace(text, ")仅负责我公司与贵单位的所有货物收货签收事宜,我公司对授托人的上述行为承担法律责任。", "", -1)
	text = strings.Replace(text, "收货委托书期限:", ",", -1)
	text = strings.Replace(text, "收货地址:", ",", -1)
	text = strings.Replace(text, "附:被授权人员身份证复印件", ",", -1)

	fmt.Printf("收货授权委托书 ocr result: %s\n", text)
	splits := strings.Split(text, ",")
	fmt.Printf("收货授权委托书 split result: %v\n", splits)
	for k, v := range splits {
		switch k {
		case 0: // 公司名
			consigneeResult.CompanyName = v
			break
		case 1: // 被授权人名
			consigneeResult.EmployeeName = v
			break
		case 2: // 被授权人身份证
			consigneeResult.EmployeeId = v
			break
		case 3: // 被授权人电话
			consigneeResult.EmployeePhone = v
			break
		case 4: // 期限
			consigneeResult.TimeLimit = v
			break
		case 5: // 收货地址
			consigneeResult.Address = v
			break
		}
	}
	return consigneeResult
}
func consignee(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := ConsigneeData{}
		if resp, err := request(consigneeMethod, r); err != nil {
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
