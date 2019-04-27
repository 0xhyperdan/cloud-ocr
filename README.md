# 腾讯云开放平台 OCR 识别能力库接口

[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

* 迁移到腾讯云平台，识别速度更快，采用 URL 的传参数形式，同时支持上传文件的方式。

# API

1. 营业执照识别 `http://127.0.0.1:6663/ocr/license` 

2. 身份证识别 `http://127.0.0.1:6663/ocr/identity` 

3. 通用识别 `http://127.0.0.1:6663/ocr/general` 

4. 手写体识别 `http://127.0.0.1:6663/ocr/write` 

5. 开票信息识别（base 通用识别） `http://127.0.0.1:6663/ocr/invoice`
 
6. 银行卡识别 `http://127.0.0.1:6663/ocr/bank` 

# Docker

* build `docker build --rm -t ocr:v1.0 .`

* run `docker run -d -p 6663:6663 --name ybr/ocr ocr:v1.0`

# Make

1. modify `common.go` init func `os.Open("/go/src/tencent-ocr/config.json")` to `os.Open("config
.json")`

2. run `make` then `./admin.sh start`