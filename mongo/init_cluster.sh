#!/usr/bin/env bash

# Инициализируем конфигсервер
docker-compose --file docker-compose-mongo-cluster.yml exec configsvr01 sh -c "mongosh < /scripts/init-configserver.js"

# Инициализируем шарды
docker-compose --file docker-compose-mongo-cluster.yml exec shard01-a sh -c "mongosh < /scripts/init-shard01.js"
docker-compose --file docker-compose-mongo-cluster.yml exec shard02-a sh -c "mongosh < /scripts/init-shard02.js"
docker-compose --file docker-compose-mongo-cluster.yml exec shard03-a sh -c "mongosh < /scripts/init-shard03.js"

# Инициализируем mongos
docker-compose --file docker-compose-mongo-cluster.yml exec router01 sh -c "mongosh < /scripts/init-router.js"
docker-compose --file docker-compose-mongo-cluster.yml exec router02 sh -c "mongosh < /scripts/init-router.js"