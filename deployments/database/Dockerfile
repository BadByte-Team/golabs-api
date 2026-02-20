FROM mariadb:11.3

HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
    CMD mariadb-admin ping -h 127.0.0.1 -u root -p"$MARIADB_ROOT_PASSWORD" || exit 1