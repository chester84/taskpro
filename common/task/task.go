package task

import (
	"fmt"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"github.com/chester84/libtools"
)

//! see: https://github.com/adonovan/gopl.io/blob/master/ch8/du4/main.go

type TaskWork0 interface {
	Start()
	Cancel()
}

type Work0Item struct {
	Describe string
	Task     TaskWork0
}

var taskWork0Map = map[string]Work0Item{}

func TaskWork0Map() map[string]Work0Item {
	return taskWork0Map
}

func Register(name, describe string, task TaskWork0) {
	if task == nil {
		panic("[task->Register] register task is nil")
	}

	if _, ok := taskWork0Map[name]; ok {
		panic("[task->Register]: register called twice for task " + name)
	}

	taskWork0Map[name] = Work0Item{
		Describe: describe,
		Task:     task,
	}
}

//!+1
var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var mutex sync.Mutex
var currentDatas map[string]interface{} = make(map[string]interface{})

func addCurrentData(key string, value interface{}) {
	mutex.Lock()

	currentDatas[key] = value

	mutex.Unlock()
}

func removeCurrentData(key string) {
	mutex.Lock()

	delete(currentDatas, key)

	mutex.Unlock()
}

func GetCurrentData() string {
	mutex.Lock()

	str := fmt.Sprintf("%v", currentDatas)

	mutex.Unlock()

	return str
}

var lastTimetag int64

func TaskHeartBeat(coon redis.Conn, lockKey string) {
	nowT := libtools.GetUnixMillis()

	if lastTimetag == 0 {
		lastTimetag = nowT
		return
	}

	if nowT-lastTimetag < int64(1000*60*10) {
		return
	}

	lastTimetag = nowT

	_, err := coon.Do("SET", lockKey, nowT)
	if err != nil {
		logs.Error("[TaskHeartBeat] set key error time:%d, err:%v", nowT, err)
	}
}

//!-1
