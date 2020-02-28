package main

//入口程序
import (
	"fmt"
	"logAgent/kafka"
	"logAgent/taillog"
	"time"
)

//taillog是从日志中读取内容
//sarama是把日志内容放入kafka

func main() {
	//1.初始化kafka的链接
	err := kafka.Init([]string{"127.0.0.1:9092"})
	if err != nil {
		fmt.Printf("init kafka failed,err:%v", err)
		return
	}
	fmt.Println("kafka init success")
	//打开日志文件准备读取日志
	err = taillog.Init("./my.log") //传入日志文件的路径
	if err != nil {
		fmt.Printf("init tail failed,err:%v", err)
		return
	}
	fmt.Println("taillog init success")
	run()
}

func run() {
	//1.读取日志
	for {
		select {
		case line := <-taillog.ReadChan():
			//2.发送到kafka
			kafka.SendToKafka("web_log", line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
