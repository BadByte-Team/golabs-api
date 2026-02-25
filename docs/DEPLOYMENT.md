# Despliegue — golabs-api

Esta guía te lleva de cero a la API corriendo, ya sea en tu máquina local o en un servidor de producción. Sigue los pasos en orden y no habrá sorpresas.

---

## Índice

1. [Prerrequisitos](#1-prerrequisitos)
2. [Quick start — todo en Docker (recomendado)](#2-quick-start--todo-en-docker-recomendado)
3. [Desarrollo local — API sin Docker + BD con Docker](#3-desarrollo-local--api-sin-docker--bd-con-docker)
4. [Producción en un servidor / VPS](#4-producción-en-un-servidor--vps)
5. [Variables de entorno — referencia completa](#5-variables-de-entorno--referencia-completa)
6. [Verificar que todo funciona](#6-verificar-que-todo-funciona)
7. [Solución de problemas frecuentes](#7-solución-de-problemas-frecuentes)

---

## 1. Prerrequisitos

Instala estas herramientas antes de continuar.

| Herramienta | Versión mínima | Verificar con |
|---|---|---|
| Docker | 24+ | `docker --version` |
| Docker Compose | v2 (plugin) | `docker compose version` |
| Go (solo para desarrollo sin Docker) | 1.22+ | `go version` |

---

## 2. Quick start — todo en Docker (recomendado)

El modo más rápido: **BD + API en un solo comando**.

```bash
# 1. Clonar el repositorio
git clone <url-del-repo>
cd golabs-api

# 2. Crear y configurar el .env
cp .env.example .env
```

Editar `.env` con un editor y cambiar **al menos** estas variables:

```env
JWT_SECRET=<cadena-aleatoria-larga>   # OBLIGATORIO cambiar
DB_PASSWORD=<tu-password>             # Recomendado cambiar
```

Generar un JWT_SECRET seguro:

```bash
openssl rand -hex 32
```

```bash
# 3. Levantar todo
make docker-up

# La API estará disponible en http://localhost:8080
# La BD estará disponible en localhost:3306
```

> [!IMPORTANT]
> Las migraciones de BD se aplican automáticamente al arrancar gracias al volumen `./deployments/database/init` montado en MariaDB.

Para detener:

```bash
make docker-down
```

---

## 3. Desarrollo local — API sin Docker + BD con Docker

Útil cuando quieres editar código y ver cambios en tiempo real.

### Paso 1: Levantar solo la base de datos

```bash
cd deployments/database

# Copiar y configurar el .env de la BD
cp .env.example .env
# Editar .env con tus credenciales

# Levantar MariaDB y aplicar migraciones
make up

cd ../..
```

### Paso 2: Configurar el .env de la API

```bash
cp .env.example .env
```

Editar `.env` apuntando a localhost:

```env
DB_HOST=localhost
DB_PORT=3306             # mismo que MARIADB_PORT en deployments/database/.env
DB_NAME=golabs           # mismo que MARIADB_DATABASE
DB_USER=golabs_user      # mismo que MARIADB_USER
DB_PASSWORD=tu_password  # mismo que MARIADB_PASSWORD

JWT_SECRET=tu_jwt_secret_aqui
```

### Paso 3: Correr la API

```bash
# Opción A: compilar y correr directamente
make run

# Opción B: hot-reload (requiere air instalado)
# go install github.com/air-verse/air@latest
make dev
```

La API estará en `http://localhost:8080`.

---

## 4. Producción en un servidor / VPS

### Requisitos del servidor

- Ubuntu 22.04 LTS o equivalente
- Docker + Docker Compose instalados
- Puerto 80/443 abierto (si se usa reverse proxy)
- Puerto 8080 abierto (si se expone directamente, **no recomendado**)

### Paso 1: Clonar el repo en el servidor

```bash
ssh usuario@tu-servidor
git clone <url-del-repo>
cd golabs-api
```

### Paso 2: Configurar variables de producción

```bash
cp .env.example .env
nano .env   # o vim .env
```

Variables **críticas** para producción:

```env
JWT_SECRET=<genera con: openssl rand -hex 32>
DB_PASSWORD=<password largo y único>
ALLOWED_ORIGINS=https://tu-frontend.com
JWT_EXP_MINUTES=15    # mantener corto en producción
```

### Paso 3: Levantar el stack

```bash
make docker-up
```

Verificar que los contenedores estén sanos:

```bash
docker compose ps
```

Deberías ver:

```
NAME          STATUS
golabs-api    running (healthy)
golabs-db     running (healthy)
```

### Paso 4: Configurar reverse proxy (nginx)

Ejemplo mínimo de configuración nginx para exponer la API por HTTPS:

```nginx
server {
    listen 80;
    server_name api.tu-dominio.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name api.tu-dominio.com;

    ssl_certificate     /etc/letsencrypt/live/api.tu-dominio.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.tu-dominio.com/privkey.pem;

    location / {
        proxy_pass         http://localhost:8080;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}
```

Obtener certificado SSL gratuito con Certbot:

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.tu-dominio.com
```

### Paso 5: Configurar reinicio automático

Los contenedores ya tienen `restart: unless-stopped` en el `docker-compose.yml`. Si el servidor se reinicia, los contenedores levantan solos automáticamente (siempre que el servicio Docker esté habilitado):

```bash
sudo systemctl enable docker
```

### Paso 6: Configurar backups automáticos de BD

```bash
# Backup automático cada 6 horas
cd deployments/database && make backup-auto
```

---

## 5. Variables de entorno — referencia completa

### `.env` (raíz del proyecto — API)

```env
# ── Base de datos ──────────────────────────────────────
DB_HOST=localhost              # Host de MariaDB (usar 'db' si la API corre en Docker)
DB_PORT=3306                   # Puerto de MariaDB
DB_NAME=golabs                 # Nombre de la base de datos
DB_USER=golabs_user            # Usuario de BD (no usar root)
DB_PASSWORD=CAMBIAR_ESTO       # Contraseña del usuario de BD

# Pool de conexiones (valores razonables para producción)
DB_MAX_OPEN_CONNS=25           # Máximo de conexiones abiertas simultáneas
DB_MAX_IDLE_CONNS=10           # Máximo de conexiones inactivas en el pool
DB_CONN_MAX_LIFETIME=5m        # Tiempo máximo de vida de una conexión

# ── JWT ────────────────────────────────────────────────
JWT_SECRET=GENERAR_CON_OPENSSL # Clave de firma. Mín. 32 chars. NUNCA commitear.
JWT_ISSUER=golabs-api          # Issuer del token (aparece en el campo 'iss')
JWT_EXP_MINUTES=15             # Duración del access token. 15min recomendado.

# ── CORS ───────────────────────────────────────────────
ALLOWED_ORIGINS=http://localhost:3000  # Orígenes permitidos. Separar por coma.
                                       # En prod: https://tu-frontend.com
```

### `deployments/database/.env` (solo para BD standalone)

```env
MARIADB_ROOT_PASSWORD=CAMBIAR_ESTO   # Password de root de MariaDB
MARIADB_DATABASE=golabs              # Base de datos a crear automáticamente
MARIADB_USER=golabs_user             # Usuario de aplicación (no root)
MARIADB_PASSWORD=CAMBIAR_ESTO        # Password del usuario de aplicación
MARIADB_PORT=3306                    # Puerto expuesto al host
```

> [!CAUTION]
> `DB_USER` en el `.env` de la API debe coincidir con `MARIADB_USER` del `.env` de la BD. Lo mismo para `DB_PASSWORD` y `MARIADB_PASSWORD`.

---

## 6. Verificar que todo funciona

Una vez levantado el stack, ejecutar estas verificaciones:

```bash
# 1. Liveness — el servidor responde
curl http://localhost:8080/healthz/live
# Esperado: {"status":"ok"}

# 2. Readiness — la BD está conectada
curl http://localhost:8080/healthz/ready
# Esperado: {"status":"ok"}

# 3. Registrar un usuario de prueba
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"secret123"}' | jq .
# Esperado: objeto de usuario con id, username, email, role="user"

# 4. Login
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"testuser","password":"secret123"}' | jq .
# Esperado: {"access_token":"...","refresh_token":"...","expires_in":900}
```

Si todos devuelven las respuestas esperadas, el despliegue es exitoso.

---

## 7. Solución de problemas frecuentes

### La API no inicia y el log dice "connection refused" o "dial tcp"

La BD no está lista. La API tiene retry automático, pero si el problema persiste:

```bash
# Ver logs de la BD
docker compose logs db

# Verificar que la BD responde
docker exec golabs-db mariadb-admin ping -h 127.0.0.1 -u root -p<DB_ROOT_PASSWORD>
```

---

### Error: "JWT_SECRET is required"

El archivo `.env` no está presente o la variable está vacía.

```bash
# Verificar que el .env existe y tiene JWT_SECRET
cat .env | grep JWT_SECRET
```

---

### Error 403 en todos los endpoints

El usuario puede estar baneado, o el JWT expiró. Verificar:

- Hacer login de nuevo para obtener un token fresco
- Comprobar que la fecha del servidor es correcta (afecta la validación del JWT)

---

### Las migraciones no se aplican (stack raíz)

El stack raíz monta `./deployments/database/init` directamente en MariaDB como scripts de inicialización. Solo se ejecutan la primera vez que se crea el volumen.

Si necesitas re-aplicar migraciones:

```bash
# Opción 1: borrar el volumen y recrear (DESTRUYE LOS DATOS)
docker compose down -v
make docker-up

# Opción 2: usar el stack de deployments/database con Flyway (recomendado)
cd deployments/database && make migrate
```

---

### Cambié el .env pero los cambios no se reflejan

Los contenedores deben reiniciarse para leer el nuevo `.env`:

```bash
make docker-down && make docker-up
```

---

### Ver logs en tiempo real

```bash
# Todos los servicios
docker compose logs -f

# Solo la API
docker compose logs -f api

# Solo la BD
docker compose logs -f db
```
