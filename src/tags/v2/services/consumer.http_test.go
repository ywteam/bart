package ydk

import (
	"fmt"
	"testing"
)
func TestHttpHttpConsumerServiceForm(t *testing.T) {
	// service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	service := HttpBinConsumer()
	resp, err := service.Form("/post", HttpRequestForm{
		UrlEncoded: map[string]string{
			"key":   "value",
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("TestHttpHttpConsumerServiceForm()", resp.Text())
}
func TestHttpHttpConsumerServiceDelete(t *testing.T) {
	// service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	service := HttpBinConsumer()
	resp, err := service.Delete("/delete", map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("TestHttpHttpConsumerServiceDelete()", resp.Text())
}
func TestHttpHttpConsumerServicePath(t *testing.T) {
	// service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	service := HttpBinConsumer()
	resp, err := service.Get("/anything", map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("TestHttpHttpConsumerServicePath()", resp.Text())
}

func TestHttpHttpConsumerServicePost(t *testing.T) {
	// service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	service := HttpBinConsumer()
	resp, err := service.Post("/post", map[string]string{
		"key":   "value",
		"index": "0",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Text()", resp.Text())
	var body = map[string]interface{}{}
	if err := resp.Json(&body); err != nil {
		t.Error(err)
	}
	// fmt.Println("resp.Json", body)
	if body["form"] == nil {
		t.Error("form not found")
	}
	if body["form"].(map[string]interface{})["key"] != "value" {
		t.Error("form.key not found")
	}
	// check header X-Amzn-Trace-Id
	if body["headers"] == nil {
		t.Error("headers not found")
	}
	if body["headers"].(map[string]interface{})["X-Amzn-Trace-Id"] == nil {
		t.Error("headers.X-Amzn-Trace-Id not found")
	}
	if AmznTraceId, ok := resp.Header("Access-Control-Allow-Credentials"); !ok {
		t.Error("headers.Access-Control-Allow-Credentials")
	} else {
		fmt.Println("X-Amzn-Trace-Id", AmznTraceId)
	}
	fmt.Println("TestHttpHttpConsumerServicePost", resp)
}
func TestHttpHttpConsumerServiceGet(t *testing.T) {
	// service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	service := HttpBinConsumer()
	resp, err := service.Get("/get", map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Error(err)
	}
	// fmt.Println("Text()", resp.Text())
	var body = map[string]interface{}{}
	if err := resp.Json(&body); err != nil {
		t.Error(err)
	}
	// fmt.Println("resp.Json", body)
	if body["args"] == nil {
		t.Error("args not found")
	}
	if body["args"].(map[string]interface{})["key"] != "value" {
		t.Error("args.key not found")
	}
	// check header X-Amzn-Trace-Id
	if body["headers"] == nil {
		t.Error("headers not found")
	}
	if body["headers"].(map[string]interface{})["X-Amzn-Trace-Id"] == nil {
		t.Error("headers.X-Amzn-Trace-Id not found")
	}
	if AmznTraceId, ok := resp.Header("Access-Control-Allow-Credentials"); !ok {
		t.Error("headers.Access-Control-Allow-Credentials")
	} else {
		fmt.Println("X-Amzn-Trace-Id", AmznTraceId)
	}
	fmt.Println("TestHttpHttpConsumerServiceGet	", resp)
}

func HttpBinConsumer() *HttpConsumerService {
	return NewHttpConsumerService(&HttpConsumerOptions{
		BaseUrl: "https://httpbin.org",
	})
}
