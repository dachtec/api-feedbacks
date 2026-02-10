---
trigger: manual
---

# Prompt: Desarrollo de Microservicios en Golang

## 1. CONFIGURACIÓN DEL ENTORNO DE DESARROLLO

### 1.1 Contenedor de Desarrollo
- Utilizar imagen base oficial de Golang con versión estable específica
- Configurar volúmenes para código fuente y caché de dependencias
- Incluir herramientas de desarrollo: linters, formatters, debuggers
- Configurar hot-reload para desarrollo ágil
- Establecer variables de entorno para desarrollo local
- Incluir herramientas de testing y profiling

### 1.2 Gestión de Dependencias
- Utilizar Go Modules para gestión de dependencias
- Mantener go.mod y go.sum versionados
- Establecer política de actualización de dependencias
- Documentar dependencias críticas y su propósito

## 2. ESTRUCTURA DEL PROYECTO

### 2.1 Organización de Directorios
- Implementar estructura estándar de proyecto Golang
- Separar código por capas: handlers, services, repositories, models
- Directorio independiente para configuración
- Directorio separado para scripts y utilidades
- Directorio para migraciones de base de datos
- Directorio para pruebas de integración

### 2.2 Separación de Responsabilidades
- Aplicar principios SOLID en la organización del código
- Mantener cohesión alta y acoplamiento bajo entre paquetes
- Definir interfaces claras entre capas
- Evitar dependencias cíclicas

## 3. PATRONES DE DISEÑO

### 3.1 Patrones Arquitectónicos
- Implementar arquitectura hexagonal o clean architecture
- Aplicar patrón Repository para acceso a datos
- Utilizar patrón Factory para creación de objetos complejos
- Implementar patrón Strategy para algoritmos intercambiables
- Aplicar patrón Adapter para integración con sistemas externos

### 3.2 Patrones de Concurrencia
- Utilizar goroutines y channels de manera eficiente
- Implementar patrón Worker Pool para procesamiento paralelo
- Aplicar patrón Pipeline para procesamiento de datos en etapas
- Utilizar context para propagación de cancelación y timeouts
- Implementar sincronización adecuada con sync package

### 3.3 Patrones de Resiliencia
- Implementar Circuit Breaker para llamadas a servicios externos
- Aplicar patrón Retry con backoff exponencial
- Implementar Timeout en todas las operaciones de red
- Utilizar patrón Bulkhead para aislamiento de recursos
- Implementar Rate Limiting para protección de endpoints

## 4. DESARROLLO DEL MICROSERVICIO

### 4.1 API y Comunicación
- Definir contratos API claros y versionados
- Implementar RESTful API siguiendo convenciones HTTP
- Utilizar formato JSON para intercambio de datos
- Implementar validación de entrada en todos los endpoints
- Documentar API con especificación OpenAPI/Swagger
- Considerar gRPC para comunicación entre microservicios

### 4.2 Manejo de Errores
- Implementar jerarquía de errores personalizada
- Propagar errores de manera consistente
- No exponer detalles internos en mensajes de error
- Registrar errores con contexto suficiente
- Implementar códigos de error estándar

### 4.3 Logging y Observabilidad
- Implementar logging estructurado
- Definir niveles de log apropiados
- Incluir correlation IDs para trazabilidad
- Evitar logging de información sensible
- Implementar métricas de negocio y técnicas
- Preparar para integración con sistemas de monitoreo

### 4.4 Configuración
- Externalizar toda configuración del código
- Utilizar variables de entorno para configuración
- Implementar valores por defecto razonables
- Validar configuración al inicio de la aplicación
- Soportar múltiples ambientes

## 5. RENDIMIENTO

### 5.1 Optimización de Recursos
- Reutilizar conexiones de base de datos mediante pool
- Implementar caching estratégico
- Minimizar allocaciones de memoria innecesarias
- Utilizar buffers y pools de objetos donde corresponda
- Optimizar consultas a base de datos

### 5.2 Procesamiento Asíncrono
- Identificar operaciones que pueden ser asíncronas
- Implementar procesamiento en background para tareas pesadas
- Utilizar message queues para desacoplamiento
- Implementar procesamiento batch cuando sea apropiado

### 5.3 Perfilado y Benchmarking
- Incluir benchmarks para funciones críticas
- Preparar endpoints de profiling
- Identificar y documentar bottlenecks conocidos
- Establecer baseline de rendimiento

## 6. SEGURIDAD

### 6.1 Autenticación y Autorización
- Implementar autenticación robusta
- Aplicar principio de mínimo privilegio
- Validar y sanitizar todas las entradas
- Implementar autorización a nivel de endpoint
- Utilizar tokens seguros con expiración

### 6.2 Protección de Datos
- Encriptar datos sensibles en tránsito y reposo
- No almacenar secretos en código o repositorio
- Utilizar gestores de secretos
- Implementar enmascaramiento de datos sensibles en logs
- Aplicar prácticas seguras de gestión de credenciales

### 6.3 Seguridad de APIs
- Implementar rate limiting por cliente
- Proteger contra inyecciones SQL y otras vulnerabilidades OWASP
- Validar content-types
- Implementar CORS apropiadamente
- Utilizar HTTPS exclusivamente en producción

