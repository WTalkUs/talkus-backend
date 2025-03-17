package main

import (
	"context"
	"log"
	"net/http"


	"github.com/gorilla/mux"


	"github.com/JuanPidarraga/talkus-backend/config"
	"github.com/JuanPidarraga/talkus-backend/internal/controllers"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
)

func main() {
	ctx := context.Background()

	// Inicializar Firebase (con credenciales definidas en la variable de entorno FIREBASE_CREDENTIALS_PATH)
	firebaseApp, err := config.InitFirebase()
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}
	defer firebaseApp.Firestore.Close()

	// Verificar conexión (opcional, para debug)
	docs, err := firebaseApp.Firestore.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Error accediendo a Firestore: %v", err)
	}
	log.Printf("Conexión exitosa. Usuarios en Firestore: %d", len(docs))

	// Inicialización de las capas según Clean Architecture:
	// 1. Repositorio
	userRepo := repositories.NewUserRepository(firebaseApp.Firestore)
	// 2. Usecase (lógica de negocio)
	userUsecase := usecases.NewUserUsecase(userRepo)
	// 3. Controller (HTTP handlers)
	userController := controllers.NewUserController(userUsecase)

	// Usar Gorilla Mux para definir rutas
	router := mux.NewRouter()
	// Registro de ruta para obtener usuario
	router.HandleFunc("/users", userController.GetUser).Methods("GET")

	// Iniciar servidor HTTP
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
