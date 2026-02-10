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
