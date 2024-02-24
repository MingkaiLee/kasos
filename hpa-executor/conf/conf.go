package conf

var ServerConf = struct {
	Port string
}{
	Port: "8080",
}

var NormalTesterConf = struct {
	DefaultClientTimeout int64
	RtTestEpoch          int
}{
	DefaultClientTimeout: 60,
	RtTestEpoch:          10,
}
