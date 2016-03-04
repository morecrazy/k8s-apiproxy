package base

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"sort"
	"strings"
	"time"
	"zhanghan/slb/api"
)

type RestApi struct {
	Domain          string
	Port            int
	HttpMethod      string
	AccessKeyId     string
	AccessKeySecret string
	ApiName         string
	EncodeParams    map[string]string
}

func NewRestApi(domain string) *RestApi {
	var r = new(RestApi)
	if api.GetAppInfo() != nil {
		r.AccessKeyId = api.GetAppInfo().AccessKeyId
		r.AccessKeySecret = api.GetAppInfo().AccessKeySecret
	}

	//r.AccessKeyId = "KM0BWn5yGIjiYW3S"
	//r.AccessKeySecret = "VIECii4MYVv7QEVl5QDJbAxGH6nqH0"
	r.Domain = domain
	r.EncodeParams = make(map[string]string)
	return r
}
func (r *RestApi) GetRequestHeader() string {
	return ""
}

func (r *RestApi) GetResponse(authrize string, timeout string) string {
	apiNames := strings.Split(r.ApiName, ".")
	api.Info("apiNames = %s", apiNames)

	var timeStamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	fmt.Println(timeStamp)

	var values = url.Values{}
	values.Add("Format", "json")
	values.Add("Version", apiNames[4])
	api.Debug("Version : %s", apiNames[4])

	values.Add("Action", apiNames[3])
	api.Debug("Action : %s", apiNames[3])

	values.Add("AccessKeyId", r.AccessKeyId)
	values.Add("SignatureVersion", "1.0")
	values.Add("SignatureMethod", "HMAC-SHA1")
	out, _ := exec.Command("uuidgen").Output()

	values.Add("SignatureNonce", string(out))
	api.Debug("sinatureNonce : %s", string(out))

	values.Add("TimeStamp", timeStamp)
	values.Add("partner_id", "1.0")

	for k, v := range r.EncodeParams {
		values.Add(k, v)
	}

	//生成签名
	signature := Sign(r.AccessKeySecret, values)
	values.Add("Signature", signature)

	jsonobj, _ := HttpRequestClient(r.Domain, "POST", values)
	return jsonobj

}

/**
func (r RestApi) GetApplicationParams() map[string]string {
	var applicationParams = map[string]string{}
	//使用reflect获取struct变量类型和值
	val := reflect.ValueOf(&r).Elem()
	t := reflect.TypeOf(r)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i).String()
		typeField := t.Field(i).Name
		applicationParams[typeField] = string(valueField)
	}
	return applicationParams
}
**/

func Sign(accessKeySecret string, parameters map[string][]string) string {
	mk := make([]string, len(parameters))
	i := 0
	for k, _ := range parameters {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	var canonicalizedQueryString = ""
	for _, v := range mk {
		canonicalizedQueryString = canonicalizedQueryString + "&" + PercentEncode(v) + "=" + PercentEncode(parameters[v][0])
	}

	var stringToSign = "POST&%2F&" + PercentEncode(canonicalizedQueryString[1:])

	var signature = ComputeHmacSha1(stringToSign, accessKeySecret)
	return signature
}

func ComputeHmacSha1(message string, secret string) string {
	secret = secret + "&"
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func PercentEncode(encodeStr string) string {
	var encodedStr string = url.QueryEscape(encodeStr)
	encodedStr = strings.Replace(encodedStr, "+", "%20", -1)
	encodedStr = strings.Replace(encodedStr, "*", "%2A", -1)
	encodedStr = strings.Replace(encodedStr, "%7E", "~", -1)
	return encodedStr
}

func HttpRequestClient(requestPath string, method string, extraParams url.Values) (string, error) {
	parseRequestUrl, _ := url.Parse(requestPath)
	//更改URL Struct中的RawQuery为Encode后的Query string
	parseRequestUrl.RawQuery = extraParams.Encode()
	var requestUrl = parseRequestUrl.String()
	//var resp *http.Response
	var err error

	/**
	switch method {
	case "GET":
		resp, err = http.Get(requestUrl)
	case "POST":
		fmt.Println("ok")
		resp, err = http.PostForm(requestPath, extraParams)
	}
	**/

	client := &http.Client{}
	r, _ := http.NewRequest("POST", requestUrl, nil) // <-- URL-encoded payload
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Connection", "Keep-Alive")

	resp, _ := client.Do(r)

	//完成后关闭 Response
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)

	return string(bodyByte), err

}
