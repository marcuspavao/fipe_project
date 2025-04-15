package models

const AnoZeroKm = 32000 

type PriceInfo struct {
	Modelo   string  `json:"modelo"`
	Valor    float64 `json:"-"`
	ValorFmt string  `json:"valorFmt"`
}

type BrandPeriodStats struct {
	Ref                string             `json:"ref"`
	TabelaId           int                `json:"-"`
	MenorPreco0km      PriceInfo          `json:"menorPreco0km"`
	MaiorPreco0km      PriceInfo          `json:"maiorPreco0km"`
	ValorMedio0km      float64            `json:"-"`
	ValorMedio0kmFmt   string             `json:"valorMedio0kmFmt"`
	TotalModelos       int                `json:"totalModelos"`
	TotalVeiculos0km   int                `json:"totalVeiculos0km"`
	SomaValores0km     float64            `json:"-"` // Exemplo, tornando explícito que é interno
	ModelosEncontrados map[int32]struct{} `json:"-"` // Exemplo
	Inicializado       bool               `json:"-"` // Exemplo
}

type PercentageDiffs struct {
	ValorMedio0km *float64 `json:"valorMedio0km,omitempty"`
	TotalModelos  *float64 `json:"totalModelos,omitempty"`
}

type DashboardBrandEntry struct {
	BrandName             string           `json:"brandName"`
	BrandCode             int32            `json:"brandCode"`
	Periodo1              BrandPeriodStats `json:"periodo1"`
	Periodo2              BrandPeriodStats `json:"periodo2"`
	DiferencasPercentuais PercentageDiffs  `json:"diferencasPercentuais"`
}


