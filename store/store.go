package store

type Store interface {
	 Get(key string) (resp string, err error)
	 Put(key string, val string) (err error)
}






