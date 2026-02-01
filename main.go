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
	// 1) Load config
	cfg := loadConfig()
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// 2) Setup DB
	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 3) Dependency Injection (wiring)
	productRepo := repositories.NewProductRepository(db)
	productSvc := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	categoryRepo := repositories.NewCategoryRepository(db)
	categorySvc := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categorySvc)

	// 4) Routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Products
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Categories
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	addr := "0.0.0.0:" + cfg.Port
	fmt.Println("Server running di", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
