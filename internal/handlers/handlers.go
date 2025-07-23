package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fipe_project/internal/database"
	"fipe_project/internal/models"
	"fipe_project/internal/utils"
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
	defer cursor.Close(ctx)

	var wg sync.WaitGroup
	tabelasChan := make(chan bson.M)
	errChan := make(chan error, 1)

	go func() {
		defer close(tabelasChan)
		for cursor.Next(ctx) {
			var tabela bson.M
			if err := cursor.Decode(&tabela); err != nil {
				errChan <- fmt.Errorf("erro ao decodificar tabela de referência: %v", err)
				return
			}
			wg.Add(1)
			go func(tabela bson.M) {
				defer wg.Done()
				codigo, ok := tabela["codigo"].(int32)
				if !ok {
					log.Printf("Código da tabela não encontrado ou não é int32: %v", tabela)
					return
				}
				temVeiculos, err := TabelaTemVeiculos(ctx, int(codigo))
				if err != nil {
					log.Printf("Erro ao verificar veículos para a tabela %d: %v", codigo, err)
					return
				}
				if temVeiculos {
					tabelasChan <- tabela
				}
			}(tabela)
		}
	}()

	var tabelasFiltradas []bson.M
	done := make(chan struct{})

	go func() {
		for tabela := range tabelasChan {
			tabelasFiltradas = append(tabelasFiltradas, tabela)
		}
		close(done)
	}()

	wg.Wait()
	select {
	case err := <-errChan:
		log.Printf("Erro no processamento das tabelas: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	default:
	}

	<-done

	if err := cursor.Err(); err != nil {
		log.Printf("Erro no cursor ao final: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tabelasFiltradas)
}

func TabelaTemVeiculos(ctx context.Context, tabelaId int) (bool, error) {
	collection := database.DB.Collection("Veiculos")
	filter := bson.M{"monthYearId": tabelaId}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("Veiculos")

	filter := bson.M{
		"monthYearId":      tabelaId,
		"models.modelCode": modeloId,
	}

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

	if len(selectedYears) == 0 {
		log.Printf("Nenhum ano encontrado para o modelo especificado.")
		http.Error(w, "Anos não encontrados para o modelo especificado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedYears)
}

// Dashboard de Marcas - de acordo com as marcas analisar para dois períodos
// infomações como carro/modelo com menor e maior preço 0km, valor médio,
// número de modelos disponíveis e as difenças em porcentagens entre esses aspectos

func getTabelaRef(ctx context.Context, tabelaId int) (string, error) {
	coll := database.DB.Collection("TabelaReferencia")
	var result struct {
		Mes string `bson:"mes"`
	}

	log.Printf("Log tabela: %d", tabelaId)

	filter := bson.M{"codigo": tabelaId}
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Sprintf("Tabela %d", tabelaId), nil
		}
		return "", fmt.Errorf("erro ao buscar ref da tabela %d: %v", tabelaId, err)
	}
	return result.Mes, nil
}

func GetDashboardMarcas(w http.ResponseWriter, r *http.Request) {

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

	var marcaIdFiltro *int32
	if marcaParam != "" {
		marcaIdTemp, errM := strconv.Atoi(marcaParam)
		if errM != nil {
			http.Error(w, "Parâmetro 'marca' inválido", http.StatusBadRequest)
			return
		}
		temp := int32(marcaIdTemp)
		marcaIdFiltro = &temp
		log.Printf("Iniciando GetDashboardMarcas para tabelas: %d, %d e FILTRANDO pela marca ID: %d", tabela1Id, tabela2Id, *marcaIdFiltro)
	} else {
		log.Printf("Iniciando GetDashboardMarcas para tabelas: %d e %d (todas as marcas)", tabela1Id, tabela2Id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	statsTabela1 := make(map[int32]*models.BrandPeriodStats)
	statsTabela2 := make(map[int32]*models.BrandPeriodStats)
	BrandInfo := make(map[int32]struct {
		Name string
		Code int32
	})

	var processErr1, processErr2 error
	var wgProcess sync.WaitGroup
	wgProcess.Add(2)

	processTable := func(tabelaId int, tabelaRef string, targetStats map[int32]*models.BrandPeriodStats, filterBrand *int32) error {
		defer wgProcess.Done()
		collection := database.DB.Collection("Veiculos")

		filter := bson.M{"monthYearId": tabelaId}
		if filterBrand != nil {
			filter["brandCode"] = *filterBrand
			log.Printf("Tabela %d: Aplicando filtro para brandCode: %d", tabelaId, *filterBrand)
		} else {
			log.Printf("Tabela %d: Buscando todas as marcas.", tabelaId)
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("erro ao buscar dados da tabela %d com filtro %v: %v", tabelaId, filter, err)
		}
		defer cursor.Close(ctx)

		processedBrands := 0
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

			log.Printf("Verificando se existe a marca %d", brandCode)

			if _, exists := BrandInfo[brandCode]; !exists {
				BrandInfo[brandCode] = struct {
					Name string
					Code int32
				}{Name: brandName, Code: brandCode}
			}

			log.Printf("Tabela %s", tabelaRef)

			if _, exists := targetStats[brandCode]; !exists {
				targetStats[brandCode] = &models.BrandPeriodStats{
					Ref:                tabelaRef,
					TabelaId:           tabelaId,
					ModelosEncontrados: make(map[int32]struct{}),
					MenorPreco0km:      models.PriceInfo{Valor: math.Inf(1)},
					MaiorPreco0km:      models.PriceInfo{Valor: math.Inf(-1)},
				}
			}
			stats := targetStats[brandCode]

			modelsRaw, exists := doc["models"]
			if !exists {
				continue
			}
			modelsList, ok := modelsRaw.(primitive.A)
			if !ok {
				continue
			}

			for _, modelRaw := range modelsList {
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

				stats.ModelosEncontrados[modelCode] = struct{}{}

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

					if yearVal == models.AnoZeroKm {
						modelosDisponiveis++
					}

					if okYV && yearVal == models.AnoZeroKm && okPS {
						price, errP := utils.ParsePrice(priceStr)
						if errP != nil {
							continue
						}
						stats.TotalVeiculos0km++
						stats.SomaValores0km += price
						if price < stats.MenorPreco0km.Valor {
							stats.MenorPreco0km = models.PriceInfo{Modelo: modelName, Valor: price, ValorFmt: utils.FormatPrice(price)}
						}
						if price > stats.MaiorPreco0km.Valor {
							stats.MaiorPreco0km = models.PriceInfo{Modelo: modelName, Valor: price, ValorFmt: utils.FormatPrice(price)}
						}
						stats.Inicializado = true
					}
				}
			}

			stats.TotalModelos = modelosDisponiveis
			if stats.TotalVeiculos0km > 0 {
				stats.ValorMedio0km = stats.SomaValores0km / float64(stats.TotalVeiculos0km)
				stats.ValorMedio0kmFmt = utils.FormatPrice(stats.ValorMedio0km)
			} else {
				stats.ValorMedio0km = math.NaN()
				stats.ValorMedio0kmFmt = "N/A"
			}
			if !stats.Inicializado {
				stats.MenorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats.MaiorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			}
		}

		log.Printf("Tabela %d: Processou %d marcas (documentos).", tabelaId, processedBrands)

		if err := cursor.Err(); err != nil {
			return fmt.Errorf("erro no cursor da tabela %d: %v", tabelaId, err)
		}
		return nil
	}

	go func() { processErr1 = processTable(tabela1Id, tabela1Ref, statsTabela1, marcaIdFiltro) }()
	go func() { processErr2 = processTable(tabela2Id, tabela2Ref, statsTabela2, marcaIdFiltro) }()
	wgProcess.Wait()

	if processErr1 != nil {
		log.Printf("Erro ao processar tabela 1 (%d): %v", tabela1Id, processErr1)
	}
	if processErr2 != nil {
		log.Printf("Erro ao processar tabela 2 (%d): %v", tabela2Id, processErr2)
	}

	var dashboardResult []models.DashboardBrandEntry

	if marcaIdFiltro != nil {
		brandCode := *marcaIdFiltro
		info, ok := BrandInfo[brandCode]
		if !ok {
			http.Error(w, "Marca não encontrada nos períodos especificados", http.StatusNotFound)
			return
		}

		stats1, ok1 := statsTabela1[brandCode]
		stats2, ok2 := statsTabela2[brandCode]

		if !ok1 {
			stats1 = &models.BrandPeriodStats{Ref: tabela1Ref, TabelaId: tabela1Id, ValorMedio0km: math.NaN()}
			stats1.MenorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			stats1.MaiorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			stats1.ValorMedio0kmFmt = "N/A"
		}
		if !ok2 {
			stats2 = &models.BrandPeriodStats{Ref: tabela2Ref, TabelaId: tabela2Id, ValorMedio0km: math.NaN()}
			stats2.MenorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			stats2.MaiorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
			stats2.ValorMedio0kmFmt = "N/A"
		}
		diffs := models.PercentageDiffs{}
		if diffAvg, ok := utils.CalculatePercentageDiff(stats1.ValorMedio0km, stats2.ValorMedio0km); ok {
			diffs.ValorMedio0km = diffAvg
		}
		if diffModels, ok := utils.CalculatePercentageDiff(float64(stats1.TotalModelos), float64(stats2.TotalModelos)); ok {
			diffs.TotalModelos = diffModels
		}
		entry := models.DashboardBrandEntry{
			BrandName:             info.Name,
			BrandCode:             info.Code,
			Periodo1:              *stats1,
			Periodo2:              *stats2,
			DiferencasPercentuais: diffs,
		}
		dashboardResult = append(dashboardResult, entry)
	} else {
		for brandCode, info := range BrandInfo {
			stats1, ok1 := statsTabela1[brandCode]
			stats2, ok2 := statsTabela2[brandCode]

			if !ok1 {
				stats1 = &models.BrandPeriodStats{Ref: tabela1Ref, TabelaId: tabela1Id, ValorMedio0km: math.NaN()}
				stats1.MenorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.MaiorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats1.ValorMedio0kmFmt = "N/A"
			}
			if !ok2 {
				stats2 = &models.BrandPeriodStats{Ref: tabela2Ref, TabelaId: tabela2Id, ValorMedio0km: math.NaN()}
				stats2.MenorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.MaiorPreco0km = models.PriceInfo{Modelo: "N/A", ValorFmt: "N/A", Valor: math.NaN()}
				stats2.ValorMedio0kmFmt = "N/A"
			}

			diffs := models.PercentageDiffs{}
			if diffAvg, ok := utils.CalculatePercentageDiff(stats1.ValorMedio0km, stats2.ValorMedio0km); ok {
				diffs.ValorMedio0km = diffAvg
			}
			if diffModels, ok := utils.CalculatePercentageDiff(float64(stats1.TotalModelos), float64(stats2.TotalModelos)); ok {
				diffs.TotalModelos = diffModels
			}

			entry := models.DashboardBrandEntry{
				BrandName:             info.Name,
				BrandCode:             info.Code,
				Periodo1:              *stats1,
				Periodo2:              *stats2,
				DiferencasPercentuais: diffs,
			}
			dashboardResult = append(dashboardResult, entry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dashboardResult); err != nil {
		log.Printf("Erro ao encodar resposta JSON: %v", err)
		http.Error(w, "Erro interno ao gerar resposta", http.StatusInternalServerError)
	}
	log.Printf("GetDashboardMarcas concluído com sucesso para tabelas: %d, %d e marca %d/n", tabela1Id, tabela2Id, *marcaIdFiltro)
}

func GetVeiculosNovos(w http.ResponseWriter, r *http.Request) {

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

	filter := bson.M{
		"monthYearId": tabelaId,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Erro ao decodificar documento da tabela %d: %v", tabelaId, err)
	}

	var selectedYears []bson.M
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Erro ao decodificar documento da tabela %d: %v", tabelaId, err)
			continue
		}

		modelsRaw, exists := doc["models"]
		if !exists {
			continue
		}

		modelsList, ok := modelsRaw.(primitive.A)
		if !ok {
			continue
		}

		for _, model := range modelsList {
			m, ok := model.(bson.M)
			if ok {
				if years, exists := m["years"].(primitive.A); exists {
					for _, year := range years {
						if yearMap, isMap := year.(bson.M); isMap {
							if yearMap["year"] == int32(32000) {
								yearMap["model"] = m["modelName"]
								selectedYears = append(selectedYears, yearMap)
							}
						}
					}
				} else {
					log.Printf("Campo 'years' não encontrado ou não é um array no modelo selecionado.")
				}
			}
		}
	}

	if len(selectedYears) == 0 {
		log.Printf("Nenhum modelo encontrado com o ano especificado.")
		http.Error(w, "Modelos não encontrados com o ano especificado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedYears)
}
