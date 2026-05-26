# Lux Facturas - Planning

## Objetivo

Crear una landing page para vender Lux Facturas, un producto de facturacion orientado a pymes y mypes del Peru. La primera version se enfocara solo en la UI, comunicando valor de forma clara, confiable y directa.

## Publico objetivo

- Pymes y mypes peruanas.
- Negocios que necesitan emitir comprobantes sin complicarse.
- Emprendedores que buscan ordenar clientes, productos, ventas y reportes.
- Usuarios que valoran soporte cercano y una experiencia simple.

## Propuesta de valor

Lux Facturas ayuda a negocios peruanos a emitir, organizar y controlar sus comprobantes de forma rapida, simple y segura.

Mensajes clave:

- Factura facil, rapido y sin complicarte.
- Pensado para pymes y mypes del Peru.
- Controla tus ventas, clientes y comprobantes desde un solo lugar.
- Ahorra tiempo y reduce errores.

## Stack recomendado

### Frontend

- Next.js
- React
- TypeScript
- Tailwind CSS
- lucide-react para iconos
- Framer Motion para microinteracciones suaves

### Backend API

- Go
- Gin como framework recomendado para la API
- PostgreSQL como base de datos
- sqlc o GORM para acceso a datos

Recomendacion:

- Usar `sqlc` si se prioriza control, rendimiento y queries explicitas.
- Usar `GORM` si se quiere avanzar mas rapido al inicio.

### Infraestructura

- Frontend: Vercel
- API Go: Render, Fly.io, Railway, DigitalOcean o AWS
- Storage: Cloudflare R2 o AWS S3
- Analytics: Vercel Analytics o Google Analytics

### Integraciones futuras

- WhatsApp Business
- Email transaccional
- Culqi, Mercado Pago o Niubiz
- SUNAT, OSE o PSE segun el modelo tributario elegido

## Arquitectura inicial

```txt
Usuario
  |
  v
Landing / App Web en Next.js
  |
  v
API en Go
  |
  v
PostgreSQL
  |
  v
Servicios externos: SUNAT, pagos, WhatsApp, email
```

## Planning de la landing UI

### 1. Hero principal

Objetivo: comunicar el beneficio principal en segundos.

Contenido:

- Headline: "Factura facil, rapido y sin complicarte"
- Subtexto orientado a pymes y mypes peruanas.
- CTA principal: "Probar gratis" o "Hablar por WhatsApp"
- CTA secundario: "Ver como funciona"
- Mockup visual del dashboard o emision de factura.

### 2. Problema

Objetivo: mostrar empatia con los dolores del negocio.

Ideas:

- Emitir comprobantes toma tiempo.
- Es facil cometer errores.
- SUNAT y los procesos tributarios pueden sentirse complejos.
- Clientes, productos y ventas terminan desordenados.

### 3. Solucion

Objetivo: presentar Lux Facturas como herramienta simple y confiable.

Contenido:

- Emision y organizacion de comprobantes.
- Clientes y productos guardados.
- Reportes claros.
- Acceso desde cualquier dispositivo.

### 4. Beneficios

Beneficios principales:

- Facturas y boletas en minutos.
- Menos errores operativos.
- Informacion ordenada.
- Reportes simples para tomar decisiones.
- Pensado para negocios peruanos.
- Soporte cercano.

### 5. Como funciona

Flujo recomendado:

1. Crea tu cuenta.
2. Registra clientes y productos.
3. Emite tu comprobante.
4. Revisa tus ventas y reportes.

### 6. Vista del producto

Objetivo: hacer tangible la herramienta.

Pantallas sugeridas:

- Dashboard de ventas.
- Formulario de emision de factura.
- Lista de clientes.
- Catalogo de productos.
- Reportes.

### 7. Planes

Planes sugeridos:

- Basico
- Pyme
- Negocio

Cada plan debe incluir:

- Precio.
- Limite o volumen de comprobantes, si aplica.
- Funciones principales.
- CTA claro.

### 8. Confianza

Elementos sugeridos:

- Testimonios.
- Mensaje de seguridad.
- Soporte por WhatsApp.
- Enfoque en Peru.
- Acompanamiento para negocios que recien empiezan.

### 9. FAQ

Preguntas sugeridas:

- Sirve para mypes?
- Puedo emitir boletas y facturas?
- Funciona desde celular?
- Hay soporte?
- Puedo probarlo gratis?
- Se integrara con SUNAT?

### 10. CTA final

Objetivo: cerrar con accion clara.

Mensaje:

- "Empieza a facturar mejor hoy"

Botones:

- Hablar por WhatsApp.
- Probar gratis.

## Direccion visual

La UI debe sentirse limpia, confiable, moderna y cercana. No debe parecer una web bancaria pesada ni una landing generica.

Paleta recomendada:

- Azul profundo para confianza.
- Verde o menta para accion y crecimiento.
- Blanco y grises suaves para claridad.
- Toques calidos sutiles para cercania local.

Estilo:

- Layout claro y escaneable.
- Secciones amplias.
- CTAs visibles.
- Iconografia simple.
- Mockups realistas del producto.
- Texto directo, sin exceso de tecnicismos.

## Fases del producto

### Fase 1 - Landing UI

- Crear landing visual.
- Capturar leads por formulario o WhatsApp.
- Mostrar beneficios, planes y FAQ.

### Fase 2 - App inicial

- Login.
- Gestion de clientes.
- Gestion de productos.
- Registro de comprobantes internos.
- Dashboard basico.

### Fase 3 - API Go formal

- API REST en Go.
- PostgreSQL.
- Autenticacion.
- Roles.
- Reportes.
- Auditoria basica.

### Fase 4 - Integraciones

- SUNAT, OSE o PSE.
- Pagos.
- WhatsApp Business.
- Email transaccional.
- Reportes avanzados.

## Recomendacion final

Usar Next.js para la UI y Go para la API es una buena combinacion para Lux Facturas. Permite lanzar una landing rapida, crecer hacia una app real y mantener un backend solido para procesos de facturacion, reportes e integraciones futuras.
