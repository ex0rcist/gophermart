package entities

import "github.com/shopspring/decimal"

// gophermart decimal
// сериализация decimal по-умолчанию идёт в строку
type GDecimal decimal.Decimal

// автотесты требуют, чтобы в ответе были числа, а не строки
func (d GDecimal) MarshalJSON() ([]byte, error) {
	return []byte(decimal.Decimal(d).String()), nil
}
