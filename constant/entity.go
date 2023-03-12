package constant

import (
	"database/sql"

	"gorm.io/gorm"
)

type (
	CreatedDeleted struct {
		CreatedAt sql.NullTime   `json:"created_at"`
		CreatedBy sql.NullString `json:"created_by" gorm:"default:SYSTEM"`
		DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"default:null"`
		DeletedBy sql.NullString `json:"deleted_by" gorm:"default:null"`
	}

	CreatedUpdatedDeleted struct {
		CreatedAt sql.NullTime   `json:"created_at"`
		CreatedBy sql.NullString `json:"created_by" gorm:"default:SYSTEM"`
		UpdatedAt sql.NullTime   `json:"updated_at" gorm:"default:null"`
		UpdatedBy sql.NullString `json:"updated_by" gorm:"default:null"`
		DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"default:null"`
		DeletedBy sql.NullString `json:"deleted_by" gorm:"default:null"`
	}
)
