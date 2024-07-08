package exchange

import (
	"errors"
	"github.com/albakov/go-currency-exchange/internal/config"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"github.com/albakov/go-currency-exchange/internal/entity"
	"github.com/albakov/go-currency-exchange/internal/services"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/storage/currencies"
	"github.com/albakov/go-currency-exchange/internal/storage/exchangerates"
	"github.com/albakov/go-currency-exchange/internal/validation"
	"net/http"
)

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

func (ce Controller) Exchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ce.commonController.ShowMethodNotAllowedError(w)

		return
	}

	validated := validation.NewExchange(
		r,
		map[string]string{"from": "", "to": "", "amount": ""},
	)
	validated.Validate()

	if !validated.IsValid() {
		ce.commonController.ShowError(w, http.StatusBadRequest, validated.ErrorMessage())

		return
	}

	baseCurrency, err := ce.storageCurrencies.ByCode(validated.Field("from"))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	targetCurrency, err := ce.storageCurrencies.ByCode(validated.Field("to"))
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesCurrencyNotFound)

			return
		}

		return
	}

	exchangeService := services.New(
		ce.storageCurrencies,
		ce.storageExchangeRates,
		baseCurrency.ID,
		targetCurrency.ID,
		validated.Amount(),
	)
	rate, err := exchangeService.Rate()
	if err != nil {
		if errors.Is(err, services.NotFoundError) {
			ce.commonController.ShowError(w, http.StatusNotFound, controller.MessageExchangeRatesPairNotFound)

			return
		}

		return
	}

	exchange := entity.Exchange{
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		Rate:            rate,
		Amount:          validated.Amount(),
		ConvertedAmount: exchangeService.ConvertedAmount(),
	}

	ce.commonController.ShowResponse(w, http.StatusOK, exchange)
}
