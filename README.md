# API Feedbacks

API REST en Golang para capturar, actualizar y consultar feedbacks de usuarios sobre la plataforma digital.

## Arquitectura

El proyecto sigue principios de **Clean Architecture** con separación clara de responsabilidades:

```
cmd/server/        → Entry point y bootstrap
internal/
  config/          → Configuración por variables de entorno
  domain/          → Entidades y reglas de negocio
  service/         → Lógica de negocio
  repository/      → Interfaces y implementaciones de persistencia
  handler/         → Controladores HTTP, DTOs, router
  middleware/      → Auth, logging, recovery, CORS, rate limiting
pkg/response/      → Helpers de respuesta JSON estandarizada
```

## Requisitos Previos

- [Docker](https://www.docker.com/) y Docker Compose

## Ejecución

### Modo producción
```bash
# Levantar la API + PostgreSQL
docker compose up -d --build

# Ver logs
docker compose logs -f app

# Apagar
docker compose down -v
```

### Modo desarrollo (hot-reload)
```bash
docker compose --profile dev up --build
```

### Ejecutar tests
```bash
docker compose run --rm --no-deps app-dev sh -c "go test ./internal/... -v -count=1"
```

### Seed de datos de ejemplo
```bash
bash scripts/seed.sh
```

## Endpoints

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `GET` | `/health` | Health check | No |
| `GET` | `/ready` | Readiness probe | No |
| `POST` | `/api/v1/feedbacks` | Crear feedback | Sí |
| `GET` | `/api/v1/feedbacks` | Listar feedbacks (con filtros) | Sí |
| `GET` | `/api/v1/feedbacks/{id}` | Obtener feedback por ID | Sí |
| `PATCH` | `/api/v1/feedbacks/{id}` | Actualizar feedback parcial | Sí |

### Autenticación

Todos los endpoints bajo `/api/v1/` requieren el header `X-API-Key`:
```
X-API-Key: my-secret-api-key
```

### Filtros (GET /api/v1/feedbacks)

| Query Param | Tipo | Descripción |
|-------------|------|-------------|
| `user_id` | string | Filtrar por ID de usuario |
| `feedback_type` | string | Filtrar por tipo: bug, suggestion, praise, question |
| `min_rating` | int | Rating mínimo (1-5) |
| `max_rating` | int | Rating máximo (1-5) |
| `created_from` | RFC3339 | Fecha de creación desde |
| `created_to` | RFC3339 | Fecha de creación hasta |
| `limit` | int | Resultados por página (default: 20, max: 100) |
| `offset` | int | Desplazamiento para paginación |

## Ejemplos cURL

### Crear feedback
```bash
curl -X POST http://localhost:8080/api/v1/feedbacks \
  -H "Content-Type: application/json" \
  -H "X-API-Key: my-secret-api-key" \
  -d '{
    "user_id": "usr-001",
    "feedback_type": "bug",
    "rating": 3,
    "comment": "El botón de pago no responde en Safari"
  }'
```

### Listar con filtros
```bash
curl -H "X-API-Key: my-secret-api-key" \
  "http://localhost:8080/api/v1/feedbacks?feedback_type=bug&min_rating=1&max_rating=3"
```

### Actualizar feedback
```bash
curl -X PATCH http://localhost:8080/api/v1/feedbacks/{ID} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: my-secret-api-key" \
  -d '{"rating": 5, "comment": "Se resolvió el problema"}'
```

## Decisiones Técnicas

Ver [tech-decisions.md](tech-decisions.md) para el detalle completo. Resumen:

1. **PostgreSQL**: Ideal para filtros por rangos, fechas e índices compuestos
2. **API Key auth**: Esquema básico para proteger endpoints sin complejidad de OAuth2
3. **chi router**: Ligero, idiomatic Go, compatible con `net/http` estándar
4. **pgx/v5**: Mayor rendimiento que `lib/pq`, connection pooling nativo
5. **Clean Architecture**: Testing facilitado por interfaces entre capas
6. **Docker-first**: Todo el desarrollo y ejecución contenedorizados

## Bibliotecas Externas

| Paquete | Propósito |
|---------|-----------|
| `go-chi/chi/v5` | HTTP router composable |
| `jackc/pgx/v5` | Driver PostgreSQL de alto rendimiento |
| `google/uuid` | Generación de UUID v4 |
| `golang.org/x/time/rate` | Rate limiting (token bucket) |
| `stretchr/testify` | Assertions para tests (dev dependency) |

## Variables de Entorno

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `SERVER_PORT` | No | `8080` | Puerto del servidor |
| `DATABASE_URL` | Sí | — | URL de conexión PostgreSQL |
| `API_KEY` | Sí | — | Clave de autenticación |
| `LOG_LEVEL` | No | `info` | Nivel de log: debug, info, warn, error |
| `CORS_ORIGINS` | No | `*` | Orígenes CORS permitidos |
| `RATE_LIMIT_RPS` | No | `100` | Límite de requests por segundo por IP |
