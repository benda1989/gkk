package gkk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gkk/code"
	"gkk/expect"
)

type ResponseJson struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

func HttpBase(url, method string, bodys any) *[]byte {
	var req *http.Request
	var err error
	body := make([]byte, 0)
	switch bodys.(type) {
	case []byte:
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodys.([]byte)))
	case string:
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(bodys.(string))))
	default:
		jsonStr, e := json.Marshal(bodys)
		expect.PBM(e != nil, method+": "+url+": bodys参数错误")
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	expect.PEMC(err, "接口请求失败: "+url, code.REMOTE_ERROR)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	return &body
}

func Http(url, method string, bodys any) string {
	return string(*HttpBase(url, method, bodys))
}

func HttpRes(url, method string, bodys interface{}) *ResponseJson {
	var tt ResponseJson
	re := *HttpBase(url, method, bodys)
	expect.PEM(json.Unmarshal(re, &tt), url+": 请求返回解析失败")
	return &tt
}

func HttpRespond(url, method string, bodys any, errFuncs ...func(int)) *ResponseJson {
	var tt ResponseJson
	re := *HttpBase(url, method, bodys)
	expect.PEM(json.Unmarshal(re, &tt), url+": 请求返回解析失败")
	if len(errFuncs) > 0 {
		for _, f := range errFuncs {
			f(tt.Code)
		}
	} else {
		expect.PBMC(tt.Code != 0, tt.Msg, tt.Code)
	}
	return &tt
}

func HttpResData(url, method string, bodys, res any, errFuncs ...func(int)) {
	re := HttpRespond(url, method, bodys, errFuncs...)
	expect.PEM(json.Unmarshal(re.Data, res), url+":返回数据data解析失败")
}

func HttpGet(url string) []byte {
	resp, err := http.Get(url)
	expect.PEMC(err, "接口请求失败: "+url, code.REMOTE_ERROR)
	expect.PBM(resp.StatusCode == 404, "接口暂时无法访问")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
