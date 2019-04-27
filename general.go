package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var generalMethod = OcrMethodData{"GeneralFastOCR", "2018-11-19"}

type GeneralData struct {
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

func (generalData GeneralData) error() bool {
	return generalData.Response.Error.Code != ""
}

func general(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := GeneralData{}
		if resp, err := request(generalMethod, r); err != nil {
			ResponseData{Code: 502, Msg: err.Error()}.response(w)
		} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			ResponseData{Code: 500, Msg: err.Error()}.response(w)
		} else if res.error() {
			ResponseData{Code: 10002, Msg: res.Response.Error.Message}.response(w)
		} else {
			ResponseData{Code: 0, Msg: "success", Data: res.Response}.response(w)
		}
	} else {
		msg := fmt.Sprintf("不支持 %s 方式请求，请使用 %s", r.Method, http.MethodPost)
		ResponseData{Code: 10000, Msg: msg}.response(w)
	}
}
