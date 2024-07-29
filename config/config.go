package config

var Version = "0.0.1"

type Config struct {
	Url          string
	Host         string
	Port         int
	Dict         string
	GoroutineNum int
	SleepTime    int
}
