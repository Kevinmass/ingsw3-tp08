package main

import (
	"log"
	"net/http"
	"os"

	"tp08-testing/internal/database"
	"tp08-testing/internal/handlers"
	"tp08-testing/internal/repository"
	"tp08-testing/internal/router"
	"tp08-testing/internal/services"
)

func main() {
	// Inicializar base de datos
	db, err := database.InitDB("./database.db")
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}
	defer db.Close()

	// Crear repositorios
	userRepo := repository.NewSQLiteUserRepository(db)
	postRepo := repository.NewSQLitePostRepository(db)

	// Crear servicios
	authService := services.NewAuthService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)

	// Crear handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)

	// Configurar rutas
	r := router.Setup(authHandler, postHandler)

	// Leer puerto de variable de entorno (Render la define automáticamente)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
