package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

//初始化etcd
var (
	cli *clientv3.Client
)

//LogEntry ...
type LogEntry struct {
	Path  string `json:"path"`  //日志存放的路径
	Topic string `json:"topic"` //日志要发往kafka的哪个topic
}

//Init ...
func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

//GetConf 从etcd中根据key获取配置项
func GetConf(key string) (LogEntryConf []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}

	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, &LogEntryConf)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed,err:%v\n", err)
			return
		}
	}
	return
}

//WatchConf ...
func WatchConf(Key string, newConfCh chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), Key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v Key:%v value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed err:%v\n", err)
					continue
				}

			}

			newConfCh <- newConf
		}
	}
}
