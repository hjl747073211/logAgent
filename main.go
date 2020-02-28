package main

//入口程序
import (
	"fmt"
	conf "logAgent/Conf"
	"logAgent/kafka"
	"logAgent/taillog"
	"time"

	"gopkg.in/ini.v1"
)

//taillog是从日志中读取内容
//sarama是把日志内容放入kafka

var cfg = new(conf.AppConf)

func main() {
	//加载配置文件,映射到cfg
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("ini load failed,err:%v", err)
		return
	}

	//1.初始化kafka的链接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Printf("init kafka failed,err:%v", err)
		return
	}
	fmt.Println("kafka init success")
	//打开日志文件准备读取日志
	err = taillog.Init(cfg.TaillogConf.FileName) //传入日志文件的路径
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
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
