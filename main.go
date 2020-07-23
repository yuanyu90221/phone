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
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()
	must(createPhoneNumbersTable(db))
}
func createPhoneNumbersTable(db *sql.DB) error {
	statement := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS  phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`)
	_, err := db.Exec(statement)
	return err
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
