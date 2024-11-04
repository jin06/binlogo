package configs

import (
	"log/slog"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

var ENV Env

// NodeName current node's name
var NodeName string

// NodeIP current node's ip
var NodeIP net.IP

var Default Config = Config{
	Test: "1111",
}

type Config struct {
	Test        string  `yaml:"test"`
	ClusterName string  `yaml:"cluster_name"`
	NodeName    string  `yaml:"node_name"`
	Console     Console `yaml:"console"`
	Roles       Roles   `yaml:"roles"`
	Monitor     Monitor `yaml:"monitor"`
	Store       Store   `yaml:"store"`
	Auth        Auth    `yaml:"auth"`
}

type Console struct {
	Port   int    `yaml:"port"`
	Listen string `yaml:"listen"`
}

type Roles struct {
	API    bool `yaml:"api"`
	Master bool `yaml:"master"`
}

type Monitor struct {
	Port int `yaml:"port"`
}

type Store struct {
	Type  string `yaml:"type"`
	Redis Redis  `yaml:"redis"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	Passwrod string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Auth struct {
	Expiration string    `yaml:"auth"`
	Type       string    `yaml:"type"`
	AuthBasic  AuthBasic `yaml:"basic"`
	AuthLDAP   AuthLDAP  `yaml:"ldap"`
}

type AuthBasic struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type AuthLDAP struct {
	Addr      string   `yaml:"addr"`
	UserName  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	BaseDN    string   `yaml:"baseDN"`
	IDAttr    string   `yaml:"idAttr"`
	Attributs []string `yaml:"attributs"`
}

func InitConfig(file string) (err error) {
	if file == "" {
		file = "./etc/binlogo.yaml"
	}
	slog.Info("config file path: %s", file)
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &Default)
	return err
}
