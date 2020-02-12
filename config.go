package argus

type db struct {
	User string
	Pass string
	Host string
	Port string
}

type Config struct {
	Db db
}
