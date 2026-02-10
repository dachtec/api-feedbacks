# PROMPTS.md — Trazabilidad del uso de IA

## Resumen

Este proyecto fue desarrollado con asistencia de IA (Claude/Antigravity). Se documentan los prompts principales, las iteraciones realizadas, y las decisiones tomadas durante el proceso.

---

## Prompt 1: Análisis del reto y plan de implementación

**Prompt**: "Actúa como desarrollador experto en Golang. Tu objetivo es analizar de forma exhaustiva el reto que se relaciona en el archivo challenge-context.md y definir el plan de implementación de acuerdo con los requerimientos técnicos definidos en tech-requirements.md."

**Iteración**: La IA generó un plan de implementación detallado con estructura de directorios, arquitectura, endpoints, dependencias y estrategia de testing. Se identificaron 3 puntos de decisión para revisión del usuario.

**Decisiones tomadas**:
- PostgreSQL como motor de persistencia (aprobado por el usuario)
- API Key como mecanismo de autenticación (el usuario solicitó implementarla; la IA inicialmente proponía omitirla)
- CI/CD y monitoreo diferidos a iteración futura (aprobado por el usuario)

---

## Prompt 2: Implementación del proyecto

**Prompt**: Ejecución del plan aprobado.

**Iteración**: Se intentó usar `go mod init` localmente, pero Go no estaba instalado. El usuario indicó que todo el desarrollo debía hacerse dentro de Docker. Se ajustó la estrategia para crear todos los archivos fuente primero y usar Docker para compilar y ejecutar.

**Validación/Corrección**:
- La IA inicialmente asumió que Go estaba instalado localmente → corregido a desarrollo 100% containerizado
- La IA inicialmente proponía no implementar autenticación → corregido a API Key básica

---

## Decisiones y Tradeoffs

### ¿Por qué esta estructura de API?
- **Clean Architecture**: Permite testear cada capa de forma independiente con mocks
- **PATCH para updates**: Soporta actualizaciones parciales, más flexible que PUT
- **Query params para filtros**: Más natural para operaciones GET que body params
- **Paginación con limit/offset**: Simple y efectiva para el volumen esperado

### ¿Por qué estas bibliotecas?
- **chi en vez de gin**: Más ligero, sin reflection, compatible con net/http estándar
- **pgx en vez de lib/pq o GORM**: Mayor rendimiento, sin overhead de ORM, control total de queries
- **slog en vez de zerolog/zap**: Parte del stdlib desde Go 1.21, reduce dependencias externas

### ¿Por qué API Key y no JWT?
- El challenge no incluye gestión de usuarios ni UI
- API Key es suficiente para demostrar el patrón de autenticación
- JWT agregaría complejidad sin valor para este caso de uso
