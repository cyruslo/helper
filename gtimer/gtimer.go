package gtimer

import (
	"time"
	//"log"
	//	LOGGER "base/logger"
	//"git.huoys.com/common/gomsg"
	"sync"
)

var (
	bStartPolling = false
	recordTicks   int64
	timer stTimerList
)

type stTimer struct {
	reg      int64
	period   int64
	f        func()
	repeated bool
}

type stTimerList struct {
	mapTimers     map[int32]*stTimer
	lock *sync.RWMutex
}

func polling() {
//	defer gomsg.Recover()
	time.AfterFunc(time.Second*1, polling)

	timer.lock.Lock()
	//LOGGER.Info("base polling")
	for id, t := range timer.mapTimers {
		//LOGGER.Info("polling recordTicks: %d,%d,%d,%d", recordTicks, id, timer.period, timer.reg)
		if recordTicks > t.reg && (recordTicks-t.reg)%t.period == 0 {
			//	 LOGGER.Warning("polling recordTicks: %d,%d,%d,%d", recordTicks, id, timer.period, timer.reg)
			t.f()
			if t.repeated == false {
				//	LOGGER.Info("===delete:id:%d.", id)
				delete(timer.mapTimers, id)
			}
		}
	}
	timer.lock.Unlock()
	recordTicks++
}

func ClearTimer() {
	timer.lock.Lock()
	timer.mapTimers = make(map[int32]*stTimer)
	timer.mapTimers = nil
	timer.lock.Unlock()
}

func RemoveTimer(id int32) {
	timer.lock.Lock()
	delete(timer.mapTimers, id)
	timer.lock.Unlock()
}

func GetTimer(id int32) bool {
	if _, ok := timer.mapTimers[id]; ok {
		return true
	}
	return false
}

func SetTimer(id int32, period int64, start bool, repeated bool, f func()) {
	if bStartPolling == false {
	//	a := new(stTimerList)
		timer.lock = new(sync.RWMutex)
		timer.mapTimers = make(map[int32]*stTimer)
		bStartPolling = true
		recordTicks = 0
		time.AfterFunc(time.Second*1, polling)
	}

	//	LOGGER.Info("========= len(mapTimers):%d,period(%d).", len(mapTimers),period)
	if period <= 0 {
		return
	}

	timer.lock.Lock()
	//已经存在的定时器
	if _, ok := timer.mapTimers[id]; ok {
		timer.mapTimers[id].period = period
		timer.mapTimers[id].repeated = repeated
		timer.mapTimers[id].f = f
		timer.mapTimers[id].reg = recordTicks
	} else {
		t := &stTimer{
			reg:      recordTicks,
			period:   period,
			f:        f,
			repeated: repeated,
		}
		timer.mapTimers[id] = t
	}
	timer.lock.Unlock()

	//马上执行
	if start {
		f()
	}
}
