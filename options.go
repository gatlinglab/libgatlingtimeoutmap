package libgatlingtimeoutmap

type OptionsGatlingTimeoutMap struct {
	DefaultExpiredSeconds  int
	AutoUpdateTimeWhenCall bool
}

var g_DefaultOptions = &OptionsGatlingTimeoutMap{
	DefaultExpiredSeconds:  18000, // 5 hours
	AutoUpdateTimeWhenCall: true,
}

func GetDefaultTimeoutMapOptions() *OptionsGatlingTimeoutMap {
	return g_DefaultOptions
}
