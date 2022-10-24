package task

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"github.com/chester84/libtools"
	"taskpro/common/lib/redis/storage"
)

type Template struct {
}

func init() {
	Register("template", "task template", &Template{})
}

func (r *Template) Cancel() {
	logs.Informational("[task->Template] cancel task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	lockKey := templateLock
	_, _ = storageClient.Do("DEL", lockKey)
}

func (r *Template) Start() {
	logs.Info("[task->Template] start task")

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	// +1 Distributed Locks {{{
	lockKey := templateLock
	lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
	if err != nil || lock == nil {
		logs.Error("[task->Template] process is working, so, I will exit, err: %v", err)
		close(done)
		return
	}

	queueName := templateQueue
	for {
		if cancelled() {
			logs.Info("[task->Template] receive exit cmd.")
			break
		}

		TaskHeartBeat(storageClient, lockKey)

		// 1. create task queue
		logs.Info("[task->Template] produceTemplateQueue")
		queueLen, err := redis.Int64(storageClient.Do("LLEN", queueName))
		logs.Debug("queueLen:", queueLen, ", err:", err)

		if err != nil {
			logs.Error("[task->Template] redis get exception `LLEN %s` error. err: %v", queueName, err)
			break
		} else if queueLen == 0 {
			for i := 0; i < 100; i++ {
				_, err = storageClient.Do("LPUSH", queueName, libtools.GetUnixMillis()+int64(i))
				if err != nil {
					logs.Error("[task->Template] redis> LPUSH %s %d, err: %v", queueName, i, err)
				}
			}
		}

		var wg sync.WaitGroup
		//  increase goroutines according the situation
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go consumeTemplateQueue(&wg, i)
		}

		wg.Wait()
	}

	// -1 exit normally, release lock }}}
	_, err = storageClient.Do("DEL", lockKey)
	if err != nil {
		logs.Error("[task->Template] redis> DEL %s, err: %v", lockKey, err)
	}

	logs.Info("[task->Template] politeness exit.")
}

func consumeTemplateQueue(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	logs.Info("It will do consumeTemplateQueue, workerID:", workerID)

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	queueName := templateQueue

	for {
		if cancelled() {
			logs.Info("[consumeTemplateQueue] receive exit cmd, workID:", workerID)
			break
		}

		qValue, err := redis.String(storageClient.Do("RPOP", queueName))
		if err != nil {
			logs.Error("[consumeTemplateQueue] RPOP error workID:%d, err:%v", workerID, err)
		}

		// no datas for consuming
		if qValue == "" {
			logs.Info("[consumeTemplateQueue] no data for consume, so exit, workID:", workerID)
			break
		}

		cmdID, _ := libtools.Str2Int64(qValue)
		if cmdID == -111 {
			logs.Info("[consumeTemplateQueue] receive exit cmd, I will exit after jobs done. workID:", workerID, ", cmdID:", cmdID)
			close(done)
			break
		}

		// do the task
		addCurrentData(qValue, cmdID)
		handleTemplate(cmdID, workerID)
		removeCurrentData(qValue)
	}
}

func handleTemplate(dataID int64, workerID int) {
	if dataID <= 0 {
		logs.Warning("[handleTemplate] get zero id, workerID: %d", workerID)
		return
	}

	//TODO do the logics here
	time.Sleep(time.Second * 2)

	logs.Info("[handleTemplate] finish handel, workerID: %d", workerID)
}
