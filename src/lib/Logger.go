package lib

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var Logger = &FileLogger{}

// Log levels to control the logging output.
const (
	LevelDebug = iota
	LevelAccess
	LevelWarn
	LevelError
)

//日志类
type FileLogger struct {
	loggerMap map[string]*log.Logger
	fdMap     map[string]*os.File
	rootPath  string
	logLevel  int
	mu        sync.Mutex
}

func (this *FileLogger) Init(rootPath string, logLevel int) {
	this.rootPath = rootPath
	this.loggerMap = make(map[string]*log.Logger)
	this.fdMap = make(map[string]*os.File)
	this.logLevel = logLevel
}

func (this *FileLogger) getLogger(logName string) (*log.Logger, error) {
	nowDate := time.Now().Format("20060102")
	//当前的日志文件名
	filePath := this.rootPath + "/" + logName + ".log." + nowDate
	fd, ok := this.fdMap[logName]
	//如果日志文件没有打开，或者日志名已经变了，就重新打开另外的日志文件
	if !ok || (fd != nil && fd.Name() != filePath) {
		this.mu.Lock()
		defer this.mu.Unlock()
		fd, ok = this.fdMap[logName]
		//双重判断，减少重复操作
		if !ok || (fd != nil && fd.Name() != filePath) {
			if fd != nil {
				fd.Close()
			}
			fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
			//fmt.Println("fd name:",  fd.Name())
			if err != nil {
				return nil, err
			}
			//创建文件的时候指定777权限不管用，所有只能在显示chmod, 奇了个怪..
			fd.Chmod(0777)
			this.loggerMap[logName] = log.New(fd, "", 0)
			this.fdMap[logName] = fd
			fmt.Println("new logger:", filePath)
		}
	}
	retLogger, ok := this.loggerMap[logName]
	return retLogger, nil
}

func (this *FileLogger) writeLog(logName string, v ...interface{}) {
	logger, err := this.getLogger(logName)
	if err != nil {
		fmt.Println("log failed", err)
		return
	}
	msgstr := ""
	for _, msg := range v {
		if msg1, ok := msg.(map[string]interface{}); ok {
			//map每次输出的顺序是随机的，以下保证每次输出的顺序一致，如果map比较大，可能有一定性能损耗
			var keys []string
			for k := range msg1 {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				msgstr = msgstr + fmt.Sprintf("%s=%+v,", k, msg1[k])
			}
		} else {
			msgstr = msgstr + fmt.Sprintf("%+v,", msg)
		}
	}
	msgstr = strings.TrimRight(msgstr, ",")
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("time[%s] msg[%s]\n", timeNow, msgstr)
}

func (this *FileLogger) Debug(logName string, v ...interface{}) {
	if this.logLevel > LevelDebug {
		return
	}
	logName = logName + "_debug"
	this.writeLog(logName, v...)
}

func (this *FileLogger) Access(logName string, v ...interface{}) {
	if this.logLevel > LevelAccess {
		return
	}
	logName = logName + "_access"
	this.writeLog(logName, v...)
}

func (this *FileLogger) Warn(logName string, v ...interface{}) {
	if this.logLevel > LevelWarn {
		return
	}
	logName = logName + "_warn"
	this.writeLog(logName, v...)
}

func (this *FileLogger) Error(logName string, v ...interface{}) {
	if this.logLevel > LevelError {
		return
	}
	logName = logName + "_error"
	this.writeLog(logName, v...)
}
