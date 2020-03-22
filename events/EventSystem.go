package events

import (
    "sync"
    "time"
)

type EventSystem struct {
    m            sync.Mutex
    minSleepMs   int
    currentList  []EventCallback
    delayList    []EventCallback
    isRunning    bool
    refMap       map[int]interface{}
    refMapM      sync.RWMutex
    OnUpdateFunc func()
}

type EventCallback struct {
    cb      func()
    delay   time.Time
    isDelay bool
}

//func init() {
//    DefaultEventSystem = NewEventSystem()
//    DefaultEventSystem.init()
//}
func NewEventSystem() *EventSystem {
    es := &EventSystem{}
    es.init()
    return es
}
//
//var (
//    DefaultEventSystem *EventSystem
//)

func (s *EventSystem) init() {
    s.minSleepMs = 33
    s.Reset()
}

func (s *EventSystem) Reset() {
    s.isRunning = true
    s.refMap = make(map[int]interface{})
}

func (s *EventSystem) OnMainThread(cb func()) {
    s.m.Lock()
    s.currentList = append(s.currentList, EventCallback{cb: cb, isDelay: false})
    s.m.Unlock()
}

func (s *EventSystem) OnMainThreadDelay(sec float64, cb func()) {
    s.m.Lock()
    s.delayList = append(s.delayList, EventCallback{cb: cb, isDelay: true, delay: time.Now().Add(time.Millisecond * time.Duration(sec*1000))})
    s.m.Unlock()
}

func (s *EventSystem) Update() {
    // on update func
    if s.OnUpdateFunc != nil {
        s.OnUpdateFunc()
    }
    // on update callback
    s.m.Lock()
    sz := len(s.currentList)
    // current list
    var runningList []EventCallback
    if sz > 0 {
        runningList = append(runningList, s.currentList...)
        s.currentList = []EventCallback{}
    }
    // delay list
    var delayAgain []EventCallback
    for _, v := range s.delayList {
        if v.delay.Sub(time.Now()) <= 0 {
            runningList = append(runningList, v)
        } else {
            delayAgain = append(delayAgain, v)
        }
    }
    s.delayList = delayAgain
    s.m.Unlock()

    if len(runningList) == 0 {
        time.Sleep(time.Duration(s.minSleepMs) * time.Millisecond)
        return
    } else {
        for _, cb := range runningList {
            cb.cb()
            if !s.IsRunning() {
                break
            }
        }
    }
}

func (s *EventSystem) Stop() {
    s.isRunning = false
}

func (s *EventSystem) IsRunning() bool {
    return s.isRunning
}

func (s *EventSystem) AddRef(ref int) {
    s.refMapM.Lock()
    defer s.refMapM.Unlock()
    s.refMap[ref] = ref
}

func (s *EventSystem) UnRef(ref int) {
    s.refMapM.Lock()
    defer s.refMapM.Unlock()
    delete(s.refMap, ref)
}

func (s *EventSystem) AllRef() []int {
    s.refMapM.RLock()
    defer s.refMapM.RUnlock()
    var rs []int
    for k, _ := range s.refMap {
        rs = append(rs, k)
    }
    return rs
}

func (s *EventSystem) SetFPS(fps int) {
    s.minSleepMs = int(1000.0 / float32(fps))
}
