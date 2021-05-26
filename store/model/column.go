package model

type Column struct {
	ID uint32 `json:"id"`
	DatabaseID uint32 `json:"database_id"`
	TableID uint32 `json:"table_id"`
	Charset string `json:"charset"`
	ColumnType string `json:"column_type"`
	EnumValues []string `json:"enum_values"`
}