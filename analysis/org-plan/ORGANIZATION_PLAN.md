# 🎯 Plan de acción — Análisis de Feedbacks

**Autor:** Daniel Chamorro
**Fecha:** 11 de Febrero, 2026

---

## Resumen Ejecutivo

![Infografía de la Propuesta](img/Infografia_propuesta_5.png)

---

## 1. Objetivo

> *"Para el cierre de Q1 2026, estabilizar los módulos de videollamadas, autenticación y soporte al cliente, elevando su calificación promedio de 1.58 a ≥ 3.0 y reduciendo la retroalimentación negativa del 58.3% al ≤ 35%, con base en los puntos de dolor identificados en el informe de inteligencia de producto."*

---

## 2. Indicadores clave (KPIs)

### KPI 1: Calificación promedio de módulos críticos

| Atributo | Detalle |
|---|---|
| **Definición** | Promedio de calificaciones (1–5) recibidas en videollamadas, autenticación y soporte |
| **Línea base** | **1.58** (ponderado: videollamadas 1.50 × 4 + autenticación 1.75 × 4 + soporte 1.50 × 2) |
| **Meta Q1** | ≥ **3.0** |
| **Meta Q2** | ≥ **4.0** |
| **Fuente** | API de retroalimentación, filtrado por tipo `bug` y `queja` en los tres módulos |


### KPI 2: Proporción de retroalimentación negativa

| Atributo | Detalle |
|---|---|
| **Definición** | Porcentaje de registros clasificados como sentimiento negativo respecto al total |
| **Línea base** | **58.3%** (14 de 24 registros) |
| **Meta Q1** | ≤ **35%** |
| **Meta Q2** | ≤ **20%** |
| **Fuente** | Análisis de sentimiento automatizado + encuesta de satisfacción dentro de la aplicación |


### KPI 3: Tasa de usuarios en riesgo de abandono

| Atributo | Detalle |
|---|---|
| **Definición** | Porcentaje de usuarios que reportan ≥ 2 calificaciones negativas (≤ 2) en 30 días |
| **Línea base** | **11.1%** (2 de 18 usuarios: u-008 con 3 bugs y u-015 con 2 quejas) |
| **Meta Q1** | ≤ **5%** |
| **Meta Q2** | ≤ **3%** |
| **Fuente** | Agrupación por `user_id` con filtro de calificación ≤ 2 |

---

## 3. Acciones prácticas

### Mitigación inmediata

#### Acción 1: Estabilización del módulo de videollamadas

| Atributo | Detalle |
|---|---|
| **Problema** | 4 bugs con calificación promedio 1.50: cortes en llamadas 1:1, fallos con +3 participantes, eco en grupales. Puntaje de severidad: 40.0 (el más alto) |
| **Qué haríamos** | Revisar los registros técnicos de los incidentes del 13 y 22 de enero. Implementar reconexión automática, degradación progresiva de calidad de video según el ancho de banda, y modo solo-audio como respaldo |
| **Responsables** | Equipo de infraestructura (especialista WebRTC + SRE) + QA |
| **Plazo** | 1–2 iteraciones (2–4 semanas) |

---

### Mejora y prevención

#### Acción 2: Monitoreo proactivo y reportes automatizados

| Atributo | Detalle |
|---|---|
| **Problema que previene** | Los picos de bugs del 13 y 22 de enero sugieren degradaciones que el monitoreo actual no detectó |
| **Qué haríamos** | (a) Tableros en tiempo real con métricas de videollamadas (fluctuación, pérdida de paquetes, latencia) y alertas automáticas. (b) Incorporar la retroalimentación negativa como señal dentro del proceso de integración continua: pruebas de regresión para módulos críticos y validación de calidad antes de cada despliegue. (c) Encuesta contextual dentro de la app después de cada videollamada y cada interacción con soporte (d) Generación de reportes automáticos con IA|
| **Responsables** | SRE + Desarrollo + Investigación UX |
| **Plazo** | Q2 2026 (diseño en iteración 1, implementación progresiva en iteraciones 2–4) |


#### Acción 3: Recuperación de usuarios en riesgo y mejora del soporte

