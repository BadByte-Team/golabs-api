.PHONY: up down migrate logs db-logs migrate-logs restart status \
        connect history backup restore repair validate info clean-flyway ps \
        check-env rotate-password seed reset-dev diff \
        backup-auto backup-list backup-clean \
        stats slow-queries size dump-schema ping version

# Extrae una variable del .env
define get_env
$(shell grep ^$(1)= .env | cut -d '=' -f2)
endef

DB_PASSWORD := $(call get_env,MARIADB_ROOT_PASSWORD)
DB_NAME     := $(call get_env,MARIADB_DATABASE)
DB_PORT     := $(call get_env,MARIADB_PORT)

# ─── Verificacion ─────────────────────────────────────────────────────────────

# Verifica que el .env tenga todas las variables necesarias
check-env:
	@echo "Verificando variables de entorno..."
	@for var in MARIADB_ROOT_PASSWORD MARIADB_DATABASE MARIADB_USER MARIADB_PASSWORD MARIADB_PORT; do \
		grep -q "^$$var=" .env || (echo "Falta $$var en .env" && exit 1); \
		val=$$(grep "^$$var=" .env | cut -d'=' -f2); \
		[ -z "$$val" ] && echo "$$var esta vacia en .env" && exit 1 || echo "  OK $$var"; \
	done
	@echo "Todas las variables estan presentes."

# ─── Servicios ────────────────────────────────────────────────────────────────

# Levanta todo y aplica migraciones automaticamente
up: check-env
	@echo "Levantando servicios..."
	docker compose up -d --build mariadb
	@echo "Aplicando migraciones..."
	docker compose up --build flyway
	@echo "Listo!"

# Baja todos los servicios
down:
	@echo "Bajando servicios..."
	docker compose down

# Reinicia todo y aplica migraciones
restart: down up

# Estado rapido de los contenedores
ps:
	docker compose ps

# Verifica si MariaDB esta respondiendo
ping:
	@docker exec golabs-mariadb mariadb-admin ping -h 127.0.0.1 -u root -p$(DB_PASSWORD) \
		&& echo "MariaDB esta respondiendo." || echo "MariaDB no responde."

# Muestra la version de MariaDB y Flyway
version:
	@echo "--- MariaDB ---"
	@docker exec golabs-mariadb mariadb -u root -p$(DB_PASSWORD) -e "SELECT VERSION();"
	@echo "--- Flyway ---"
	@docker compose run --rm flyway -v

# ─── Migraciones ──────────────────────────────────────────────────────────────

# Solo corre las migraciones
migrate:
	@echo "Corriendo migraciones..."
	docker compose up --build flyway
	@echo "Migraciones aplicadas!"

# Muestra migraciones pendientes sin aplicarlas
diff:
	@echo "Migraciones pendientes:"
	docker compose run --rm flyway info | grep "Pending"

# Muestra el estado de todas las migraciones
info:
	docker compose run --rm flyway info

# Verifica que los SQLs no fueron modificados despues de aplicarse
validate:
	docker compose run --rm flyway validate

# Repara migraciones que fallaron a medias
repair:
	docker compose run --rm flyway repair

# Reinicia solo el contenedor de Flyway
clean-flyway:
	docker compose rm -f flyway

# ─── Desarrollo ───────────────────────────────────────────────────────────────

# Carga datos de prueba desde seed.sql
seed:
	@if [ ! -f ./seed.sql ]; then echo "No se encontro seed.sql en la raiz del proyecto"; exit 1; fi
	@echo "Cargando datos de prueba..."
	docker exec -i golabs-mariadb mariadb -u root -p$(DB_PASSWORD) $(DB_NAME) < ./seed.sql
	@echo "Seed aplicado."

# Borra y recrea la DB completa (solo para desarrollo local)
reset-dev:
	@echo "Advertencia: esto eliminara todos los datos de $(DB_NAME)."
	@read -p "Escribi CONFIRMAR para continuar: " confirm && [ "$$confirm" = "CONFIRMAR" ] || exit 1
	docker exec golabs-mariadb mariadb -u root -p$(DB_PASSWORD) \
		-e "DROP DATABASE IF EXISTS $(DB_NAME); CREATE DATABASE $(DB_NAME);"
	@echo "Base de datos recreada. Aplicando migraciones..."
	docker compose up --build flyway
	@echo "Listo."

