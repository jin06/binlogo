package basedata

type Database struct {
	Name string
	Charset string
	CaseSensitivity bool
}

func (d *Database) getName() string {
	return d.Name
}