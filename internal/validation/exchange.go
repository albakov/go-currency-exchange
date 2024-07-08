package validation

import (
	"fmt"
	"github.com/albakov/go-currency-exchange/internal/controller"
	"net/http"
	"strconv"
)

type RequestExchange struct {
	r            *http.Request
	fields       map[string]string
	errorMessage string
	amount       float64
}

func NewExchange(r *http.Request, fields map[string]string) *RequestExchange {
	return &RequestExchange{
		r:      r,
		fields: fields,
	}
}

func (re *RequestExchange) Validate() {
	query := re.r.URL.Query()

	for field := range re.fields {
		v := query.Get(field)

		if v == "" {
			re.errorMessage = fmt.Sprintf(controller.MessageFieldEmpty, field)

			return
		}

		re.fields[field] = v
	}

	amount, err := strconv.ParseFloat(re.fields["amount"], 64)
	if err != nil || amount <= 0 {
		re.errorMessage = fmt.Sprintf(controller.MessageFieldIncorrectError, "amount")

		return
	}

	re.amount = amount
}

func (re *RequestExchange) IsValid() bool {
	return re.errorMessage == ""
}

func (re *RequestExchange) ErrorMessage() string {
	return re.errorMessage
}

func (re *RequestExchange) Field(field string) string {
	return re.fields[field]
}

func (re *RequestExchange) Amount() float64 {
	return re.amount
}
