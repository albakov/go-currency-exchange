package validation

import (
	"fmt"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"net/http"
)

type RequestCurrenciesAdd struct {
	r            *http.Request
	fields       map[string]string
	errorMessage string
}

func NewCurrencies(r *http.Request, fields map[string]string) *RequestCurrenciesAdd {
	return &RequestCurrenciesAdd{
		r:      r,
		fields: fields,
	}
}

func (cc *RequestCurrenciesAdd) Validate() {
	for field := range cc.fields {
		v := cc.r.FormValue(field)

		if v == "" {
			cc.errorMessage = fmt.Sprintf(controller.MessageFieldEmpty, field)

			return
		}

		cc.fields[field] = v
	}
}

func (cc *RequestCurrenciesAdd) IsValid() bool {
	return cc.errorMessage == ""
}

func (cc *RequestCurrenciesAdd) ErrorMessage() string {
	return cc.errorMessage
}

func (cc *RequestCurrenciesAdd) Field(field string) string {
	return cc.fields[field]
}
