package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/JuanPidarraga/talkus-backend/config"
	"github.com/JuanPidarraga/talkus-backend/internal/controllers"
	"github.com/JuanPidarraga/talkus-backend/internal/handlers"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
	"github.com/JuanPidarraga/talkus-backend/internal/service"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
)

func main() {

	// Inicializar Firebase (con credenciales definidas en la variable de entorno FIREBASE_CREDENTIALS_PATH)
	firebaseApp, err := config.InitFirebase()
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}
	defer firebaseApp.Firestore.Close()

	authService := service.NewAuthService(firebaseApp)
    authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Usar Gorilla Mux para definir rutas
	router := mux.NewRouter()

	publicRouter := router.PathPrefix("/public").Subrouter()

	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)
	protectedRouter.HandleFunc("/profile", authHandler.GetUserProfile)

	
	

	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Acccept", "Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders: []string{"Content-Length","Content-Type"},
		AllowCredentials: true,
		MaxAge: 300,
	}

	handler := cors.New(corsOptions).Handler(router)
	serverPort := ":8080"

	// 1. Repositorio
	userRepo := repositories.NewUserRepository(firebaseApp.Firestore)
	// 2. Usecase (lógica de negocio)
	userUsecase := usecases.NewUserUsecase(userRepo)
	// 3. Controller (HTTP handlers)
	userController := controllers.NewUserController(userUsecase)

	
	// Registro de ruta para obtener usuario
	publicRouter.HandleFunc("/users", userController.GetUser).Methods("GET")

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Ruta registrada:", path)
		}
		return nil
	})
	

	// Iniciar servidor HTTP
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
