#!/bin/bash

DB_HOST="db"
DB_NAME="$POSTGRES_DB"
DB_USER="$POSTGRES_USER"
DB_PASSWORD="$POSTGRES_PASSWORD"
BACKUP_NAME="backup_$(date +'%Y-%m-%d_%H-%M-%S').sql.gz"
LOCAL_BACKUP_PATH="/backups/$BACKUP_NAME"

PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | gzip > "$LOCAL_BACKUP_PATH"