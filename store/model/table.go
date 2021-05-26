package model

type Table struct {
	ID uint32 `json:"id"`
	Name uint32 `json:"name"`
	DatabaseID uint32 `json:"database_id"`
	Charset string `json:"charset"`
	PrimaryKey string `json:"primary_key"`
}
