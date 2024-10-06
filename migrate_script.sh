#!/sh
echo "Запуск миграции"

DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${POSTGRES_DB:-"postgres_db"}
DB_USER=${POSTGRES_USER:-"username"}
DB_PASSWORD=${POSTGRES_PASSWORD:-"password"}

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << EOF
CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_done BOOLEAN DEFAULT FALSE
);
EOF

echo "Миграция завершена"
