package api

import (
	"encoding/json"
	"fmt"
	"go-api-framework/src/lib"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strconv"
	"time"
)

//所有Controller类的父类
type BaseController struct {
	rw           http.ResponseWriter
	r            *http.Request
	startTime    time.Time
	OutputDirect bool //是否直接输出到http
}

func (this *BaseController) Init(rw http.ResponseWriter, r *http.Request) {
	this.startTime = time.Now()
	this.rw = rw
	this.r = r
	this.OutputDirect = true
}

//如果有异常，server会回调这个方法
func (this *BaseController) OutputError(err interface{}) {
	errno := lib.ERR_SYSTEM.ErrorNo
	errmsg := lib.ERR_SYSTEM.ErrorMsg

	switch errinfo := err.(type) {
	case string:
		errmsg = errinfo
	case lib.Err:
		errno = errinfo.ErrorNo
		errmsg = errinfo.Error()
	case lib.ErrorInfo:
		errno = errinfo.ErrorNo
		errmsg = errinfo.Error()
	case error:
		errmsg = errinfo.Error()
	default:
		errmsg = fmt.Sprint(errinfo)
	}
	ret := map[string]interface{}{
		"errno":  errno,
		"errmsg": errmsg,
		"data":   "",
	}

	this.toJson(ret)
	if _, ok := err.(error); ok {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		dump, err := httputil.DumpRequest(this.r, true)
		if err != nil {
			dump = []byte("failed get request")
		}
		lib.Logger.Error("qgo", string(dump), ret, string(buf))
		return
	}
	lib.Logger.Error("qgo", this.genLog(), ret)
	return

}

func (this *BaseController) GetString(key string, defaultValue string) string {
	ret := this.r.FormValue(key)
	if ret == "" {
		ret = defaultValue
	}
	return ret
}
func (this *BaseController) InputData() map[string]interface{} {
	this.r.ParseForm()
	ret := make(map[string]interface{})
	if this.r.Form == nil {
		return ret
	}
	for key, value := range this.r.Form {
		if len(value) > 1 {
			ret[key] = value
		} else {
			ret[key] = value[0]
		}
	}
	return ret
}

func (this *BaseController) GetStrings(key string) []string {
	if this.r.Form == nil {
		return []string{}
	}
	vs := this.r.Form[key]
	return vs
}

func (this *BaseController) GetInt(key string, defaultValue int64) int64 {
	ret, err := strconv.ParseInt(this.r.FormValue(key), 10, 64)
	if err != nil {
		ret = defaultValue
	}
	return ret
}

func (this *BaseController) GetBool(key string, defaultValue bool) bool {
	ret, err := strconv.ParseBool(this.r.FormValue(key))
	if err != nil {
		ret = defaultValue
	}
	return ret
}

func (this *BaseController) Output(data interface{}) {
	this.toJson(data)
	lib.Logger.Access("go-api-framework", this.genLog())
}

func (this *BaseController) OutputString(data string) {
	this.writeToWriter([]byte(data))
	lib.Logger.Access("go-api-framework", this.genLog())
}

func (this *BaseController) toJson(data interface{}) {
	//content, err := json.MarshalIndent(data, "", "  ")
	content, err := json.Marshal(data)
	if err != nil {
		panic(lib.ErrorInfo{lib.ERR_OUTPUT, err})
	}
	this.rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.writeToWriter(content)
}

//获取需要打印到日志的信息
func (this *BaseController) genLog() map[string]interface{} {
	ret := make(map[string]interface{})
	for key, value := range this.r.Form {
		if len(value) > 1 {
			ret[key] = value
		} else {
			ret[key] = value[0]
		}
	}
	//访问ip
	ret["user_ip"] = this.r.RemoteAddr
	//请求路径
	ret["req_uri"] = this.r.URL.Path
	//时间消耗
	ret["time_cost"] = this.cost()
	return ret
}

func (this *BaseController) cost() int64 {
	return time.Now().Sub(this.startTime).Nanoseconds() / 1000 / 1000
}

func (this *BaseController) writeToWriter(rb []byte) {
	//this.rw.Header().Set("Content-Length", strconv.Itoa(len(rb)))
	//this.rw.Header().Set("Content-Type", "text/html")
	if this.OutputDirect {
		this.rw.Write(rb)
	}
}
