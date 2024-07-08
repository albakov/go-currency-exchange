package main

import (
	"github.com/albakov/go-currency-exchange/dbinit"
	"github.com/albakov/go-currency-exchange/internal/app"
	"github.com/albakov/go-currency-exchange/internal/config"
)

func main() {
	c := config.MustNew()
	dbinit.New(c).MustCreateDatabaseIfNotExists()
	app.New(c).MustStart()
}
