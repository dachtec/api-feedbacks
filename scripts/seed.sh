#!/bin/bash
# seed.sh - Seeds sample feedback data via the API
# Usage: bash scripts/seed.sh

API_URL="${API_URL:-http://localhost:8080}"
API_KEY="${API_KEY:-my-secret-api-key}"

echo "🌱 Seeding feedback data..."

curl -s -X POST "$API_URL/api/v1/feedbacks" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "user_id": "usr-001",
    "feedback_type": "bug",
    "rating": 2,
    "comment": "El botón de pago no responde en Safari"
  }' | jq .

curl -s -X POST "$API_URL/api/v1/feedbacks" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "user_id": "usr-002",
    "feedback_type": "praise",
    "rating": 5,
    "comment": "Excelente experiencia de usuario, muy intuitivo"
  }' | jq .

curl -s -X POST "$API_URL/api/v1/feedbacks" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "user_id": "usr-003",
    "feedback_type": "suggestion",
    "rating": 4,
    "comment": "Sería genial tener modo oscuro en la plataforma"
  }' | jq .

curl -s -X POST "$API_URL/api/v1/feedbacks" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "user_id": "usr-001",
    "feedback_type": "question",
    "rating": 3,
    "comment": "¿Cómo exporto mis datos en formato CSV?"
  }' | jq .

curl -s -X POST "$API_URL/api/v1/feedbacks" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "user_id": "usr-004",
    "feedback_type": "bug",
    "rating": 1,
    "comment": "La página de perfil tarda más de 10 segundos en cargar"
  }' | jq .

echo ""
echo "✅ Seeding complete! Use GET /api/v1/feedbacks to verify."
