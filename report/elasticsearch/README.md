# Отчёт по домашнему заданию №7

## Настройка кластера

Перейдём в директорию elasticsearch
```bash
cd elasticsearch
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Создадим данные
```bash
curl -X PUT "localhost:9200/mama/_bulk?pretty" -H 'Content-Type: application/json' -d'
{ "create": { } }
{ "say": "моя мама мыла посуду а кот жевал сосиски"}
{ "create": { } }
{ "say": "рама была отмыта и вылизана котом"}
{ "create": { } }
{ "say": "мама мыла раму"}
'
```

Выполним запрос на поиск
```bash
curl -X GET "localhost:9200/mama/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query" : {
    "match" : { "say": "мама ела сосиски" }
  }
}
'
{
  "took" : 18,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 2,
      "relation" : "eq"
    },
    "max_score" : 1.241674,
    "hits" : [
      {
        "_index" : "mama",
        "_id" : "71iH1IQB5YitY4ALaZW4",
        "_score" : 1.241674,
        "_source" : {
          "say" : "моя мама мыла посуду а кот жевал сосиски"
        }
      },
      {
        "_index" : "mama",
        "_id" : "8ViH1IQB5YitY4ALaZW5",
        "_score" : 0.5820575,
        "_source" : {
          "say" : "мама мыла раму"
        }
      }
    ]
  }
}
```
Посомтрим результат.
Найдено 2 совпадения: "моя мама мыла посуду а кот жевал сосиски" и "моя мама мыла посуду а кот жевал сосиски"

Теперь выполним запрос нечёткого поиска и посмотрим что получится
```bash
curl -X GET "localhost:9200/mama/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query" : {
    "match" : { 
      "say" : {
        "query": "мама ела сосиски",
        "fuzziness": "auto"
      }
    }
  }
}
'
{
  "took" : 7,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 3,
      "relation" : "eq"
    },
    "max_score" : 1.241674,
    "hits" : [
      {
        "_index" : "mama",
        "_id" : "NfGiDIUBZcQ-g7fz4z7w",
        "_score" : 1.241674,
        "_source" : {
          "say" : "моя мама мыла посуду а кот жевал сосиски"
        }
      },
      {
        "_index" : "mama",
        "_id" : "N_GiDIUBZcQ-g7fz4z7z",
        "_score" : 0.5820575,
        "_source" : {
          "say" : "мама мыла раму"
        }
      },
      {
        "_index" : "mama",
        "_id" : "NvGiDIUBZcQ-g7fz4z7z",
        "_score" : 0.34421936,
        "_source" : {
          "say" : "рама была отмыта и вылизана котом"
        }
      }
    ]
  }
}
```
Найдено 3 совпадения: "моя мама мыла посуду а кот жевал сосиски", "моя мама мыла посуду а кот жевал сосиски" и "рама была отмыта и вылизана котом"
Третяя строка попала в результат так как слово рама отличается от мама на одну букву