| Atributo | Detalle |
|---|---|
| **Problema que previene** | u-008 reportó 3 bugs sin resolución; u-015 sufrió un fallo técnico y luego una mala experiencia con soporte. El soporte está amplificando la frustración en lugar de contenerla |
| **Qué haríamos** | (a) Sistema de alerta temprana: detectar usuarios con ≥ 2 quejas en 30 días y contactarlos de forma proactiva con un canal de atención prioritario. (b) Acuerdos de nivel de servicio: primera respuesta en ≤ 4 horas, resolución de bugs críticos en ≤ 24 horas. (c) Invitar a usuarios comprometidos (u-008, u-012) a un programa de pruebas anticipadas para convertir su frustración en colaboración. (d) Revisión integral del flujo de autenticación: renovación silenciosa de sesión, respaldo de 2FA por correo/TOTP, tolerancia a cambios de red |
| **Responsables** | Producto + Soporte + UX + Desarrollo |
| **Plazo** | Q2 2026 (acuerdos de nivel de servicio en iteración 1; programa de alerta temprana en iteraciones 2–3) |


#### Acción 4: Asistente virtual con IA para soporte de primer nivel

| Atributo | Detalle |
|---|---|
| **Problema que previene** | El soporte tiene calificación 1.50 y amplifica la frustración: u-015 esperó 7 días sin respuesta por un fallo de 2FA que se resuelve en minutos con una guía |
| **Qué haríamos** | Construir un asistente conversacional integrado en la app que resuelva incidentes típicos sin intervención humana. La arquitectura usaría RAG (generación aumentada por recuperación): se indexa la documentación de producto, las guías de resolución y los tickets resueltos en una base vectorial, y el modelo consulta esa base para generar respuestas contextualizadas. |
| **Responsables** | Desarrollo (especialista IA + backend) + Soporte (curación de la base de conocimiento y validación de respuestas) + Producto (priorización de casos de uso) |
| **Plazo** | Q3 2026 (2 a 3 iteraciones) |

---

## 4. Gestión de riesgos

| Riesgo | Impacto | Mitigación |
|---|---|---|
| **Volumen de datos insuficiente para medir avance** — Solo 24 registros en 16 días; los KPIs podrían no moverse visiblemente en Q1 | Los indicadores no reflejan el progreso real del equipo, generando frustración | Acelerar la encuesta dentro de la app (Acción 2c) desde la primera iteración para aumentar el volumen de datos cuanto antes |
| **Dependencia de un solo especialista WebRTC** — Si la estabilización de videollamadas depende de una persona, cualquier ausencia bloquea la acción más crítica | Retraso en la acción con mayor puntaje de severidad (40.0) | Documentar hallazgos desde el día 1 y asignar un segundo desarrollador como respaldo en la auditoría |
| **Resistencia del equipo de soporte a los acuerdos de nivel de servicio** — Un compromiso de respuesta en ≤ 4 horas puede percibirse como imposición | Los acuerdos se definen pero no se cumplen, erosionando la credibilidad del plan | Involucrar a soporte desde la definición de los acuerdos, revisar su capacidad real y ajustar plazos si es necesario |
| **Regresión en funcionalidades que hoy funcionan bien** — Mientras se corrigen videollamadas y autenticación, podríamos introducir errores en módulos estables (notificaciones, app móvil) | Perdemos las fortalezas que hoy sostienen la satisfacción (rating 5.00) | Ejecutar pruebas de humo obligatorias para las fortalezas del producto antes de cada despliegue |

---

## 5. Organización del equipo

### ¿Cómo comunicaría los resultados y prioridades?

1. **Kick-off con datos**: Reunión de inicio de trimestre con todo el equipo donde presento los hallazgos del reporte con datos concretos: la distribución bimodal, los usuarios en riesgo y los comentarios más representativos. Los datos generan urgencia compartida sin necesidad de imponer mandatos.

2. **Resumen visible**: El objetivo, los 3 KPIs y las 3 acciones se documentan en una página accesible (wiki interna) que se consulta en cada reunión diaria y de planificación.

3. **Tableros en vivo**: Los KPIs se proyectan en herramientas de monitoreo (Grafana/Datadog) visibles para todo el equipo. Ver cómo se mueven las métricas hace tangible el impacto de cada mejora.

