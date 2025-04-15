package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ParsePrice(priceStr string) (float64, error) {
	if priceStr == "" { return 0, fmt.Errorf("preço vazio") }
	cleaned := strings.ReplaceAll(priceStr, "R$", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	cleaned = strings.TrimSpace(cleaned)
	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil { return 0, fmt.Errorf("erro ao converter '%s' para float: %v", priceStr, err) }
	return price, nil
}

func FormatPrice(price float64) string {
	if math.IsNaN(price) || math.IsInf(price, 0) { return "N/A" }
	p := message.NewPrinter(language.BrazilianPortuguese)
	return p.Sprintf("R$ %.2f", price)
}

func CalculatePercentageDiff(v1, v2 float64) (*float64, bool) {
	if v2 == 0 || math.IsNaN(v1) || math.IsNaN(v2) || math.IsInf(v1, 0) || math.IsInf(v2, 0) {
		return nil, false // Não é possível calcular
	}
	diff := ((v1 / v2) - 1) * 100
	return &diff, true
}