package entity

type ExchangeRates struct {
	ID             int64    `json:"id"`
	BaseCurrency   Currency `json:"baseCurrency"`
	TargetCurrency Currency `json:"targetCurrency"`
	Rate           float64  `json:"rate"`
}
