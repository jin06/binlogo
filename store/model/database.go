package model

type Database struct {
	ID uint32 `json:"id"`
	Name string `json:"name"`
	SchemasID uint32`json:"schema_id"`
	Charset string `json:"charset"`
}