### 6.4 Dependencias y Vulnerabilidades
- Escanear dependencias regularmente
- Mantener dependencias actualizadas
- Documentar excepciones de seguridad conocidas
- Implementar proceso de parching

## 7. TESTING

### 7.1 Estrategia de Testing
- Implementar pruebas unitarias con cobertura mínima definida
- Crear pruebas de integración para flujos críticos
- Implementar pruebas de contrato para APIs
- Considerar pruebas de carga para endpoints críticos
- Utilizar mocks e interfaces para facilitar testing

### 7.2 Calidad de Código
- Configurar linters y aplicar sus recomendaciones
- Mantener código formateado consistentemente
- Aplicar análisis estático de código
- Implementar revisión de código
- Documentar código complejo

## 8. RESILIENCIA

### 8.1 Manejo de Fallos
- Implementar graceful shutdown
- Manejar señales del sistema operativo apropiadamente
- Implementar health checks y readiness probes
- Preparar para recuperación automática de errores
- Implementar fallbacks para operaciones críticas

### 8.2 Tolerancia a Fallos
- No asumir disponibilidad de servicios externos
- Implementar timeouts en todas las operaciones de red
- Manejar errores transitorios con reintentos
- Implementar degradación controlada de funcionalidad
- Documentar comportamiento en caso de fallo de dependencias

## 9. ESCALABILIDAD

### 9.1 Diseño Stateless
- Evitar estado compartido en instancias del servicio
- Externalizar sesiones y cache
- Diseñar para escalamiento horizontal
- Evitar dependencias en filesystem local

### 9.2 Gestión de Recursos
- Implementar límites de recursos configurables
- Utilizar connection pooling
- Implementar backpressure cuando sea necesario
- Monitorear consumo de recursos
- Documentar requerimientos de recursos

### 9.3 Base de Datos
- Implementar estrategias de particionado si es necesario
- Utilizar índices apropiadamente
- Implementar paginación en consultas que retornan múltiples resultados
- Considerar read replicas para operaciones de lectura
- Optimizar queries frecuentes

## 10. CONTENERIZACIÓN

### 10.1 Dockerfile de Producción
- Utilizar multi-stage builds
- Usar imagen base mínima para producción
- Implementar usuario no-root
- Minimizar layers y tamaño de imagen
- Incluir solo dependencias necesarias para runtime
- Establecer labels informativos

### 10.2 Configuración del Contenedor
- Exponer solo puertos necesarios
- Configurar health checks
- Establecer límites de recursos
- Implementar graceful shutdown
- Definir variables de entorno requeridas

### 10.3 Seguridad del Contenedor
- Escanear imágenes por vulnerabilidades
- No incluir secretos en la imagen
- Utilizar image digest en vez de tags
- Mantener imágenes base actualizadas
- Aplicar principio de mínimo privilegio

## 11. DOCUMENTACIÓN

### 11.1 Documentación Técnica
- README con instrucciones claras de setup
- Documentar decisiones arquitectónicas importantes
- Incluir diagramas de arquitectura y flujo
- Documentar variables de entorno
- Incluir guía de troubleshooting

### 11.2 Documentación de API
- Mantener especificación OpenAPI actualizada
- Documentar ejemplos de request/response
- Documentar códigos de error posibles
- Incluir información de autenticación
- Documentar rate limits y restricciones

## 12. CI/CD

### 12.1 Pipeline de Integración
- Ejecutar linters y análisis estático
- Ejecutar suite completa de tests
- Verificar cobertura de código
- Construir imagen de contenedor
- Escanear vulnerabilidades
- Validar manifiestos de despliegue

### 12.2 Despliegue
- Implementar estrategia de versionamiento semántico
- Generar changelog automático
- Tagear imágenes apropiadamente
- Implementar rollback automático en caso de fallo
- Mantener artefactos de build

## 13. MONITOREO Y MÉTRICAS

### 13.1 Métricas de Aplicación
- Exponer métricas en formato estándar (Prometheus)
- Incluir métricas de negocio relevantes
- Implementar métricas de latencia por endpoint
- Monitorear tasa de errores
- Incluir métricas de recursos

### 13.2 Logging Centralizado
- Estructurar logs para agregación
- Incluir metadata relevante en logs
- Implementar niveles de log apropiados
- Preparar para integración con sistemas centralizados
- Establecer políticas de retención

### 13.3 Tracing Distribuido
- Preparar instrumentación para tracing
- Propagar trace context
- Identificar transacciones críticas para tracing
- Documentar integración con sistemas de tracing

## 14. ENTREGA

### 14.1 Artefactos Requeridos
- Código fuente completo y organizado
- Dockerfile para desarrollo y producción
- Docker Compose para ambiente local
- Documentación técnica completa
- Scripts de utilidad necesarios
- Configuraciones de ejemplo
- Colección de pruebas de API

### 14.2 Checklist de Calidad
- Código pasa todos los linters configurados
- Cobertura de tests cumple umbral mínimo
- Documentación está completa y actualizada
- Imagen de contenedor está optimizada
- Vulnerabilidades conocidas están documentadas o resueltas
- Health checks funcionan correctamente
- Variables de entorno están documentadas