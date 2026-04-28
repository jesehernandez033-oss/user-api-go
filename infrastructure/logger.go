package infrastructure

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening log file:", err)
	}

	Logger = log.New(file, "API_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}