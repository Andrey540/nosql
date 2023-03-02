# Отчёт по домашнему заданию №14

Создадим кластер mongodb в yandex cloud самой простой конфигурации

Загрузим сертификат

```bash
mkdir --parents ~/.mongodb && \
wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
    --output-document ~/.mongodb/root.crt && \
chmod 0644 ~/.mongodb/root.crt
```

Соберём приложение на go
```bash
make
```

Перейдём в директорию bin и запутим приложение для mongo-cloud
Приложение записыват в 10 потоков по 1000 записей в каждом
А также аналогично читает 10 потоков по 1000 записей в каждом
```bash
cd bin
./mongo-cloud
Write time: 32063
Read time: 38811
```

Время записи составило 32 секунды, а чтения 38 секунд.
Обработали 10000 запросов в 10 потоков

Дальше протестируем другие базы, но уже локально с помощью докера.

Перейдём в директорию cassandra
```bash
cd ..
cd cassandra
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Перейдём в директорию bin и запутим приложение для cassandra
При этом количество записей увеличим в 100 раз
Приложение записыват в 10 потоков по 100000 записей в каждом
А также аналогично читает 10 потоков по 100000 записей в каждом
```bash
cd ..
cd bin
./cassandra
Write time: 59589
Read time: 99087
```

Время записи составило 60 секунд, а чтения 100 секунд.
Обработали 1 млн запросов в 10 потоков

Остановим контейнеры
```bash
cd ..
cd cassandra
docker-compose down
```

Перейдём в директорию scylla
```bash
cd ..
cd scylla
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Перейдём в директорию bin и запутим приложение для scylla
```bash
cd ..
cd bin
./scylla
Write time: 114775
Read time: 106714
```

Запись оконло двух минуут - 114 секунд
Чтение 106 секунд
Время выше, чем у cassandra, что странно, тк scylla вроде по тестам должна быть быстрее cassandra, возможно такой результат из-за дефолтных настроек

Остановим контейнеры
```bash
cd ..
cd scylla
docker-compose down
```

Перейдём в директорию mongo
```bash
cd ..
cd mongo
```

Поднимем контейнеры
```bash
docker-compose --file docker-compose-mongo-simple-cluster.yml up -d
```

Инициализируем реплики
```bash
docker-compose --file docker-compose-mongo-simple-cluster.yml exec mongodb_1 sh -c "mongosh < /scripts/init-shards.js"
```

Создадим коллекцию и индкс по ней
```bash
use testing
db.users.createIndex({ login: 1 })
```

Перейдём в директорию bin и запутим приложение для mongo
```bash
cd ..
cd bin
./mongo
Write time: 476494
Read time: 42475
```

Запись порядком больше времени занимает 476 секунд, а вот чтение раза в 2 меньше - 42 секунды

Остановим контейнеры
```bash
cd ..
cd mongo
docker-compose --file docker-compose-mongo-simple-cluster.yml down
```

Перейдём в директорию mongo
```bash
cd ..
cd couchbase
```

Поднимем контейнеры
```bash
docker-compose --file docker-compose-couchbase.yml up -d
```

Перейдём в директорию bin и запутим приложение для couchbase
```bash
cd ..
cd bin
./couchbase
Write time: 49995
Read time: 8211
```

Время записи составило 50 секунд, а время чтения 8 секунд

На всех базах записывалось 1 млн записей в 10 потоков, а также читалось 1 млн записей также в 10 потоков
Результаты по всем базам приведены ниже в мл секундах

|            | Cassandra | Scylla  | Mongo   | Couchbase |
|------------|-----------|---------|---------|-----------|
| Write (ms) | 59589     | 114775  | 476494  | 49995     |
| Read (ms)  | 99087     | 106714  | 42475   | 8211      |

Самой быстрой по записи и по чтению оказался couchbase.
Среди остальных самой быстрой за запись оказалась cassandra, scylla аказалась в 2 раза медленнее, возможно из-за дефолтных настроек, ожидалось, что она будет раза в 2 быстрее
Самой быстрой по чтению среди остальных стала mongo.