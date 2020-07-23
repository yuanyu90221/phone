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

type phone struct {
	id     int
	number string
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
	_, err = insertPhone(db, "1234567890")
	must(err)
	_, err = insertPhone(db, "123 456 7891")
	must(err)
	id, err := insertPhone(db, "(123) 456 7892")
	must(err)
	_, err = insertPhone(db, "(123) 456-7893")
	must(err)
	_, err = insertPhone(db, "123-456-7894")
	must(err)
	_, err = insertPhone(db, "123-456-7890")
	must(err)
	_, err = insertPhone(db, "1234567892")
	must(err)
	_, err = insertPhone(db, "(123)456-7892")
	must(err)

	number, err := getPhone(db, id)
	must(err)
	fmt.Println("Number is ", number)
	phones, err := allPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("working on... %+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating or removing ...", number)
		} else {
			fmt.Println("No changes required")
		}
	}
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}
func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id=$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", nil
	}
	return number, nil
}
func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
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