# Cambia la password de root y actualiza el .env
rotate-password:
	@read -s -p "Nueva password: " newpass && echo "" && \
	docker exec golabs-mariadb mariadb -u root -p$(DB_PASSWORD) \
		-e "ALTER USER 'root'@'%' IDENTIFIED BY '$$newpass'; FLUSH PRIVILEGES;" && \
	sed -i "s/^MARIADB_ROOT_PASSWORD=.*/MARIADB_ROOT_PASSWORD=$$newpass/" .env && \
	echo "Password actualizada en la DB y en el .env."

# ─── Acceso y consulta ────────────────────────────────────────────────────────

# Entra a la consola de MariaDB
connect:
	docker exec -it golabs-api-db-1 mariadb -u root -p$(DB_PASSWORD) $(DB_NAME)

# Muestra el historial de migraciones aplicadas
history:
	docker exec -it golabs-mariadb mariadb -u root -p$(DB_PASSWORD) $(DB_NAME) \
		-e "SELECT version, description, type, installed_on, success FROM flyway_schema_history ORDER BY installed_rank;"

# Muestra el tamaño de la DB y cada tabla
size:
	docker exec -it golabs-mariadb mariadb -u root -p$(DB_PASSWORD) \
		-e "SELECT table_name, ROUND((data_length + index_length) / 1024 / 1024, 2) AS 'Tamaño (MB)' \
		    FROM information_schema.tables WHERE table_schema='$(DB_NAME)' ORDER BY (data_length + index_length) DESC;"

# Muestra las queries mas lentas
slow-queries:
	docker exec -it golabs-mariadb mariadb -u root -p$(DB_PASSWORD) \
		-e "SELECT query_time, lock_time, rows_examined, sql_text \
		    FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10;"

# Exporta solo la estructura sin datos
dump-schema:
	@mkdir -p ./backups
	docker exec golabs-mariadb mariadb-dump \
		-u root -p$(DB_PASSWORD) --no-data $(DB_NAME) \
		> ./backups/schema_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Schema exportado en ./backups/"

# ─── Backup y restore ─────────────────────────────────────────────────────────

# Exporta un backup al volumen nombrado golabs-backups
backup:
	@echo "Generando backup en el volumen golabs-backups..."
	docker exec golabs-mariadb bash -c \
		"mariadb-dump -u root -p$(DB_PASSWORD) $(DB_NAME) > /backups/golabs_$$(date +%Y%m%d_%H%M%S).sql"
	@echo "Backup guardado en el volumen golabs-backups."

# Lista todos los backups disponibles en el volumen con fecha y tamaño
backup-list:
	@echo "Backups disponibles en el volumen golabs-backups:"
	@docker exec golabs-mariadb bash -c "ls -lh /backups/*.sql 2>/dev/null | awk '{print \$$5, \$$9}'" \
		|| echo "No hay backups disponibles."

# Elimina backups con mas de N dias (por defecto 7): make backup-clean days=30
backup-clean:
	@days=$(or $(days),7); \
	echo "Eliminando backups con mas de $$days dias del volumen..."; \
	docker exec golabs-mariadb bash -c "find /backups -name '*.sql' -mtime +$$days -print -delete"
	@echo "Limpieza completada."

# Configura backup automatico cada 6 horas via cron en el contenedor
backup-auto:
	@echo "Configurando backup automatico cada 6 horas..."
	docker exec golabs-mariadb bash -c \
		'echo "0 */6 * * * mariadb-dump -u root -p$$MARIADB_ROOT_PASSWORD $$MARIADB_DATABASE > /backups/auto_\$$(date +\%Y\%m\%d_\%H\%M\%S).sql" | crontab -'
	@echo "Backup automatico configurado."

# Restaura desde un backup del volumen: make restore file=golabs_20240101.sql
restore:
	@if [ -z "$(file)" ]; then echo "Uso: make restore file=nombre_del_backup.sql"; exit 1; fi
	@echo "Restaurando desde /backups/$(file)..."
	docker exec -i golabs-mariadb bash -c \
		"mariadb -u root -p$(DB_PASSWORD) $(DB_NAME) < /backups/$(file)"
	@echo "Restore completado desde $(file)."

# ─── Monitoreo ────────────────────────────────────────────────────────────────

# Muestra uso de CPU, memoria y disco en tiempo real
stats:
	docker stats golabs-mariadb

# ─── Logs ─────────────────────────────────────────────────────────────────────

# Logs de todos los servicios
logs:
	docker compose logs -f

# Logs solo de MariaDB
db-logs:
	docker compose logs -f mariadb

# Logs solo de Flyway
migrate-logs:
	docker compose logs flyway
