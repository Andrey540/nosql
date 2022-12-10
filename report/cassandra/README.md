# Отчёт по домашнему заданию №6

## Настройка кластера

Перейдём в директорию cassandra
```bash
cd cassandra
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Каждая нода ограничена в ресурсах памяти, по 256 Mb на ноду

Проверим статусы нод
```bash
docker exec -it cassandra-1 nodetool status
docker exec -it cassandra-2 nodetool status
docker exec -it cassandra-3 nodetool status

Datacenter: datacenter1
=======================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens  Owns (effective)  Host ID                               Rack 
UN  172.22.0.2  137.76 KiB  16      66.3%             309f2071-2359-4410-84a3-dd13d4acea26  rack1
UN  172.22.0.3  101.62 KiB  16      63.8%             aa33fcff-9fb8-4632-885a-023ad3000d03  rack1
UN  172.22.0.4  142.81 KiB  16      70.0%             dde1a67a-9add-4d8f-98e6-2cd61fb33911  rack1
```

Зайдмн на ноду
```bash
docker exec -it cassandra-1 bash
```

Запустим нагрузачный тест на запись
```bash
/opt/cassandra/tools/bin/cassandra-stress write n=1000000

Results:
Op rate                   :   31,586 op/s  [WRITE: 31,586 op/s]
Partition rate            :   31,586 pk/s  [WRITE: 31,586 pk/s]
Row rate                  :   31,586 row/s [WRITE: 31,586 row/s]
Latency mean              :    6.2 ms [WRITE: 6.2 ms]
Latency median            :    2.9 ms [WRITE: 2.9 ms]
Latency 95th percentile   :   24.3 ms [WRITE: 24.3 ms]
Latency 99th percentile   :   54.2 ms [WRITE: 54.2 ms]
Latency 99.9th percentile :   97.9 ms [WRITE: 97.9 ms]
Latency max               :  357.3 ms [WRITE: 357.3 ms]
Total partitions          :  1,000,000 [WRITE: 1,000,000]
Total errors              :          0 [WRITE: 0]
Total GC count            : 0
Total GC memory           : 0.000 KiB
Total GC time             :    0.0 seconds
Avg GC time               :    NaN ms
StdDev GC time            :    0.0 ms
Total operation time      : 00:00:31
```
Среднее время записи = 6 ms
95% запросов на запись уложились в 24.3 ms


Запустим нагрузачный тест на чтение
```bash
/opt/cassandra/tools/bin/cassandra-stress read n=1000000

Running with 4 threadCount
Running READ with 4 threads for 1000000 iteration

Results:
Op rate                   :   14,304 op/s  [READ: 14,304 op/s]
Partition rate            :   14,304 pk/s  [READ: 14,304 pk/s]
Row rate                  :   14,304 row/s [READ: 14,304 row/s]
Latency mean              :    0.3 ms [READ: 0.3 ms]
Latency median            :    0.2 ms [READ: 0.2 ms]
Latency 95th percentile   :    0.4 ms [READ: 0.4 ms]
Latency 99th percentile   :    0.6 ms [READ: 0.6 ms]
Latency 99.9th percentile :    5.4 ms [READ: 5.4 ms]
Latency max               :   37.4 ms [READ: 37.4 ms]
Total partitions          :  1,000,000 [READ: 1,000,000]
Total errors              :          0 [READ: 0]
Total GC count            : 0
Total GC memory           : 0.000 KiB
Total GC time             :    0.0 seconds
Avg GC time               :    NaN ms
StdDev GC time            :    0.0 ms
Total operation time      : 00:01:09

Running with 8 threadCount
Running READ with 8 threads for 1000000 iteration

Results:
Op rate                   :   18,419 op/s  [READ: 18,419 op/s]
Partition rate            :   18,419 pk/s  [READ: 18,419 pk/s]
Row rate                  :   18,419 row/s [READ: 18,419 row/s]
Latency mean              :    0.4 ms [READ: 0.4 ms]
Latency median            :    0.3 ms [READ: 0.3 ms]
Latency 95th percentile   :    0.6 ms [READ: 0.6 ms]
Latency 99th percentile   :    1.3 ms [READ: 1.3 ms]
Latency 99.9th percentile :   17.6 ms [READ: 17.6 ms]
Latency max               :   47.6 ms [READ: 47.6 ms]
Total partitions          :  1,000,000 [READ: 1,000,000]
Total errors              :          0 [READ: 0]
Total GC count            : 0
Total GC memory           : 0.000 KiB
Total GC time             :    0.0 seconds
Avg GC time               :    NaN ms
StdDev GC time            :    0.0 ms
Total operation time      : 00:00:54

Improvement over 4 threadCount: 29%

Running with 16 threadCount
Running READ with 16 threads for 1000000 iteration

Results:
Op rate                   :   23,950 op/s  [READ: 23,950 op/s]
Partition rate            :   23,950 pk/s  [READ: 23,950 pk/s]
Row rate                  :   23,950 row/s [READ: 23,950 row/s]
Latency mean              :    0.6 ms [READ: 0.6 ms]
Latency median            :    0.5 ms [READ: 0.5 ms]
Latency 95th percentile   :    1.2 ms [READ: 1.2 ms]
Latency 99th percentile   :    3.7 ms [READ: 3.7 ms]
Latency 99.9th percentile :   21.5 ms [READ: 21.5 ms]
Latency max               :   54.5 ms [READ: 54.5 ms]
Total partitions          :  1,000,000 [READ: 1,000,000]
Total errors              :          0 [READ: 0]
Total GC count            : 0
Total GC memory           : 0.000 KiB
Total GC time             :    0.0 seconds
Avg GC time               :    NaN ms
StdDev GC time            :    0.0 ms
Total operation time      : 00:00:41

Improvement over 8 threadCount: 30%

Running with 24 threadCount
Running READ with 24 threads for 1000000 iteration

Results:
Op rate                   :   24,923 op/s  [READ: 24,923 op/s]
Partition rate            :   24,923 pk/s  [READ: 24,923 pk/s]
Row rate                  :   24,923 row/s [READ: 24,923 row/s]
Latency mean              :    0.9 ms [READ: 0.9 ms]
Latency median            :    0.6 ms [READ: 0.6 ms]
Latency 95th percentile   :    1.9 ms [READ: 1.9 ms]
Latency 99th percentile   :    6.7 ms [READ: 6.7 ms]
Latency 99.9th percentile :   26.7 ms [READ: 26.7 ms]
Latency max               :   54.0 ms [READ: 54.0 ms]
Total partitions          :  1,000,000 [READ: 1,000,000]
Total errors              :          0 [READ: 0]
Total GC count            : 0
Total GC memory           : 0.000 KiB
Total GC time             :    0.0 seconds
Avg GC time               :    NaN ms
StdDev GC time            :    0.0 ms
Total operation time      : 00:00:40
```

При увеличении количества потоков, среднее время ответа на запрос увеличивается, но количество прочитанных строк в секунду тоже увеличивается,
за счёт чего пропускная способность повышается, то есть один и тот же объём данных вычитывается быстрее.