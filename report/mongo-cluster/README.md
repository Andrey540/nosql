# Отчёт по домашнему заданию №3

## Настройки кластера

Перейдём в директорию mondo
```bash
cd mongo
```

Поднимем контейнеры
```bash
docker-compose --file docker-compose-mongo-cluster.yml up -d
```

Инициализируем кластер
```bash
.init_cluster.sh
```

## Шардирование

Создадим коллекцию
```bash
docker exec -it router-01 bash
mongosh
use order
sh.enableSharding("order")
```

Рсшардируем её
```bash
sh.enableSharding("order")
```

Выставим размер чанка равным 1Мб
```bash
use config
db.settings.updateOne(
   { _id: "chunksize" },
   { $set: { _id: "chunksize", value: 1 } },
   { upsert: true }
)
```

Заполним коллекцию
```bash
use order
for (var i=0; i<200000; i++) { db.orders.insertOne({_id: UUID(), user_id: UUID(), amount: Math.random()*100}) }
```

Создадим индекс по полю user_id
```bash
db.orders.createIndex({user_id: 1})
```

Посомтрим на статус шардирования
```bash
sh.status()
    database: {
      _id: 'order',
      primary: 'rs-shard-02',
      partitioned: false,
      version: {
        uuid: UUID("eed2eba7-9fc6-4563-9d94-6082c47c286b"),
        timestamp: Timestamp({ t: 1668348888, i: 20 }),
        lastMod: 1
      }
    },
```
Видим, что наша коллекция ещё не шардирована

Расшардируем её по полю user_id. Ключ шардироания выбрал user_id, тк данная коллекция хранит заказы пользователей. А пользователю нужно видеть все свои заказы.
Поэтому, чтобы данные пользователя находились на одном шарде выбран этот ключ. Конечно, в реальных условиях у разных пользователей количество заказов может сильно отливаться,
что может привести не совсем к равномерному распределению данных. В данном случае на каждую запись создан свой uuid, то есть считается, что у одного пользователя один заказ для упрощения.
```bash
sh.shardCollection("order.orders",{user_id: 1})
```

