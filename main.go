package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/JuanPidarraga/talkus-backend/config"
	_ "github.com/JuanPidarraga/talkus-backend/docs"
	"github.com/JuanPidarraga/talkus-backend/internal/controllers"
	"github.com/JuanPidarraga/talkus-backend/internal/handlers"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
	"github.com/JuanPidarraga/talkus-backend/internal/service"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/cloudinary/cloudinary-go/v2"
)

func main() {

	// Inicializar Firebase (con credenciales definidas en la variable de entorno FIREBASE_CREDENTIALS_PATH)
	firebaseApp, err := config.InitFirebase()
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}
	defer firebaseApp.Firestore.Close()

	// Inicializar Cloudinary (con credenciales definidas en la variable de entorno CLOUDINARY_URL)
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Error iniciando Cloudinary: %v", err)
	}

	authService := service.NewAuthService(firebaseApp)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	userRepo := repositories.NewUserRepository(firebaseApp.Firestore)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userController := controllers.NewUserController(userUsecase)

	// Post layer
	postRepo := repositories.NewPostRepository(firebaseApp.Firestore)
	postUsecase := usecases.NewPostUsecase(postRepo)
	postController := controllers.NewPostController(postUsecase, cld)

	// Usar Gorilla Mux para definir rutas
	router := mux.NewRouter()

	publicRouter := router.PathPrefix("/public").Subrouter()
	publicRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	publicRouter.HandleFunc("/users", userController.GetUser).Methods("GET")
	publicRouter.HandleFunc("/forgot-password", handlers.ForgotPasswordHandler(authService)).Methods("POST")
	publicRouter.HandleFunc("/posts", postController.GetAll).Methods("GET")
	publicRouter.HandleFunc("/posts", postController.Create).Methods("POST")
	publicRouter.HandleFunc("/post/{id}", postController.GetByID).Methods("GET")

	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)
	protectedRouter.HandleFunc("/profile", authHandler.GetUserProfile)
	protectedRouter.HandleFunc("/posts", postController.Delete).Methods("DELETE")
	protectedRouter.HandleFunc("/posts", postController.Edit).Methods("PUT")
	protectedRouter.HandleFunc("/posts/{id}/react", postController.React).Methods("POST")

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Acccept", "Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	handler := cors.New(corsOptions).Handler(router)
	serverPort := ":8080"

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Ruta registrada:", path)
		}
		return nil
	})

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Iniciar servidor HTTP
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	// Puerto en el que se está ejecutando el servidor
	log.Println("📚 Swagger UI en http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(serverPort, handler))

}
