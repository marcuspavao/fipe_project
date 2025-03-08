package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"fipe_project/database"

	// [CHANGED] Certifique-se de ter todos os imports do mongo driver
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// [CHANGED - Definimos o tipo ModelAPI para um modelo individual]
type ModelAPI struct {
	Label string `json:"Label"`
	Value int    `json:"Value"`
}

// [CHANGED - Definimos o tipo ModelosAPI, que contém o array de ModelAPI]
type ModelosAPI struct {
	Modelos []ModelAPI `json:"Modelos"`
}

// -----------------------------------------------------------------------------

func LoadData() error {
	ctx := context.Background()
	veiculos := database.DB.Collection("Veiculos")
	tabRef := database.DB.Collection("TabelaReferencia")

	// 1. Tabela de referência
	codigoTabela, err := fetchLatestTabelaReferencia(ctx, &database.CollectionWrapper{Collection: tabRef})
	if err != nil {
		return fmt.Errorf("falha em fetchLatestTabelaReferencia: %v", err)
	}

	// 2. Busca todas as marcas na API
	payloadMarcas := map[string]interface{}{
		"codigoTabelaReferencia": codigoTabela,
		"codigoTipoVeiculo":      1,
	}
	dataMarcas, err := makePostRequest("https://veiculos.fipe.org.br/api/veiculos/ConsultarMarcas", payloadMarcas)
	if err != nil {
		return err
	}
	var marcasAPI []struct {
		Label string `json:"Label"`
		Value string `json:"Value"`
	}
	if err := json.Unmarshal(dataMarcas, &marcasAPI); err != nil {
		return err
	}

	// 3. Loop nas marcas da API
	for _, marca := range marcasAPI {
		codMarca, _ := strconv.Atoi(marca.Value)

		// 3.1. Busca quantos modelos existem na API pra essa marca
		modPayload := map[string]interface{}{
			"codigoTabelaReferencia": codigoTabela,
			"codigoTipoVeiculo":      1,
			"codigoMarca":            codMarca,
		}
		modData, err := makePostRequest("https://veiculos.fipe.org.br/api/veiculos/ConsultarModelos", modPayload)
		if err != nil {
			log.Printf("Erro ao consultar modelos da marca %s: %v", marca.Label, err)
			continue
		}

		// [CHANGED - Em vez de struct anônimo, usamos ModelosAPI]
		var modStruct ModelosAPI
		if err := json.Unmarshal(modData, &modStruct); err != nil {
			log.Printf("Erro parse modelos da marca %s: %v", marca.Label, err)
			continue
		}
		totalModelosAPI := len(modStruct.Modelos)

		// 3.2. Verifica no banco se a marca já existe
		var brandDoc bson.M
		errFind := veiculos.FindOne(ctx, bson.M{
			"brandCode":   codMarca,
			"monthYearId": codigoTabela,
		}).Decode(&brandDoc)

		// Se não existe no banco ou houve "not found", processa tudo normalmente
		if errFind != nil {
			// Insere a marca e todos os modelos (já que não existe nada)
			log.Printf("Marca %s (%d) não encontrada ou erro => processando do zero", marca.Label, codMarca)
			errIns := processarMarca(ctx, veiculos, codigoTabela, codMarca, marca.Label, modStruct)
			if errIns != nil {
				log.Printf("Erro ao processar marca %s: %v", marca.Label, errIns)
			}
			continue
		}
		// [ADICIONADO] Exibe todo o documento recuperado para debug
		//log.Printf("brandDoc recuperado do BD para marca %s (%d): %v", marca.Label, codMarca, brandDoc)

		// 3.3. Se a marca existe, checar se já temos o mesmo total de modelos gravados
		modelsBsonA, ok := brandDoc["models"].(bson.A)
		modelsBD := []interface{}(modelsBsonA)
		if !ok {
			// Se "models" não for array ou algo deu errado, processa do zero
			log.Printf("Marca %s encontrada, mas sem campo 'models' ou não é array => processando tudo do zero", marca.Label)
			errIns := processarMarca(ctx, veiculos, codigoTabela, codMarca, marca.Label, modStruct)
			if errIns != nil {
				log.Printf("Erro ao processar marca %s: %v", marca.Label, errIns)
			}
			continue
		}

		// [ADICIONADO] Exibe o conteúdo de modelsBD
		log.Printf("Marca %s => modelsBD (já no banco)", marca.Label)

		// 3.4. Se o total de modelos na API == total de modelos no BD => pula
		if len(modelsBD) >= totalModelosAPI {
			log.Printf("Marca %s (%d) já possui todos os modelos (%d). Pulando.", marca.Label, codMarca, totalModelosAPI)
			continue
		}

		// 3.5. Se faltam modelos, processa só o que falta (comparando um a um)

		log.Printf("Marca %s (%d) incompleta: BD tem %d de %d modelos. Processando os faltantes...",
			marca.Label, codMarca, len(modelsBD), totalModelosAPI)

		errIns := processarModelosFaltantes(ctx, veiculos, codigoTabela, codMarca, marca.Label, modelsBD, modStruct)
		if errIns != nil {
			log.Printf("Erro ao processar modelos faltantes da marca %s: %v", marca.Label, errIns)
		}
	}

	log.Println("Carregamento concluído.")
	return nil
}

