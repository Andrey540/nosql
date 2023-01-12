# Отчёт по домашнему заданию №10

Перейдём в директорию etcd
```bash
cd etcd
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Зайдём в контейнер
```bash
docker exec -it etcd-1 bash
```

Запишем данные
```bash
etcdctl put foo bar
OK
```

Получим данные
```bash
etcdctl get foo 
foo
bar
```

Проверим состояние кластера
```bash
etcdctl --write-out=table --endpoints=localhost:2379 member list
+------------------+---------+--------+--------------------+--------------------+------------+
|        ID        | STATUS  |  NAME  |     PEER ADDRS     |    CLIENT ADDRS    | IS LEARNER |
+------------------+---------+--------+--------------------+--------------------+------------+
| 88d11e2649dad027 | started | etcd-2 | http://etcd-2:2380 | http://etcd-2:2379 |      false |
| b8c6addf901e4e46 | started | etcd-1 | http://etcd-1:2380 | http://etcd-1:2379 |      false |
| c3697a4fd7a20dcd | started | etcd-3 | http://etcd-3:2380 | http://etcd-3:2379 |      false |
+------------------+---------+--------+--------------------+--------------------+------------+
```

Остановим третюю ноду
```bash
docker stop etcd-3
```

Попробуем опять зписать данные
```bash
etcdctl put foo bar1
OK
```

Получим данные
```bash
etcdctl get foo 
foo
bar1
```

Кластер продолжает работать.
Остановим теперь ещё и вторую ноду
```bash
docker stop etcd-2
```

Попробуем теперь зписать данные
```bash
etcdctl put foo bar2
{"level":"warn","ts":"2023-01-08T17:09:55.326Z","logger":"etcd-client","caller":"v3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000146a80/127.0.0.1:2379","attempt":0,"error":"rpc error: code = Unknown desc = context deadline exceeded"}
Error: rpc error: code = Unknown desc = context deadline exceeded
```

Получили ошибку, тк нет достаточно участников для консенсуса
Попробуем прочитать данные
```bash
etcdctl get foo 
{"level":"warn","ts":"2023-01-08T17:11:10.771Z","logger":"etcd-client","caller":"v3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc0003d2a80/127.0.0.1:2379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = context deadline exceeded"}
Error: context deadline exceeded
```

Аналогично не работат чтение
Поднимем остановленные ноды
```bash
docker-compose up -d
```

Попробуем прочитать данные
```bash
etcdctl get foo 
foo
bar1
```

И записать
```bash
etcdctl put foo bar2
OK
```

Проверим что записалось
```bash
etcdctl get foo 
foo
bar2
```

Запустим consul
```bash
cd ../consul
docker-compose up -d
```

Зайдём на сервер
```bash
docker exec -it consul1 /bin/sh
```

Проверим кластер
```bash
consul members
Node     Address          Status  Type    Build   Protocol  DC   Partition  Segment
consul1  172.20.0.2:8301  alive   server  1.12.3  2         dc1  default    <all>
consul2  172.20.0.4:8301  alive   server  1.12.3  2         dc1  default    <all>
consul3  172.20.0.3:8301  alive   server  1.12.3  2         dc1  default    <all>
```

Сохраним ключ - значений
```bash
consul kv put otus consul
Success! Data written to: otus
```

Прочитаем значение
```bash
consul kv get otus
consul
```

Остановим одну ноду
```bash
docker stop consul3
```

Проверим кластер
```bash
consul members
Node     Address          Status  Type    Build   Protocol  DC   Partition  Segment
consul1  172.20.0.2:8301  alive   server  1.12.3  2         dc1  default    <all>
consul2  172.20.0.4:8301  alive   server  1.12.3  2         dc1  default    <all>
consul3  172.20.0.3:8301  failed  server  1.12.3  2         dc1  default    <all>
```
Одна нода не доступна

Прочитаем значение
```bash
consul kv get otus
consul
```

Остановим ещё одну ноду
```bash
docker stop consul2
```

Проверим кластер
```bash
consul members
Node     Address          Status  Type    Build   Protocol  DC   Partition  Segment
consul1  172.20.0.2:8301  alive   server  1.12.3  2         dc1  default    <all>
consul2  172.20.0.4:8301  failed  server  1.12.3  2         dc1  default    <all>
consul3  172.20.0.3:8301  left    server  1.12.3  2         dc1  default    <all>
```
Две ноды не доступны

Прочитаем значение
```bash
consul kv get otus
Error querying Consul agent: Unexpected response code: 500 (No cluster leader)
```

Кластер отвечает ошибкой, тк нет косенсуса
Аналогично ошибка будет выдаваться на запись

Восстановим кластер
```bash
docker-compose up -d
```

Проверим кластер
```bash
consul members
Node     Address          Status  Type    Build   Protocol  DC   Partition  Segment
consul1  172.20.0.2:8301  alive   server  1.12.3  2         dc1  default    <all>
consul2  172.20.0.3:8301  alive   server  1.12.3  2         dc1  default    <all>
consul3  172.20.0.4:8301  alive   server  1.12.3  2         dc1  default    <all>
```
Все ноды активны

Прочитаем значение
```bash
consul kv get otus
consul
```