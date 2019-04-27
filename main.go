package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	if file, err := os.Open("config.json"); err != nil {
		//if file, err := os.Open("/go/src/tencent-ocr/config.json"); err != nil {
		fmt.Printf("open config file error: %s\n", err.Error())
	} else if err := json.NewDecoder(file).Decode(&conf); err != nil {
		fmt.Printf("decode config file error: %s\n", err.Error())
	} else {
		defer file.Close()
	}
}

func main() {
	http.HandleFunc("/ocr/license", license)
	http.HandleFunc("/ocr/write", write)
	http.HandleFunc("/ocr/general", general)
	http.HandleFunc("/ocr/invoice", invoice)
	http.HandleFunc("/ocr/identity", identity)
	http.HandleFunc("/ocr/bank", bank)
	log.Fatal(http.ListenAndServe(":6663", nil))
}
