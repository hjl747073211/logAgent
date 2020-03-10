package taillog

//专门从日志文件收集日志内容的模块

import (
	"context"
	"fmt"
	"logAgent/kafka"

	"github.com/hpcloud/tail"
)

//TailTask 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	//为了实现退出t.run()
	ctx        context.Context
	cancelFunc context.CancelFunc
}

//NewTailTask ...
func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init()
	return
}

func (t *TailTask) init() {
	var err error
	t.instance, err = tail.TailFile(t.path, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)

	}
	go t.run()
}

func (t *TailTask) run() {
	for {
		select {
		//收到ctx的信号退出
		case <-t.ctx.Done():
			fmt.Printf("tail task 退出了：%s_%s\n", t.path, t.topic)
			return
		//1.一行一行读取日志
		case line := <-t.instance.Lines:
			//2.发送到kafka
			kafka.SendToChan(t.topic, line.Text)
			// fmt.Println(line.Text)

		}
	}
}
