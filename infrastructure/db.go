package infrastructure

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func ConnectDB() {
	// 🔐 Leer conexión desde variable de entorno
	connString := os.Getenv("DB_CONN")

	if connString == "" {
		log.Fatal("DB_CONN not set")
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
