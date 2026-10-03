package main

// Documentación Swagger (OpenAPI 3) del microservicio.
//
//	GET /docs          -> interfaz Swagger UI
//	GET /openapi.json  -> especificación OpenAPI

import "net/http"

const openapiJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Microservicio Go - 3D Organ Inc",
    "version": "1.0.0",
    "description": "Microservicio CRUD de modelos 3D de 3D-Organ-Inc. Lee y escribe en la tabla 'modelos' de Supabase (PostgreSQL): listar, obtener, crear, editar y eliminar."
  },
  "servers": [
    {
      "url": "/"
    }
  ],
  "tags": [
    {
      "name": "Estado",
      "description": "Verificar que el servicio está vivo"
    },
    {
      "name": "Modelos",
      "description": "Listar, obtener, crear, editar y eliminar modelos 3D"
    }
  ],
  "paths": {
    "/": {
      "get": {
        "tags": [
          "Estado"
        ],
        "summary": "Verifica que el servicio está vivo",
        "responses": {
          "200": {
            "description": "Servicio activo",
            "content": {
              "application/json": {
                "example": {
                  "servicio": "microservicio-modelos-3d-go",
                  "estado": "ok"
                }
              }
            }
          }
        }
      }
    },
    "/api/modelos": {
      "get": {
        "tags": [
          "Modelos"
        ],
        "summary": "Listar modelos 3D",
        "parameters": [
          {
            "name": "categoria",
            "in": "query",
            "required": false,
            "description": "Filtra por categoría exacta",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Lista de modelos",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Modelo"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Error de base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          }
        }
      },
      "post": {
        "tags": [
          "Modelos"
        ],
        "summary": "Crear un modelo 3D",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ModeloEntrada"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Modelo creado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Modelo"
                }
              }
            }
          },
          "400": {
            "description": "JSON inválido",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "422": {
            "description": "Faltan 'nombre' o 'categoria'",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "500": {
            "description": "Error de base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          }
        }
      }
    },
    "/api/modelos/{id}": {
      "get": {
        "tags": [
          "Modelos"
        ],
        "summary": "Obtener un modelo 3D",
        "parameters": [
          {
            "$ref": "#/components/parameters/IdModelo"
          }
        ],
        "responses": {
          "200": {
            "description": "Modelo encontrado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Modelo"
                }
              }
            }
          },
          "400": {
            "description": "Id no numérico",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "404": {
            "description": "Modelo no encontrado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "500": {
            "description": "Error de base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          }
        }
      },
      "put": {
        "tags": [
          "Modelos"
        ],
        "summary": "Editar un modelo 3D",
        "description": "Solo se actualizan los campos enviados (nombre, descripcion, categoria, activo). Debe enviarse al menos uno.",
        "parameters": [
          {
            "$ref": "#/components/parameters/IdModelo"
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ModeloCambios"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Modelo actualizado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Modelo"
                }
              }
            }
          },
          "400": {
            "description": "Id no numérico, JSON inválido o sin campos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "404": {
            "description": "Modelo no encontrado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "500": {
            "description": "Error de base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Modelos"
        ],
        "summary": "Eliminar un modelo 3D",
        "parameters": [
          {
            "$ref": "#/components/parameters/IdModelo"
          }
        ],
        "responses": {
          "204": {
            "description": "Eliminado (sin contenido)"
          },
          "400": {
            "description": "Id no numérico",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "404": {
            "description": "Modelo no encontrado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          },
          "500": {
            "description": "Error de base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "parameters": {
      "IdModelo": {
        "name": "id",
        "in": "path",
        "required": true,
        "description": "Id numérico del modelo",
        "schema": {
          "type": "integer",
          "example": 1
        }
      }
    },
    "schemas": {
      "ModeloEntrada": {
        "type": "object",
        "required": [
          "nombre",
          "categoria"
        ],
        "properties": {
          "nombre": {
            "type": "string",
            "example": "Cóclea 3D"
          },
          "descripcion": {
            "type": "string",
            "example": "Modelo anatómico de la cóclea"
          },
          "categoria": {
            "type": "string",
            "example": "🦻 Modelos del Oído"
          },
          "activo": {
            "type": "boolean",
            "default": true
          }
        }
      },
      "ModeloCambios": {
        "type": "object",
        "properties": {
          "nombre": {
            "type": "string",
            "example": "Cóclea 3D v2"
          },
          "descripcion": {
            "type": "string"
          },
          "categoria": {
            "type": "string",
            "example": "🤖 Modelos de Prueba"
          },
          "activo": {
            "type": "boolean",
            "example": false
          }
        }
      },
      "Modelo": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer",
            "example": 1
          },
          "nombre": {
            "type": "string"
          },
          "descripcion": {
            "type": "string"
          },
          "categoria": {
            "type": "string"
          },
          "fecha_registro": {
            "type": "string",
            "format": "date-time"
          },
          "activo": {
            "type": "boolean"
          }
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "detail": {
            "type": "string",
            "example": "Modelo 3D no encontrado"
          }
        }
      }
    }
  }
}`

const swaggerHTML = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="utf-8">
<title>Documentación API</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = function () {
  SwaggerUIBundle({ url: "/openapi.json", dom_id: "#swagger-ui" });
};
</script>
</body>
</html>
`

func registrarDocs() {
	http.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(openapiJSON))
	})
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerHTML))
	})
}
