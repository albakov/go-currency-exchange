package exchangerates

import (
	"errors"
	"github.com/albakov/go-currency-exchange/internal/config"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"github.com/albakov/go-currency-exchange/internal/entity"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/storage/currencies"
	"github.com/albakov/go-currency-exchange/internal/storage/exchangerates"
	"github.com/albakov/go-currency-exchange/internal/util"
	"github.com/albakov/go-currency-exchange/internal/validation"
	"net/http"
	"strings"
)

const f = "exchangerates.Controller"

type Controller struct {
	commonController     controller.ServerResponse
	storageExchangeRates exchangerates.StorageExchangeRates
	storageCurrencies    currencies.StorageCurrencies
}

func New(config *config.Config, commonController controller.ServerResponse) *Controller {
	return &Controller{
		commonController:     commonController,
		storageExchangeRates: exchangerates.New(config.PathToDB),
		storageCurrencies:    currencies.New(config.PathToDB),
	}
}

func (ce *Controller) ExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ce.exchangeRatesGetHandler(w)

		return
	}

	if r.Method == http.MethodPost {
		ce.exchangeRatesAddHandler(w, r)

		return
	}

	ce.commonController.ShowMethodNotAllowedError(w)
}

func (ce *Controller) ExchangeRatesPairHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		ce.commonController.ShowReadyToPatch(w)

		return
	}

	if r.Method == http.MethodGet {
		ce.exchangeRatesPairGetHandler(w, r)

		return
	}

	if r.Method == http.MethodPatch {
		ce.exchangeRatesPairUpdateHandler(w, r)

		return
	}

	ce.commonController.ShowMethodNotAllowedError(w)
}

func (ce *Controller) exchangeRatesGetHandler(w http.ResponseWriter) {
	ce.commonController.ShowResponse(w, http.StatusOK, ce.storageExchangeRates.All())
}

func (ce *Controller) exchangeRatesAddHandler(w http.ResponseWriter, r *http.Request) {
	const op = "exchangeRatesAddHandler"

	validated := validation.NewExchangeRates(
		r,
		map[string]string{"baseCurrencyCode": "", "targetCurrencyCode": "", "rate": ""},
	)
	validated.Validate()

	if !validated.IsValid() {
		ce.commonController.ShowError(w, http.StatusBadRequest, validated.ErrorMessage())

		return
	}

	baseCurrency, err := ce.storageCurrencies.ByCode(validated.Field("baseCurrencyCode"))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	targetCurrency, err := ce.storageCurrencies.ByCode(validated.Field("targetCurrencyCode"))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	exchangeRates := entity.ExchangeRates{
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           validated.Rate(),
	}

	id, err := ce.storageExchangeRates.Add(exchangeRates)
	if err != nil {
		if errors.Is(err, storage.EntityAlreadyExistsError) {
			ce.commonController.ShowError(w, http.StatusConflict, controller.MessageExchangeRatesAlreadyExists)

			return
		}

		util.LogError(f, op, err)
		ce.commonController.ShowError(w, http.StatusInternalServerError, controller.MessageServerError)

		return
	}

	exchangeRates.ID = id

	ce.commonController.ShowResponse(w, http.StatusCreated, exchangeRates)
}

func (ce *Controller) exchangeRatesPairGetHandler(w http.ResponseWriter, r *http.Request) {
	pair := r.PathValue("pair")
	if pair == "" {
		ce.commonController.ShowError(w, http.StatusBadRequest, controller.MessageExchangeRatesPairEmpty)

		return
	}

	currenciesCodes := strings.Split(strings.ToUpper(pair), "")

	if len(currenciesCodes) != 6 {
		ce.commonController.ShowError(w, http.StatusBadRequest, controller.MessageExchangeRatesCurrencyNotFound)

		return
	}

	baseCurrency, err := ce.storageCurrencies.ByCode(strings.Join(currenciesCodes[0:3], ""))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesPairCurrencyNotFound)

			return
		}

		return
	}

	targetCurrency, err := ce.storageCurrencies.ByCode(strings.Join(currenciesCodes[3:], ""))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesPairCurrencyNotFound)

			return
		}

		return
	}

	exchangeRate, err := ce.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesPairNotFound)

			return
		}

		return
	}

	ce.commonController.ShowResponse(w, http.StatusOK, exchangeRate)
}

func (ce *Controller) exchangeRatesPairUpdateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "exchangeRatesPairUpdateHandler"

	pair := r.PathValue("pair")
	if pair == "" {
		ce.commonController.ShowError(w, http.StatusBadRequest, controller.MessageExchangeRatesPairEmpty)

		return
	}

	currenciesCodes := strings.Split(strings.ToUpper(pair), "")

	if len(currenciesCodes) != 6 {
		ce.commonController.ShowError(w, http.StatusBadRequest, controller.MessageExchangeRatesCurrencyNotFound)

		return
	}

	baseCurrency, err := ce.storageCurrencies.ByCode(strings.Join(currenciesCodes[0:3], ""))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	targetCurrency, err := ce.storageCurrencies.ByCode(strings.Join(currenciesCodes[3:], ""))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	validated := validation.NewExchangeRates(r, map[string]string{"rate": ""})
	validated.Validate()

	if !validated.IsValid() {
		ce.commonController.ShowError(w, http.StatusBadRequest, validated.ErrorMessage())

		return
	}

	exchangeRate, err := ce.storageExchangeRates.ByBaseCurrencyIdAndTargetCurrencyId(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesPairNotFound)

			return
		}

		return
	}

	exchangeRate.Rate = validated.Rate()

	err = ce.storageExchangeRates.UpdateRate(exchangeRate)
	if err != nil {
		util.LogError(f, op, err)
		ce.commonController.ShowError(w, http.StatusInternalServerError, controller.MessageServerError)

		return
	}

	ce.commonController.ShowResponse(w, http.StatusOK, exchangeRate)
}