Посмотрим что получилось
```bash
sh.status()
    database: {
      _id: 'order',
      primary: 'rs-shard-02',
      partitioned: false,
      version: {
        uuid: UUID("eed2eba7-9fc6-4563-9d94-6082c47c286b"),
        timestamp: Timestamp({ t: 1668348888, i: 20 }),
        lastMod: 1
      }
    },
    collections: {
      'order.orders': {
        shardKey: { user_id: 1 },
        unique: false,
        balancing: true,
        chunkMetadata: [
          { shard: 'rs-shard-01', nChunks: 4 },
          { shard: 'rs-shard-02', nChunks: 11 },
          { shard: 'rs-shard-03', nChunks: 3 }
        ],
        chunks: [
          { min: { user_id: MinKey() }, max: { user_id: UUID("0e90acb9-1ba4-4d66-9b12-cb27c6316d8b") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 2, i: 0 }) },
          { min: { user_id: UUID("0e90acb9-1ba4-4d66-9b12-cb27c6316d8b") }, max: { user_id: UUID("1cf91f67-9350-4e2a-998f-c0da09ce98f5") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 3, i: 0 }) },
          { min: { user_id: UUID("1cf91f67-9350-4e2a-998f-c0da09ce98f5") }, max: { user_id: UUID("2bc9888e-ca14-4b3e-b3f0-080c1f6da471") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 4, i: 0 }) },
          { min: { user_id: UUID("2bc9888e-ca14-4b3e-b3f0-080c1f6da471") }, max: { user_id: UUID("3a093503-8ddf-4fb0-ab2a-79fff78b2f04") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 5, i: 0 }) },
          { min: { user_id: UUID("3a093503-8ddf-4fb0-ab2a-79fff78b2f04") }, max: { user_id: UUID("4892ca77-d445-4e9a-a879-460d2963101a") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 6, i: 0 }) },
          { min: { user_id: UUID("4892ca77-d445-4e9a-a879-460d2963101a") }, max: { user_id: UUID("56e89e9e-b4ad-4a31-b321-5e794c3d6b33") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 7, i: 0 }) },
          { min: { user_id: UUID("56e89e9e-b4ad-4a31-b321-5e794c3d6b33") }, max: { user_id: UUID("6533fcfa-848d-42e0-bedf-0a0db64afa84") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 8, i: 0 }) },
          { min: { user_id: UUID("6533fcfa-848d-42e0-bedf-0a0db64afa84") }, max: { user_id: UUID("73be0332-ca0e-4811-bb32-66d2a09ba7a2") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 8, i: 1 }) },
          { min: { user_id: UUID("73be0332-ca0e-4811-bb32-66d2a09ba7a2") }, max: { user_id: UUID("825818f8-0e98-4a4b-9cb5-1cb0d5424021") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 8 }) },
          { min: { user_id: UUID("825818f8-0e98-4a4b-9cb5-1cb0d5424021") }, max: { user_id: UUID("90991783-46c7-45eb-b9cf-529802344281") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 9 }) },
          { min: { user_id: UUID("90991783-46c7-45eb-b9cf-529802344281") }, max: { user_id: UUID("9edefe07-b792-4cd6-ae91-296413fb76da") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 10 }) },
          { min: { user_id: UUID("9edefe07-b792-4cd6-ae91-296413fb76da") }, max: { user_id: UUID("ad4b5b4a-aa21-4ddd-927b-0425e3466fbf") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 11 }) },
          { min: { user_id: UUID("ad4b5b4a-aa21-4ddd-927b-0425e3466fbf") }, max: { user_id: UUID("bb9e80a6-3daf-4fca-ba57-047c008e1322") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 12 }) },
          { min: { user_id: UUID("bb9e80a6-3daf-4fca-ba57-047c008e1322") }, max: { user_id: UUID("ca040334-f96c-4121-9002-4e5320f34517") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 13 }) },
          { min: { user_id: UUID("ca040334-f96c-4121-9002-4e5320f34517") }, max: { user_id: UUID("d7952777-581c-4edb-b5ba-1c40e7828114") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 14 }) },
          { min: { user_id: UUID("d7952777-581c-4edb-b5ba-1c40e7828114") }, max: { user_id: UUID("e527cdd3-f441-4710-85a7-3a989b473d04") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 15 }) },
          { min: { user_id: UUID("e527cdd3-f441-4710-85a7-3a989b473d04") }, max: { user_id: UUID("f27252d4-7ead-44ad-a986-81bebbb41408") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 16 }) },
          { min: { user_id: UUID("f27252d4-7ead-44ad-a986-81bebbb41408") }, max: { user_id: MaxKey() }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 17 }) }
        ],
        tags: []
      }
```
Можно видеть, что получилось 18 чанков, 11 из них ещё на изначальном втором шарде, но некоторые уже начали перебалансироваться

