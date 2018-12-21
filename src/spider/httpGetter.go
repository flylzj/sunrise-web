package spider

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"math/rand"
	"model"
	"net/http"
	"sort"
	"strings"
	"time"
)

func GetJsonData(url string, method string, headers map[string]string, body string) (*simplejson.Json, error){
	//fmt.Println(url)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 4.4.2; HUAWEI MLA-AL10 Build/HUAWEIMLA-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 Html5Plus/1.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		model.Error.Println("http error", err.Error())
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		model.Error.Println("io error", err.Error())
		return nil, err
	}

	jsonDate, err := simplejson.NewJson(responseBody)
	if err != nil {
		model.Error.Println("json error", err.Error())
		return nil, err
	}
	return jsonDate, nil
}

func GetToken() string{
	url := "http://srmemberapp.srgow.com/sys/token"
	appsecret := "e1d0b361201e4324b37c968fb71f0d3c"
	appid := "sunrise_member"
	nonce := fmt.Sprintf("%d", rand.Intn(9000) + 1001)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	_array := []string{appsecret, nonce, timestamp}
	sort.Strings(_array)
	_tmp := strings.Join(_array, "")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(_tmp))
	cipherStr := md5Ctx.Sum(nil)
	signature := hex.EncodeToString(cipherStr)
	data := fmt.Sprintf("{\"appid\": \"%s\", \"appsecret\": \"%s\", \"timestamp\": \"%s\", \"signature\": \"%s\", \"nonce\": \"%s\"}",
		appid, appsecret, timestamp, strings.ToUpper(signature), nonce)
	jsonData, err := GetJsonData(url, "POST", map[string]string{"Content-Type": "application/json"}, data)
	if err != nil{
		model.Error.Println("get token error", err.Error())
		return ""
	}
	token, _ := jsonData.Get("data").Get("token").String()
	return 	token
}
