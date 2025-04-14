package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"fipe_project/internal/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language" // Para formatação de moeda
	"golang.org/x/text/message"  // Para formatação de moeda
)

func GetTabelasReferencia(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("TabelaReferencia")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Erro ao buscar tabelas de referência: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	var tabelas []bson.M
	if err := cursor.All(ctx, &tabelas); err != nil {
		log.Printf("Erro ao decodificar tabelas de referência: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tabelas)
}

func GetMarcas(w http.ResponseWriter, r *http.Request) {
	tabelaParam := r.URL.Query().Get("tabela")
	if tabelaParam == "" {
		http.Error(w, "Parâmetro 'tabela' é obrigatório", http.StatusBadRequest)
		return
	}
	tabelaId, err := strconv.Atoi(tabelaParam)
	if err != nil {
		http.Error(w, "Parâmetro 'tabela' inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("Veiculos")
	filter := bson.M{"monthYearId": tabelaId}
	projection := options.Find().SetProjection(bson.M{"brandName": 1, "brandCode": 1, "_id": 0})
	cursor, err := collection.Find(ctx, filter, projection)
	if err != nil {
		log.Printf("Erro ao buscar marcas: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	var marcas []bson.M
	if err := cursor.All(ctx, &marcas); err != nil {
		log.Printf("Erro ao decodificar marcas: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marcas)
}

func GetModelos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	marcaParam := vars["marca"]
	codMarca, err := strconv.Atoi(marcaParam)
	if err != nil {
		http.Error(w, "Código de marca inválido", http.StatusBadRequest)
		return
	}
	tabelaParam := r.URL.Query().Get("tabela")
	if tabelaParam == "" {
		http.Error(w, "Parâmetro 'tabela' é obrigatório", http.StatusBadRequest)
		return
	}
	tabelaId, err := strconv.Atoi(tabelaParam)
	if err != nil {
		http.Error(w, "Parâmetro 'tabela' inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("Veiculos")
	var brand bson.M
	filter := bson.M{"brandCode": codMarca, "monthYearId": tabelaId}
	if err := collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		log.Printf("Erro ao buscar marca %d: %v", codMarca, err)
		http.Error(w, "Marca não encontrada", http.StatusNotFound)
		return
	}
	models, ok := brand["models"]
	if !ok {
		http.Error(w, "Modelos não encontrados", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

func GetVeiculos(w http.ResponseWriter, r *http.Request) {

	modeloParam := r.URL.Query().Get("modelo")
	if modeloParam == "" {
		http.Error(w, "Parâmetro 'modelo' é obrigatório", http.StatusBadRequest)
		return
	}
	tabelaParam := r.URL.Query().Get("tabela")
	if tabelaParam == "" {
		http.Error(w, "Parâmetro 'tabela' é obrigatório", http.StatusBadRequest)
		return
	}

	tabelaId, err := strconv.Atoi(tabelaParam)
	if err != nil {
		http.Error(w, "Parâmetro 'tabela' inválido", http.StatusBadRequest)
		return
	}

	modeloId, err := strconv.Atoi(modeloParam)
	if err != nil {
		http.Error(w, "Parâmetro 'modelo' inválido", http.StatusBadRequest)
		return
	}

	//log.Printf("Parâmetros recebidos: tabelaId=%d, modeloId=%d", tabelaId, modeloId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("Veiculos")

	filter := bson.M{
		"monthYearId":      tabelaId,
		"models.modelCode": modeloId,
	}

	//log.Printf("Filtro usado na consulta: %v", filter)

	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("Erro ao buscar veículos: %v", err)
		http.Error(w, "Veículo não encontrado", http.StatusNotFound)
		return
	}

	modelsRaw, exists := result["models"]
	if !exists {
		log.Printf("Campo 'models' não encontrado no documento")
		http.Error(w, "Modelos não encontrados", http.StatusInternalServerError)
		return
	}

	models, ok := modelsRaw.(primitive.A)
	if !ok {
		log.Printf("Campo 'models' não é um array, tipo: %T", modelsRaw)
		http.Error(w, "Modelos não encontrados", http.StatusInternalServerError)
		return
	}

	var selectedYears []bson.M
	for _, model := range models {
		m, ok := model.(bson.M)
		if ok && m["modelCode"].(int32) == int32(modeloId) {
			if years, exists := m["years"].(primitive.A); exists {
				for _, year := range years {
					if yearMap, isMap := year.(bson.M); isMap {
						yearMap["model"] = m["modelName"]
						selectedYears = append(selectedYears, yearMap)
					}
				}
			} else {
				log.Printf("Campo 'years' não encontrado ou não é um array no modelo selecionado.")
			}
			break
		}
	}

	//log.Printf("Anos: %v", selectedYears)

	if len(selectedYears) == 0 {
		log.Printf("Nenhum ano encontrado para o modelo especificado.")
		http.Error(w, "Anos não encontrados para o modelo especificado", http.StatusNotFound)
		return
	}

	//log.Printf("Anos selecionados para o modelo: %v", selectedYears)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedYears)
}

// Dashboard de Marcas - de acordo com as marcas analisar para dois períodos
// infomações como carro/modelo com menor e maior preço 0km, valor médio,
// número de modelos disponíveis e as difenças em porcentagens entre esses aspectos

const anoZeroKm = 32000

// Estrutura para guardar informações de preço (min/max)
type PriceInfo struct {
	Modelo   string  `json:"modelo"`
	Valor    float64 `json:"-"`        // Valor numérico para cálculos
	ValorFmt string  `json:"valorFmt"` // Valor formatado para exibição
}

// Estrutura para guardar as estatísticas de uma marca para um período
type BrandPeriodStats struct {
	Ref                string             `json:"ref"` // Ex: "Abril/2024"
	TabelaId           int                `json:"-"`   // ID da tabela para referência interna
	MenorPreco0km      PriceInfo          `json:"menorPreco0km"`
	MaiorPreco0km      PriceInfo          `json:"maiorPreco0km"`
	ValorMedio0km      float64            `json:"-"`
	ValorMedio0kmFmt   string             `json:"valorMedio0kmFmt"`
	TotalModelos       int                `json:"totalModelos"`
	TotalVeiculos0km   int                `json:"totalVeiculos0km"` // Contagem de veículos 0km encontrados
	somaValores0km     float64            // Usado internamente para calcular média
	modelosEncontrados map[int32]struct{} // Usado internamente para contar modelos únicos
	inicializado       bool               // Flag para saber se já iniciamos min/max
}

// Estrutura para guardar as diferenças percentuais
type PercentageDiffs struct {
	ValorMedio0km *float64 `json:"valorMedio0km,omitempty"`
	TotalModelos  *float64 `json:"totalModelos,omitempty"`
}

// Estrutura final para cada entrada no dashboard
type DashboardBrandEntry struct {
	BrandName             string           `json:"brandName"`
	BrandCode             int32            `json:"brandCode"`
	Periodo1              BrandPeriodStats `json:"periodo1"`
	Periodo2              BrandPeriodStats `json:"periodo2"`
	DiferencasPercentuais PercentageDiffs  `json:"diferencasPercentuais"`
}

// --- Funções Auxiliares ---

// parsePrice converte string de preço (ex: "R$ 123.456,78") para float64
func parsePrice(priceStr string) (float64, error) {
	if priceStr == "" {
		return 0, fmt.Errorf("preço vazio")
	}
	// Remove "R$", espaços, e troca vírgula por ponto
	cleaned := strings.ReplaceAll(priceStr, "R$", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")  // Remove separador de milhar
	cleaned = strings.ReplaceAll(cleaned, ",", ".") // Troca separador decimal
	cleaned = strings.TrimSpace(cleaned)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter '%s' para float: %v", priceStr, err)
	}
	return price, nil
}

// formatPrice converte float64 para string de preço formatada (ex: "R$ 123.456,78")
func formatPrice(price float64) string {
	if math.IsNaN(price) || math.IsInf(price, 0) {
		return "N/A"
	}
	p := message.NewPrinter(language.BrazilianPortuguese)
	return p.Sprintf("R$ %.2f", price)
}

// calculatePercentageDiff calcula a diferença percentual ((v1/v2) - 1) * 100
func calculatePercentageDiff(v1, v2 float64) (*float64, bool) {
	if v2 == 0 || math.IsNaN(v1) || math.IsNaN(v2) || math.IsInf(v1, 0) || math.IsInf(v2, 0) {
		return nil, false // Não é possível calcular
	}
	diff := ((v1 / v2) - 1) * 100
	return &diff, true
}

// getTabelaRef busca o nome do mês/ano da tabela
func getTabelaRef(ctx context.Context, tabelaId int) (string, error) {
	coll := database.DB.Collection("TabelasReferencia")
	var result struct {
		Mes string `bson:"mes"`
	}
	filter := bson.M{"codigo": tabelaId}
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Sprintf("Tabela %d", tabelaId), nil // Retorna ID se não achar nome
		}
		return "", fmt.Errorf("erro ao buscar ref da tabela %d: %v", tabelaId, err)
	}
	return result.Mes, nil
}

// --- Handler Principal ---

func GetDashboardMarcas(w http.ResponseWriter, r *http.Request) {
	// 1. Obter e validar parâmetros tabela1, tabela2 e o opcional marca
	tabela1Param := r.URL.Query().Get("tabela1")
	tabela2Param := r.URL.Query().Get("tabela2")
	marcaParam := r.URL.Query().Get("marca")

	if tabela1Param == "" || tabela2Param == "" {
		http.Error(w, "Parâmetros 'tabela1' e 'tabela2' são obrigatórios", http.StatusBadRequest)
		return
	}

	tabela1Id, err1 := strconv.Atoi(tabela1Param)
	tabela2Id, err2 := strconv.Atoi(tabela2Param)
	if err1 != nil || err2 != nil {
		http.Error(w, "Parâmetros 'tabela1' ou 'tabela2' inválidos", http.StatusBadRequest)
		return
	}

	if tabela1Id == tabela2Id {
		http.Error(w, "Os períodos de comparação devem ser diferentes", http.StatusBadRequest)
		return
	}

	var marcaIdFiltro *int32 // <<< NOVO: Ponteiro para armazenar o ID da marca a filtrar (nil se não filtrar)
	if marcaParam != "" {
		marcaIdTemp, errM := strconv.Atoi(marcaParam)
		if errM != nil {
			http.Error(w, "Parâmetro 'marca' inválido", http.StatusBadRequest)
			return
		}
		temp := int32(marcaIdTemp)
		marcaIdFiltro = &temp // Atribui o endereço da variável temporária
		log.Printf("Iniciando GetDashboardMarcas para tabelas: %d, %d e FILTRANDO pela marca ID: %d", tabela1Id, tabela2Id, *marcaIdFiltro)
	} else {
		log.Printf("Iniciando GetDashboardMarcas para tabelas: %d e %d (todas as marcas)", tabela1Id, tabela2Id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Buscar referências das tabelas (sem alterações aqui)
	var tabela1Ref, tabela2Ref string
	var refErr1, refErr2 error
	var wgRefs sync.WaitGroup
	wgRefs.Add(2)
	go func() { defer wgRefs.Done(); tabela1Ref, refErr1 = getTabelaRef(ctx, tabela1Id) }()
	go func() { defer wgRefs.Done(); tabela2Ref, refErr2 = getTabelaRef(ctx, tabela2Id) }()
	wgRefs.Wait()
	if refErr1 != nil {
		log.Printf("Erro ao buscar referência da tabela %d: %v", tabela1Id, refErr1)
	}
	if refErr2 != nil {
		log.Printf("Erro ao buscar referência da tabela %d: %v", tabela2Id, refErr2)
	}

	// 3. Processar dados de cada tabela concorrentemente
	statsTabela1 := make(map[int32]*BrandPeriodStats)
	statsTabela2 := make(map[int32]*BrandPeriodStats)
	allBrandInfo := make(map[int32]struct {
		Name string
		Code int32
	})

	var processErr1, processErr2 error // Capturar erros das goroutines
	var wgProcess sync.WaitGroup
	wgProcess.Add(2)

	// Função para processar dados de uma tabela (agora recebe marcaIdFiltro)
	processTable := func(tabelaId int, tabelaRef string, targetStats map[int32]*BrandPeriodStats, filterBrand *int32) error { // <<< NOVO: Parâmetro filterBrand
		defer wgProcess.Done()
		collection := database.DB.Collection("Veiculos")

		// <<< MODIFICADO: Adiciona filtro de marca se fornecido
		filter := bson.M{"monthYearId": tabelaId}
		if filterBrand != nil {
			filter["brandCode"] = *filterBrand
			log.Printf("Tabela %d: Aplicando filtro para brandCode: %d", tabelaId, *filterBrand)
		} else {
			log.Printf("Tabela %d: Buscando todas as marcas.", tabelaId)
		}
		// --- Fim da Modificação do Filtro ---

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("erro ao buscar dados da tabela %d com filtro %v: %v", tabelaId, filter, err)
		}
		defer cursor.Close(ctx)

		processedBrands := 0 // Contador para log
		for cursor.Next(ctx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				log.Printf("Erro ao decodificar documento da tabela %d: %v", tabelaId, err)
				continue
			}

			brandCode, okBC := doc["brandCode"].(int32)
			brandName, okBN := doc["brandName"].(string)
			if !okBC || !okBN {
				continue
			}

			processedBrands++

			modelosDisponiveis := 0

			// Armazena informações da marca globalmente (importante mesmo filtrando, para obter o nome)
			if _, exists := allBrandInfo[brandCode]; !exists {
				allBrandInfo[brandCode] = struct {
					Name string
					Code int32
				}{Name: brandName, Code: brandCode}
			}

			// Inicializa estatísticas (igual antes)
			if _, exists := targetStats[brandCode]; !exists {
				targetStats[brandCode] = &BrandPeriodStats{
					Ref:                tabelaRef,
					TabelaId:           tabelaId,
					modelosEncontrados: make(map[int32]struct{}),
					MenorPreco0km:      PriceInfo{Valor: math.Inf(1)},
					MaiorPreco0km:      PriceInfo{Valor: math.Inf(-1)},
				}
			}
			stats := targetStats[brandCode]

			// Processamento de modelos e anos (lógica interna igual antes)
			modelsRaw, exists := doc["models"]
			if !exists {
				continue
			}
			models, ok := modelsRaw.(primitive.A)
			if !ok {
				continue
			}

			for _, modelRaw := range models {
				model, okM := modelRaw.(bson.M)
				if !okM {
					continue
				}
				modelCode, okMC := model["modelCode"].(int32)
				if !okMC {
					continue
				}
				modelName, okMN := model["modelName"].(string)
				if !okMN {
					continue
				}

				stats.modelosEncontrados[modelCode] = struct{}{}

				yearsRaw, existsY := model["years"]
				if !existsY {
					continue
				}
				years, okY := yearsRaw.(primitive.A)
				if !okY {
					continue
				}

				for _, yearRaw := range years {
					yearData, okYD := yearRaw.(bson.M)
					if !okYD {
						continue
					}
					yearVal, okYV := yearData["year"].(int32)
					priceStr, okPS := yearData["price"].(string)

					if yearVal == anoZeroKm {
						modelosDisponiveis++
					}

					if okYV && yearVal == anoZeroKm && okPS {
						price, errP := parsePrice(priceStr)
						if errP != nil {
							continue
						}
						stats.TotalVeiculos0km++
						stats.somaValores0km += price
						if price < stats.MenorPreco0km.Valor {
							stats.MenorPreco0km = PriceInfo{Modelo: modelName, Valor: price, ValorFmt: formatPrice(price)}
						}
						if price > stats.MaiorPreco0km.Valor {
							stats.MaiorPreco0km = PriceInfo{Modelo: modelName, Valor: price, ValorFmt: formatPrice(price)}
						}
						stats.inicializado = true
					}
				}
			}

			// Finaliza cálculos (igual antes)
			stats.TotalModelos = modelosDisponiveis
			if stats.TotalVeiculos0km > 0 {
				stats.ValorMedio0km = stats.somaValores0km / float64(stats.TotalVeiculos0km)
				stats.ValorMedio0kmFmt = formatPrice(stats.ValorMedio0km)
			} else {
				stats.ValorMedio0km = math.NaN()
				stats.ValorMedio0kmFmt = "N/A"
			}
			if !stats.inicializado {
				stats.MenorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats.MaiorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			}
		} // Fim loop cursor

		log.Printf("Tabela %d: Processou %d marcas (documentos).", tabelaId, processedBrands)

		if err := cursor.Err(); err != nil {
			return fmt.Errorf("erro no cursor da tabela %d: %v", tabelaId, err)
		}
		return nil
	}

	// <<< MODIFICADO: Passa marcaIdFiltro para as goroutines
	go func() { processErr1 = processTable(tabela1Id, tabela1Ref, statsTabela1, marcaIdFiltro) }()
	go func() { processErr2 = processTable(tabela2Id, tabela2Ref, statsTabela2, marcaIdFiltro) }()
	wgProcess.Wait() // Espera as goroutines terminarem

	// Verifica erros do processamento (igual antes)
	if processErr1 != nil {
		log.Printf("Erro ao processar tabela 1 (%d): %v", tabela1Id, processErr1) /* Tratar erro */
	}
	if processErr2 != nil {
		log.Printf("Erro ao processar tabela 2 (%d): %v", tabela2Id, processErr2) /* Tratar erro */
	}

	// 4. Combinar resultados e calcular diferenças
	var dashboardResult []DashboardBrandEntry

	// <<< MODIFICADO: Lógica diferente se estiver filtrando por marca
	if marcaIdFiltro != nil {
		// Caso filtrado: Processa apenas a marca especificada
		brandCode := *marcaIdFiltro
		info, brandInfoOk := allBrandInfo[brandCode] // Verifica se a marca foi encontrada em algum período

		if brandInfoOk {
			stats1, ok1 := statsTabela1[brandCode]
			stats2, ok2 := statsTabela2[brandCode]

			// Cria entradas padrão se a marca não existir em um dos períodos (igual antes)
			if !ok1 {
				stats1 = &BrandPeriodStats{Ref: tabela1Ref, TabelaId: tabela1Id, ValorMedio0km: math.NaN()}
				stats1.MenorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.MaiorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.ValorMedio0kmFmt = "N/A"
			}
			if !ok2 {
				stats2 = &BrandPeriodStats{Ref: tabela2Ref, TabelaId: tabela2Id, ValorMedio0km: math.NaN()}
				stats2.MenorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.MaiorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.ValorMedio0kmFmt = "N/A"
			}

			// Calcular diferenças percentuais (igual antes)
			diffs := PercentageDiffs{}
			if diffAvg, ok := calculatePercentageDiff(stats1.ValorMedio0km, stats2.ValorMedio0km); ok {
				diffs.ValorMedio0km = diffAvg
			}
			if diffModels, ok := calculatePercentageDiff(float64(stats1.TotalModelos), float64(stats2.TotalModelos)); ok {
				diffs.TotalModelos = diffModels
			}

			entry := DashboardBrandEntry{
				BrandName:             info.Name,
				BrandCode:             info.Code,
				Periodo1:              *stats1,
				Periodo2:              *stats2,
				DiferencasPercentuais: diffs,
			}
			dashboardResult = append(dashboardResult, entry) // Adiciona a única entrada
		} else {
			log.Printf("Marca filtrada (%d) não encontrada em nenhum dos períodos.", brandCode)
			// Retorna um array vazio, o que é apropriado
		}

	} else {
		// Caso não filtrado: Processa todas as marcas (lógica original)
		brandCodes := make([]int32, 0, len(allBrandInfo))
		for code := range allBrandInfo {
			brandCodes = append(brandCodes, code)
		}
		sort.Slice(brandCodes, func(i, j int) bool { return brandCodes[i] < brandCodes[j] }) // Ordena por código

		for _, brandCode := range brandCodes {
			info := allBrandInfo[brandCode]
			stats1, ok1 := statsTabela1[brandCode]
			stats2, ok2 := statsTabela2[brandCode]

			if !ok1 {
				stats1 = &BrandPeriodStats{Ref: tabela1Ref, TabelaId: tabela1Id, ValorMedio0km: math.NaN()}
				stats1.MenorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.MaiorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.ValorMedio0kmFmt = "N/A"
			}
			if !ok2 {
				stats2 = &BrandPeriodStats{Ref: tabela2Ref, TabelaId: tabela2Id, ValorMedio0km: math.NaN()}
				stats2.MenorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.MaiorPreco0km = PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.ValorMedio0kmFmt = "N/A"
			}

			diffs := PercentageDiffs{}
			if diffAvg, ok := calculatePercentageDiff(stats1.ValorMedio0km, stats2.ValorMedio0km); ok {
				diffs.ValorMedio0km = diffAvg
			}
			if diffModels, ok := calculatePercentageDiff(float64(stats1.TotalModelos), float64(stats2.TotalModelos)); ok {
				diffs.TotalModelos = diffModels
			}

			entry := DashboardBrandEntry{
				BrandName:             info.Name,
				BrandCode:             info.Code,
				Periodo1:              *stats1,
				Periodo2:              *stats2,
				DiferencasPercentuais: diffs,
			}
			dashboardResult = append(dashboardResult, entry)
		}

		// Ordenar resultado final pelo nome da marca (opcional, apenas no caso não filtrado)
		sort.Slice(dashboardResult, func(i, j int) bool { return dashboardResult[i].BrandName < dashboardResult[j].BrandName })
	} // --- Fim da Modificação da Lógica de Combinação ---

	// 5. Enviar resposta JSON (sem alterações aqui)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dashboardResult); err != nil {
		log.Printf("Erro ao encodar resposta JSON: %v", err)
		http.Error(w, "Erro interno ao gerar resposta", http.StatusInternalServerError)
	}
	if marcaIdFiltro != nil {
		log.Printf("GetDashboardMarcas concluído com sucesso para tabelas: %d, %d e marca %d", tabela1Id, tabela2Id, *marcaIdFiltro)
	} else {
		log.Printf("GetDashboardMarcas concluído com sucesso para tabelas: %d, %d (todas as marcas)", tabela1Id, tabela2Id)
	}
}
