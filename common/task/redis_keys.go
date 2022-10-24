package task

import "taskpro/common/types"

const (
	eventLock  = "taskpro:lock:event"
	eventQueue = types.EventTaskRdsKey

	templateLock  = "taskpro:lock:template"
	templateQueue = "taskpro:queue:template"
)
