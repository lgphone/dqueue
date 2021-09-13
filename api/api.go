package api

import (
	"dqueue/deley_queue"
	"dqueue/job"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ResponseBody 响应Body格式
type ResponseBody struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func PushJob(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(404)
		resp.Write([]byte("404 page not found\n"))
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	requestBody := &job.Job{}
	if err := readBody(req, requestBody); err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	requestBody.Task = strings.TrimSpace(requestBody.Task)
	if requestBody.Task == "" {
		resp.Write(ParamsErrResp("task 参数非法！"))
		return
	}

	newJob := job.NewJobWithId(requestBody.Task, requestBody.Delay, requestBody.TTL, requestBody.Args...)
	if err := deley_queue.PushJob(newJob); err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	resp.Write(SuccessResp("提交成功!"))
}

func DeleteJob(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(404)
		resp.Write([]byte("404 page not found\n"))
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	requestBody := struct {
		ID string `json:"id"`
	}{}
	if err := readBody(req, &requestBody); err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	if requestBody.ID == "" {
		resp.Write(ParamsErrResp("id 参数非法！"))
		return
	}
	if err := job.DeleteJob(requestBody.ID); err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	resp.Write(SuccessResp("删除成功!"))
}

func GetJob(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(404)
		resp.Write([]byte("404 page not found\n"))
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	requestBody := struct {
		ID string `json:"id"`
	}{}
	if err := readBody(req, &requestBody); err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	if requestBody.ID == "" {
		resp.Write(ParamsErrResp("id 参数非法！"))
		return
	}
	jobInfo, err := job.GetJob(requestBody.ID)
	if err != nil {
		resp.Write(UnKnowErrResp(err.Error()))
		return
	}
	resp.Write(SuccessResp(jobInfo))
}

func readBody(req *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("读取body错误: %s", err.Error()))
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return errors.New(fmt.Sprintf("解析json失败: %s", err.Error()))
	}

	return nil
}

func SuccessResp(data interface{}) []byte {
	return responseBody(1000, nil, data)
}

func ParamsErrResp(msg string) []byte {
	return responseBody(1001, msg, nil)
}

func AuthErrResp(msg string) []byte {
	return responseBody(1002, msg, nil)
}

func LogicErrResp(msg string) []byte {
	return responseBody(1003, msg, nil)
}

func UnKnowErrResp(msg string) []byte {
	return responseBody(1005, msg, nil)
}

func responseBody(code int, msg, data interface{}) []byte {
	body := &ResponseBody{}
	body.Code = code
	body.Message = msg
	body.Data = data

	bytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("生成response body,转换json失败#%s", err.Error())
		return []byte(`{"code":"5000", "message": "生成响应body异常", "data":[]}`)
	}
	return bytes
}
