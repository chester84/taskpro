# taskpro

## distributed locks
```
# code segment

lockKey := templateLock
lock, err := storageClient.Do("SET", lockKey, libtools.GetUnixMillis(), "NX")
if err != nil || lock == nil {
  logs.Error("[task->Template] process is working, so, I will exit, err: %v", err)
  close(done)
  return
}

```

## redis list
```
# code segment

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
```

## goroutine channel
```
# code segment

var wg sync.WaitGroup
//  increase goroutines according the situation
for i := 0; i < 2; i++ {
    wg.Add(1)
    go consumeTemplateQueue(&wg, i)
}  
```

## close channels gracefully
```
# code segment

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
```

## usage
```
go run cmd/task/main.go

Options:
  -h	show usage and exit
  -n task-name
    	crontab, cli or backend task-name, need assign.
  -name task-name
    	crontab, cli or backend task-name, need assign.
  -v	show version and exit
```

## example
```
go run cmd/task/main.go -n=template

#logs as below

2022/10/25 10:20:31.801 [D] [main.go:115]  runmode:  dev , config: {"filename":"./logs/task.log","rotate":false,"separate":["emergency","alert","critical","error","warning","notice","info","debug"]}
2022/10/25 10:20:31.801 [I] [template.go:32]  [task->Template] start task
2022/10/25 10:20:31.801 [I] [template.go:56]  [task->Template] produceTemplateQueue
2022/10/25 10:20:31.801 [D] [template.go:58]  queueLen: 88 , err: <nil>
2022/10/25 10:20:31.801 [I] [template.go:94]  It will do consumeTemplateQueue, workerID: 1
2022/10/25 10:20:31.801 [I] [template.go:94]  It will do consumeTemplateQueue, workerID: 0

```
