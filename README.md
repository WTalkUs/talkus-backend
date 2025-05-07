# TalkUs Backend

Este repositorio contiene el backend para la aplicación **TalkUs**, una plataforma que permite a los usuarios registrarse, iniciar sesión, crear publicaciones y gestionar contenido. Está desarrollado en **Go** y utiliza servicios como Firebase y Cloudinary para la autenticación, almacenamiento de datos y manejo de imágenes.

## Características

- **Autenticación**: Registro, inicio de sesión y recuperación de contraseñas utilizando Firebase Authentication.
- **Gestión de usuarios**: Recuperación de perfiles de usuario desde Firestore.
- **Publicaciones**: Creación y recuperación de publicaciones con soporte para imágenes subidas a Cloudinary.
- **Swagger**: Documentación de la API generada automáticamente con Swagger.
- **Middleware**: Autenticación de rutas protegidas mediante middleware.
- **CORS**: Configuración de políticas de acceso entre dominios.

## Tecnologías utilizadas

- **Go**: Lenguaje de programación principal.
- **Firebase**: Para autenticación y base de datos Firestore.
- **Cloudinary**: Para almacenamiento de imágenes.
- **Gorilla Mux**: Enrutador para manejar las rutas HTTP.
- **Swagger**: Para la documentación de la API.
- **CORS**: Configuración de políticas de acceso entre dominios.

## Configuración del entorno

### Variables de entorno

El proyecto utiliza un archivo `.env` para configurar las credenciales necesarias. Asegúrate de incluir las siguientes variables en tu archivo `.env`:

```properties
FIREBASE_CREDENTIALS=firebaseCredentials.json
FIREBASE_WEB_API_KEY=tu_api_key_de_firebase

CLOUDINARY_CLOUD_NAME=tu_nombre_de_cloudinary
CLOUDINARY_API_KEY=tu_api_key_de_cloudinary
CLOUDINARY_API_SECRET=tu_api_secret_de_cloudinary
```

### Instalación

1. Clona este repositorio:
   ```bash
   git clone https://github.com/JuanPidarraga/talkus-backend.git
   cd talkus-backend
   ```

2. Instala las dependencias:
   ```bash
   go mod tidy
   ```

3. Configura las variables de entorno en un archivo `.env` en la raíz del proyecto.

4. Asegúrate de tener las credenciales de Firebase en el archivo `firebaseCredentials.json` en la raíz del proyecto.

5. Ejecuta el servidor:
   ```bash
   go run main.go
   ```

## Endpoints principales

### Autenticación

- **POST** `/public/register`: Registrar un nuevo usuario.
- **POST** `/public/forgot-password`: Enviar un enlace de recuperación de contraseña.

### Usuarios

- **GET** `/public/users`: Obtener un usuario por ID.

### Publicaciones

- **GET** `/public/posts`: Obtener todas las publicaciones.
- **POST** `/public/posts`: Crear una nueva publicación.

### Swagger

La documentación de la API está disponible en [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

## Estructura del proyecto

```
talkus-backend/
├── config/                 # Configuración de Firebase
├── docs/                   # Documentación Swagger
├── internal/
│   ├── controllers/        # Controladores HTTP
│   ├── handlers/           # Manejadores de rutas
│   ├── middleware/         # Middleware para autenticación
│   ├── models/             # Modelos de datos
│   ├── repositories/       # Repositorios para Firestore
│   ├── service/            # Servicios de negocio
│   └── usecases/           # Casos de uso
├── .env                    # Variables de entorno
├── .gitignore              # Archivos ignorados por Git
├── go.mod                  # Dependencias del proyecto
├── main.go                 # Punto de entrada de la aplicación
└── README.md               # Documentación del proyecto
```

## Contribuciones

Si deseas contribuir a este proyecto, por favor sigue los pasos:

1. Haz un fork del repositorio.
2. Crea una nueva rama para tu funcionalidad o corrección de errores:
   ```bash
   git checkout -b feature/nueva-funcionalidad
   ```
3. Realiza tus cambios y haz un commit:
   ```bash
   git commit -m "Agrega nueva funcionalidad"
   ```
4. Sube tus cambios a tu fork:
   ```bash
   git push origin feature/nueva-funcionalidad
   ```
5. Abre un Pull Request en este repositorio.

## Licencia

Este proyecto está bajo la licencia MIT. Consulta el archivo `LICENSE` para más detalles.