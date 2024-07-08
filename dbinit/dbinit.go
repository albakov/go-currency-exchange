package dbinit

import (
	"database/sql"
	"github.com/albakov/go-currency-exchange/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type DBInit struct {
	pathToDb string
}

func New(config *config.Config) *DBInit {
	return &DBInit{
		pathToDb: config.PathToDB,
	}
}

func (d *DBInit) MustCreateDatabaseIfNotExists() {
	if !d.mustCheckIsDatabaseExists() {
		_, err := os.Create(d.pathToDb)
		if err != nil {
			panic(err)
		}
	}

	d.mustCreateTablesIfNotExists()
}

func (d *DBInit) mustCheckIsDatabaseExists() bool {
	_, err := os.Stat(d.pathToDb)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	return true
}

func (d *DBInit) mustCreateTablesIfNotExists() {
	db, err := sql.Open("sqlite3", d.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Currencies (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    	Code VARCHAR(255) NOT NULL UNIQUE, 
    	FullName VARCHAR(255) NOT NULL, 
    	Sign VARCHAR(255) NOT NULL)`,
	)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS ExchangeRates (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    	BaseCurrencyId INT NOT NULL, 
    	TargetCurrencyId INT NOT NULL, 
    	Rate DECIMAL(6) NOT NULL,
    	FOREIGN KEY (BaseCurrencyId) REFERENCES Currencies (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	FOREIGN KEY (TargetCurrencyId) REFERENCES Currencies (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	UNIQUE(BaseCurrencyId, TargetCurrencyId) ON CONFLICT ABORT)`,
	)
	if err != nil {
		panic(err)
	}
}
