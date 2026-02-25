# deployments/database

Infraestructura de base de datos para **golabs-api**.

Levanta **MariaDB 11** con gestión de migraciones automáticas vía **Flyway 10**, usando Docker Compose. Incluye herramientas de backup, restauración, monitoreo y acceso rápido a la consola.

---

## Contenido

```
deployments/database/
├── Dockerfile              # Imagen MariaDB 11.3 con healthcheck
├── docker-compose.yaml     # MariaDB + Flyway en redes aisladas
├── Makefile                # Todos los comandos de operación
├── .env                    # Variables de entorno (no commitear)
├── init/                   # Migraciones SQL versionadas con Flyway
│   ├── V1__crate_users_table.sql
│   ├── V2__ban_added.sql
│   ├── V3__create_events_table.sql
│   ├── V4__create_event_teams.sql
│   ├── V5__create_challenges_tables.sql
│   ├── V6__unique_username.sql
│   ├── V7__indexes.sql
│   └── V8__refresh_tokens.sql
└── volumes/
    └── mariadb/            # Datos persistentes de MariaDB (ignorado por git)
```

---

## Configuración

Crear el archivo `.env` en esta carpeta con las siguientes variables:

```env
MARIADB_ROOT_PASSWORD=tu_password_segura
MARIADB_DATABASE=golabs
MARIADB_USER=golabs
MARIADB_PASSWORD=password_del_usuario
MARIADB_PORT=3306
```

> [!CAUTION]
> Nunca commitear el archivo `.env`. Está incluido en `.gitignore`.

---

## Arquitectura de red

El `docker-compose.yaml` define dos redes Docker:

| Red | Tipo | Propósito |
|---|---|---|
| `golabs-db-internal` | Bridge interno (sin internet) | Comunicación privada entre MariaDB y Flyway |
| `golabs-db-net` | Bridge estándar | Conexión del backend Go a MariaDB |

MariaDB está conectada a **ambas** redes. Flyway solo a `golabs-db-internal`.

---

## Volúmenes

| Volumen | Descripción |
|---|---|
| `./volumes/mariadb` | Datos persistentes de MariaDB (bind mount local) |
| `golabs-backups` | Volumen nombrado Docker para backups — persiste aunque se baje el contenedor |

---

## Comandos Make

### Servicios

| Target | Descripción |
|---|---|
| `make up` | Levanta MariaDB y aplica las migraciones pendientes |
| `make down` | Detiene todos los contenedores |
| `make restart` | `down` + `up` |
| `make ps` | Estado rápido de los contenedores |
| `make ping` | Verifica que MariaDB está respondiendo |
| `make version` | Muestra la versión de MariaDB y Flyway |
| `make check-env` | Valida que todas las variables del `.env` estén definidas |

### Migraciones (Flyway)

| Target | Descripción |
|---|---|
| `make migrate` | Aplica las migraciones pendientes |
| `make diff` | Muestra las migraciones pendientes sin aplicarlas |
| `make info` | Estado completo de todas las migraciones |
| `make validate` | Verifica que los SQL no fueron modificados post-aplicación |
| `make repair` | Repara migraciones fallidas a medias |
| `make clean-flyway` | Elimina el contenedor de Flyway |

### Desarrollo

| Target | Descripción |
|---|---|
| `make connect` | Abre la consola interactiva de MariaDB |
| `make history` | Historial de migraciones aplicadas |
| `make seed` | Carga datos de prueba desde `./seed.sql` |
| `make reset-dev` | Elimina y recrea la DB completa (requiere confirmar) |
| `make rotate-password` | Rota la contraseña de root y actualiza el `.env` |

### Backup y restauración

| Target | Descripción |
|---|---|
| `make backup` | Genera un dump de la DB en el volumen `golabs-backups` |
| `make backup-list` | Lista los backups disponibles con fecha y tamaño |
| `make backup-clean [days=N]` | Elimina backups con más de N días (default: 7) |
| `make backup-auto` | Configura backup automático cada 6 horas vía cron en el contenedor |
| `make restore file=<nombre>.sql` | Restaura desde un backup del volumen |

### Consulta y monitoreo

| Target | Descripción |
|---|---|
| `make size` | Tamaño de la DB y cada tabla |
| `make slow-queries` | Las 10 queries más lentas |
| `make dump-schema` | Exporta el schema sin datos a `./backups/` |
| `make stats` | Uso de CPU, memoria y disco en tiempo real |
| `make logs` | Logs de todos los servicios |
| `make db-logs` | Logs solo de MariaDB |
| `make migrate-logs` | Logs solo de Flyway |

---

## Migraciones SQL

Las migraciones siguen la convención de Flyway: `V{version}__{descripcion}.sql`.

| Versión | Descripción |
|---|---|
| V1 | Tabla `users` con campos de identidad y puntos |
| V2 | Campo `banned` y `banned_at` en `users` |
| V3 | Tabla `events` con ciclo de vida y `max_team_size` |
| V4 | Tablas `event_teams` y `event_team_members` |
| V5 | Tablas `challenges`, `flags` y `solves` |
| V6 | Constraint `UNIQUE` en `users.username` |
| V7 | Índices de performance en FK y campos de búsqueda |
| V8 | Tabla `refresh_tokens` con soporte de revocación |

> [!NOTE]
> Flyway registra cada migración en `flyway_schema_history`. No se puede modificar un SQL ya aplicado; para cambios en esquema, crear una nueva migración `V9__...sql`.

---

## Uso típico de desarrollo

```bash
# Primera vez
make up        # Levanta MariaDB, aplica todas las migraciones

# Agregar nueva migración
echo "ALTER TABLE users ADD COLUMN ..." > init/V9__nueva_columna.sql
make migrate   # Aplica solo la nueva migración

# Inspeccionar el estado
make info      # Ver todas las migraciones y su estado

# Reset total (destructivo)
make reset-dev # Pide confirmación antes de borrar datos
```
