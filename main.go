package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func loadConfig() Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	return Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}
}

func main() {
	cfg := loadConfig()
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBConn == "" {
		log.Fatal("DB_CONN kosong. Pastikan .env kebaca.")
	}

	// Init DB pool (pgxpool)
	dbPool, err := database.InitDBPool(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbPool.Close()

	// DI
	productRepo := repositories.NewProductRepository(dbPool)
	productSvc := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	categoryRepo := repositories.NewCategoryRepository(dbPool)
	categorySvc := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categorySvc)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(dbPool)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Report
	reportRepo := repositories.NewReportRepository(dbPool)
	reportSvc := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportSvc)

	// Routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleHariIni)
	http.HandleFunc("/api/report", reportHandler.HandleReportRange) // optional

	addr := "0.0.0.0:" + cfg.Port
	fmt.Println("Server running di", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
