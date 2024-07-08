package controller

const (
	MessageServerError                       = "Ошибка на сервере"
	MessageMethodNotAllowed                  = "Метод не доступен"
	MessageFieldEmpty                        = "Отсутствует нужное поле: %s"
	MessageFieldIncorrectError               = "Некорректно указано поле %s"
	MessageCurrencyCodeEmpty                 = "Код валюты отсутствует в адресе"
	MessageCurrencyAlreadyExists             = "Валюта с таким кодом уже существует"
	MessageCurrencyNotFound                  = "Валюта не найдена"
	MessageExchangeRatesAlreadyExists        = "Валютная пара с таким кодом уже существует"
	MessageExchangeRatesCurrencyNotFound     = "Одна (или обе) валюты из валютной пары не существует в БД"
	MessageExchangeRatesPairEmpty            = "Коды валют пары отсутствуют в адресе"
	MessageExchangeRatesPairCurrencyNotFound = "Обменный курс для пары не найден"
	MessageExchangeRatesPairNotFound         = "Валютная пара не найдена"
)
