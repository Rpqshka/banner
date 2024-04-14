#!/bin/sh
# wait-for-postgres.sh

set -e
host="$1"
dbname="$2"

shift 2
cmd="$@"

until PGPASSWORD="admin" psql -h "$host" -U "postgres" -d "$dbname" -c '\q'; do
  >&2 echo "Postgres loading - sleeping"
  sleep 5
done

>&2 echo "Postgres is up - run banner service"

exec $cmd