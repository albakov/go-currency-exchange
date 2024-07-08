package services

import (
	"errors"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/storage/currencies"
	"github.com/albakov/go-currency-exchange/internal/storage/exchangerates"
	"github.com/albakov/go-currency-exchange/internal/util"
	"math"
)

const f = "services.Exchange"

var NotFoundError = errors.New("not found")

type Exchange struct {
	baseCurrencyId, targetCurrencyId int64
	isReversed                       bool
	rate, amount                     float64
	storageCurrencies                currencies.StorageCurrencies
	storageExchangeRates             exchangerates.StorageExchangeRates
}

func New(
	storageCurrencies currencies.StorageCurrencies,
	storageExchangeRates exchangerates.StorageExchangeRates,
	baseCurrencyId,
	targetCurrencyId int64,
	amount float64,
) *Exchange {
	return &Exchange{
		storageCurrencies:    storageCurrencies,
		storageExchangeRates: storageExchangeRates,
		baseCurrencyId:       baseCurrencyId,
		targetCurrencyId:     targetCurrencyId,
		amount:               amount,
	}
}

func (e *Exchange) Rate() (float64, error) {
	err := e.calculate()
	if err != nil {
		return 0, err
	}

	return e.round(e.rate), nil
}

func (e *Exchange) ConvertedAmount() float64 {
	convertedAmount := 0.0
	rate := e.round(e.rate)

	if e.amount <= 0 || rate <= 0 {
		return 0
	}

	if e.isReversed {
		convertedAmount = e.amount / rate
	} else {
		convertedAmount = e.amount * rate
	}

	return e.round(convertedAmount)
}

func (e *Exchange) calculate() error {
	err := e.direct()
	if err == nil {
		return nil
	}

	if !errors.Is(err, NotFoundError) {
		return err
	}

	err = e.reverse()
	if err == nil {
		e.isReversed = true

		return nil
	}

	if !errors.Is(err, NotFoundError) {
		return err
	}

	err = e.cross()
	if err == nil {
		return nil
	}

	return err
}

func (e *Exchange) direct() error {
	const op = "direct"

	exchangeRate, err := e.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(
		e.baseCurrencyId,
		e.targetCurrencyId,
	)

	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			return NotFoundError
		}

		util.LogError(f, op, err)

		return err
	}

	e.rate = exchangeRate.Rate

	return nil
}

func (e *Exchange) reverse() error {
	const op = "reverse"

	exchangeRate, err := e.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(
		e.targetCurrencyId,
		e.baseCurrencyId,
	)

	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			return NotFoundError
		}

		util.LogError(f, op, err)

		return err
	}

	e.rate = exchangeRate.Rate

	return nil
}

func (e *Exchange) cross() error {
	const op = "cross"

	usdCurrency, err := e.storageCurrencies.ByCode("USD")
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			return NotFoundError
		}

		util.LogError(f, op, err)

		return err
	}

	exchangeRateA, err := e.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(
		usdCurrency.ID,
		e.baseCurrencyId,
	)
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			return NotFoundError
		}

		util.LogError(f, op, err)

		return err
	}

	exchangeRateB, err := e.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(
		usdCurrency.ID,
		e.targetCurrencyId,
	)
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			return NotFoundError
		}

		util.LogError(f, op, err)

		return err
	}

	e.rate = exchangeRateB.Rate / exchangeRateA.Rate

	return nil
}

func (e *Exchange) round(value float64) float64 {
	if value <= 0 {
		return 0
	}

	return math.Round(value*100) / 100
}
