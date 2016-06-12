package main

import (
	"flag"
	"fmt"
	"go-api-framework/src/controllers/api"
	"go-api-framework/src/lib"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func main() {
	defer func() {
		managePid(false)
		fmt.Println("Server Exit")
	}()
	envInit()
	managePid(true)
	startServer()
}

func envInit() {
	os.Chdir(path.Dir(os.Args[0]))
	confiPath := flag.String("f", "../conf/app.conf", "config file")
	flag.Parse()
	if *confiPath == "" {
		panic("config file missing")
	}
	lib.Conf.Init(*confiPath)
	lib.Conf.InitHttp("../conf/http.yaml")
	lib.Logger.Init(lib.Conf.Get("log_root"), lib.Conf.GetInt("log_level"))
}

func startServer() {
	server := lib.NewHttpServer("", lib.Conf.GetInt("http_port"),
		lib.Conf.GetInt("http_timeout"),
		lib.Conf.GetBool("pprof_enable"))
	server.AddController(&api.XxxController{})
	fmt.Println("Server Start")
	server.Run()
}

func managePid(create bool) {
	pidFile := lib.Conf.Get("app_pid_file")
	if create {
		pid := os.Getpid()
		pidString := strconv.Itoa(pid)
		ioutil.WriteFile(pidFile, []byte(pidString), 0777)
	} else {
		os.Remove(pidFile)
	}
}
