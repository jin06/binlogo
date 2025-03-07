package configs

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// NodeName current node's name
var NodeName string

// NodeIP current node's ip
var NodeIP net.IP

var Default Config

func init() {
	hostname, _ := os.Hostname()
	Default = Config{
		NodeName:    hostname,
		ClusterName: CLUSTER_NAME,
		Console: Console{
			Listen: CONSOLE_LISTEN,
			Port:   8081,
		},
		Roles: Roles{
			API:    true,
			Master: true,
		},
		Monitor: Monitor{
			Port: 8085,
		},
	}
}

type Config struct {
	Test        string  `yaml:"test"`
	ClusterName string  `yaml:"cluster_name"`
	NodeName    string  `yaml:"node_name"`
	LogLevel    string  `yaml:"log_level"`
	Console     Console `yaml:"console"`
	Roles       Roles   `yaml:"roles"`
	Monitor     Monitor `yaml:"monitor"`
	Store       Store   `yaml:"store"`
	Auth        Auth    `yaml:"auth"`
	Profile     bool    `yaml:"profile"`
	ProfilePort uint    `yaml:"profile_port"`
	EventExpire int     `yaml:"event_expire"`
}

func (cfg *Config) GetNodeName() string {
	if len(cfg.NodeName) == 0 {
		host, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		return host
	}
	return cfg.NodeName
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
	Password string `yaml:"password"`
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
	logrus.Infof("config file path: %s", file)
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &Default)
	return err
}

func GetNodeName() string {
	return Default.GetNodeName()
}
