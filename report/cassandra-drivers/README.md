# Отчёт по домашнему заданию №8

## Подготовка к сравнению драйверов

Сравнивать будет приложения на go и на php
Соберём исходник go
```bash
make
```

Соберём контейнер для него
```bash
docker build -f Dockerfile.cassandra -t cassandra-drivers .
```

Перейдём в директорию cassandra
```bash
cd cassandra
```

Выкачаем репозиторй внешней библиотеки php
```bash
cd php
git clone https://github.com/duoshuo/php-cassandra
cd ..
```

## Снятие снапшотов и восстановление

Тк в предыдущей работу использовал docker-compose, то для снятию снапшота буду использовать nodetool и делать это буду вручную.

Поднимем контейнеры
```bash
docker-compose up -d
```

Зайдём на ноду
```bash
docker exec -it cassandra-1 bash
```

Запустим клиент
```bash
cqlsh
```

Создадим базу данных
```cassandraql
CREATE KEYSPACE IF NOT EXISTS test_keyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
```

Выберем базу
```cassandraql
USE test_keyspace;
```

Создадим простую таблицу user
```cassandraql
CREATE TABLE user (
   id uuid,
   login text,
   created_at timestamp,
   PRIMARY KEY (id)
);
```

И таблицу user_statistic
```cassandraql
CREATE TABLE user_statistic (
   id uuid,
   data text,
   created_at timestamp,
   PRIMARY KEY ((id, data), created_at)
);
```

Выйдем из клиента
```bash
exit
```

Заполним таблицу
```bash
for ((i=1; i<=1000; i++))
do
  UUID=$(cat /proc/sys/kernel/random/uuid)
  cqlsh -e "INSERT INTO test_keyspace.user (id, login, created_at) VALUES ($UUID, 'test-$i', dateof(now()));";
  for ((j=1; j<=10; j++))
  do
    let "DATA = $j / 2"
    cqlsh -e "INSERT INTO test_keyspace.user_statistic (id, data, created_at) VALUES ($UUID, 'data-$DATA', dateof(now()));";
  done
done
```

Проверим количество записей
```cassandraql
SELECT COUNT(*) FROM user;

1000

SELECT COUNT(*) FROM user_statistic;

10000
```

Сбросим данные на диск
```bash
nodetool flush test_keyspace
```

Создадим снапшот данных на каждом узле
```bash
nodetool snapshot --tag test-ks-snapshot test_keyspace
```

Проделаем данную процедуру на всех узлах

Найдём файлы созданных снапшотов
```bash
find -name test-ks-snapshot

./var/lib/cassandra/data/test_keyspace/user_statistic-e4a255307c8211eda6412929f36a54c7/snapshots/test-ks-snapshot
./var/lib/cassandra/data/test_keyspace/user-893036507c8b11ed9a8811ea5add240e/snapshots/test-ks-snapshot
```

Очистим таблицы
```cassandraql
truncate user;
truncate user_statistic;
```

Создадим папки для восстановления таблиц
```bash
mkdir /var/lib/cassandra/data/test_keyspace/user
mkdir /var/lib/cassandra/data/test_keyspace/user_statistic
```

И скопируем туда содержимое снапшотов
```bash
cd /var/lib/cassandra/data/test_keyspace/user-893036507c8b11ed9a8811ea5add240e/snapshots/test-ks-snapshot
cp * /var/lib/cassandra/data/test_keyspace/user
cd /var/lib/cassandra/data/test_keyspace/user
```

Выполним восстановление таблицы user
```bash
sstableloader -d 172.22.0.2,172.22.0.3,172.22.0.4 /var/lib/cassandra/data/test_keyspace/user
```

Аналогично для таблицы user_statistic

```bash
cd /var/lib/cassandra/data/test_keyspace/user_statistic-e4a255307c8211eda6412929f36a54c7/snapshots/test-ks-snapshot
cp * /var/lib/cassandra/data/test_keyspace/user_statistic
cd /var/lib/cassandra/data/test_keyspace/user_statistic
```

```bash
sstableloader -d 172.22.0.2,172.22.0.3,172.22.0.4 /var/lib/cassandra/data/test_keyspace/user_statistic
```

Проверим количество записей после восстановления
```cassandraql
SELECT COUNT(*) FROM user;

1000

SELECT COUNT(*) FROM user_statistic;

10000
```

## Сравнение драйверов

Зайдём на ноду php
```bash
docker exec -it php bash
cd /var/code
```

Включим модуль socket
```bash
docker-php-ext-install sockets
```

Запустим php
```bash
php cassandra.php

exitstring(31) "time: 0.024019002914429 seconds"
string(27) "memory: 0.82025146484375 Mb"
```
Время выполнения 24 млсек
Память 0.8 Mb

Зайдём на ноду go
```bash
docker exec -it go bash
cd /var/code
```

Запустим go
```bash
./bin/cassandra

Time: 32
Memory: 0.361
```
Время выполнения 32 млсек
Память 0.36 Mb

Время примерно такое же, чуть больше, памяти в 2 раза меньше потребляет решение на go
