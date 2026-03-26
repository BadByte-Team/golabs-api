#!/bin/bash

# Salir inmediatamente si un comando falla
set -e

echo "🛑 Deteniendo los contenedores actuales..."
docker compose down

echo "🧹 Eliminando volúmenes (incluyendo la configuración en caché de MariaDB)..."
docker compose down -v

echo "🚀 Levantando el entorno nuevamente en segundo plano..."
docker compose up -d

echo "✅ ¡Entorno reiniciado con éxito! Estado actual:"
docker compose ps
