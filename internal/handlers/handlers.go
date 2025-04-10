package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"fipe_project/internal/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	log.Printf("Parâmetros recebidos: tabelaId=%d, modeloId=%d", tabelaId, modeloId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("Veiculos")


	filter := bson.M{
		"monthYearId":      tabelaId,
		"models.modelCode": modeloId,
	}

	log.Printf("Filtro usado na consulta: %v", filter)

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

	log.Printf("Anos: %v", selectedYears)

	if len(selectedYears) == 0 {
		log.Printf("Nenhum ano encontrado para o modelo especificado.")
		http.Error(w, "Anos não encontrados para o modelo especificado", http.StatusNotFound)
		return
	}

	log.Printf("Anos selecionados para o modelo: %v", selectedYears)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedYears)
}
