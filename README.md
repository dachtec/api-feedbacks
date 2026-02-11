<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.24"/>
  <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL 16"/>
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"/>
  <img src="https://img.shields.io/badge/Architecture-Clean-blueviolet?style=for-the-badge" alt="Clean Architecture"/>
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="MIT License"/>
</p>

# 📣 API Feedbacks

> **Backend RESTful en Golang para capturar, gestionar y consultar feedbacks de usuarios sobre una plataforma digital.**

API diseñada como solución al [challenge técnico](docs/challenge-context.md), cuyo objetivo es construir una base sólida para análisis internos, dashboards y detección de problemas de experiencia, desempeño y percepción de los usuarios.

---

## 📑 Tabla de Contenidos

- [🎯 Contexto del Proyecto](#-contexto-del-proyecto)
- [🏗️ Arquitectura y Diseño](#️-arquitectura-y-diseño)
- [⚙️ Decisiones Técnicas](#️-decisiones-técnicas)
- [📌 Supuestos y Limitaciones](#-supuestos-y-limitaciones)
- [🔧 Requisitos Previos](#-requisitos-previos)
- [🚀 Instalación y Ejecución](#-instalación-y-ejecución)
- [🧪 Pruebas con Postman](#-pruebas-con-postman)
- [📖 Referencia Rápida de Endpoints](#-referencia-rápida-de-endpoints)
- [💡 Ideas para Evolucionar el Proyecto](#-ideas-para-evolucionar-el-proyecto)
- [📂 Estructura del Proyecto](#-estructura-del-proyecto)

---

## 🎯 Contexto del Proyecto

Este proyecto nace de un [challenge técnico](docs/challenge-context.md) que plantea la necesidad de construir un backend en **Golang** (sin interfaz gráfica) capaz de:

| Requerimiento | Estado |
|---|---|
| Crear feedbacks de usuarios | ✅ Implementado |
| Actualizar feedbacks (rating, comment, feedback_type) | ✅ Implementado |
| Consultar feedbacks con filtros combinados | ✅ Implementado |
| Filtros: `user_id`, `feedback_type`, `rating`, fechas | ✅ Implementado |
| Autenticación API Key | ✅ Implementado |
| Pruebas vía Postman/cURL | ✅ [Colección Postman](docs/API_Feedbacks.postman_collection.json) |
| Documentación de decisiones técnicas | ✅ [Documentadas](docs/tech-decisions.md) |
| Trazabilidad del uso de IA | ✅ [PROMPTS.md](docs/PROMPTS.md) y [sessions/](docs/sessions) |

### ¿Qué es un Feedback?

```json
{
  "feedback_id": "f-0001",
  "user_id": "u-001",
  "feedback_type": "sugerencia",
  "rating": 4,
  "comment": "La navegación podría ser más intuitiva",
  "created_at": "2026-02-10T14:30:00Z",
  "updated_at": "2026-02-10T14:30:00Z"
}
```

**Tipos válidos:** `bug` · `sugerencia` · `elogio` · `duda` · `queja`

---

## 🏗️ Arquitectura y Diseño

El proyecto implementa **Clean Architecture** con separación estricta de responsabilidades en 4 capas:

```
┌─────────────────────────────────────────────────┐
│                   HTTP Layer                     │
│          (chi router + middlewares)              │
├─────────────────────────────────────────────────┤
│                 Handler Layer                    │
│        (request/response, validación)            │
├─────────────────────────────────────────────────┤
│                 Service Layer                    │
│          (lógica de negocio, IDs)                │
├─────────────────────────────────────────────────┤
│               Repository Layer                   │
│         (PostgreSQL via pgx/pgxpool)             │
├─────────────────────────────────────────────────┤
│                 Domain Layer                     │
│     (entidades, interfaces, validaciones)        │
└─────────────────────────────────────────────────┘
```

> 📐 **Principios aplicados:** SOLID, inyección de dependencias, interfaces para desacoplamiento, y DTOs para transformación de datos.

### Stack Tecnológico

| Componente | Tecnología | Justificación |
|---|---|---|
| **Lenguaje** | Go 1.24 | Rendimiento, concurrencia nativa, tipado estático |
| **Router HTTP** | chi v5 | Ligero, idiomático, compatible con `net/http` |
| **Base de Datos** | PostgreSQL 16 | Índices compuestos, filtros por rangos, `TIMESTAMPTZ` |
| **Driver DB** | pgx v5 + pgxpool | Mayor rendimiento que `lib/pq`, connection pooling nativo |
| **Contenedores** | Docker + Compose | Entorno reproducible, sin dependencias locales |
| **Rate Limiting** | golang.org/x/time | Protección contra abuso de la API |

> 📄 Todas las decisiones técnicas están documentadas con contexto y justificación en [`tech-decisions.md`](docs/tech-decisions.md).

---

## ⚙️ Decisiones Técnicas Destacadas

| ID | Decisión | Resumen |
|---|---|---|
| TD-001 | PostgreSQL | Soporte nativo para filtros complejos y tipos ricos |
| TD-002 | API Key Auth | Mecanismo simple pero funcional vía header `X-API-Key` |
| TD-003 | CI/CD diferido | Priorización de API funcional sobre infraestructura de pipelines |
| TD-004 | Router chi v5 | Ligero, idiomático, compatible con `net/http` estándar |
| TD-005 | Driver pgx v5 | Mayor rendimiento que `lib/pq`, connection pooling nativo |
| TD-006 | Clean Architecture | Capas: domain → repository → service → handler |
| TD-007 | Colección Postman | Scripts de test automatizados para todos los escenarios |
| TD-008 | IDs secuenciales `f-####` | Legibles, predecibles, thread-safe con `sync/atomic` |
| TD-009 | Validación `u-###` | Regex en capa de dominio para consistencia |
| TD-010 | Tipos en español | Alineados con el dominio de negocio |
| TD-011 | Timestamps sin ms | Formato `RFC3339` truncado a segundos |
| TD-012 | Corrección migración | Alineación de migración inline con modelo de dominio actual |
| TD-013 | Seed dinámico | Carga de datos semilla desde `seed-data.json` con `jq` |
| TD-014 | Go 1.24 + air pinning | Actualización de Go, pinning de `air` y perfiles Docker Compose |

> 📋 Ver detalle completo en [`tech-decisions.md`](docs/tech-decisions.md)

---

## 📌 Supuestos y Limitaciones

### Supuestos
- El sistema será consumido exclusivamente vía API REST (sin UI).
- La autenticación con API Key estática es suficiente para el alcance del challenge.
- Los IDs secuenciales (`f-####`) se generan en memoria; en producción multi-instancia se migraría a secuencias PostgreSQL.
- Un único usuario/API Key gestiona todos los feedbacks (no hay multi-tenancy).

### Limitaciones Actuales
- ❌ **Sin CI/CD** — Pipeline de integración continua diferido a iteración futura (TD-003).
- ❌ **Sin monitoreo** — Métricas Prometheus y tracing distribuido no implementados.
- ❌ **Sin paginación cursor-based** — Se usa offset/limit básico.
- ❌ **Sin rate limiting por usuario** — Rate limit global, no por API Key individual.

---

## 🔧 Requisitos Previos

Antes de ejecutar el proyecto, asegúrate de tener instalado:

| Herramienta | Versión mínima | Verificación |
|---|---|---|
| **Docker** | 20.10+ | `docker --version` |
| **Docker Compose** | 2.0+ | `docker compose version` |
| **Git** | 2.0+ | `git --version` |
| **Postman** *(opcional)* | Última versión | Para pruebas con la colección incluida |

> [!NOTE]
> **No necesitas Go instalado localmente.** Todo el desarrollo y ejecución se realiza dentro de contenedores Docker.

---

## 🚀 Instalación y Ejecución

### 1. Clonar el repositorio

```bash
git clone https://github.com/dachtec/api-feedbacks.git
cd api-feedbacks
```

### 2. Ejecutar en modo producción

```bash
# Construir y levantar todos los servicios (API + PostgreSQL)
make run
```

Esto ejecutará `docker compose up -d --build`, levantando:
- 🟢 **API** en `http://localhost:8080`
- 🟢 **PostgreSQL** en `localhost:5432`

### 3. Verificar que la API está activa

```bash
curl http://localhost:8080/health
```

**Respuesta esperada:**
```json
{ "status": "ok" }
```

### 4. Realizar una primera petición autenticada

```bash
curl -X POST http://localhost:8080/api/v1/feedbacks \
  -H "Content-Type: application/json" \
  -H "X-API-Key: my-secret-api-key" \
  -d '{
    "user_id": "u-001",
    "feedback_type": "elogio",
    "rating": 5,
    "comment": "Excelente plataforma, muy intuitiva"
  }'
```

### Comandos útiles (`Makefile`)

| Comando | Descripción |
|---|---|
| `make run` | Levanta los contenedores en modo producción |
| `make dev` | Levanta en modo desarrollo con hot-reload |
| `make test` | Ejecuta tests unitarios dentro de Docker |
| `make test-cover` | Tests con reporte de cobertura |
| `make lint` | Ejecuta `go vet` en el código |
| `make logs` | Muestra logs de la aplicación en tiempo real |
| `make seed` | Carga datos de ejemplo |
| `make docker-down` | Detiene los contenedores |
| `make docker-clean` | Detiene y elimina contenedores + volúmenes |
| `make clean` | Limpieza total (contenedores, imágenes, temporales) |

### Variables de Entorno

| Variable | Valor por defecto | Descripción |
|---|---|---|
| `SERVER_PORT` | `8080` | Puerto del servidor HTTP |
| `DATABASE_URL` | *(ver docker-compose)* | URL de conexión a PostgreSQL |
| `API_KEY` | `my-secret-api-key` | Clave de autenticación |
| `LOG_LEVEL` | `info` | Nivel de logging (`debug`, `info`, `warn`, `error`) |
| `CORS_ORIGINS` | `*` | Orígenes permitidos para CORS |
| `RATE_LIMIT_RPS` | `100` | Requests por segundo permitidos |

---

## 🧪 Pruebas con Postman

El proyecto incluye una **colección Postman completa** con scripts de test automatizados para validar todos los escenarios de la API.

### Importar la colección

1. Abre **Postman** y haz clic en **Import**.
2. Selecciona el archivo:
   ```
   docs/API_Feedbacks.postman_collection.json
   ```
3. La colección se importará con todas las variables pre-configuradas.

### Variables de la colección

| Variable | Valor | Descripción |
|---|---|---|
| `base_url` | `http://localhost:8080` | URL base de la API |
| `api_key` | `my-secret-api-key` | API Key para autenticación |
| `feedback_id` | *(se auto-genera)* Ej. `f-0001` | ID del primer feedback creado (capturado en scripts) |
| `feedback_id_2` | *(se auto-genera)* Ej. `f-0002`| ID del segundo feedback creado (capturado en scripts) |

### Ejecutar las pruebas

#### Opción A: Request por request
Navega por las carpetas de la colección y ejecuta cada request individualmente. Los **scripts de test** validan automáticamente:
- ✅ Status codes correctos
- ✅ Estructura de la respuesta
- ✅ Reglas de negocio (formatos, rangos, tipos válidos)

#### Opción B: Collection Runner (ejecución completa)
1. Haz clic derecho en la colección → **Run Collection**.
2. Asegúrate de que el orden de ejecución sea el correcto (crear antes de consultar/actualizar).
3. Haz clic en **Run** y observa los resultados de todos los tests.

### Escenarios cubiertos

| Categoría | Escenarios |
|---|---|
| **Crear feedback** | Happy path, tipos inválidos, rating fuera de rango, user_id inválido, campos faltantes |
| **Obtener por ID** | Existente, no encontrado, formato inválido |
| **Actualizar** | Parcial, completa, feedback inexistente, valores inválidos |
| **Listar con filtros** | Por `user_id`, `feedback_type`, rango de `rating`, rango de fechas, combinaciones, paginación |
| **Autenticación** | Sin API Key, API Key inválida |

---

## 📖 Referencia Rápida de Endpoints

Todas las rutas están bajo el prefijo `/api/v1` y requieren el header `X-API-Key`.

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/health` | Health check (sin auth) |
| `POST` | `/api/v1/feedbacks` | Crear un feedback |
| `GET` | `/api/v1/feedbacks` | Listar feedbacks con filtros |
| `GET` | `/api/v1/feedbacks/{id}` | Obtener feedback por ID |
| `PUT` | `/api/v1/feedbacks/{id}` | Actualizar un feedback |

### Parámetros de filtro (`GET /api/v1/feedbacks`)

| Parámetro | Tipo | Ejemplo |
|---|---|---|
| `user_id` | string | `?user_id=u-001` |
| `feedback_type` | string | `?feedback_type=bug` |
| `min_rating` | int | `?min_rating=3` |
| `max_rating` | int | `?max_rating=5` |
| `created_from` | datetime | `?created_from=2026-01-01T00:00:00Z` |
| `created_to` | datetime | `?created_to=2026-12-31T23:59:59Z` |
| `limit` | int | `?limit=10` |
| `offset` | int | `?offset=0` |

> 📄 Documentación OpenAPI completa disponible en [`docs/openapi.yaml`](docs/openapi.yaml)

---

## 💡 Ideas para Evolucionar el Proyecto

### 🔜 Corto Plazo
- **Paginación cursor-based** — Más eficiente que offset/limit para datasets grandes.
- **Soft delete** — Marcar feedbacks como eliminados sin borrarlos físicamente.
- **Endpoint DELETE** — Permitir la eliminación (lógica) de feedbacks.
- **Validaciones enriquecidas** — Longitud mínima/máxima de `comment`, sanitización de HTML/XSS.

### 🔮 Mediano Plazo
- **Pipeline CI/CD** — GitHub Actions con build, test, lint y deploy automático.
- **Métricas y Observabilidad** — Prometheus para métricas, Grafana para dashboards, tracing con OpenTelemetry.
- **Autenticación JWT/OAuth2** — Reemplazar API Key por tokens con scopes y expiración.
- **Rate limiting por usuario** — Limitar requests por API Key individual en lugar de global.
- **Caché con Redis** — Cache de consultas frecuentes para reducir carga en PostgreSQL.

### 🚀 Largo Plazo
- **Análisis de sentimiento** — Integrar un modelo de NLP para clasificar automáticamente el sentimiento de los comentarios.
- **Dashboard en tiempo real** — WebSockets o SSE para notificaciones de nuevos feedbacks.
- **Multi-tenancy** — Soporte para múltiples plataformas/equipos con aislamiento de datos.
- **Búsqueda full-text** — PostgreSQL full-text search o Elasticsearch para buscar en comentarios.
- **Event sourcing** — Almacenar el historial completo de cambios para auditoría.
- **API GraphQL** — Alternativa a REST para consultas más flexibles por parte de frontends.

---


<p align="center">
  Desarrollado con 🤖 asistencia de IA + supervisión humana<br/>
  <sub>Ver trazabilidad completa en <a href="docs/PROMPTS.md">PROMPTS.md</a> y <a href="docs/sessions">sessions/</a></sub>
</p>
