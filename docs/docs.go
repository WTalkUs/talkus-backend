// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/public/forgot-password": {
            "post": {
                "description": "Envía un enlace de recuperación de contraseña al correo electrónico proporcionado.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Recuperación de contraseña",
                "parameters": [
                    {
                        "description": "Correo electrónico del usuario",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Enlace de recuperación enviado correctamente",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Solicitud inválida: el email es obligatorio",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno al enviar el enlace de recuperación",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/public/posts": {
            "get": {
                "description": "Obtiene una lista de todas las publicaciones ordenadas por fecha de creación.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Obtener todas las publicaciones",
                "responses": {
                    "200": {
                        "description": "Lista de publicaciones",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Post"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Permite crear una nueva publicación con un título, contenido y una imagen opcional. La imagen se sube a Cloudinary y se guarda la URL en la publicación.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Crear una nueva publicación",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Título de la publicación",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contenido de la publicación",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Imagen para la publicación",
                        "name": "image",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Publicación creada exitosamente",
                        "schema": {
                            "$ref": "#/definitions/models.Post"
                        }
                    },
                    "400": {
                        "description": "Solicitud inválida, título o contenido faltante",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno al crear la publicación",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/public/register": {
            "post": {
                "description": "Permite registrar un nuevo usuario con su correo y contraseña",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Registrar un nuevo usuario",
                "parameters": [
                    {
                        "description": "Datos del usuario a registrar",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Usuario creado exitosamente",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Solicitud incorrecta: los datos no son válidos",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno al registrar el usuario",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/public/users": {
            "get": {
                "description": "Recupera un usuario de la base de datos utilizando su ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Obtener un usuario por ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del usuario a recuperar",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Usuario encontrado",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Parámetro 'id' faltante",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ForgotPasswordRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "handlers.RegisterRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Post": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "dislikes": {
                    "type": "integer"
                },
                "forum_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "is_flagged": {
                    "type": "boolean"
                },
                "likes": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
