package main

import (
    "log"

    "fipe_project/database"
    "fipe_project/services"
)

func main() {
    // Conecta ao MongoDB
    if err := database.ConnectMongoDB(); err != nil {
        log.Fatalf("Erro ao conectar no MongoDB: %v", err)
    }

    // Executa a rotina de carregamento FIPE
    if err := services.LoadData(); err != nil {
        log.Printf("Erro ao carregar dados FIPE: %v", err)
    } else {
        log.Println("Dados FIPE carregados com sucesso!")
    }

    // Mantém a aplicação ativa (sem servidor HTTP neste exemplo)
    select {}
}