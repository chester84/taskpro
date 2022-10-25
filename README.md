# taskpro

## distributed locks

## redis list

## goroutine channel

## produce consume

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
