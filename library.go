package library

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func HttpClientPost(timeout time.Duration, connection_timeout time.Duration, gurl string, data map[string]interface{}, headers []interface{}) (string, error) {

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	var clusterinfo = url.Values{}
	for k, v := range data {
		if k == "fix_n" {
			k = "n"
		}
		clusterinfo.Add(k, v.(string))
	}
	dataStr := clusterinfo.Encode()

	request, err := http.NewRequest("GET", gurl, strings.NewReader(dataStr))
	if err != nil {
		return "", err
	}
	/*
		cookie := &http.Cookie{Name: "userId", Value: strconv.Itoa(12345)}
		request.AddCookie(cookie)
	*/
	for _, v := range headers {
		v := v.(map[string]interface{})
		if v["name"].(string) == "host" {
			request.Host = v["value"].(string) //草
		} else {
			request.Header.Set(strings.Title(v["name"].(string)), v["value"].(string))
		}
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
	}
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(r), nil

}
func HttpClientGet(timeout time.Duration, connection_timeout time.Duration, gurl string, data map[string]interface{}, headers []interface{}) (string, error) {

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connection_timeout)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	var clusterinfo = url.Values{}
	for k, v := range data {
		if k == "fix_n" {
			k = "n"
		}
		clusterinfo.Add(k, v.(string))
	}
	dataStr := clusterinfo.Encode()
	gurl = gurl + "?" + dataStr

	request, err := http.NewRequest("GET", gurl, nil)
	if err != nil {
		return "", err
	}
	/*
		cookie := &http.Cookie{Name: "userId", Value: strconv.Itoa(12345)}
		request.AddCookie(cookie)

	*/
	fmt.Println(gurl)
	for _, v := range headers {
		v := v.(map[string]interface{})
		if v["name"].(string) == "host" {
			request.Host = v["value"].(string) //草
		} else {
			request.Header.Set(strings.Title(v["name"].(string)), v["value"].(string))
		}

	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
	}
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(r), nil

}
func GetIdc() string {
	host_name, err := os.Hostname()
	idc := ""
	if err != nil {
		idc = "default"
	}
	hosts := strings.Split(host_name, ".")
	if len(hosts) >= 3 {
		idc = hosts[2]
	}
	idc = "default"
	return idc
}
func GetLocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	var ip string = ""
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			ip = ipnet.IP.String()
			if ip != "127.0.0.1" {
			}
		}
	}
	return ip
}

func MapMerge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func SerializePhp(data map[string]interface{}) string {
	ret := fmt.Sprintf("a:%d:{", len(data))
	for key, value := range data {
		ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(key), key)
		if valuemap, ok := value.(map[string]interface{}); ok {
			ret = ret + SerializePhp(valuemap)
		} else {
			valuestr := value.(string)
			ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(valuestr), valuestr)
		}
	}
	ret = ret + "}"
	return ret
}
