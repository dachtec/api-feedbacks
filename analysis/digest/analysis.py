#!/usr/bin/env python3
"""
Comprehensive Feedback Analysis - Product Intelligence Report
Senior Data Scientist & UX Research Specialist Analysis
"""

import json
import os
import re
import math
from collections import Counter, defaultdict
from datetime import datetime

# ─────────────────────────────────────────────
# 1. LOAD DATA
# ─────────────────────────────────────────────
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_PATH = os.path.join(SCRIPT_DIR, "seed-data.json")

with open(DATA_PATH, "r", encoding="utf-8") as f:
    data = json.load(f)

print(f"Total records loaded: {len(data)}")

# ─────────────────────────────────────────────
# 2. DATA STRUCTURE OVERVIEW
# ─────────────────────────────────────────────
print("\n" + "=" * 60)
print("2. DATA STRUCTURE OVERVIEW")
print("=" * 60)

columns = list(data[0].keys())
print(f"Columns ({len(columns)}): {columns}")
print(f"Total records: {len(data)}")

# Type inference
sample = data[0]
for col in columns:
    val = sample[col]
    print(f"  - {col}: {type(val).__name__} (example: {val})")

# Check for nulls / missing
missing_count = {col: 0 for col in columns}
for row in data:
    for col in columns:
        if col not in row or row[col] is None or row[col] == "":
            missing_count[col] += 1

print(f"\nMissing values: {missing_count}")

# ─────────────────────────────────────────────
# 3. QUANTITATIVE ANALYSIS
# ─────────────────────────────────────────────
print("\n" + "=" * 60)
print("3. QUANTITATIVE ANALYSIS")
print("=" * 60)

# 3a. Rating Distribution
ratings = [d["rating"] for d in data]
rating_dist = Counter(ratings)
print("\n--- Rating Distribution ---")
for r in sorted(rating_dist.keys()):
    pct = rating_dist[r] / len(data) * 100
    print(f"  Rating {r}: {rating_dist[r]} ({pct:.1f}%)")

