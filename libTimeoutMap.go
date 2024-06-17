package libgatlingtimeoutmap

type ITimeoutMap interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
}

func NewTimeoutMap() ITimeoutMap {
	return NewTimeoutMapWithOptions(GetDefaultTimeoutMapOptions())
}
func NewTimeoutMapWithOptions(options *OptionsGatlingTimeoutMap) ITimeoutMap {
	return newmap(options)
}
