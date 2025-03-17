package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// FirebaseApp encapsula los clientes de Firebase que usaremos.
type FirebaseApp struct {
	Auth      *auth.Client
	Firestore *firestore.Client
}

// InitFirebase inicializa Firebase y crea el cliente de Firestore.

func InitFirebase() (*FirebaseApp, error) {
	// Cargar variables de entorno desde .env (si existe)
	_ = godotenv.Load()

	// Obtener la ruta de credenciales desde la variable de entorno
	credFile := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credFile == "" {
		log.Fatalf("❌ No se encontró FIREBASE_CREDENTIALS_PATH en las variables de entorno")
		return nil, nil
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(credFile)

	// Inicializar la app de Firebase
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("❌ Error inicializando Firebase: %v", err)
		return nil, err
	}

	// Crear cliente de autenticación (opcional, si lo necesitas)
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("❌ Error inicializando Auth: %v", err)
		return nil, err
	}

	// Crear cliente de Firestore
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("❌ Error inicializando Firestore: %v", err)
		return nil, err
	}

	log.Println("🔥 Firebase inicializado correctamente")
	return &FirebaseApp{
		Auth:      authClient,
		Firestore: firestoreClient,
	}, nil
}
