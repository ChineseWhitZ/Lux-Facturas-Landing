# Deployament - Lux Facturas

Este documento describe como se desplegaran los proyectos de Lux Facturas segun el alcance definido en `PLANNING.md`.

## Resumen de despliegue

Lux Facturas se desplegara por fases:

1. Landing UI en Next.js para captar leads y validar la propuesta de valor.
2. App inicial en Next.js para login, clientes, productos, comprobantes internos y dashboard basico.
3. API formal en Go con PostgreSQL para la logica de negocio, autenticacion, roles, reportes y auditoria basica.

Arquitectura objetivo:

```txt
Usuario
  |
  v
Frontend Next.js en Vercel
  |
  v
API Go en Render, Fly.io, Railway, DigitalOcean o AWS
  |
  v
PostgreSQL administrado
  |
  v
Servicios externos: SUNAT, pagos, WhatsApp, email y storage
```

## Entornos

Se recomienda manejar tres entornos:

- `development`: entorno local para desarrollo.
- `staging`: entorno publico de pruebas antes de produccion.
- `production`: entorno real para clientes.

Variables comunes:

```sh
APP_ENV=production
PORT=8080
DATABASE_URL=postgres://usuario:password@host:5432/lux_facturas
```

El frontend tambien debera manejar variables publicas y privadas:

```sh
NEXT_PUBLIC_APP_URL=https://luxfacturas.pe
NEXT_PUBLIC_API_URL=https://api.luxfacturas.pe
```

## Fase 1 - Landing UI

Objetivo: publicar una landing rapida, confiable y clara para vender Lux Facturas, capturar leads y dirigir usuarios a WhatsApp.

### Plataforma recomendada

- Vercel para el proyecto Next.js.
- Vercel Analytics o Google Analytics para medicion.
- Dominio principal: `luxfacturas.pe` o el dominio comercial elegido.

### Pasos de despliegue

1. Crear el proyecto frontend con Next.js, TypeScript y Tailwind CSS.
2. Subir el proyecto a un repositorio Git.
3. Conectar el repositorio en Vercel.
4. Configurar el framework como Next.js.
5. Configurar variables de entorno:

```sh
NEXT_PUBLIC_APP_URL=https://luxfacturas.pe
NEXT_PUBLIC_WHATSAPP_URL=https://wa.me/51XXXXXXXXX
NEXT_PUBLIC_ANALYTICS_ID=valor_si_aplica
```

6. Configurar dominio personalizado en Vercel.
7. Verificar que los CTAs funcionen:

- `Probar gratis`
- `Hablar por WhatsApp`
- `Ver como funciona`

### Validacion antes de produccion

- La pagina carga correctamente en desktop y mobile.
- Los textos comunican la propuesta de valor para pymes y mypes del Peru.
- Los formularios o enlaces de contacto funcionan.
- Analytics registra visitas y eventos principales.
- El Lighthouse score es aceptable para performance, accesibilidad y SEO.

## Fase 2 - App inicial

Objetivo: desplegar la primera version funcional de la app con login, clientes, productos, comprobantes internos y dashboard basico.

### Plataforma recomendada

- Vercel para la app Next.js.
- API Go separada para datos y reglas de negocio.
- PostgreSQL administrado para persistencia.

### Estrategia de rutas

Opcion simple:

```txt
https://luxfacturas.pe          -> Landing
https://app.luxfacturas.pe      -> App web
https://api.luxfacturas.pe      -> API Go
```

Opcion monorepo:

```txt
/frontend   -> Landing y app Next.js
/backend    -> API Go
```

### Pasos de despliegue

1. Crear vistas protegidas para la app inicial.
2. Configurar autenticacion.
3. Conectar el frontend con la API usando `NEXT_PUBLIC_API_URL`.
4. Crear una base PostgreSQL para `staging` y otra para `production`.
5. Ejecutar migraciones de base de datos antes de publicar cambios.
6. Configurar CORS en la API para permitir solo los dominios necesarios.
7. Publicar primero en `staging`.
8. Validar flujos principales.
9. Promover a `production`.

### Validacion funcional

- Login y cierre de sesion.
- Registro y edicion de clientes.
- Registro y edicion de productos.
- Registro de comprobantes internos.
- Dashboard basico con ventas y resumen.
- Manejo correcto de errores de API.

## Fase 3 - API Go formal

Objetivo: desplegar una API REST en Go con PostgreSQL, autenticacion, roles, reportes y auditoria basica.

El backend actual esta en `backend/` y expone:

```txt
GET /health
GET /api/v1/status
```

### Plataformas recomendadas

Se puede elegir una de estas opciones:

