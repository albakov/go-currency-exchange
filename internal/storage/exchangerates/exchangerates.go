package exchangerates

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-currency-exchange/internal/entity"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/util"
	"strings"
)

const f = "storage.ExchangeRatesHandler"

type StorageExchangeRates interface {
	All() []entity.ExchangeRates
	Add(exchangeRates entity.ExchangeRates) (int64, error)
	ByBaseCurrencyIdAndTargetCurrencyId(baseCurrencyId int64, targetCurrencyId int64) (entity.ExchangeRates, error)
	UpdateRate(exchangeRates entity.ExchangeRates) error
}

type ExchangeRates struct {
	pathToDb string
}

func New(pathToDb string) *ExchangeRates {
	return &ExchangeRates{
		pathToDb: pathToDb,
	}
}

func (c *ExchangeRates) All() []entity.ExchangeRates {
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

	stmt, err := db.Query(
		`SELECT ExchangeRates.ID, ExchangeRates.Rate, 
       		BaseCurrency.ID as BaseCurrencyID,
       		BaseCurrency.Code as BaseCurrencyCode,
       		BaseCurrency.FullName as BaseCurrencyFullName,
       		BaseCurrency.Sign as BaseCurrencySign,
       		TargetCurrency.ID as TargetCurrencyID,
       		TargetCurrency.Code as TargetCurrencyCode,
       		TargetCurrency.FullName as TargetCurrencyFullName,
       		TargetCurrency.Sign as TargetCurrencySign
       	FROM ExchangeRates
		LEFT JOIN Currencies as BaseCurrency ON BaseCurrency.Id = ExchangeRates.BaseCurrencyId 
		LEFT JOIN Currencies as TargetCurrency ON TargetCurrency.Id = ExchangeRates.TargetCurrencyId`,
	)
	if err != nil {
		return []entity.ExchangeRates{}
	}
	defer func(stmt *sql.Rows) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	currencies := []entity.ExchangeRates{}

	for stmt.Next() {
		exchangeRates := entity.ExchangeRates{}
		baseCurrency := entity.Currency{}
		targetCurrency := entity.Currency{}

		err := stmt.Scan(
			&exchangeRates.ID,
			&exchangeRates.Rate,
			&baseCurrency.ID,
			&baseCurrency.Code,
			&baseCurrency.FullName,
			&baseCurrency.Sign,
			&targetCurrency.ID,
			&targetCurrency.Code,
			&targetCurrency.FullName,
			&targetCurrency.Sign,
		)
		if err != nil {
			util.LogError(f, op, err)

			return []entity.ExchangeRates{}
		}

		exchangeRates.BaseCurrency = baseCurrency
		exchangeRates.TargetCurrency = targetCurrency

		currencies = append(currencies, exchangeRates)
	}

	return currencies
}

func (c *ExchangeRates) Add(exchangeRates entity.ExchangeRates) (int64, error) {
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

	stmt, err := db.Prepare(
		"INSERT INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate) VALUES (?, ?, ?)",
	)
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

	exec, err := stmt.Exec(exchangeRates.BaseCurrency.ID, exchangeRates.TargetCurrency.ID, exchangeRates.Rate)
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

func (c *ExchangeRates) ByBaseCurrencyIdAndTargetCurrencyId(
	baseCurrencyId int64,
	targetCurrencyId int64,
) (entity.ExchangeRates, error) {
	const op = "ByBaseCurrencyIdAndTargetCurrencyId"

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

	row := db.QueryRow(
		`SELECT ExchangeRates.ID, ExchangeRates.Rate, 
       		BaseCurrency.ID as BaseCurrencyID,
       		BaseCurrency.Code as BaseCurrencyCode,
       		BaseCurrency.FullName as BaseCurrencyFullName,
       		BaseCurrency.Sign as BaseCurrencySign,
       		TargetCurrency.ID as TargetCurrencyID,
       		TargetCurrency.Code as TargetCurrencyCode,
       		TargetCurrency.FullName as TargetCurrencyFullName,
       		TargetCurrency.Sign as TargetCurrencySign
       	FROM ExchangeRates
		LEFT JOIN Currencies as BaseCurrency ON BaseCurrency.Id = ExchangeRates.BaseCurrencyId 
		LEFT JOIN Currencies as TargetCurrency ON TargetCurrency.Id = ExchangeRates.TargetCurrencyId 
		WHERE ExchangeRates.BaseCurrencyId = ? AND ExchangeRates.TargetCurrencyId = ?`,
		baseCurrencyId,
		targetCurrencyId,
	)
	if row.Err() != nil {
		util.LogError(f, op, err)

		return entity.ExchangeRates{}, row.Err()
	}

	exchangeRates := entity.ExchangeRates{}
	baseCurrency := entity.Currency{}
	targetCurrency := entity.Currency{}

	err = row.Scan(
		&exchangeRates.ID,
		&exchangeRates.Rate,
		&baseCurrency.ID,
		&baseCurrency.Code,
		&baseCurrency.FullName,
		&baseCurrency.Sign,
		&targetCurrency.ID,
		&targetCurrency.Code,
		&targetCurrency.FullName,
		&targetCurrency.Sign,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ExchangeRates{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.ExchangeRates{}, err
	}

	exchangeRates.BaseCurrency = baseCurrency
	exchangeRates.TargetCurrency = targetCurrency

	return exchangeRates, nil
}

func (c *ExchangeRates) UpdateRate(exchangeRates entity.ExchangeRates) error {
	const op = "UpdateRate"

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

	stmt, err := db.Prepare("UPDATE ExchangeRates SET Rate = ? WHERE ID = ?")
	if err != nil {
		util.LogError(f, op, err)

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	_, err = stmt.Exec(exchangeRates.Rate, exchangeRates.ID)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	return nil
}
