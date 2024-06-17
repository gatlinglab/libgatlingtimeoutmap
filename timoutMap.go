package libgatlingtimeoutmap

import (
	"sync"
	"time"
)

type cGatlingTimeoutMapValue struct {
	endTime    int64
	key        interface{}
	data       interface{}
	mapTimeAdd int
}

type CTimeoutMap struct {
	mapOptions *OptionsGatlingTimeoutMap
	lockMain   *sync.Mutex
	mapData    map[interface{}]*cGatlingTimeoutMapValue
	mapTime    map[int64]*cGatlingTimeoutMapValue
	lastid     int64
}

func newmap(options *OptionsGatlingTimeoutMap) *CTimeoutMap {
	return &CTimeoutMap{
		mapOptions: options,
		lockMain:   new(sync.Mutex),
		mapData:    make(map[interface{}]*cGatlingTimeoutMapValue),
		mapTime:    make(map[int64]*cGatlingTimeoutMapValue),
		lastid:     0,
	}
}

func (pInst *CTimeoutMap) Set(key, value interface{}) {
	timeNow := time.Now().UnixMilli()
	endStamp := timeNow + int64(pInst.mapOptions.DefaultExpiredSeconds)*1000
	pInst.lockMain.Lock()
	defer pInst.lockMain.Unlock()
	dataNew := &cGatlingTimeoutMapValue{
		endTime:    endStamp,
		data:       value,
		key:        key,
		mapTimeAdd: 0,
	}
	pInst.mapData[key] = dataNew
	pInst.addToTimeMap(dataNew)

	pInst.checkDelete(timeNow)
	//fmt.Println("add: ", fmt.Sprintf("%v", key), timeNow, endStamp, pInst.lastid)
}
func (pInst *CTimeoutMap) Get(key interface{}) (retData interface{}) {
	retData = nil
	pInst.lockMain.Lock()
	defer pInst.lockMain.Unlock()
	value, exists := pInst.mapData[key]
	if !exists {
		return nil
	}
	timeNow := time.Now().UnixMilli()
	value.endTime = timeNow + int64(pInst.mapOptions.DefaultExpiredSeconds)*1000
	retData = value.data
	pInst.addToTimeMap(value)

	pInst.checkDelete(timeNow)

	return retData
}

func (pInst *CTimeoutMap) addToTimeMap(dataNew *cGatlingTimeoutMapValue) {
	if dataNew.mapTimeAdd > 1 {
		return
	}
	if dataNew.endTime <= pInst.lastid {
		pInst.lastid++
	} else {
		pInst.lastid = dataNew.endTime
	}
	pInst.mapTime[pInst.lastid] = dataNew
	dataNew.mapTimeAdd++
}
func (pInst *CTimeoutMap) checkDelete(timeNow int64) {
	for k, v := range pInst.mapTime {
		if k < timeNow {
			v.mapTimeAdd--
			if v.endTime < timeNow {
				delete(pInst.mapData, v.key)
			}
			delete(pInst.mapTime, k)
		}
		break
	}
}
