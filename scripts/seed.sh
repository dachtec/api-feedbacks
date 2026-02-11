#!/bin/bash
# seed.sh - Seeds sample feedback data via the API
# Usage: bash scripts/seed.sh
# Loads data from docs/seed-data.json

set -e

API_URL="${API_URL:-http://localhost:8080}"
API_KEY="${API_KEY:-my-secret-api-key}"

# Resolve the path to seed-data.json relative to the project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SEED_FILE="$PROJECT_ROOT/docs/seed-data.json"

if [ ! -f "$SEED_FILE" ]; then
  echo "❌ Seed file not found: $SEED_FILE"
  exit 1
fi

# Validate that jq is available
if ! command -v jq &> /dev/null; then
  echo "❌ jq is required but not installed. Install it with: brew install jq"
  exit 1
fi

TOTAL=$(jq length "$SEED_FILE")
echo "🌱 Seeding $TOTAL feedbacks from docs/seed-data.json..."
echo ""

SUCCESS=0
FAIL=0

for i in $(seq 0 $((TOTAL - 1))); do
  # Extract only the fields accepted by the API (exclude feedback_id, created_at, updated_at)
  PAYLOAD=$(jq -c ".[$i] | {user_id, feedback_type, rating, comment}" "$SEED_FILE")

  FEEDBACK_ID=$(jq -r ".[$i].feedback_id" "$SEED_FILE")
  echo -n "  [$((i + 1))/$TOTAL] Sending $FEEDBACK_ID... "

  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/v1/feedbacks" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d "$PAYLOAD")

  HTTP_CODE=$(echo "$RESPONSE" | tail -1)
  BODY=$(echo "$RESPONSE" | sed '$d')

  if [ "$HTTP_CODE" -eq 201 ]; then
    echo "✅"
    SUCCESS=$((SUCCESS + 1))
  else
    echo "❌ (HTTP $HTTP_CODE)"
    echo "    $BODY"
    FAIL=$((FAIL + 1))
  fi
done

echo ""
echo "📊 Results: $SUCCESS succeeded, $FAIL failed (of $TOTAL total)"
echo "✅ Seeding complete! Use GET /api/v1/feedbacks to verify."
