package currencies

import (
	"errors"
	"github.com/albakov/go-currency-exchange/internal/config"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"github.com/albakov/go-currency-exchange/internal/entity"
	"github.com/albakov/go-currency-exchange/internal/storage"
	"github.com/albakov/go-currency-exchange/internal/storage/currencies"
	"github.com/albakov/go-currency-exchange/internal/util"
	"github.com/albakov/go-currency-exchange/internal/validation"
	"net/http"
)

const f = "currencies.Controller"

type Controller struct {
	storageCurrencies currencies.StorageCurrencies
	commonController  controller.ServerResponse
}

func New(config *config.Config, commonController controller.ServerResponse) *Controller {
	return &Controller{
		storageCurrencies: currencies.New(config.PathToDB),
		commonController:  commonController,
	}
}

func (cc *Controller) CurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.currenciesGetHandler(w)

		return
	}

	if r.Method == http.MethodPost {
		cc.currenciesAddHandler(w, r)

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w)
}

func (cc *Controller) CurrencyCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cc.commonController.ShowMethodNotAllowedError(w)

		return
	}

	const op = "currencyCodeHandler"

	code := r.PathValue("code")
	if code == "" {
		cc.commonController.ShowError(w, http.StatusBadRequest, controller.MessageCurrencyCodeEmpty)

		return
	}

	currency, err := cc.storageCurrencies.ByCode(code)
	if err != nil {
		if errors.Is(err, storage.EntitiesNotFoundError) {
			cc.commonController.ShowError(w, http.StatusNotFound, controller.MessageCurrencyNotFound)

			return
		}

		util.LogError(f, op, err)
		cc.commonController.ShowError(w, http.StatusInternalServerError, controller.MessageServerError)

		return
	}

	cc.commonController.ShowResponse(w, http.StatusOK, currency)
}

func (cc *Controller) currenciesGetHandler(w http.ResponseWriter) {
	cc.commonController.ShowResponse(w, http.StatusOK, cc.storageCurrencies.All())
}

func (cc *Controller) currenciesAddHandler(w http.ResponseWriter, r *http.Request) {
	const op = "currenciesAddHandler"

	validated := validation.NewCurrencies(r, map[string]string{"name": "", "code": "", "sign": ""})
	validated.Validate()

	if !validated.IsValid() {
		cc.commonController.ShowError(w, http.StatusBadRequest, validated.ErrorMessage())

		return
	}

	currency := entity.Currency{
		Code:     validated.Field("code"),
		FullName: validated.Field("name"),
		Sign:     validated.Field("sign"),
	}

	id, err := cc.storageCurrencies.Add(currency)
	if err != nil {
		if errors.Is(err, storage.EntityAlreadyExistsError) {
			cc.commonController.ShowError(w, http.StatusConflict, controller.MessageCurrencyAlreadyExists)

			return
		}

		util.LogError(f, op, err)
		cc.commonController.ShowError(w, http.StatusInternalServerError, controller.MessageServerError)

		return
	}

	currency.ID = id

	cc.commonController.ShowResponse(w, http.StatusCreated, currency)
}
