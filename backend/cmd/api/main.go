package main

import (
	"log"
	"net/http"
	"os"

	"ingsw3-tp7-tp8-integrated/backend/internal/database"
	"ingsw3-tp7-tp8-integrated/backend/internal/handlers"
	"ingsw3-tp7-tp8-integrated/backend/internal/repository"
	"ingsw3-tp7-tp8-integrated/backend/internal/router"
	"ingsw3-tp7-tp8-integrated/backend/internal/services"
)

func main() {
	// Obtener URL de base de datos desde variable de entorno
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Inicializar base de datos
	db, err := database.InitDB(databaseURL)
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}
	defer db.Close()

	// Crear repositorios
	userRepo := repository.NewPostgreSQLUserRepository(db)
	postRepo := repository.NewPostgreSQLPostRepository(db)

	// Crear servicios
	authService := services.NewAuthService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)

	// Crear handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)

	// Configurar rutas
	r := router.Setup(authHandler, postHandler)

	// Definir puerto desde variable de entorno o default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
