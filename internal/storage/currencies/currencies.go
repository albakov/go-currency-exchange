package currencies

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-currency-exchange/internal/entity"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/util"
	"strings"
)

const f = "storage.Currencies"

type StorageCurrencies interface {
	All() []entity.Currency
	ByCode(code string) (entity.Currency, error)
	Add(currency entity.Currency) (int64, error)
}

type Currencies struct {
	pathToDb string
}

func New(pathToDb string) *Currencies {
	return &Currencies{
		pathToDb: pathToDb,
	}
}

func (c *Currencies) All() []entity.Currency {
	const op = "All"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	stmt, err := db.Query("SELECT ID, Code, FullName, Sign FROM Currencies")
	if err != nil {
		return []entity.Currency{}
	}
	defer func(stmt *sql.Rows) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	currencies := []entity.Currency{}

	for stmt.Next() {
		currency := entity.Currency{}

		err := stmt.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign)
		if err != nil {
			util.LogError(f, op, err)

			return []entity.Currency{}
		}

		currencies = append(currencies, currency)
	}

	return currencies
}

func (c *Currencies) ByCode(code string) (entity.Currency, error) {
	const op = "Code"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	currency := entity.Currency{}

	row := db.QueryRow("SELECT ID, Code, FullName, Sign FROM Currencies WHERE Code = ?", code)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return currency, row.Err()
	}

	err = row.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Currency{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Currency{}, err
	}

	return currency, nil
}

func (c *Currencies) Add(currency entity.Currency) (int64, error) {
	const op = "Add"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	stmt, err := db.Prepare("INSERT INTO Currencies (Code, FullName, Sign) VALUES (?, ?, ?)")
	if err != nil {
		util.LogError(f, op, err)

		return 0, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	exec, err := stmt.Exec(currency.Code, currency.FullName, currency.Sign)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return 0, storage.EntityAlreadyExistsError
		}

		util.LogError(f, op, err)

		return 0, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		util.LogError(f, op, err)

		return 0, err
	}

	return id, nil
}
