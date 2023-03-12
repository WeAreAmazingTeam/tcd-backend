package constant

import (
	"os"
	"strconv"
)

var (
	DB_USER    string
	DB_PASS    string
	DB_NAME    string
	RW_HOST    string
	RO_HOST    string
	DB_PORT    string
	STR_DSN    string
	DB_CACHING bool
	DELETED_BY bool
)

func InitDBConstant() {
	dbCaching, _ := strconv.ParseBool(os.Getenv("DB_CACHING"))
	dbDeletedBy, _ := strconv.ParseBool(os.Getenv("DELETED_BY"))
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
	RW_HOST = os.Getenv("RW_HOST")
	RO_HOST = os.Getenv("RW_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	STR_DSN = os.Getenv("STR_DSN")
	DB_CACHING = dbCaching
	DELETED_BY = dbDeletedBy
}
