package lib

import (
)
//错误类
type Err struct {
    ErrorNo  int
    ErrorMsg string
}

func (this *Err) Error() string{
	return this.ErrorMsg
}

//需要传递原始错误消息的使用这个类
type ErrorInfo struct {
    Err
    Raw error
}

func (this *ErrorInfo) Error() string{
	rawError := this.Raw.Error()
	errstring := this.ErrorMsg + ":" + rawError
	return errstring
}


//在这里添加错误码
var(
ERR_SUC 	= Err{ErrorNo:0, ErrorMsg:"OK"}
ERR_SYSTEM  = Err{ErrorNo:100, ErrorMsg:"system error"}
ERR_INPUT  	= Err{ErrorNo:101, ErrorMsg:"input param error"}
ERR_OUTPUT  = Err{ErrorNo:102, ErrorMsg:"output error"}
ERR_MONGODB = Err{ErrorNo:102, ErrorMsg:"mongodb error"}
)