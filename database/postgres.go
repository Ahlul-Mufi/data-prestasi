package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DB      *sql.DB
	MongoDB *mongo.Database
)

func ConnectDB() {
	connectPostgres()
	connectMongoDB()
}

func connectPostgres() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "your_password")
	dbname := getEnv("DB_NAME", "prestasi_db")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open PostgreSQL:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping PostgreSQL:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("✅ Connected to PostgreSQL:", dbname)
}

func connectMongoDB() {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGODB_DATABASE", "prestasi_db")

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	err = client.Ping(ctx2, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoDB = client.Database(dbName)

	fmt.Println("✅ Connected to MongoDB:", dbName)
}

func DisconnectDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Println("Error closing PostgreSQL:", err)
		} else {
			fmt.Println("✅ PostgreSQL disconnected")
		}
	}

	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := MongoDB.Client().Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		} else {
			fmt.Println("✅ MongoDB disconnected")
		}
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	if MongoDB == nil {
		log.Fatal("MongoDB is not connected")
	}
	return MongoDB.Collection(collectionName)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