avg_rating = sum(ratings) / len(ratings)
median_rating = sorted(ratings)[len(ratings) // 2]
print(f"\n  Mean rating: {avg_rating:.2f}")
print(f"  Median rating: {median_rating}")
std_dev = math.sqrt(sum((r - avg_rating) ** 2 for r in ratings) / len(ratings))
print(f"  Std deviation: {std_dev:.2f}")

# 3b. Feedback Type Distribution
types = [d["feedback_type"] for d in data]
type_dist = Counter(types)
print("\n--- Feedback Type Distribution ---")
for t, count in type_dist.most_common():
    pct = count / len(data) * 100
    print(f"  {t}: {count} ({pct:.1f}%)")

# 3c. Rating by Feedback Type
print("\n--- Average Rating by Feedback Type ---")
type_ratings = defaultdict(list)
for d in data:
    type_ratings[d["feedback_type"]].append(d["rating"])

for t, r_list in sorted(type_ratings.items(), key=lambda x: sum(x[1]) / len(x[1])):
    avg = sum(r_list) / len(r_list)
    print(f"  {t}: avg={avg:.2f}, n={len(r_list)}, range=[{min(r_list)}-{max(r_list)}]")

# 3d. Severity & Prioritization (based on frequency × inverse rating)
print("\n--- Severity & Prioritization ---")
# Group by type, compute severity score = count * (6 - avg_rating) to prioritize low-rated high-volume issues
severity_scores = {}
for t, r_list in type_ratings.items():
    avg = sum(r_list) / len(r_list)
    severity = len(r_list) * (6 - avg)
    severity_scores[t] = {"count": len(r_list), "avg_rating": avg, "severity_score": round(severity, 2)}

for t, info in sorted(severity_scores.items(), key=lambda x: x[1]["severity_score"], reverse=True):
    print(f"  {t}: score={info['severity_score']}, count={info['count']}, avg_rating={info['avg_rating']:.2f}")

# 3e. User Behavior Analysis
print("\n--- User Behavior (Anonymized) ---")
user_feedback = defaultdict(list)
for d in data:
    user_feedback[d["user_id"]].append(d)

repeat_users = {uid: entries for uid, entries in user_feedback.items() if len(entries) > 1}
print(f"  Total unique users: {len(user_feedback)}")
print(f"  Users with single feedback: {len(user_feedback) - len(repeat_users)}")
print(f"  Users with multiple feedback: {len(repeat_users)}")

for uid, entries in sorted(repeat_users.items()):
    ftypes = [e["feedback_type"] for e in entries]
    ratings = [e["rating"] for e in entries]
    print(f"    {uid}: {len(entries)} feedbacks, types={ftypes}, ratings={ratings}")

# 3f. Temporal Trend Analysis
print("\n--- Temporal Trend ---")
dates_parsed = []
for d in data:
    dt = datetime.fromisoformat(d["created_at"].replace("Z", "+00:00"))
    dates_parsed.append((dt, d))

dates_parsed.sort(key=lambda x: x[0])
print(f"  Date range: {dates_parsed[0][0].strftime('%Y-%m-%d')} to {dates_parsed[-1][0].strftime('%Y-%m-%d')}")

# Daily volume
daily_volume = Counter()
daily_ratings = defaultdict(list)
for dt, d in dates_parsed:
    day = dt.strftime("%Y-%m-%d")
    daily_volume[day] += 1
    daily_ratings[day].append(d["rating"])

print("\n  Daily Volume & Avg Rating:")
for day in sorted(daily_volume.keys()):
    avg_r = sum(daily_ratings[day]) / len(daily_ratings[day])
    print(f"    {day}: {daily_volume[day]} feedback(s), avg_rating={avg_r:.1f}")

# Weekly aggregation
weekly_volume = Counter()
for dt, d in dates_parsed:
    week = dt.strftime("%Y-W%W")
    weekly_volume[week] += 1

print("\n  Weekly Volume:")
for week in sorted(weekly_volume.keys()):
    print(f"    {week}: {weekly_volume[week]} feedback(s)")

# 3g. Correlations & Contradictions
print("\n--- Correlations & Contradictions ---")
# Comment length vs rating
comment_lengths = [(d["rating"], len(d["comment"])) for d in data]
print("\n  Comment Length vs Rating:")
for r in sorted(set(ratings)):
    lengths = [cl for rat, cl in comment_lengths if rat == r]
    avg_len = sum(lengths) / len(lengths) if lengths else 0
    print(f"    Rating {r}: avg_comment_length={avg_len:.0f} chars (n={len(lengths)})")

# Contradiction: same user praises and complains
print("\n  User Contradictions:")
for uid, entries in user_feedback.items():
    types_set = set(e["feedback_type"] for e in entries)
    if "elogio" in types_set and ("queja" in types_set or "bug" in types_set):
        print(f"    {uid}: Has both praise and complaints")
        for e in entries:
            print(f"      [{e['feedback_type']}, rating={e['rating']}] {e['comment'][:60]}...")

# u-016 specific contradiction
if "u-016" in user_feedback:
    entries_016 = user_feedback["u-016"]
    print(f"\n  ⚠️ Notable: u-016 gave rating=1 (queja: soporte) AND rating=5 (elogio: app móvil)")

# ─────────────────────────────────────────────
# 4. QUALITATIVE ANALYSIS (NLP)
# ─────────────────────────────────────────────
print("\n" + "=" * 60)
print("4. QUALITATIVE ANALYSIS (NLP)")
print("=" * 60)

# Spanish stop words
STOP_WORDS = {
    "la", "el", "en", "de", "y", "a", "que", "es", "los", "las", "un", "una",
    "se", "no", "me", "mi", "muy", "por", "para", "con", "al", "del", "lo",
    "le", "si", "ya", "más", "mas", "pero", "su", "sus", "e", "o", "u",
    "este", "esta", "estos", "estas", "son", "ser", "ha", "he", "hay",
    "cuando", "como", "todo", "todos", "toda", "todas", "otro", "otra",
    "otros", "bien", "mal", "vez", "después", "despues", "debo", "siempre",
    "mucho", "hacer", "nunca", "así", "asi", "veces", "manera", "cada",
    "sola", "puedo", "pueden", "tengo"
}

# 4a. Keyword extraction
all_comments = " ".join(d["comment"] for d in data)
# Clean and tokenize
words = re.findall(r'[a-záéíóúñü]+', all_comments.lower())
filtered_words = [w for w in words if w not in STOP_WORDS and len(w) > 2]
word_freq = Counter(filtered_words)

print("\n--- Top 20 Keywords ---")
for word, count in word_freq.most_common(20):
    print(f"  {word}: {count}")

# 4b. Thematic Clustering (rule-based)
THEMES = {
    "Videollamadas/Comunicación": ["videollamada", "videollamadas", "corta", "cortan", "eco", "cortes", "personas"],
    "Autenticación/Login": ["login", "sesión", "sesion", "verificación", "verificacion", "código", "codigo", "sms", "iniciar", "desconecta"],
    "Rendimiento/Estabilidad": ["tarda", "cargar", "reinicia", "reiniciar", "actualización", "actualizacion", "lenta", "rápida", "rapida", "ágil", "agil"],
    "Búsqueda": ["búsqueda", "busqueda", "resultados", "incorrectos"],
    "Notificaciones": ["notificaciones", "notificación", "novedades", "email"],
    "Soporte al cliente": ["soporte", "tickets", "atención", "atencion", "demora", "resolver", "responde"],
    "UX/Interfaz": ["navegación", "navegacion", "interfaz", "intuitiva", "tema", "oscuro", "textos", "configuración", "configuracion", "fácil", "facil"],
    "Integraciones": ["integración", "integracion", "google", "drive", "exportar", "reportes"],
    "Trabajo remoto": ["remoto", "sincronización", "sincronizacion", "organizarme"],
}

print("\n--- Thematic Clustering ---")
theme_assignments = defaultdict(list)
for d in data:
    comment_lower = d["comment"].lower()
    assigned = False
    for theme, keywords in THEMES.items():
        if any(kw in comment_lower for kw in keywords):
            theme_assignments[theme].append(d)
            assigned = True
    if not assigned:
        theme_assignments["Otros"].append(d)

for theme, entries in sorted(theme_assignments.items(), key=lambda x: len(x[1]), reverse=True):
    avg_r = sum(e["rating"] for e in entries) / len(entries)
    print(f"  {theme}: {len(entries)} feedbacks, avg_rating={avg_r:.2f}")
    for e in entries:
        print(f"    [{e['feedback_type']}, r={e['rating']}] {e['comment'][:70]}")

# 4c. Sentiment Analysis (rule-based Spanish)
POSITIVE_WORDS = {
    "excelente", "encanta", "genial", "útil", "util", "rápida", "rapida", "intuitiva",
    "estable", "fácil", "facil", "felicidades", "ideal", "ágil", "agil",
    "organizarme", "ayudan", "útiles", "utiles", "gustaria", "increíble"
}
NEGATIVE_WORDS = {
    "falla", "corta", "cortan", "imposible", "nunca", "mala", "demora",
    "incorrectos", "tarda", "reinicia", "desconecta", "echo", "cortes",
    "críticos", "criticos", "problema", "problemas"
}

print("\n--- Sentiment Analysis ---")
sentiments = []
for d in data:
    comment_lower = d["comment"].lower()
    words_in_comment = set(re.findall(r'[a-záéíóúñü]+', comment_lower))
    pos_count = len(words_in_comment & POSITIVE_WORDS)
    neg_count = len(words_in_comment & NEGATIVE_WORDS)

    if pos_count > neg_count:
        sentiment = "Positivo"
    elif neg_count > pos_count:
        sentiment = "Negativo"
    else:
        # Use rating as tiebreaker
        if d["rating"] >= 4:
            sentiment = "Positivo"
        elif d["rating"] <= 2:
            sentiment = "Negativo"
        else:
            sentiment = "Neutral"

    sentiments.append(sentiment)
    d["sentiment"] = sentiment

sentiment_dist = Counter(sentiments)
print(f"  Positivo: {sentiment_dist.get('Positivo', 0)} ({sentiment_dist.get('Positivo', 0)/len(data)*100:.1f}%)")
print(f"  Neutral: {sentiment_dist.get('Neutral', 0)} ({sentiment_dist.get('Neutral', 0)/len(data)*100:.1f}%)")
print(f"  Negativo: {sentiment_dist.get('Negativo', 0)} ({sentiment_dist.get('Negativo', 0)/len(data)*100:.1f}%)")

# Sentiment by type
print("\n  Sentiment by Feedback Type:")
type_sentiments = defaultdict(Counter)
for d in data:
    type_sentiments[d["feedback_type"]][d["sentiment"]] += 1

for t in sorted(type_sentiments.keys()):
    print(f"    {t}: {dict(type_sentiments[t])}")

# Avg rating by sentiment
print("\n  Avg Rating by Sentiment:")
sentiment_ratings = defaultdict(list)
for d in data:
    sentiment_ratings[d["sentiment"]].append(d["rating"])
for s in ["Positivo", "Neutral", "Negativo"]:
    if sentiment_ratings[s]:
        avg = sum(sentiment_ratings[s]) / len(sentiment_ratings[s])
        print(f"    {s}: avg_rating={avg:.2f}, n={len(sentiment_ratings[s])}")

print("\n\n✅ Analysis Complete. Data exported for visualization.")

# ─────────────────────────────────────────────
# EXPORT STRUCTURED DATA FOR VISUALIZATION
# ─────────────────────────────────────────────
export_data = {
    "rating_distribution": dict(sorted(rating_dist.items())),
    "type_distribution": dict(type_dist.most_common()),
    "sentiment_distribution": dict(sentiment_dist),
    "theme_counts": {t: len(e) for t, e in theme_assignments.items()},
    "theme_avg_ratings": {t: round(sum(e["rating"] for e in entries)/len(entries), 2) for t, entries in theme_assignments.items()},
    "daily_volume": {k: v for k, v in sorted(daily_volume.items())},
    "daily_avg_rating": {k: round(sum(v)/len(v), 2) for k, v in sorted(daily_ratings.items())},
    "weekly_volume": {k: v for k, v in sorted(weekly_volume.items())},
    "keyword_freq": dict(word_freq.most_common(25)),
    "severity_scores": severity_scores,
    "avg_rating": round(avg_rating, 2),
    "median_rating": median_rating,
    "std_dev": round(std_dev, 2),
    "total_records": len(data),
    "unique_users": len(user_feedback),
    "repeat_users": len(repeat_users),
    "comment_length_by_rating": {},
}

for r in sorted(set(ratings)):
    lengths = [cl for rat, cl in comment_lengths if rat == r]
    export_data["comment_length_by_rating"][r] = round(sum(lengths)/len(lengths), 1) if lengths else 0

with open(os.path.join(SCRIPT_DIR, "analysis_results.json"), "w", encoding="utf-8") as f:
    json.dump(export_data, f, indent=2, ensure_ascii=False)

print("Results saved to analysis_results.json")
