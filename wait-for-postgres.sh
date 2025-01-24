#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
port="${2:-5432}"
shift 2
cmd="$@"

if [ -z "$DB_PASSWORD" ]; then
  echo "Error: DB_PASSWORD environment variable is not set!"
  exit 1
fi

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -p "$port" -U "postgres" -d postgres -c '\q' >/dev/null 2>&1; do
  >&2 echo "Postgres $host:$port is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd