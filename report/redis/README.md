# Отчёт по домашнему заданию №9

Перейдём в директорию cassandra
```bash
cd redis
```

Создадим файлы с данными

JSON
```bash
JSON="{"
for ((i=1; i<=120000; i++))
do
  JSON+="'key_$i':'Test data $i','key_1_$i':'Test data 1 $i','key_2_$i':'Test data 2 $i','key_3_$i':'Test data 3 $i','key_4_$i':'Test data 4 $i',"
done
JSON+="'key_end':'Test data end'}"
echo $JSON > BIG_JSON.json
```

HSET
```bash
JSON=""
for ((i=1; i<=140000; i++))
do
  JSON+=" key_$i 'Test data $i' key_1_$i 'Test data $i 1' key_2_$i 'Test data $i 2' key_3_$i 'Test data $i 3' key_4_$i 'Test data $i 4'"
done
echo $JSON > BIG_HSET_JSON.json
```

ZSET
```bash
JSON=""
for ((i=1; i<=800000; i=i+10))
do
  TWO=$(($i+1))
  THREE=$(($i+2))
  FOUR=$(($i+3))
  FIVE=$(($i+4))
  SIX=$(($i+5))
  SEVEN=$(($i+6))
  EIGHT=$(($i+7))
  NINE=$(($i+8))
  TEN=$(($i+9))
  JSON+=" $i 'Test data $i' $TWO 'Test data $TWO' $THREE 'Test data $THREE' $FOUR 'Test data $FOUR' $FIVE 'Test data $FIVE' $SIX 'Test data $SIX' $SEVEN 'Test data $SEVEN' $EIGHT 'Test data $EIGHT' $NINE 'Test data $NINE' $TEN 'Test data $TEN'"
done
echo $JSON > BIG_ZSET_JSON.json
```

LPUSH
```bash
JSON=""
for ((i=1; i<=110000; i++))
do
  JSON+=" 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i' 'Test data $i'"
done
echo $JSON > BIG_LPUSH_JSON.json
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Зайдём на мастер ноду
```bash
docker exec -it redis-node-1 bash
```

Выставим маленькое пороговое значение для slowlog
```bash
echo "config set slowlog-log-slower-than 100" | redis-cli
```

Зальём большую строку как значение
```bash
val=$(cat /tmp/BIG_JSON.json)

echo "set big_json \"${val}\"" |  redis-cli
```

Проверим время выполнения
```redis
SLOWLOG GET
1) (integer) 478
2) (integer) 1672002532
3) (integer) 1147
4) 1) "set"
   2) "big_json"
   3) "{'key_1':'Test data 1','key_1_1':'Test data 1 1','key_2_1':'Test data 2 1','key_3_1':'Test data 3 1','key_4_1':'Test data 4 1','... (20008849 more bytes)"
5) "127.0.0.1:56048"
6) ""
```
1 мл сек

Зальём большой объкет в виде списка ключей и их значений
```bash
val=$(cat /tmp/BIG_HSET_JSON.json)

echo "hset big_hset ${val}" |  redis-cli
(integer) 700000
(0.79s)
```
Вставка заняла 790 мл сек

Зальём большой объкет в виде отсортированного множества
```bash
val=$(cat /tmp/BIG_ZSET_JSON.json)

echo "zadd big_zset ${val}" |  redis-cli
(integer) 800000
(0.95s)
```
Вставка заняла 950 мл сек. Если орентироватья на количество записей, то сапоставимо с hset.

Зальём большой список
```bash
val=$(cat /tmp/BIG_LPUSH_JSON.json)

echo "lpush big_lpush ${val}" |  redis-cli
```

Проверим время выполнения
```redis
SLOWLOG GET
1) (integer) 195
2) (integer) 1672085502
3) (integer) 94736
4)  1) "lpush"
    2) "big_lpush"
    3) "Test data 1"
    4) "Test data 1"
    5) "Test data 1"
    ...
    32) "... (1099971 more arguments)"
5) "127.0.0.1:39456"
6) ""
```
94 мл сек

Протестируем теперь скорость чтения строковых значений
```redis
get big_json
(2.27s)
```
2.27 секунды получение большой строки

Пар ключ значение
```redis
hgetall big_hset
(4.51s)
```
4.51 секунды на получение всех значений. Если же нужно получить значение по конкрутному ключу, то время обработки будет окого 1 микро секунды.

Отсортированного множества
```redis
ZRANGE big_zset 0 -1 WITHSCORES
(4.60s)
```
4.60 секунды на получение всех значений. Однако получение конкретного значения полагаю что будет существенно меньше.
Например, команда ниже выполняется 32 микросекунды
```redis
ZRANGE big_zset 1 2 WITHSCORES
```

Список
```redis
lrange big_lpush 0 -1
(3.41s)
```
3.41 секунды на получение всех значений.
Получение конкретного значения
```redis
lrange big_lpush 100 100
```
24 микросекунды


Поднимем кластреное решение redis. Конфиг редиса находится в файле redis.conf
```bash
docker-compose --file docker-compose-cluster.yml up -d
```

И подключимся к одной из нод
```bash
docker exec -it redis-1 bash
redis-cli
```

Проверим состояние кластера
```redis
cluster nodes
83cf89faafea95d92a9595e92367dd59d9da1ac5 172.25.0.4:6379@16379 master - 0 1672003317066 2 connected 5461-10922
cc847c9a95244ea0d3f12c27865bb65c20d0054e 172.25.0.6:6379@16379 slave e536ff0b783630e1ba58b16928da3022188a7ae0 0 1672003315558 3 connected
e536ff0b783630e1ba58b16928da3022188a7ae0 172.25.0.5:6379@16379 master - 0 1672003316061 3 connected 10923-16383
a2e702a334990aa95cc869c67128f4a7de310739 172.25.0.8:6379@16379 slave 83cf89faafea95d92a9595e92367dd59d9da1ac5 0 1672003315558 2 connected
77b63ab4e49c3546f2dbca9b7a62608704a267bd 172.25.0.7:6379@16379 slave 203cd172bfc4853aeb7acc91ec4f723cec174e6c 0 1672003315000 1 connected
203cd172bfc4853aeb7acc91ec4f723cec174e6c 172.25.0.3:6379@16379 myself,master - 0 1672003316000 1 connected 0-5460
```
3 мастера, 3 слейва

Попробуем записать значения
```bash
redis-cli -c
```

```redis
set key_1 1
OK
```

```redis
set key_2 2
-> Redirected to slot [7869] located at 172.25.0.4:6379
OK
```

```redis
set key_3 3
-> Redirected to slot [3740] located at 172.25.0.3:6379
OK
```

```redis
get key_1
-> Redirected to slot [11998] located at 172.25.0.5:6379
"1"
```

```redis
get key_2
-> Redirected to slot [7869] located at 172.25.0.4:6379
"2"
```

```redis
get key_3
-> Redirected to slot [3740] located at 172.25.0.3:6379
"3"
```