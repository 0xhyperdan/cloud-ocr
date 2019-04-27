package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	requestMethod = "POST"
	service       = "ocr"
	host          = "ocr.tencentcloudapi.com"
	region        = "ap-beijing"
	algorithm     = "TC3-HMAC-SHA256"
	contentType   = "application/json"
	tc3           = "tc3_request"
	signedHeaders = "content-type;host"
	RequestUrl    = "https://ocr.tencentcloudapi.com"
)

var conf Config

type OcrMethodData struct {
	Name    string
	Version string
}

// 旧版本签名
func getSign() string {
	tm := time.Now().Unix()
	origin := url.Values{}
	origin.Add("a", conf.AppId)
	origin.Add("b", "b")
	origin.Add("k", conf.SecretId)
	origin.Add("e", strconv.FormatInt(tm+2592000, 10))
	origin.Add("t", strconv.FormatInt(tm, 10))
	origin.Add("r", strconv.Itoa(getRandom()))
	origin.Add("f", "")
	//fmt.Printf("origin sign: %s\n", origin)
	strOrigin := origin.Encode()
	byteSign := hmac1Secret([]byte(conf.SecretKey), strOrigin)
	return base64.StdEncoding.EncodeToString(append(byteSign, []byte(strOrigin)...))
}

// 随机数字
func getRandom() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000000)
}

// 获取请求的 Header
func getHeader(payload []byte, action, version string) http.Header {
	header := http.Header{}
	tm := time.Now()
	timeSpan := strconv.FormatInt(tm.Unix(), 10)
	header.Add("Authorization", getAuth(tm, payload))
	header.Add("Host", host)
	header.Add("Content-Type", contentType)
	header.Add("X-TC-Action", action)
	header.Add("X-TC-Timestamp", timeSpan)
	header.Add("X-TC-Version", version)
	header.Add("X-TC-Region", region)
	return header
}

// 1. 拼接规范请求串
func getSpliceData(payload []byte) string {
	h := sha256.New()
	h.Write(payload)
	hs := h.Sum(nil)

	httpRequestMethod := requestMethod
	uri := "/"
	query := ""
	headers := "content-type:" + contentType + "\nhost:" + host + "\n"
	hashedRequestPayload := hex.EncodeToString(hs)

	request := httpRequestMethod + "\n" +
		uri + "\n" +
		query + "\n" +
		headers + "\n" +
		signedHeaders + "\n" +
		hashedRequestPayload

	//fmt.Printf("getSpliceData result: %s\n\n", request)
	return request
}

// 2. 拼接待签名字符串
func getStringToSign(splice string, tm time.Time) string {
	h := sha256.New()
	h.Write([]byte(splice))
	hs1 := h.Sum(nil)

	hashedRequest := hex.EncodeToString(hs1)
	scope := tm.Format("2006-01-02") + "/" + service + "/" + tc3
	timeSpan := strconv.FormatInt(tm.Unix(), 10)

	stringToSign := algorithm + "\n" + timeSpan + "\n" + scope + "\n" + hashedRequest
	//fmt.Printf("getStringToSign result: %s\n\n", stringToSign)
	return stringToSign
}

// 	3. 计算签名
func getSignature(stringToSign string, tm time.Time) string {
	secretDate := hmac256Secret([]byte("TC3"+conf.SecretKey), tm.Format("2006-01-02"))
	secretService := hmac256Secret(secretDate, service)
	secretSigning := hmac256Secret(secretService, tc3)
	signature := hex.EncodeToString(hmac256Secret(secretSigning, stringToSign))
	//fmt.Printf("getSignature result: %s\n\n", signature)
	return signature
}

// 4. 拼接 Authorization
func getAuth(tm time.Time, payload []byte) string {
	scope := tm.Format("2006-01-02") + "/" + service + "/" + tc3
	auth := algorithm + " " + "Credential=" + conf.SecretId + "/" + scope + ", " +
		"SignedHeaders=" + signedHeaders + ", " + "Signature=" +
		getSignature(getStringToSign(getSpliceData(payload), tm), tm)
	//fmt.Printf("auth result: %s\n\n", auth)
	return auth
}

// 哈希256摘要
func hmac256Secret(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	sha := mac.Sum(nil)
	return sha
}

// 哈希1摘要
func hmac1Secret(key []byte, content string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	sha := mac.Sum(nil)
	return sha
}

// 发送请求
func request(ocr OcrMethodData, r *http.Request) (*http.Response, error) {
	params := make(map[string]string, 2)
	if ocr.Name == "IDCardOCR" {
		params["CardSide"] = "FRONT" //FRONT为身份证有照片的一面（正面） BACK为身份证有国徽的一面（反面）
		//params["Config"] = ""   //可选字段，根据需要选择是否请求对应字段。目前包含的字段为： CropIdCard-身份证照片裁剪， CropPortrait-人像照片裁剪， CopyWarn-复印件告警， ReshootWarn-翻拍告警。
	}
	imageUrl := r.FormValue("imageUrl")
	if imageUrl == "" {
		base64Data := getBase64File(r)
		if base64Data == "" {
			return nil, fmt.Errorf("image error: %s", "参数错误，请上传图片，或图片URL")
		}
		params["ImageBase64"] = base64Data
	} else {
		params["ImageUrl"] = imageUrl
	}
	client := &http.Client{}
	js, _ := json.Marshal(params)
	fmt.Printf("params: %s\n\n", string(js))
	req, _ := http.NewRequest(http.MethodPost, RequestUrl, bytes.NewReader(js))
	req.Header = getHeader(js, ocr.Name, ocr.Version)
	return client.Do(req)
}

// 发送请求
func requestOld(requestUrl string, r *http.Request) (*http.Response, error) {
	params := make(map[string]string, 2)
	params["appid"] = conf.AppId
	imageUrl := r.FormValue("imageUrl")
	if imageUrl == "" {
		base64Data := getBase64File(r)
		if base64Data == "" {
			return nil, fmt.Errorf("image error: %s", "参数错误，请上传图片，或图片URL")
		}
		params["image"] = base64Data
	} else {
		params["url"] = imageUrl
	}
	client := &http.Client{}
	js, _ := json.Marshal(params)
	//fmt.Printf("params: %s\n\n", string(js))
	req, _ := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(js))
	req.Header = http.Header{}
	req.Header.Add("authorization", getSign())
	req.Header.Add("content-type", "application/json")
	req.Header.Add("host", "recognition.image.myqcloud.com")
	return client.Do(req)
}

// 配置文件
type Config struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	AppId     string `json:"app_id"`
}

// 响应请求数据结构
type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 结构化
func (response ResponseData) toJson() ([]byte, error) {
	return json.Marshal(response)
}

// 响应请求
func (response ResponseData) response(w http.ResponseWriter) {
	if js, err := response.toJson(); err != nil {
		fmt.Printf("reponse toJson error: %s\n", err.Error())
	} else if _, err := w.Write(js); err != nil {
		fmt.Printf("reponse errors: %s\n", err.Error())
	}
}

// 文件转 Base64
func getBase64File(r *http.Request) string {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Print(err)
		return ""
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)
	data := base64.StdEncoding.EncodeToString(content)
	timeStr := time.Now().Format("2006-01-02") + string(os.PathSeparator)
	fmt.Printf("%s ，上传成功 => %v\n", timeStr, handler.Header)
	return data
}
