#! /bin/sh
# ./execute-sql .env "SELECT * FROM properties;"
if [[ -z "$1" || ( "$1" != "test" && "$1" != "dev" ) ]]; then
	echo "Env is required: (e.g. ./execute-sql test)"
	exit 0
fi

ENVFILE=".env.test"
DB_CONTAINER_NAME="properties-seeker-test-db-1"
if [ $1 = "dev" ]; then
	ENVFILE=".env"
	DB_CONTAINER_NAME="properties-seeker-dev-db-1"
fi

source $ENVFILE
docker exec -i ${DB_CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NAME} -c "\d+"
docker exec -i ${DB_CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NAME} -c "$2"

# SELECT conrelid::regclass AS table_name,
#        conname AS foreign_key,
#        pg_get_constraintdef(oid)
# FROM   pg_constraint
# WHERE  contype = 'f'
# AND    connamespace = 'public'::regnamespace
# ORDER  BY conrelid::regclass::text, contype DESC;
