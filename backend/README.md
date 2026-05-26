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
- `LEADS_FILE`: archivo JSON donde se guardan los leads. Por defecto `data/leads.json`.
- `ADMIN_USER`: usuario para entrar al dashboard. Por defecto `admin`.
- `ADMIN_PASSWORD`: clave para entrar al dashboard. En produccion debe cambiarse.
- `SESSION_KEY`: secreto para firmar la cookie de sesion. En produccion debe ser largo y privado.

## Ejecutar

```sh
go run ./cmd/api
```

## Dashboard

El backend tambien sirve una pantalla interna para revisar leads:

```txt
http://localhost:8080/dashboard
```

Primero debes iniciar sesion en:

```txt
http://localhost:8080/login
```

En desarrollo, si no defines variables, las credenciales son:

```txt
usuario: admin
clave: lux-admin
```

En produccion configura `ADMIN_USER`, `ADMIN_PASSWORD` y `SESSION_KEY` en Vercel antes de publicar el backend.

La raiz redirige automaticamente al dashboard:

```txt
http://localhost:8080
```

## Endpoints iniciales

```txt
GET /
GET /login
POST /login
POST /logout
GET /dashboard
POST /dashboard/leads/{id}/status
GET /health
GET /api/v1/status
POST /api/v1/leads
GET /api/v1/leads
GET /api/v1/leads/{id}
PATCH /api/v1/leads/{id}/status
GET /openapi.yaml
GET /swagger
```

## Crear un lead

```sh
curl -X POST http://localhost:8080/api/v1/leads \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Maria Lopez",
    "company_name": "Bodega San Jose",
    "document_number": "10456789123",
    "phone": "987654321",
    "email": "maria@bodega.pe",
    "selected_plan": "kit_basico",
    "license_type": "mensual",
    "customer_type": "quiere_cotizar",
    "business_category": "bodega",
    "approx_product_quantity": 250,
    "city": "Lima",
    "message": "Quiero informacion para mi tienda"
  }'
```

Valores permitidos:

- `selected_plan`: `kit_basico`, `estandar`, `completo`, `software`.
- `license_type`: `mensual`, `permanente`.
- `customer_type`: `ya_tiene_mac`, `necesita_kit_completo`, `quiere_cotizar`.
- `status`: `nuevo`, `contactado`, `cotizado`, `cerrado`, `descartado`.

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
