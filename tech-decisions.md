# Decisiones Técnicas

## TD-001: Base de datos PostgreSQL
**Contexto**: Se necesita persistir feedbacks con filtros por rangos (rating, fechas), búsqueda por tipo y usuario.
**Decisión**: Usar PostgreSQL como motor de persistencia.
**Justificación**: Soporte nativo para índices compuestos, filtros por rangos, tipos de datos ricos (TIMESTAMPTZ, UUID), y excelente rendimiento para consultas con múltiples filtros simultáneos.

## TD-002: Autenticación con API Key
**Contexto**: Los requerimientos técnicos exigen autenticación robusta, pero el challenge no incluye UI ni gestión de usuarios como objetivo principal.
**Decisión**: Implementar un modelo básico de autenticación mediante API Key en el header `X-API-Key`.
**Justificación**: Proporciona un mecanismo de autenticación funcional y probeable vía cURL/Postman sin la complejidad de OAuth2/JWT. Es suficiente para proteger los endpoints y demostrar el patrón de seguridad.

## TD-003: CI/CD y monitoreo diferidos
**Contexto**: Los requerimientos técnicos incluyen pipelines CI/CD, métricas Prometheus y tracing distribuido.
**Decisión**: Diferir la implementación de CI/CD y monitoreo avanzado a una iteración futura. Solo se documentará la estrategia.
**Justificación**: El alcance del challenge se centra en la API funcional. Se priorizó la calidad del código, testing y documentación sobre infraestructura de CI/CD.

## TD-004: Router chi
**Contexto**: Se necesita un router HTTP para montar los endpoints REST.
**Decisión**: Usar `github.com/go-chi/chi/v5`.
**Justificación**: Ligero, idiomático en Go, compatible con el `net/http` estándar, y middleware composable sin dependencias pesadas.

## TD-005: Driver pgx para PostgreSQL
**Contexto**: Se necesita un driver para conectarse a PostgreSQL desde Go.
**Decisión**: Usar `github.com/jackc/pgx/v5` con `pgxpool`.
**Justificación**: Mayor rendimiento que `lib/pq`, soporte nativo de connection pooling, y manejo eficiente de tipos PostgreSQL.

## TD-006: Clean Architecture
**Contexto**: Los requerimientos técnicos exigen arquitectura hexagonal o clean architecture con separación de responsabilidades.
**Decisión**: Implementar Clean Architecture con capas: domain → repository → service → handler.
**Justificación**: Facilita testing (mock de interfaces entre capas), mantiene baja dependencia entre componentes, y cumple con principios SOLID.

## TD-007: Colección Postman como herramienta de pruebas de integración
**Contexto**: El challenge requiere que la API pueda probarse fácilmente vía Postman/cURL.
**Decisión**: Crear una colección Postman completa en `docs/API_Feedbacks.postman_collection.json` con scripts de test automatizados y variables de colección.
**Justificación**: Permite ejecutar todos los escenarios (happy path y excepciones) de forma reproducible. Los scripts de test validan automáticamente status codes, estructura de respuesta y reglas de negocio. Las variables de colección (`base_url`, `api_key`, `feedback_id`) facilitan la configuración y el encadenamiento de requests.

## TD-008: Formato de feedback_id secuencial (f-####)
**Contexto**: Se requiere que el identificador de feedback tenga formato `f-####` (e.g. `f-0001`) en lugar de UUID.
**Decisión**: Usar un contador atómico (`sync/atomic`) en el servicio para generar IDs secuenciales con formato `f-####`. Se eliminó la dependencia de `github.com/google/uuid`.
**Justificación**: Proporciona IDs legibles y predecibles. El contador atómico es thread-safe y no requiere consultas adicionales a la base de datos. En un escenario de producción multi-instancia, se podría migrar a una secuencia PostgreSQL.

## TD-009: Validación de formato user_id (u-###)
**Contexto**: Se requiere que el campo `user_id` siga el formato `u-###` (e.g. `u-001`, `u-015`).
**Decisión**: Agregar validación con regex `^u-\d{3}$` en la capa de dominio.
**Justificación**: Garantiza consistencia en los datos desde la capa más interna, independiente del punto de entrada (API, CLI, etc.).

## TD-010: Valores de feedback_type en español
**Contexto**: Se requiere cambiar los valores permitidos de feedback_type a: bug, sugerencia, elogio, duda, queja.
**Decisión**: Actualizar las constantes del dominio, la validación, el constraint SQL y la documentación OpenAPI.
**Justificación**: Alinea los valores del sistema con el idioma del dominio de negocio. Se mantiene `bug` como término técnico universal.

## TD-011: Timestamps sin milisegundos
**Contexto**: Los campos `created_at` y `updated_at` deben presentarse sin milisegundos (e.g. `2026-01-10T09:12:11Z`).
**Decisión**: Truncar timestamps a segundos con `time.Truncate(time.Second)` al crearlos, y formatear con `time.RFC3339` en el DTO de respuesta.
**Justificación**: El truncado en la capa de servicio elimina los milisegundos desde la generación. El formateo en el DTO asegura que la serialización JSON produzca el formato correcto sin depender del marshaler por defecto de `time.Time`.

## TD-012: Corrección de migración inline en main.go
**Contexto**: La función `runMigrations` en `cmd/server/main.go` tenía una migración SQL hardcodeada que no coincidía con el modelo de dominio actual: usaba `id UUID` como PK (en lugar de `feedback_id VARCHAR(10)`) y tipos en inglés (`suggestion`, `praise`, `question`) en lugar de los valores en español definidos en TD-010.
**Decisión**: Actualizar la migración inline para que coincida con el archivo de migración canónico (`001_create_feedbacks.sql`) y el modelo de dominio.
**Justificación**: La inconsistencia causaba errores 500 en todas las operaciones que interactuaban con la base de datos (Create, GetByID, Update, List). El `CREATE TABLE IF NOT EXISTS` impedía que la migración se auto-corrigiera al reiniciar la aplicación.
