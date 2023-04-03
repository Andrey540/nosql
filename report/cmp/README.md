# Отчёт по проектной работе

Создадим кластер mongodb в yandex cloud самой простой конфигурации

Загрузим сертификат

```bash
./tarantool 
Write time: 110253
Read time: 24452
```

```bash
./redis 
Write time: 205381
Read time: 123207
```

```bash
./couchbase 
Write time: 494009
Read time: 29673
```

```bash
./couchbase-single
Write time: 221072
Read time: 40784
```

создадим индекс в couchbase по полю first_name
```bash
CREATE INDEX idx_users_first_name ON `travel-sample`.tenant_agent_01.users(first_name);
```

создадим индекс в tarantool по полю first_name
```bash
space:create_index('first_name_idx', {
type = 'tree',
unique = false,
parts = {'first_name'},
if_not_exists = true
})
```

Соберём приложения на go для тестирования нагрузки
```bash
make
```

redis entity
https://overload.yandex.net/589330
line(5, 100000, 2m) const(100000,2m)
RAM 7.5 GB
CPU 80%
rps 30000
latency 98 - 52ms; 95 - 46ms; 90 - 42ms
avg time - 30ms

tarantool entity
https://overload.yandex.net/589331
line(5, 100000, 2m) const(100000,2m)
RAM 6 GB
CPU 80%
rps 45000
latency 98 - 40ms; 95 - 35ms; 90 - 32ms
avg time - 21ms

couchbase entity
https://overload.yandex.net/589332
line(5, 100000, 2m) const(100000,2m)
RAM 12 GB
CPU 80%
rps 43000
latency 98 - 47ms; 95 - 41ms; 90 - 36ms
avg time - 24ms

tarantool name
https://overload.yandex.net/589335
line(5, 500, 2m) const(500,2m)
RAM 9 GB
CPU 70% волнами
rps 220
latency 98 - 5303ms; 95 - 5185ms; 90 - 5117ms
avg time - 4912ms

couchbase name
https://overload.yandex.net/589338
line(5, 500, 2m) const(500,2m)
RAM 13.5 GB
CPU 80% волнами
rps 60
latency 98 - 11000ms; 95 - 11000ms; 90 - 11000ms
avg time - 6600ms
При 180 rps стал 500-ками сыпать


couchbase-single name
https://overload.yandex.net/591009
line(5, 500, 2m) const(500,2m)
RAM 9.3 GB
CPU 90% волнами
rps 110
latency 98 - 7214ms; 95 - 7109ms; 90 - 7048ms
avg time - 4439ms
При 180 rps стал 500-ками сыпать


couchbase-single entity
https://overload.yandex.net/591010
line(5, 100000, 2m) const(100000,2m)
RAM 7 GB
CPU 90% волнами
rps 42000
latency 98 - 42ms; 95 - 37ms; 90 - 33ms
avg time - 24ms