Проверим ещё раз
```bash
sh.status()
     database: {
      _id: 'order',
      primary: 'rs-shard-02',
      partitioned: false,
      version: {
        uuid: UUID("eed2eba7-9fc6-4563-9d94-6082c47c286b"),
        timestamp: Timestamp({ t: 1668348888, i: 20 }),
        lastMod: 1
      }
    },
    collections: {
      'order.orders': {
        shardKey: { user_id: 1 },
        unique: false,
        balancing: true,
        chunkMetadata: [
          { shard: 'rs-shard-01', nChunks: 6 },
          { shard: 'rs-shard-02', nChunks: 6 },
          { shard: 'rs-shard-03', nChunks: 6 }
        ],
        chunks: [
          { min: { user_id: MinKey() }, max: { user_id: UUID("0e90acb9-1ba4-4d66-9b12-cb27c6316d8b") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 2, i: 0 }) },
          { min: { user_id: UUID("0e90acb9-1ba4-4d66-9b12-cb27c6316d8b") }, max: { user_id: UUID("1cf91f67-9350-4e2a-998f-c0da09ce98f5") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 3, i: 0 }) },
          { min: { user_id: UUID("1cf91f67-9350-4e2a-998f-c0da09ce98f5") }, max: { user_id: UUID("2bc9888e-ca14-4b3e-b3f0-080c1f6da471") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 4, i: 0 }) },
          { min: { user_id: UUID("2bc9888e-ca14-4b3e-b3f0-080c1f6da471") }, max: { user_id: UUID("3a093503-8ddf-4fb0-ab2a-79fff78b2f04") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 5, i: 0 }) },
          { min: { user_id: UUID("3a093503-8ddf-4fb0-ab2a-79fff78b2f04") }, max: { user_id: UUID("4892ca77-d445-4e9a-a879-460d2963101a") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 6, i: 0 }) },
          { min: { user_id: UUID("4892ca77-d445-4e9a-a879-460d2963101a") }, max: { user_id: UUID("56e89e9e-b4ad-4a31-b321-5e794c3d6b33") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 7, i: 0 }) },
          { min: { user_id: UUID("56e89e9e-b4ad-4a31-b321-5e794c3d6b33") }, max: { user_id: UUID("6533fcfa-848d-42e0-bedf-0a0db64afa84") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 8, i: 0 }) },
          { min: { user_id: UUID("6533fcfa-848d-42e0-bedf-0a0db64afa84") }, max: { user_id: UUID("73be0332-ca0e-4811-bb32-66d2a09ba7a2") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 9, i: 0 }) },
          { min: { user_id: UUID("73be0332-ca0e-4811-bb32-66d2a09ba7a2") }, max: { user_id: UUID("825818f8-0e98-4a4b-9cb5-1cb0d5424021") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 10, i: 0 }) },
          { min: { user_id: UUID("825818f8-0e98-4a4b-9cb5-1cb0d5424021") }, max: { user_id: UUID("90991783-46c7-45eb-b9cf-529802344281") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 11, i: 0 }) },
          { min: { user_id: UUID("90991783-46c7-45eb-b9cf-529802344281") }, max: { user_id: UUID("9edefe07-b792-4cd6-ae91-296413fb76da") }, 'on shard': 'rs-shard-01', 'last modified': Timestamp({ t: 12, i: 0 }) },
          { min: { user_id: UUID("9edefe07-b792-4cd6-ae91-296413fb76da") }, max: { user_id: UUID("ad4b5b4a-aa21-4ddd-927b-0425e3466fbf") }, 'on shard': 'rs-shard-03', 'last modified': Timestamp({ t: 13, i: 0 }) },
          { min: { user_id: UUID("ad4b5b4a-aa21-4ddd-927b-0425e3466fbf") }, max: { user_id: UUID("bb9e80a6-3daf-4fca-ba57-047c008e1322") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 13, i: 1 }) },
          { min: { user_id: UUID("bb9e80a6-3daf-4fca-ba57-047c008e1322") }, max: { user_id: UUID("ca040334-f96c-4121-9002-4e5320f34517") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 13 }) },
          { min: { user_id: UUID("ca040334-f96c-4121-9002-4e5320f34517") }, max: { user_id: UUID("d7952777-581c-4edb-b5ba-1c40e7828114") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 14 }) },
          { min: { user_id: UUID("d7952777-581c-4edb-b5ba-1c40e7828114") }, max: { user_id: UUID("e527cdd3-f441-4710-85a7-3a989b473d04") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 15 }) },
          { min: { user_id: UUID("e527cdd3-f441-4710-85a7-3a989b473d04") }, max: { user_id: UUID("f27252d4-7ead-44ad-a986-81bebbb41408") }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 16 }) },
          { min: { user_id: UUID("f27252d4-7ead-44ad-a986-81bebbb41408") }, max: { user_id: MaxKey() }, 'on shard': 'rs-shard-02', 'last modified': Timestamp({ t: 1, i: 17 }) }
        ],
        tags: []
      }
```
Чанки перебалансировалсиь

## Отказоустойчивость

