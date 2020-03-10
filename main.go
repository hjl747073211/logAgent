package main

//入口程序
import (
	"fmt"
	"logAgent/conf"
	"logAgent/etcd"
	"logAgent/kafka"
	"logAgent/taillog"
	"sync"
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
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed,err:%v", err)
		return
	}
	fmt.Println("kafka init success")

	//2.初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("etcd init success")

	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("etcd getconf failed, err:%v\n", err)
		return
	}
	fmt.Printf("etcd getconf success,%v\n", logEntryConf)

	//3.从etcd中获取配置项
	taillog.Init(logEntryConf)

	//开启watch，检测配置项的变化
	newConfChan := taillog.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key, newConfChan)
	wg.Wait()

}
