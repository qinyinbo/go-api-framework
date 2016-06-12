package lib

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func NewHttpServer(addr string, port int, timout int, pprof bool) *HttpApiServer {
	ret := &HttpApiServer{
		HttpAddr: addr,
		HttpPort: port,
		Timeout:  timout,
		hanlder:  &httpApiHandler{routMap: make(map[string]map[string]reflect.Type), enablePprof: pprof},
	}
	return ret
}

//http服务监听,路由
type HttpApiServer struct {
	HttpAddr string
	HttpPort int
	Timeout  int
	hanlder  *httpApiHandler
}

func (this *HttpApiServer) AddController(c interface{}) {
	this.hanlder.addController(c)
}

func (this *HttpApiServer) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := fmt.Sprintf("%s:%d", this.HttpAddr, this.HttpPort)
	s := &http.Server{
		Addr:         addr,
		Handler:      this.hanlder,
		ReadTimeout:  time.Duration(this.Timeout) * time.Millisecond,
		WriteTimeout: time.Duration(this.Timeout) * time.Millisecond,
	}
	fmt.Println("HttpApiServer Listen: ", addr)
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

//controller中以此结尾的方法会参与路由
const METHOD_EXPORT_TAG = "Action"

type httpApiHandler struct {
	routMap     map[string]map[string]reflect.Type //key:controller: {key:method value:reflect.type}
	enablePprof bool
}

func (this *httpApiHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ServeHTTP: ", err)
			http.Error(rw, fmt.Sprintln(err), http.StatusInternalServerError)
		}
	}()
	rw.Header().Set("Server", "GoServer")
	r.ParseMultipartForm(1024)

	var cname, mname string
	//如果开启了pprof, 相关请求走DefaultServeMux
	if this.enablePprof && strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		cname = "Monitor"
		mname = "Pprof"
	} else if strings.HasPrefix(r.URL.Path, "/status.go") { //用于lvs监控
		cname = "Monitor"
		mname = "Status"
	} else {
		//处理路由
		router := strings.Trim(r.URL.Path, "/")
		routers := strings.Split(router, "/")
		pathCounts := len(routers)
		/*
			fmt.Println(pathCounts)
			for i := 0; i < len(routers); i++ {
				fmt.Println(fmt.Sprintf("%d", i) + ":" + routers[i])
			}
		*/
		if router == "" || routers == nil || pathCounts < 2 || routers[0] == "" || routers[1] == "" {
			http.NotFound(rw, r)
			return
		}

		//不处理favicon.ico路由
		if router == "/favicon.ico" {
			return
		}

		if pathCounts == 2 {

			cname = "test" + "_" + strings.Title(routers[0])
			mname = strings.Title(routers[1])

		} else if pathCounts == 3 {

			cname = strings.ToLower(routers[0]) + "_" + strings.Title(routers[1])
			mname = strings.Title(routers[2])

		} else {

			http.NotFound(rw, r)
			return

		}
	}
	mname = mname + METHOD_EXPORT_TAG
	canhandler := false
	var contollerType reflect.Type
	if cname != "" && mname != "" {
		if methodMap, ok := this.routMap[cname]; ok {
			if contollerType, ok = methodMap[mname]; ok {
				canhandler = true
			}
		}
	}

	if !canhandler {
		http.NotFound(rw, r)
		return
	}

	vc := reflect.New(contollerType)
	var in []reflect.Value
	var method reflect.Value
	defer func() {
		if err := recover(); err != nil {
			in = []reflect.Value{reflect.ValueOf(err)}
			method := vc.MethodByName("OutputError")
			method.Call(in)
		}
	}()
	in = make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(rw)
	in[1] = reflect.ValueOf(r)
	method = vc.MethodByName("Init")
	method.Call(in)

	in = make([]reflect.Value, 0)
	method = vc.MethodByName(mname)
	method.Call(in)
}

func (this *httpApiHandler) addController(c interface{}) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	pkgPath := ct.PkgPath()
	pkgPaths := strings.Split(pkgPath, "/")
	pathCounts := len(pkgPaths)
	pkgName := pkgPaths[pathCounts-1]
	firstParam := strings.ToLower(pkgName) + "_" + strings.TrimSuffix(ct.Name(), "Controller")
	if _, ok := this.routMap[firstParam]; ok {
		return
	} else {
		this.routMap[firstParam] = make(map[string]reflect.Type)
	}
	var mname string
	for i := 0; i < rt.NumMethod(); i++ {
		mname = rt.Method(i).Name
		if strings.HasSuffix(mname, METHOD_EXPORT_TAG) {
			this.routMap[firstParam][rt.Method(i).Name] = ct
		}
	}
}
