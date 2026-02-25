# deployments/

Infraestructura de despliegue de **golabs-api**.

---

## Estructura

```
deployments/
└── database/    # MariaDB + Flyway (migraciones y operaciones de BD)
```

Para instrucciones detalladas de la base de datos, ver [database/README.md](database/README.md).

---

## Despliegue completo (API + BD)

Para levantar todo el stack desde la raíz del proyecto (API + BD juntos), usar el `docker-compose.yml` de la raíz:

```bash
# Desde la raíz del proyecto
make docker-up
```

Para levantar solo la base de datos de forma independiente (útil cuando la API corre localmente):

```bash
cd deployments/database
make up
```
