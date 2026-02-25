# Guía de contribución — golabs-api

## Configuración del entorno

```bash
# 1. Clonar el repo
git clone <repo-url> && cd golabs-api

# 2. Instalar herramientas de desarrollo
go install github.com/air-verse/air@latest          # hot reload
go install honnef.co/go/tools/cmd/staticcheck@latest # análisis estático

# 3. Levantar la BD
cd deployments/database && make up && cd ../..

# 4. Configurar env
cp .env.example .env    # completar JWT_SECRET, credenciales de BD, etc.

# 5. Correr la API con hot reload
make dev
```

---

## Convenciones de código

### Idioma de comentarios

| Elemento | Idioma |
|---|---|
| Package-doc y godoc de tipos/funciones públicas | **Español** |
| Comentarios inline (explicar el *por qué*) | **Español** |
| Mensajes de `slog` (logs) | **Inglés** |
| Strings de error retornados al cliente | **Inglés** |
| Términos técnicos sin traducción natural (UUID, hash, token...) | **Inglés** |

### Estilo general

- Seguir `gofmt` y `go vet` — ejecutar `make lint` antes de cada commit
- Preferir errores centinela de `apperrors` sobre errores genéricos de `errors.New`
- Envolver errores con contexto: `fmt.Errorf("%w: contexto adicional", apperrors.ErrNotFound)`
- No usar `panic` fuera del bootstrap del servidor
- Todos los timestamps deben ser UTC: `time.Now().UTC()`

---

## Estructura de un módulo

Cada módulo de negocio sigue esta estructura de cuatro capas. Al agregar un módulo nuevo, respetar exactamente esta organización:

```
internal/<modulo>/
├── domain/
│   ├── <entidad>.go        # struct de la entidad + constructor + reglas de negocio
│   └── repository.go       # interfaz Repository (puerto de salida)
├── application/
│   └── <caso_de_uso>.go    # un archivo por caso de uso
├── interfaces/
│   ├── dto.go              # structs de request/response con json tags y validate tags
│   ├── handler.go          # métodos HTTP (un método por endpoint)
│   └── routes.go           # RegisterRoutes() — instancia repositorios, use cases y handler
└── infraestructure/        # (sí, así se llama el paquete en este proyecto)
    └── mysql_repository.go # implementación de la interfaz Repository
```

---

## Cómo agregar un módulo nuevo

### 1. Definir el dominio

```go
// internal/mymodulo/domain/myentity.go
package domain

type MyEntity struct {
    ID        uuid.UUID
    // ... campos
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewMyEntity(...) (*MyEntity, error) {
    // validaciones de negocio
}
```

```go
// internal/mymodulo/domain/repository.go
package domain

type Repository interface {
    Save(e *MyEntity) error
    GetByID(id uuid.UUID) (*MyEntity, error)
    // ... otros métodos
}
```

### 2. Implementar el repositorio MySQL

```go
// internal/mymodulo/infraestructure/mysql_repository.go
package infrastructure

type MySQLMyModuloRepository struct { db *sql.DB }

func NewMyModuloRepository(db *sql.DB) domain.Repository {
    return &MySQLMyModuloRepository{db: db}
}

// Implementar todos los métodos de domain.Repository
```

### 3. Crear los casos de uso

Cada caso de uso en un archivo separado:

```go
// internal/mymodulo/application/create_entity.go
package application

type CreateEntityUseCase struct { repo domain.Repository }

func NewCreateEntityUseCase(repo domain.Repository) *CreateEntityUseCase { ... }

func (uc *CreateEntityUseCase) Execute(...) (*domain.MyEntity, error) {
    // lógica de negocio
}
```

### 4. Crear DTOs, handler y rutas

```go
// internal/mymodulo/interfaces/dto.go
// internal/mymodulo/interfaces/handler.go
// internal/mymodulo/interfaces/routes.go — implementar RegisterRoutes()
```

### 5. Registrar las rutas en `main.go`

```go
// cmd/api/main.go
mymodulointerfaces.RegisterRoutes(router, db, jwtSvc)
```

### 6. Añadir la migración SQL

```sql
-- deployments/database/init/V9__create_mymodulo_table.sql
CREATE TABLE mymodulo (...);
```

```bash
cd deployments/database && make migrate
```

---

## Agregar un endpoint a un módulo existente

1. Agregar el método al handler en `interfaces/handler.go`
2. Registrar la ruta en `interfaces/routes.go`
3. Si hay nuevo input/output, agregar el DTO en `interfaces/dto.go`
4. Si se necesita acceder a BD, agregar el método a la interfaz `domain/repository.go` y a `infraestructure/mysql_repository.go`
5. Crear el caso de uso en `application/`

---

## Manejo de errores

```go
// En un use case:
if err != nil {
    return nil, fmt.Errorf("%w: contexto del error", apperrors.ErrNotFound)
}

// En un handler:
if err != nil {
    apperrors.RespondError(w, err)  // mapea automáticamente al código HTTP correcto
    return
}
```

Mapa de errores centinela → HTTP:

| Error centinela | HTTP |
|---|---|
| `apperrors.ErrNotFound` | 404 |
| `apperrors.ErrConflict` | 409 |
| `apperrors.ErrUnauthorized` | 401 |
| `apperrors.ErrForbidden` | 403 |
| `apperrors.ErrBadRequest` (default) | 400 |

---

## Comandos útiles

```bash
make build     # compilar
make run       # compilar y correr
make dev       # hot reload (requiere air)
make test      # tests
make lint      # go vet + staticcheck
make tidy      # go mod tidy + verify
```

---

## Antes de abrir un PR

- [ ] `make lint` pasa sin errores
- [ ] `go build ./cmd/... ./internal/...` compila sin errores
- [ ] `make test` pasa (si existen tests relevantes)
- [ ] Nuevas funciones públicas tienen godoc en español
- [ ] Mensajes de error al cliente en inglés
- [ ] Logs de `slog` en inglés
- [ ] No se commitean archivos `.env` ni `config.local.yaml`
