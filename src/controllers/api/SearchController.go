package api

import (
	"fmt"
	"go-api-framework/src/models"
)

type SearchController struct {
	BaseController
}

type Rep struct {
	Errno     string      `json:"errno"`
	Total     string      `json:"total"`
	End_state string      `json:"end_state"`
	Data      interface{} `json:"data"`
}

func (this *SearchController) AppNovelContentByCidAction() map[string]interface{} {
	inputData := this.InputData()
	searhAppNovel := models.NewSearhAppNovel()
	data := searhAppNovel.AppNovelContentByCid(inputData)
	ret := this.makeMap(0, data)
	this.Output(ret)
	return data
}

func (this *SearchController) makeMap(errno int, res map[string]interface{}) interface{} {

	total := 0
	count := 0
	start := 0
	ret_count := 0
	end := 0
	data := map[string]interface{}{}
	if d := res["total"]; d != nil {
		total = res["total"].(int)
	}
	if d := res["count"]; d != nil {
		count = res["count"].(int)
	}
	if d := res["start"]; d != nil {
		start = res["start"].(int)
	}
	if d := res["data"]; d != nil {
		data = res["data"].(map[string]interface{})
	}
	ret_count = len(data)
	if ret_count > count {
		count = ret_count
	}
	if total <= start+count {
		end = 1
	}

	rep := Rep{
		fmt.Sprint(errno),
		fmt.Sprint(total),
		fmt.Sprint(end),
		data,
	}
	return rep

}
