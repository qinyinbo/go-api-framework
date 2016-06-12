package models

import (
	"encoding/json"
	"errors"
	"go-api-framework/src/gophpserialize"
	"go-api-framework/src/lib"
	"time"
)

type SearhAppNovel struct {
	retry_max_num int
	retry_times   int
}

func NewSearhAppNovel() *SearhAppNovel {
	ret := &SearhAppNovel{3, 0}
	return ret
}

func (this *SearhAppNovel) AppNovelContentByCid(input_data map[string]interface{}) map[string]interface{} {

	http_method, time_out, connection_timeout, gurl, params, headers := this.processParams()
	params = lib.MapMerge(params, input_data)
	param_conf := map[string]interface{}{
		"_method": "xxxx",
		"_model":  "yyy",
		"client":  params["src"].(string),
	}
	params = lib.MapMerge(params, param_conf)

	data, err := this.httpNet(http_method, time_out, connection_timeout, gurl, params, headers)
	if err != nil && this.retry_max_num > 0 && this.retry_times > 0 {
		data, err = this.retry(http_method, time_out, connection_timeout, gurl, params, headers)
	}
	/*
		data, err := this.httpNet(http_method, time_out, connection_timeout, "http://42.96.174.93:2012/test_get", params, headers)
		if err != nil && this.retry_max_num > 0 && this.retry_times > 0 {
			data, err = this.retry(http_method, time_out, connection_timeout, "http://42.96.174.93:2012/test_get", params, headers)
		}
	*/
	if err != nil {
		panic(err)
	}
	final_data := this.dataUnmarshal(data)
	return final_data
}

func (this *SearhAppNovel) dataUnmarshal(data string) map[string]interface{} {

	final_data := make(map[string]interface{})
	decodeData := gophpserialize.Unmarshal([]byte(data))
	if decodeData != nil {
		final_data["data"] = decodeData
	} else {
		var json_interface interface{}
		err := json.Unmarshal([]byte(data), &json_interface)
		if err != nil {
			final_data["data"] = data
		}
		final_data["data"] = json_interface
	}
	return final_data

}

func (this *SearhAppNovel) processParams() (string, time.Duration, time.Duration, string, map[string]interface{}, []interface{}) {
	httpConfig := lib.Conf.GetAllHttpConfig()
	idc := lib.GetIdc()
	hkey := "yyyyyyyyyyyy.host_port." + idc + ".hostname"
	pkey := "yyyyyyyyyyy.host_port." + idc + ".port"
	http_protocol := ""
	path := ""
	url := ""
	http_method := ""
	total_timeout := time.Millisecond
	connection_timeout := time.Millisecond

	hostname, err := httpConfig.String(hkey)
	port, err := httpConfig.String(pkey)

	http_conf_default_fields, err := httpConfig.Map("xxxxxxxxxxxxxxxx")
	zeus_data_service, err := httpConfig.Map("yyyyyyyyyy")
	mergeConfig := lib.MapMerge(http_conf_default_fields, zeus_data_service)

	if v := mergeConfig["retry_times"]; v != nil {
		this.retry_times = mergeConfig["retry_times"].(int)
	}
	if v := mergeConfig["http_protocol"]; v != nil {
		http_protocol = mergeConfig["http_protocol"].(string)
	}
	if v := mergeConfig["path"]; v != nil {
		path = mergeConfig["path"].(string)
	}
	if v := mergeConfig["http_method"]; v != nil {
		http_method = mergeConfig["http_method"].(string)
	}
	if v := mergeConfig["total_timeout"]; v != nil {
		total_timeout = time.Duration(mergeConfig["total_timeout"].(int)) * time.Millisecond
	}
	if v := mergeConfig["connection_timeout"]; v != nil {
		connection_timeout = time.Duration(mergeConfig["connection_timeout"].(int)) * time.Millisecond
	}
	if err != nil || http_protocol != "https" {
		url = "http://" + hostname + ":" + port + path
	} else {
		url = "https://" + hostname + ":" + path
	}
	params := mergeConfig["params"].(map[string]interface{})
	headers := mergeConfig["header"].([]interface{})
	return http_method, total_timeout, connection_timeout, url, params, headers
}
func (this *SearhAppNovel) httpNet(method string, timeout time.Duration, connection_timeout time.Duration, gurl string, params map[string]interface{}, headers []interface{}) (string, error) {
	data := ""
	err := errors.New("")
	if method == "get" {
		data, err = lib.HttpClientGet(timeout, connection_timeout, gurl, params, headers)
	} else {
		data, err = lib.HttpClientPost(timeout, connection_timeout, gurl, params, headers)
	}
	return data, err
}
func (this *SearhAppNovel) retry(method string, timeout time.Duration, connection_timeout time.Duration, gurl string, params map[string]interface{}, headers []interface{}) (string, error) {
	data := ""
	err := errors.New("")
	for i := 1; i < this.retry_max_num; i++ {
		data, err = this.httpNet(method, timeout, connection_timeout, gurl, params, headers)
		if i > this.retry_max_num || i > this.retry_times || err == nil {
			break
		}
	}
	return data, err
}
