package expense

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

type Api struct {
	Db *sql.DB
}

func InitDB() {
	var err error
	Db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Cannot connect to DB server", err)
	}

	var createTable string = `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err = Db.Exec(createTable)

	if err != nil {
		log.Fatal("cannot create expense table", err)
	}
}
