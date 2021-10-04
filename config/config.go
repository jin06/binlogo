package config

var Cfg Config

type Config struct {
	Cluster struct {
		Name string `config:"name"`
	} `config:"cluster"`
	Node struct {
		Name string `config:"name"`
	} `config:"node"`
	Api struct {
		Port uint `config:"port"`
	} `config:"api"`
	Store struct {
		Type string `config:"type"`
		Etcd struct {
			Endpoints []string `config:"endpoints"`
		} `config:"etcd"`
	} `config:"store"`
}

//func (c *Config) String() string {
//	return fmt.Sprintf(
//		"%v", c,
//	)
//}

//func InitCfg(path string) {
	//loader := confita.NewLoader(file.NewBackend(path))
	//err := loader.Load(context.Background(), &Cfg)
	//fmt.Println(Cfg.String())
	//fmt.Println(Cfg)
	//if err != nil {
	//	panic(err)
	//}
//}
