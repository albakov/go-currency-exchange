package entity

type Exchange struct {
	BaseCurrency    Currency `json:"baseCurrency"`
	TargetCurrency  Currency `json:"targetCurrency"`
	Rate            float64  `json:"rate"`
	Amount          float64  `json:"amount"`
	ConvertedAmount float64  `json:"convertedAmount"`
}
