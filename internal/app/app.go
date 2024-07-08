package app

import (
	"fmt"
	"github.com/albakov/go-currency-exchange/internal/config"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"github.com/albakov/go-currency-exchange/internal/controller/currencies"
	"github.com/albakov/go-currency-exchange/internal/controller/exchange"
	"github.com/albakov/go-currency-exchange/internal/controller/exchangerates"
	"net/http"
)

type App struct {
	mux                     *http.ServeMux
	config                  *config.Config
	exchangeController      *exchange.Controller
	currenciesController    *currencies.Controller
	exchangeRatesController *exchangerates.Controller
}

func New(config *config.Config) *App {
	commonController := controller.New()

	return &App{
		mux:                     http.NewServeMux(),
		config:                  config,
		exchangeController:      exchange.New(config, commonController),
		currenciesController:    currencies.New(config, commonController),
		exchangeRatesController: exchangerates.New(config, commonController),
	}
}

func (a *App) MustStart() {
	a.SetRoutes()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.config.Host, a.config.Port), a)
	if err != nil {
		panic(err)
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.setCORS(w)
	a.mux.ServeHTTP(w, r)
}

func (a *App) SetRoutes() {
	a.mux.HandleFunc("/exchange", a.exchangeController.Exchange)
	a.mux.HandleFunc("/currencies", a.currenciesController.CurrenciesHandler)
	a.mux.HandleFunc("/currency/{code}", a.currenciesController.CurrencyCodeHandler)
	a.mux.HandleFunc("/exchangeRates", a.exchangeRatesController.ExchangeRatesHandler)
	a.mux.HandleFunc("/exchangeRate/{pair}", a.exchangeRatesController.ExchangeRatesPairHandler)
}

func (a *App) setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", a.config.AccessControlAllowOrigin)
	w.Header().Set("Access-Control-Allow-Headers", a.config.AccessControlAllowHeaders)
	w.Header().Set("Access-Control-Allow-Methods", a.config.AccessControlAllowMethods)
}
