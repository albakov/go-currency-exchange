dev:
	go build cmd/main.go && mv main currency_exchange
	./currency_exchange

build:
	go build cmd/main.go && mv main currency_exchange