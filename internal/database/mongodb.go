package database

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

type CollectionWrapper struct {
    *mongo.Collection
}

func ConnectMongoDB() error {
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://mongo:27017"
    }
    dbName := os.Getenv("MONGO_DATABASE")
    if dbName == "" {
        dbName = "fipe_db"
    }

    clientOpts := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)
    client, err := mongo.NewClient(clientOpts)
    if err != nil {
        log.Printf("Erro ao criar cliente: %v", err)
        return err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err = client.Connect(ctx); err != nil {
        log.Printf("Erro ao conectar: %v", err)
        return err
    }
    if err = client.Ping(ctx, nil); err != nil {
        log.Printf("Ping falhou: %v", err)
        return err
    }

    DB = client.Database(dbName)
    log.Println("Conexão com MongoDB ok:", uri, dbName)
    return nil
}

func GetCollection(name string) *CollectionWrapper {
    return &CollectionWrapper{DB.Collection(name)}
}