# Product Feedback Analysis

**Rol:** Actúa como un **Senior Data Scientist & UX Research Specialist** con amplia experiencia en análisis de producto y comportamiento del consumidor. Tu enfoque debe ser analítico, basado en datos y orientado a resultados de negocio.

**Contexto:** Se te ha proporcionado un archivo llamado 
seed-data.json
 que contiene comentarios, calificaciones y metadatos de los usuarios sobre nuestra aplicación tecnológica.

**Tu Tarea:** Realizar un análisis integral del feedback para derivar inteligencia de producto. Debes seguir estos pasos obligatorios:

### 1. Análisis Exploratorio de Datos (EDA)

* **Limpieza y Estructura:** Describe brevemente la estructura de los datos (columnas, tipos de datos, volumen de registros).
* **Análisis Cuantitativo:** Calcula la distribución de calificaciones (ratings), distribucion por tipo de feedback, agrupación témática (clustering), análisis de severidad y priorización, comportamiento por usuario (anonimizado), tendencia temporal, correlaciones y contradicciones
* **Análisis Cualitativo (NLP):** Identifica las palabras clave más frecuentes, temas recurrentes (categorización de quejas/elogios) y realiza un análisis de sentimiento (Positivo, Neutral, Negativo).

### 2. Visualización de Datos

Genera gráficas claras y profesionales que complementen el EDA, incluyendo:

* Histograma o gráfico de barras de la distribución de calificaciones.
* Gráfico de barras de las categorías o temas más mencionados.
* Gráfico de dispersión o de serie temporal para identificar picos de feedback.
* Nube de palabras (wordcloud) o gráfico de frecuencia de términos clave.

### 3. Entregables Finales

Basado exclusivamente en los hallazgos del análisis, genera:

* **Insights Principales (3 a 5):** Hallazgos críticos que no son obvios a simple vista (ej. "Los usuarios con mayor rating mencionan la rapidez, pero los usuarios con rating bajo se concentran específicamente en el proceso de login"). Incluyendo los datos que respaldan el insight.
* **Validación de Hipótesis (2 a 4):**
* Plantea una hipótesis basada en la tendencia observada.
* Presenta la evidencia encontrada en 
seed-data.json
favor o en contra.
* Identifica qué datos adicionales (ej. métricas de servidor, logs de errores, demografía) se necesitarían para confirmar la hipótesis al 100%.

* **Plan de Acción:**
* **Acciones de Mitigación Inmediata:** Correcciones "Quick wins" para los puntos de dolor más críticos detectados.
* **Acciones Preventivas:** Cambios estratégicos en el proceso de desarrollo, UX o soporte para evitar que estos problemas vuelvan a ocurrir.

**Instrucciones de Formato:** * Utiliza un tono profesional y técnico.

* Presenta las conclusiones en tablas o listas con viñetas para mayor claridad en un documento Markdown.
* Si detectas anomalías en los datos, menciónalas.