package db

import (
	"aat-manager/utils"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// Singleton instance
var (
	instance *sql.DB
	once     sync.Once
)

// pgConnect Connect to the DB, this function is safe for concurrent use.
func pgConnect() *sql.DB {
	once.Do(func() {
		pgConnString := utils.ReadEnvOrPanic(utils.PGRESCONNSTRING)
		db, err := sql.Open("postgres", pgConnString)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// configure the db connection pool here
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Minute * 5)

		// ping the DB to ensure connection is valid
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		log.Println("Connected to the database")

		createTables(db)

		instance = db
	})

	return instance
}

// The function takes a *sql.DB as the input parameter and creates the following tables if they do not already exist:
// - tokens: This table stores encrypted tokens, with columns name and value.
// - The table and column names have appropriate comments assigned to them for better understanding.
// The function iterates through the list of queries and executes each query using the provided DB connection.
// If there is an error during query execution, the error along with the corresponding query is logged.
func createTables(db *sql.DB) {
	queries := []string{
		`create table if not exists tokens
(
    name  varchar not null
        constraint tokens_pk
            primary key,
    value varchar not null
);

comment on table tokens is 'Tokens table.
All tokens are encrypted';

comment on column tokens.name is 'Token meaningful name, must be unique';

comment on column tokens.value is 'Encrypted token';

`,
	}

	// Actually create all table in db if not exists
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("error executing query: %s: %v", query, err)
		}
	}
}
