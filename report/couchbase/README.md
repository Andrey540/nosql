# Отчёт по домашнему заданию №4

## Настройк кластера

Перейдём в директорию couchbase
```bash
cd couchbase
```

Поднимем контейнеры
```bash
docker-compose --file docker-compose-couchbase.yml up -d
```

Инициализируем кластер по адресу http://localhost:8091
Добавим ноды couchbase-1, couchbase-2 и couchbase-3 в кластер с ip 172.18.0.3, 172.18.0.2, 172.18.0.4 соответственно
Первые три ноды сделаем для хранения данных, червёртую для всего остального
После чего перебалансируем кластер, чтобы трафик мог идти на все ноды

Создадим бакет из примеров - travel-sample

Выполним запрос на вставку
```bash
INSERT INTO `travel-sample`.inventory.hotel (KEY, VALUE)
VALUES ("key1", { "type" : "hotel", "name" : "new hotel" });
{
  "results": []
}
```

Выполним запрос на вставку ещё раз
```bash
INSERT INTO `travel-sample`.inventory.hotel (KEY, VALUE)
VALUES ("key1", { "type" : "hotel", "name" : "new hotel" });
[
  {
    "code": 12009,
    "msg": "DML Error, possible causes include concurrent modification. Failed to perform INSERT on key key1 - cause: Duplicate Key: key1",
    "reason": {
      "caller": "couchbase:1985",
      "code": 17012,
      "key": "dml.statement.duplicatekey",
      "message": "Duplicate Key: key1"
    },
    "retry": false
  }
]
```

Выполним запрос на впоиск созданной записи
```bash
SELECT * FROM `travel-sample`.inventory.hotel where name = "new hotel";
[
  {
    "hotel": {
      "name": "new hotel",
      "type": "hotel"
    }
  }
]
```

Узнаем сколько всего отелей
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 919
  }
]
```

Найдём все отели укоторых в названии есть слово hotel
```bash
SELECT name FROM `travel-sample`.inventory.hotel where name like "%hotel%";
[
  {
    "name": "The Falcondale hotel and restaurant"
  },
  {
    "name": "myhotel Chelsea"
  },
  {
    "name": "pentahotel Birmingham"
  },
  {
    "name": "new hotel"
  }
]
```

Остановим первую ноду и посмотрим что получилось
```bash
docker stop couchbase-1
```

Первая нода стала недоступной

Узнаем сколько всего отелей
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 0
  }
]
```

Кластер стал выдавать ошибку
Однако через какое-то время он автматически выполнил failover и работа восстановилась
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 919
  }
]
```
Вставка также заработала

Остановим вторую ноду и посмотрим что получилось
```bash
docker stop couchbase-2
```

Вторая нода стала недоступной
Узнаем сколько всего отелей
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 0
  }
]
```
Кластер опять стал выдавать ошибку
Ручной failover лучше не выполнять, пишут что это может привести к потере данных

Если выполнять ручной failover, то кластер станет отдавать данные, но не все
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 598
  }
]
```

Автоматический failover же при двух упавших нодах с данными из трёх не наступает

Попробуем поднять обратно ноды
```bash
docker-compose --file docker-compose-couchbase.yml up -d
```

Узнаем сколько сейчас отелей
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 919
  }
]
```
Данные восстановились
Однако вторая нода появилась, но кластер трафик на неё не направляет из-за автоматического failover-а
Сделаем частичное восстановление ноды и перебалансировку

После этого нода снова в строю


Остановим четвёртую ноду, которая обрабатывает запросы, и посмотрим что получилось
```bash
docker stop couchbase-4
```

```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
{"status": "Failure contacting server."}
```

Операция не доступна

Поднимем обратно ноду
```bash
docker-compose --file docker-compose-couchbase.yml up -d
```

Заросы снова заработали
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 919
  }
]
```

Попробуем поднять пятую ноду тоже на запросы с ip 172.18.0.6
Добавим её в кластер и сделаем перебалансировку
После чего повторно остановим червёртую ноду
```bash
docker stop couchbase-4
```

Ситуация повторилась
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
{"status": "Failure contacting server."}
```

Попробуем выполнить failover этой ноды
После чего запросы стали отрабатывать
```bash
SELECT count(*) FROM `travel-sample`.inventory.hotel;
[
  {
    "$1": 919
  }
]
```

Поднимем обратно ноду
```bash
docker-compose --file docker-compose-couchbase.yml up -d
```
Вернём её в кластер и перебалансируем.
Кластер восстановился