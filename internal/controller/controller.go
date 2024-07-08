package controller

import (
	"encoding/json"
	"fmt"
	"github.com/albakov/go-currency-exchange/internal/util"
	"net/http"
)

const f = "controller.Controller"

type ServerResponse interface {
	ShowResponse(w http.ResponseWriter, statusCode int, msg interface{})
	ShowError(w http.ResponseWriter, statusCode int, message string)
	ShowMethodNotAllowedError(w http.ResponseWriter)
	ShowReadyToPatch(w http.ResponseWriter)
}

type Controller struct {
}

type errorResponse struct {
	Message string `json:"message"`
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) ShowResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	const op = "ShowResponse"

	c.setHeaders(w)

	response, err := json.Marshal(msg)
	if err != nil {
		util.LogError(f, op, fmt.Errorf("convert response to json: %v", err))
		c.ShowError(w, http.StatusInternalServerError, MessageServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = fmt.Fprint(w, string(response))
	if err != nil {
		util.LogError(f, op, fmt.Errorf("response to json: %v", err))

		return
	}
}

func (c *Controller) ShowError(w http.ResponseWriter, statusCode int, message string) {
	const op = "ShowError"

	c.setHeaders(w)

	response := errorResponse{Message: message}
	res, err := json.Marshal(response)
	if err != nil {
		util.LogError(f, op, fmt.Errorf("convert response to json: %v", err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = fmt.Fprint(w, string(res))
	if err != nil {
		util.LogError(f, op, fmt.Errorf("response to json: %v", err))

		return
	}
}

func (c *Controller) ShowMethodNotAllowedError(w http.ResponseWriter) {
	c.ShowError(w, http.StatusMethodNotAllowed, MessageMethodNotAllowed)
}

func (c *Controller) ShowReadyToPatch(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
