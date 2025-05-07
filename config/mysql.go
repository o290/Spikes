package config

type Mysql struct {
	Host     string
	Port     string
	DB       string
	User     string
	Pwd      string
	LogLevel string
	Config   string
}

// 写连接地址
func (m Mysql) Dsn() string {
	return m.User + ":" + m.Pwd + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DB + "?" + m.Config
}
