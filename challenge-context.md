# Análisis de Feedbacks de la Plataforma

## Contexto del Challenge
Queremos construir una base para capturar y consultar feedbacks de usuarios sobre nuestra plataforma digital. Esto servirá para análisis internos, dashboards y detección de problemas de experiencia, desempeño y percepción de los usuarios.

## Objetivo
Construir un backend en Golang (sin interfaz gráfica) que permita:
* Crear, actualizar y consultar feedbacks de usuarios sobre la plataforma.
* Exponer endpoints que soporten filtros y casos típicos de consulta.
* Debe poder probarse fácilmente (Postman/cURL).

### ¿Qué es un "feedback"?
Un feedback debe contener al menos:
* **user_id**: identificador del usuario que envió el feedback.
* **feedback_type**: tipo (ej. bug, sugerencia, elogio, duda).
* **rating**: entero de 1 a 5 sobre la experiencia.
* **comment**: texto descriptivo del feedback.
* **Metadata**:
    * created_at
    * updated_at

## ¿Qué esperamos que entregue tu solución?

### Una API usable
Define los endpoints (los que consideres necesarios) para:
* Crear un feedback.
* Actualizar un feedback (por ejemplo, rating, comment o feedback_type).
* Consultar feedbacks con filtros.

### Consultas por filtros
Debe existir al menos una consulta que permita filtrar (puede ser vía query params, body o como prefieras). Ejemplos de filtros:
* user_id
* feedback_type
* rango de rating (min_rating, max_rating)
* rango de fechas (created_from, created_to)

### Decisiones técnicas
* README con decisiones técnicas, cómo ejecutar y cómo probar la solución.

* Uso de IA (obligatorio, con trazabilidad)
* Entrega un archivo **PROMPTS.md** con:
	* Prompts principales utilizados.
	* Cómo iteraste y por qué.
	* Ejemplos donde validaste/corregiste sugerencias de la IA.
	* Decisiones y tradeoffs (por qué tu API/estructura es así).

## Restricciones
* Puedes usar bibliotecas externas en Go, pero justifícalo brevemente en el README.
* No necesitas UI.