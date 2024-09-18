#!/bin/bash

docker run --name postgres-test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d -p 6432:5432 -d postgres:alpine
echo "Waiting for postgres to start"
sleep 3

docker exec -it postgres-test psql -U postgres  -d postgres -c "CREATE DATABASE productapp;"
sleep 3 
echo "Database created successfully"

docker exec -it postgres-test psql -U postgres -d productapp -c "
CREATE TABLE if not exists products 
(
    id BIGSERIAL not null PRIMARY KEY, 
    name VARCHAR(255) not null, 
    price double precision not null,
    discount double precision,
    store VARCHAR(255) not null
);
"
sleep 3 
echo "Table created successfully"