// --------------------------------------------------------------------
// Se marca não existe no BD => Upsert da marca + processarModelosFaltantes
func processarMarca(
	ctx context.Context,
	veiculos *mongo.Collection,
	codigoTabela, codMarca int,
	nomeMarca string,
	modStruct ModelosAPI,
) error {

	// Marca nova: Upsert com "models": []
	_, err := veiculos.UpdateOne(
		ctx,
		bson.M{"brandCode": codMarca, "monthYearId": codigoTabela},
		bson.M{"$set": bson.M{
			"brandName":   nomeMarca,
			"brandCode":   codMarca,
			"monthYearId": codigoTabela,
			"models":      []bson.M{},
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir marca %s (%d): %v", nomeMarca, codMarca, err)
	}

	// Chama a mesma função de inserir modelos, mas sem nenhum modelo já existente (modelsBD = nil)
	return processarModelosFaltantes(ctx, veiculos, codigoTabela, codMarca, nomeMarca, nil, modStruct)
}

// --------------------------------------------------------------------
// Se a marca existe parcialmente => processar apenas os modelos faltantes
func processarModelosFaltantes(
	ctx context.Context,
	veiculos *mongo.Collection,
	codigoTabela, codMarca int,
	nomeMarca string,
	modelsBD []interface{}, // se nil ou vazio, insere tudo
	modStruct ModelosAPI,
) error {

	// 1) Monta um map com os modelCode existentes (se houverem)
	existing := make(map[int]bool)

	for _, m := range modelsBD {
		// Log do elemento bruto
		// Verifique se é do tipo bson.M
		doc, ok := m.(bson.M)
		if !ok {
			log.Printf("Skipping non-bson.M element: %v", m)
			continue
		}
		// Tente acessar modelCode como int
		rawCode := doc["modelCode"]
		switch v := rawCode.(type) {
		case int32:
			existing[int(v)] = true
		default:
			log.Printf("Tipo inesperado para modelCode: %T, Valor: %v", rawCode, rawCode)
		}
	}

	// 2) Percorre os modelos da API
	for _, modelo := range modStruct.Modelos {
		// Se já existe no BD, pula
		if existing[modelo.Value] {
			continue
		}
		log.Printf("Inserindo modelo %s", modelo.Label)

		// (A) CONSULTAR anos do modelo
		anosPayload := map[string]interface{}{
			"codigoTabelaReferencia": codigoTabela,
			"codigoTipoVeiculo":      1,
			"codigoMarca":            codMarca,
			"codigoModelo":           modelo.Value,
		}
		anosData, err := makePostRequest("https://veiculos.fipe.org.br/api/veiculos/ConsultarAnoModelo", anosPayload)
		if err != nil {
			log.Printf("Erro ao consultar anos do modelo %s (%d): %v", modelo.Label, modelo.Value, err)
			continue
		}
		var anos []struct {
			Label string `json:"Label"`
			Value string `json:"Value"`
		}
		if err := json.Unmarshal(anosData, &anos); err != nil {
			log.Printf("Erro parse anos do modelo %s: %v", modelo.Label, err)
			continue
		}

		// (B) Para cada ano, consultar valor
		var anosArr []bson.M
		for _, ano := range anos {

			if len(ano.Value) < 3 {
				continue
			}
			anoModelo := ano.Value[:len(ano.Value)-2]
			codComb := ano.Value[len(ano.Value)-1:]
			anoInt, _ := strconv.Atoi(anoModelo)

			valorPayload := map[string]interface{}{
				"codigoTabelaReferencia": codigoTabela,
				"codigoTipoVeiculo":      1,
				"codigoMarca":            codMarca,
				"codigoModelo":           modelo.Value,
				"anoModelo":              anoModelo,
				"codigoTipoCombustivel":  codComb,
				"ano":                    ano.Value,
				"tipoConsulta":           "tradicional",
			}
			valData, err := makePostRequest("https://veiculos.fipe.org.br/api/veiculos/ConsultarValorComTodosParametros", valorPayload)
			if err != nil {
				log.Printf("Erro preco ano %s modelo %s: %v", ano.Value, modelo.Label, err)
				continue
			}
			var val struct {
				Modelo        string `json:"Modelo"`
				Valor         string `json:"Valor"`
				AnoModelo     int    `json:"AnoModelo"`
				MesReferencia string `json:"MesReferencia"`
			}
			if err := json.Unmarshal(valData, &val); err != nil {
				log.Printf("Erro parse preco ano %s: %v", ano.Value, err)
				continue
			}

			// Monta item do array de anos
			anosArr = append(anosArr, bson.M{
				"yearCode":       ano.Value,
				"year":           anoInt,
				"price":          val.Valor,
				"monthReference": val.MesReferencia,
			})
		}

		// (C) Documento do modelo, incluindo os anos
		modeloDoc := bson.M{
			"modelCode": modelo.Value,
			"modelName": modelo.Label,
			"years":     anosArr,
		}

		// (D) Push no array "models" da marca
		_, err = veiculos.UpdateOne(
			ctx,
			bson.M{"brandCode": codMarca, "monthYearId": codigoTabela},
			bson.M{"$push": bson.M{"models": modeloDoc}},
		)
		if err != nil {
			log.Printf("Erro ao inserir modelo %s na marca %s: %v", modelo.Label, nomeMarca, err)
		}
	}
	return nil
}

// makePostRequest com retry básico
func makePostRequest(url string, payload map[string]interface{}) ([]byte, error) {
	maxRetries := 5
	backoffBase := 1 * time.Second
	for attempt := 1; attempt <= maxRetries; attempt++ {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("erro ao codificar payload: %v", err)
		}
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Tentativa %d falhou: %v", attempt)
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusOK {
				log.Printf("Tentativa %d: status %d => Ok", attempt, resp.StatusCode)
				time.Sleep(backoffBase)
				return body, nil
			}
			log.Printf("Tentativa %d: status %d", attempt, resp.StatusCode)
		}
		time.Sleep(time.Duration(attempt) * backoffBase)
	}
	return nil, fmt.Errorf("todas as tentativas falharam")
}

// fetchLatestTabelaReferencia retorna o codigo de referencia atual
func fetchLatestTabelaReferencia(ctx context.Context, col *database.CollectionWrapper) (int, error) {
	url := "https://veiculos.fipe.org.br/api/veiculos/ConsultarTabelaDeReferencia"
	data, err := makePostRequest(url, map[string]interface{}{})
	if err != nil {
		return 0, err
	}
	var tabelas []struct {
		Codigo int    `json:"Codigo"`
		Mes    string `json:"Mes"`
	}
	if err := json.Unmarshal(data, &tabelas); err != nil {
		return 0, err
	}
	if len(tabelas) == 0 {
		return 0, fmt.Errorf("nenhuma tabela encontrada")
	}
	_, err = col.UpdateOne(
		ctx,
		bson.M{"codigo": tabelas[0].Codigo},
		bson.M{"$set": bson.M{
			"codigo": tabelas[0].Codigo,
			"mes":    tabelas[0].Mes,
		}},
		options.Update().SetUpsert(true),
	)
	return tabelas[0].Codigo, err
}
