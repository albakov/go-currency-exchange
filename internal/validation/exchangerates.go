package validation

import (
	"fmt"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"net/http"
	"strconv"
)

type RequestExchangeRatesAdd struct {
	r            *http.Request
	fields       map[string]string
	errorMessage string
	rate         float64
}

func NewExchangeRates(r *http.Request, fields map[string]string) *RequestExchangeRatesAdd {
	return &RequestExchangeRatesAdd{
		r:      r,
		fields: fields,
	}
}

func (er *RequestExchangeRatesAdd) Validate() {
	for field := range er.fields {
		v := er.r.FormValue(field)

		if v == "" {
			er.errorMessage = fmt.Sprintf(controller.MessageFieldEmpty, field)

			return
		}

		er.fields[field] = v
	}

	rate, err := strconv.ParseFloat(er.fields["rate"], 64)
	if err != nil || rate <= 0 {
		er.errorMessage = fmt.Sprintf(controller.MessageFieldIncorrectError, "rate")

		return
	}

	er.rate = rate
}

func (er *RequestExchangeRatesAdd) IsValid() bool {
	return er.errorMessage == ""
}

func (er *RequestExchangeRatesAdd) ErrorMessage() string {
	return er.errorMessage
}

func (er *RequestExchangeRatesAdd) Rate() float64 {
	return er.rate
}

func (er *RequestExchangeRatesAdd) Field(field string) string {
	return er.fields[field]
}
