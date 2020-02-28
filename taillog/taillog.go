package taillog

//专门从日志文件收集日志内容的模块

import (
	"fmt"

	"github.com/hpcloud/tail"
)

var (
	tailObj *tail.Tail
)

//Init ...
func Init(filename string) (err error) {
	tailObj, err = tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	return
}

//ReadChan ...
func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}