### ¿Cómo dividiría y haría seguimiento de las tareas?

| Mecanismo | Detalle |
|---|---|
| **Épicas por acción** | Cada acción se convierte en una épica en Jira con responsable claro y fecha límite del trimestre |
| **Iteraciones de 2 semanas** | Cada épica se descompone en historias estimadas; la acción 1 (videollamadas) entra desde la iteración 1 por su severidad |
| **Daily focalizado** | Destino 2 minutos del daily para revisar avance de las acciones del plan. Los bloqueos se escalan ahí mismo |
| **Revisión semanal de KPIs** | Cada semana reviso los 3 KPIs contra la meta con el equipo y ajustamos prioridades según datos nuevos |
| **Retrospectiva con voz del usuario** | Cada cierre de iteración incorpora los comentarios recientes como insumo de la retrospectiva |

### ¿Cómo involucraría a los roles clave?

| Rol | Responsabilidad | Cómo lo involucro |
|---|---|---|
| **Desarrollo** | Auditoría WebRTC, corrección de 2FA, persistencia de sesión, pruebas de regresión | Trabajo en pares durante la auditoría técnica; comparto el contexto completo de los registros y comentarios de usuarios para que entiendan el *por qué* |
| **SRE / Infra** | Tableros de monitoreo, alertas, análisis de registros del 13 y 22 de enero | Participan desde el arranque para cruzar la retroalimentación con métricas de infraestructura |
| **Producto** | Priorización del trabajo pendiente, evaluación de integración con Google Drive, programa de pruebas anticipadas | Co-owner del objetivo; participa en la revisión quincenal de KPIs y en decisiones de priorización |
| **QA** | Pruebas de regresión, validación de calidad antes de despliegue, pruebas de humo para fortalezas | Define criterios de aceptación basados en escenarios reales (ej: "videollamada con +3 personas no debe cortarse") |
| **UX** | Encuesta contextual, auditoría del panel de configuración, diseño del flujo de reconexión | Lidera el diseño de la encuesta de satisfacción y traduce los comentarios de usuarios en mejoras de interfaz |
| **Soporte** | Acuerdos de nivel de servicio, canal prioritario, contacto proactivo con usuarios en riesgo | Tiene visibilidad de los usuarios detectados por el sistema de alerta temprana y acceso al tablero de retroalimentación |

### ¿Cómo aseguraría que el aprendizaje genere mejoras continuas?

1. **Análisis recurrente, no puntual**: Automatizo la categorización de retroalimentación (por tipo, sentimiento, tema) y genero un reporte semanal. Cada planificación de iteración arranca con los 3 comentarios más relevantes de la semana como insumo.

2. **Cerrar el ciclo con el usuario**: Cuando resolvemos un problema reportado (ej: u-008), le avisamos y le pedimos que pruebe de nuevo. Eso convierte una queja en una oportunidad de recuperar confianza.

3. **Hipótesis → Experimento → Dato**: Las 4 hipótesis del reporte tienen un responsable que busca los datos necesarios (registros de WebRTC, satisfacción neta, tasa de abandono post-soporte) y reporta en la revisión quincenal. No dejamos supuestos sin validar.

4. **Cuidar lo que funciona bien**: Las notificaciones (5.00), la app móvil (5.00) y el trabajo remoto (5.00) son nuestro diferencial. Defino pruebas de humo y métricas de regresión para estas áreas. Si alguna cae, es una alerta tan crítica como un bug nuevo.

---

> **Nota:** Este plan parte de una muestra de 24 registros en 16 días. Las metas cuantitativas deben recalibrarse conforme aumente el volumen de datos con la encuesta dentro de la aplicación. La dirección, sin embargo, es clara: los puntos de dolor están identificados, las causas raíz son alcanzables, y las fortalezas del producto nos dan una base sólida para construir.

## Referencias:
- [Digest - Informe de inteligencia de producto](/analysis/digest/feedback-analysis-report.md)
- [Digest - Visualizaciones interactivas](/analysis/digest/visualizations.html)
- [Plan de acción - Galería de Infografías](/analysis/org-plan/img/)
