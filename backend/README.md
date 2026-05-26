# Lux Facturas API

Backend inicial para Lux Facturas, construido en Go.

## Requisitos

- Go 1.22 o superior

## Configuracion

Copia las variables de ejemplo:

```sh
cp .env.example .env
```

Variables disponibles:

- `APP_ENV`: entorno de ejecucion. Por defecto `development`.
- `PORT`: puerto HTTP. Por defecto `8080`.
- `DATABASE_URL`: cadena de conexion de PostgreSQL para futuras integraciones.
- `CORS_ORIGIN`: origen permitido para el frontend. Por defecto `http://localhost:3000`.

## Ejecutar

```sh
go run ./cmd/api
```

## Endpoints iniciales

```txt
GET /health
GET /api/v1/status
GET /openapi.yaml
GET /swagger
```

## Swagger

La documentacion Swagger queda disponible cuando la API esta corriendo:

```txt
http://localhost:8080/swagger
```

La especificacion OpenAPI se sirve desde:

```txt
http://localhost:8080/openapi.yaml
```

## Tests

```sh
go test ./...
```

## Siguiente paso recomendado

Agregar modulos de dominio para:

- Empresas
- Usuarios
- Clientes
- Productos
- Comprobantes
- Reportes

Cuando la API necesite PostgreSQL, se recomienda elegir entre `sqlc` para queries explicitas o `GORM` para avanzar mas rapido.
