package config

import "fmt"

type System struct {
	Name string
	Host string
	Port int
	Env  string
}

func (s System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
