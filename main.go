package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	passwd := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, passwd)
	// db, err := sql.Open("postgres", psqlInfo)
	// must(err)
	// db.Close()
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Ping())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