Проверим состояние первого шарда
Для этого зайдём на первый шард
```bash
docker exec -it shard-01-node-c bash
```

И проверим статус реплики
```bash
mongosh
rs.status().members
[
  {
    _id: 0,
    name: 'shard01-a:27017',
    health: 1,
    state: 1,
    stateStr: 'PRIMARY',
    uptime: 614,
    optime: { ts: Timestamp({ t: 1668369532, i: 1 }), t: Long("11") },
    optimeDate: ISODate("2022-11-13T19:58:52.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    lastDurableWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    electionTime: Timestamp({ t: 1668368952, i: 1 }),
    electionDate: ISODate("2022-11-13T19:49:12.000Z"),
    configVersion: 1,
    configTerm: 11,
    self: true,
    lastHeartbeatMessage: ''
  },
  {
    _id: 1,
    name: 'shard01-b:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 601,
    optime: { ts: Timestamp({ t: 1668369532, i: 1 }), t: Long("11") },
    optimeDurable: { ts: Timestamp({ t: 1668369532, i: 1 }), t: Long("11") },
    optimeDate: ISODate("2022-11-13T19:58:52.000Z"),
    optimeDurableDate: ISODate("2022-11-13T19:58:52.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    lastDurableWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    lastHeartbeat: ISODate("2022-11-13T19:59:00.680Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T19:59:01.094Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: '',
    syncSourceHost: 'shard01-a:27017',
    syncSourceId: 0,
    infoMessage: '',
    configVersion: 1,
    configTerm: 11
  },
  {
    _id: 2,
    name: 'shard01-c:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 600,
    optime: { ts: Timestamp({ t: 1668369532, i: 1 }), t: Long("11") },
    optimeDurable: { ts: Timestamp({ t: 1668369532, i: 1 }), t: Long("11") },
    optimeDate: ISODate("2022-11-13T19:58:52.000Z"),
    optimeDurableDate: ISODate("2022-11-13T19:58:52.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    lastDurableWallTime: ISODate("2022-11-13T19:58:52.433Z"),
    lastHeartbeat: ISODate("2022-11-13T19:59:00.680Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T19:59:01.094Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: '',
    syncSourceHost: 'shard01-a:27017',
    syncSourceId: 0,
    infoMessage: '',
    configVersion: 1,
    configTerm: 11
  }
]
```

Видим, что мастером является shard01-a
Попробуем его погасить

```bash
docker stop shard-01-node-a
```

Проверим статус
```bash
rs.status().members
[
  {
    _id: 0,
    name: 'shard01-a:27017',
    health: 0,
    state: 8,
    stateStr: '(not reachable/healthy)',
    uptime: 0,
    optime: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDurable: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDate: ISODate("1970-01-01T00:00:00.000Z"),
    optimeDurableDate: ISODate("1970-01-01T00:00:00.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:02:20.648Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:02:20.648Z"),
    lastHeartbeat: ISODate("2022-11-13T20:02:51.829Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:02:29.654Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: 'Error connecting to shard01-a:27017 :: caused by :: Could not find address for shard01-a:27017: SocketException: Host not found (non-authoritative), try again later',
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    configVersion: 1,
    configTerm: 12
  },
  {
    _id: 1,
    name: 'shard01-b:27017',
    health: 1,
    state: 1,
    stateStr: 'PRIMARY',
    uptime: 831,
    optime: { ts: Timestamp({ t: 1668369770, i: 1 }), t: Long("12") },
    optimeDurable: { ts: Timestamp({ t: 1668369770, i: 1 }), t: Long("12") },
    optimeDate: ISODate("2022-11-13T20:02:50.000Z"),
    optimeDurableDate: ISODate("2022-11-13T20:02:50.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:02:50.655Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:02:50.655Z"),
    lastHeartbeat: ISODate("2022-11-13T20:02:51.707Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:02:52.660Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: '',
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    electionTime: Timestamp({ t: 1668369740, i: 1 }),
    electionDate: ISODate("2022-11-13T20:02:20.000Z"),
    configVersion: 1,
    configTerm: 12
  },
  {
    _id: 2,
    name: 'shard01-c:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 847,
    optime: { ts: Timestamp({ t: 1668369770, i: 1 }), t: Long("12") },
    optimeDate: ISODate("2022-11-13T20:02:50.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:02:50.655Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:02:50.655Z"),
    syncSourceHost: 'shard01-b:27017',
    syncSourceId: 1,
    infoMessage: '',
    configVersion: 1,
    configTerm: 12,
    self: true,
    lastHeartbeatMessage: ''
  }
]
```

