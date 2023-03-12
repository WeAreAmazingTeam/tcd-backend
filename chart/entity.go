package chart

import "database/sql"

type (
	Chart struct {
		ID        int          `json:"id"`
		Name      string       `json:"name"`
		Content   string       `json:"content"`
		CreatedAt sql.NullTime `json:"created_at"`
	}
)
