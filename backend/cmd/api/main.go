package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"ecommerce-grupo3/backend/internal/config"
	"ecommerce-grupo3/backend/internal/database"
	"ecommerce-grupo3/backend/internal/products"
	"ecommerce-grupo3/backend/internal/promotions"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var pool *pgxpool.Pool
	dbStatus := "Memoria (Local Fallback)"

	p, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("[ADVERTENCIA] No se pudo conectar a PostgreSQL (%v).", err)
		log.Println("[INFO] Modo de desarrollo activado: usando repositorio en memoria con datos semilla precargados.")
	} else {
		pool = p
		dbStatus = "PostgreSQL Conectado"
		defer pool.Close()
	}

	// Repositorios y Handlers para los 2 CRUDs
	productRepo := products.NewRepository(pool)
	productHandler := products.NewHandler(productRepo)

	promotionRepo := promotions.NewRepository(pool)
	promotionHandler := promotions.NewHandler(promotionRepo, productRepo)

	mux := http.NewServeMux()

	// Endpoint de salud
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":    "ok",
			"proyecto":  "E-commerce Grupo 3",
			"database":  dbStatus,
			"patrones":  []string{"Composite (Productos/Combos)", "Decorator (Promociones/Descuentos)"},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Registrar rutas de los dos CRUDs
	productHandler.RegisterRoutes(mux)
	promotionHandler.RegisterRoutes(mux)

	// Arquitectura Monolítica: Si existe la carpeta dist del frontend, el backend Go lo sirve directamente
	frontendDist := findFrontendDist()
	if frontendDist != "" {
		log.Printf("[INFO] Servidor web estático habilitado desde: %s", frontendDist)
		fileServer := http.FileServer(http.Dir(frontendDist))
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.NotFound(w, r)
				return
			}
			targetPath := filepath.Join(frontendDist, filepath.Clean(r.URL.Path))
			if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join(frontendDist, "index.html"))
		})
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           cors(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("==================================================")
		log.Printf(" Servidor E-commerce Grupo 3 listo en http://localhost:%s", cfg.Port)
		log.Printf(" Base de Datos: %s", dbStatus)
		log.Printf(" CRUD 1: http://localhost:%s/api/productos (Composite)", cfg.Port)
		log.Printf(" CRUD 2: http://localhost:%s/api/promociones (Decorator)", cfg.Port)
		if frontendDist != "" {
			log.Printf(" Frontend Monolito: http://localhost:%s/", cfg.Port)
		}
		log.Printf("==================================================")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error del servidor: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}

func findFrontendDist() string {
	candidates := []string{
		"../frontend/dist",
		"./frontend/dist",
		"../../frontend/dist",
		"./dist",
	}
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return ""
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