Видим, что shard01-a стал неактивным, а мастером стал shard01-b

Попробуем погасить второй шард
```bash
docker stop shard-01-node-b
```

Проверим статус
```bash
rs.status().members
[
  {
    _id: 0,
    name: 'shard01-a:27017',
    health: 0,
    state: 8,
    stateStr: '(not reachable/healthy)',
    uptime: 0,
    optime: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDurable: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDate: ISODate("1970-01-01T00:00:00.000Z"),
    optimeDurableDate: ISODate("1970-01-01T00:00:00.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:02:20.648Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:02:20.648Z"),
    lastHeartbeat: ISODate("2022-11-13T20:04:54.439Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:02:29.654Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: 'Error connecting to shard01-a:27017 :: caused by :: Could not find address for shard01-a:27017: SocketException: Host not found (non-authoritative), try again later',
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    configVersion: 1,
    configTerm: 12
  },
  {
    _id: 1,
    name: 'shard01-b:27017',
    health: 0,
    state: 8,
    stateStr: '(not reachable/healthy)',
    uptime: 0,
    optime: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDurable: { ts: Timestamp({ t: 0, i: 0 }), t: Long("-1") },
    optimeDate: ISODate("1970-01-01T00:00:00.000Z"),
    optimeDurableDate: ISODate("1970-01-01T00:00:00.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:04:10.657Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:04:10.657Z"),
    lastHeartbeat: ISODate("2022-11-13T20:04:54.438Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:04:20.691Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: 'Error connecting to shard01-b:27017 :: caused by :: Could not find address for shard01-b:27017: SocketException: Host not found (non-authoritative), try again later',
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    configVersion: 1,
    configTerm: 12
  },
  {
    _id: 2,
    name: 'shard01-c:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 968,
    optime: { ts: Timestamp({ t: 1668369850, i: 1 }), t: Long("12") },
    optimeDate: ISODate("2022-11-13T20:04:10.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:04:10.657Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:04:10.657Z"),
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    configVersion: 1,
    configTerm: 12,
    self: true,
    lastHeartbeatMessage: ''
  }
]
```

Видим, что shard01-a и shard01-b неактивны, шард shard01-c активен, но является ведомым, тк большинства он полочить не смог, поэтому шард находится в нерабочем состоянии

Поднимем обратно контейнеры
```bash
docker-compose --file docker-compose-mongo-cluster.yml up -d
```