- Render: simple para empezar, con deploy desde Git y PostgreSQL administrado.
- Fly.io: buena opcion si se quiere mayor control y despliegue cercano a usuarios.
- Railway: rapido para prototipos y primeras versiones.
- DigitalOcean App Platform: opcion estable y simple para produccion.
- AWS: recomendable cuando el producto necesite mayor control, escalabilidad y servicios administrados.

Para la primera version productiva se recomienda Render o Railway por simplicidad.

### Variables de entorno de la API

```sh
APP_ENV=production
PORT=8080
DATABASE_URL=postgres://usuario:password@host:5432/lux_facturas
```

Variables futuras:

```sh
JWT_SECRET=valor_seguro
CORS_ALLOWED_ORIGINS=https://luxfacturas.pe,https://app.luxfacturas.pe
S3_BUCKET=lux-facturas
S3_REGION=auto
S3_ENDPOINT=https://endpoint-r2-o-s3
WHATSAPP_TOKEN=valor_si_aplica
EMAIL_API_KEY=valor_si_aplica
PAYMENT_PROVIDER_API_KEY=valor_si_aplica
```

### Build y start

Desde `backend/`:

```sh
go mod download
go build -o bin/api ./cmd/api
./bin/api
```

En plataformas como Render se puede configurar:

```sh
Build Command: go build -o bin/api ./cmd/api
Start Command: ./bin/api
```

### Health check

Usar el endpoint:

```txt
/health
```

La plataforma de hosting debe verificar este endpoint para detectar si el servicio esta activo.

### Base de datos

Para PostgreSQL:

1. Crear una base administrada para `staging`.
2. Crear una base administrada separada para `production`.
3. Guardar la cadena de conexion en `DATABASE_URL`.
4. Activar backups automaticos.
5. Definir estrategia de migraciones.

Herramientas sugeridas para migraciones:

- `golang-migrate/migrate`
- `goose`
- Migraciones gestionadas por el proveedor si aplica.

### Acceso a datos

Segun el planning:

- Usar `sqlc` si se prioriza control, rendimiento y queries explicitas.
- Usar `GORM` si se quiere avanzar mas rapido al inicio.

La decision debe tomarse antes de crear los primeros modulos persistentes para evitar cambios grandes despues.

## Storage

Para archivos, logos, comprobantes generados o adjuntos futuros:

- Cloudflare R2 si se quiere bajo costo y buena integracion con CDN.
- AWS S3 si se prefiere ecosistema AWS.

Recomendacion inicial: Cloudflare R2 para reducir costos operativos.

## Analytics y monitoreo

Frontend:

- Vercel Analytics para medicion simple.
- Google Analytics si se necesita integracion comercial o marketing.

Backend:

- Logs estructurados en JSON.
- Health check `/health`.
- Alertas de caida del servicio.
- Metricas basicas de latencia, errores y uso.

Opciones recomendadas:

- Logs del proveedor al inicio.
- Sentry para errores de frontend y backend cuando el producto crezca.
- UptimeRobot, Better Stack o el monitor del proveedor para uptime.

## Seguridad

Antes de produccion:

- Usar HTTPS en todos los dominios.
- Configurar CORS con dominios explicitos.
- Guardar secretos solo en variables de entorno del proveedor.
- No exponer `DATABASE_URL`, tokens ni llaves privadas en el frontend.
- Activar backups de PostgreSQL.
- Agregar rate limiting en endpoints sensibles.
- Registrar auditoria basica para acciones importantes.

## Checklist de produccion

- Dominio configurado.
- Certificados HTTPS activos.
- Variables de entorno configuradas.
- Base PostgreSQL creada y respaldada.
- Migraciones ejecutadas.
- Health check activo.
- Logs revisados.
- Landing probada en mobile y desktop.
- CTAs y formularios probados.
- API probada desde el frontend.
- CORS configurado.
- Analytics activo.
- Backups habilitados.

## Flujo recomendado de release

1. Desarrollar cambios en una rama Git.
2. Abrir pull request.
3. Ejecutar pruebas y build.
4. Desplegar automaticamente a `staging`.
5. Validar flujo principal.
6. Fusionar a la rama principal.
7. Desplegar a `production`.
8. Verificar logs, health check y analytics.

## Roadmap de despliegue

### Primer despliegue

- Landing Next.js en Vercel.
- WhatsApp como CTA principal.
- Analytics basico.
- Dominio principal.

### Segundo despliegue

- App inicial protegida.
- API Go en Render o Railway.
- PostgreSQL administrado.
- Modulos de clientes, productos y comprobantes internos.

### Tercer despliegue

- Autenticacion robusta.
- Roles.
- Reportes.
- Auditoria.
- Storage para archivos.
- Preparacion para integraciones con SUNAT, pagos, WhatsApp y email.

