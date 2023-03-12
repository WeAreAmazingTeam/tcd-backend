package helper

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// handle null string
func HNString(param sql.NullString) any {
	if param.Valid {
		return param.String
	}
	return nil
}

// handle null time
func HNTime(param sql.NullTime) any {
	if param.Valid {
		return param.Time.Format("2006-01-02 15:04:05")
	}
	return nil
}

// handle null time gorm
func HNTimeGDeletedAt(param gorm.DeletedAt) any {
	if param.Valid {
		return param.Time.Format("2006-01-02 15:04:05")
	}
	return nil
}

func ConvertToInLineQuery(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "  ", " ")

	re := regexp.MustCompile(` +\r?\n +`)
	return re.ReplaceAllString(s, " ")
}

func IsErrNoRows(s string) bool {
	return strings.Contains(s, "no rows")
}

// set null string
func SetNS(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  true,
	}
}

// set now gorm null time
func SetNowNT() *gorm.DeletedAt {
	return &gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}