Поверим статус
```bash
rs.status().members
[
  {
    _id: 0,
    name: 'shard01-a:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 30,
    optime: { ts: Timestamp({ t: 1668370089, i: 1 }), t: Long("14") },
    optimeDurable: { ts: Timestamp({ t: 1668370089, i: 1 }), t: Long("14") },
    optimeDate: ISODate("2022-11-13T20:08:09.000Z"),
    optimeDurableDate: ISODate("2022-11-13T20:08:09.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    lastHeartbeat: ISODate("2022-11-13T20:08:15.138Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:08:13.643Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: '',
    syncSourceHost: 'shard01-c:27017',
    syncSourceId: 2,
    infoMessage: '',
    configVersion: 1,
    configTerm: 14
  },
  {
    _id: 1,
    name: 'shard01-b:27017',
    health: 1,
    state: 2,
    stateStr: 'SECONDARY',
    uptime: 28,
    optime: { ts: Timestamp({ t: 1668370089, i: 1 }), t: Long("14") },
    optimeDurable: { ts: Timestamp({ t: 1668370089, i: 1 }), t: Long("14") },
    optimeDate: ISODate("2022-11-13T20:08:09.000Z"),
    optimeDurableDate: ISODate("2022-11-13T20:08:09.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    lastHeartbeat: ISODate("2022-11-13T20:08:15.139Z"),
    lastHeartbeatRecv: ISODate("2022-11-13T20:08:13.643Z"),
    pingMs: Long("0"),
    lastHeartbeatMessage: '',
    syncSourceHost: 'shard01-c:27017',
    syncSourceId: 2,
    infoMessage: '',
    configVersion: 1,
    configTerm: 14
  },
  {
    _id: 2,
    name: 'shard01-c:27017',
    health: 1,
    state: 1,
    stateStr: 'PRIMARY',
    uptime: 1169,
    optime: { ts: Timestamp({ t: 1668370089, i: 1 }), t: Long("14") },
    optimeDate: ISODate("2022-11-13T20:08:09.000Z"),
    lastAppliedWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    lastDurableWallTime: ISODate("2022-11-13T20:08:09.138Z"),
    syncSourceHost: '',
    syncSourceId: -1,
    infoMessage: '',
    electionTime: Timestamp({ t: 1668370069, i: 1 }),
    electionDate: ISODate("2022-11-13T20:07:49.000Z"),
    configVersion: 1,
    configTerm: 14,
    self: true,
    lastHeartbeatMessage: ''
  }
]
```

Все ноды активны, лидером же стала shard01-c

С configsvr можно проделать аналогичные операции, результат будут такой же

Если потушить две ноды второго шарда
```bash
docker stop shard-02-node-a shard-02-node-b
```

И выполнить запрос к коллекции
```bash
b.orders.countDocuments()
MongoServerError: Could not find host matching read preference { mode: "primary" } for set rs-shard-02
```
То получим ошибку.
Аналогичная картина будет, если потушить все три ноды второго щарда

## Аутентификация и ролевой доступ

Настроим внутриннюю аутентификацию и настроим роли пользователей

Создадим файл с ключом доступа для нод кластера
```bash
openssl rand -base64 756 > keyfile
```

Выставим права ключа доступа
```bash
chmod 600 keyfile
sudo chown 999:999 keyfile
```

Перезапустим кластер
Зайдём в кластер
```bash
docker exec -it router-01 bash
```

Зайдём в mongo без авторизации чтобы создать первого пользователя
```bash
mongosh
```

Создадим пользователя с рутовыми правами
```bash
mongosh
db.getSiblingDB("admin").createUser(
{
  user: 'admin',
  pwd: 'Jez346Cy',
  roles: [ { role: 'root', db: 'admin' } ]
});
```

Перезайдём в mongo от администратора
```bash
exit
mongosh -u admin -p Jez346Cy --authenticationDatabase admin
```

```bash
db.getSiblingDB("admin").createUser(
{
  user: 'clusterAdmin',
  pwd: '98Z29Asd',
  roles: [ { role: 'clusterAdmin', db: 'admin' } ]
});
```

Создадим пользователя для работы с базой order
```bash
db.getSiblingDB("order").createUser({      
     user: 'order',      
     pwd: 'wh0Z6e0h',      
     roles: [ { role: 'readWrite', db: 'order' } ] 
});
```

Проверим, что у пользователя order есть права на базу order
```bash
exit
mongosh -u order -p wh0Z6e0h --authenticationDatabase order
```

Выполним запрос
```bash
db.orders.countDocuments()
200000
```
Получили ответ

Попробуем выполнить запрос, для которого нужны админские права
```bash
use admin
db.system.users.find()
MongoServerError: not authorized on admin to execute command { find: "system.users", filter: {}, lsid: { id: UUID("8dd97ba4-66c2-4200-90d5-7cbde5a7aff2") }, $clusterTime: { clusterTime: Timestamp(1668538052, 1), signature: { hash: BinData(0, 2ECE885583BE6600035948241ED1EEA508CDE6BE), keyId: 7165237903478489106 } }, $db: "admin" }
```
Получили ошибку прав доступа