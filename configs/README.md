# configs/

Configuración YAML del servidor **golabs-api**.

---

## Archivos

| Archivo | Propósito |
|---|---|
| `config.yaml` | Configuración base cargada en todos los entornos |
| `config.local.yaml` | Sobreescrituras para desarrollo local (ignorado por git) |

---

## Estructura

```yaml
server:
  port: 8080          # Puerto HTTP del servidor

app:
  name: golabs-api    # Nombre de la aplicación (aparece en logs)
  env: development    # Entorno: development | production
```

---

## Uso

La API carga `configs/config.yaml` automáticamente al iniciar. Si existe `configs/config.local.yaml`, sus valores sobreescriben los del archivo base.

Para desarrollo local, crear el archivo de override:

```bash
cp configs/config.yaml configs/config.local.yaml
# Editar config.local.yaml con los valores deseados
```

> [!NOTE]
> Las credenciales de base de datos y el JWT secret se configuran exclusivamente mediante variables de entorno (`.env`), nunca en estos archivos YAML. Esto evita commitear secretos accidentalmente.
