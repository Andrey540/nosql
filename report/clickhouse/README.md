# Отчёт по домашнему заданию №5

## Разворачивание clickhouse

Перейдём в директорию clickhouse
```bash
cd clickhouse
```

Скачаем тестовые базы данных
```bash
curl https://datasets.clickhouse.com/hits/tsv/hits_v1.tsv.xz | unxz --threads=`nproc` > hits_v1.tsv
curl https://datasets.clickhouse.com/visits/tsv/visits_v1.tsv.xz | unxz --threads=`nproc` > visits_v1.tsv
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Зайдём в контейнер клиента clickhouse
```bash
docker exec -it clickhouse-client bash
```

Создадим тестовую базу
```bash
clickhouse-client --host clickhouse-server --query "CREATE DATABASE IF NOT EXISTS tutorial"
```

Зайдём в clickhouse-client
```bash
clickhouse-client --host clickhouse-server
```

И создадим таблицы
```clickhouse
CREATE TABLE tutorial.hits_v1
(
    `WatchID` UInt64,
    `JavaEnable` UInt8,
    `Title` String,
    `GoodEvent` Int16,
    `EventTime` DateTime,
    `EventDate` Date,
    `CounterID` UInt32,
    `ClientIP` UInt32,
    `ClientIP6` FixedString(16),
    `RegionID` UInt32,
    `UserID` UInt64,
    `CounterClass` Int8,
    `OS` UInt8,
    `UserAgent` UInt8,
    `URL` String,
    `Referer` String,
    `URLDomain` String,
    `RefererDomain` String,
    `Refresh` UInt8,
    `IsRobot` UInt8,
    `RefererCategories` Array(UInt16),
    `URLCategories` Array(UInt16),
    `URLRegions` Array(UInt32),
    `RefererRegions` Array(UInt32),
    `ResolutionWidth` UInt16,
    `ResolutionHeight` UInt16,
    `ResolutionDepth` UInt8,
    `FlashMajor` UInt8,
    `FlashMinor` UInt8,
    `FlashMinor2` String,
    `NetMajor` UInt8,
    `NetMinor` UInt8,
    `UserAgentMajor` UInt16,
    `UserAgentMinor` FixedString(2),
    `CookieEnable` UInt8,
    `JavascriptEnable` UInt8,
    `IsMobile` UInt8,
    `MobilePhone` UInt8,
    `MobilePhoneModel` String,
    `Params` String,
    `IPNetworkID` UInt32,
    `TraficSourceID` Int8,
    `SearchEngineID` UInt16,
    `SearchPhrase` String,
    `AdvEngineID` UInt8,
    `IsArtifical` UInt8,
    `WindowClientWidth` UInt16,
    `WindowClientHeight` UInt16,
    `ClientTimeZone` Int16,
    `ClientEventTime` DateTime,
    `SilverlightVersion1` UInt8,
    `SilverlightVersion2` UInt8,
    `SilverlightVersion3` UInt32,
    `SilverlightVersion4` UInt16,
    `PageCharset` String,
    `CodeVersion` UInt32,
    `IsLink` UInt8,
    `IsDownload` UInt8,
    `IsNotBounce` UInt8,
    `FUniqID` UInt64,
    `HID` UInt32,
    `IsOldCounter` UInt8,
    `IsEvent` UInt8,
    `IsParameter` UInt8,
    `DontCountHits` UInt8,
    `WithHash` UInt8,
    `HitColor` FixedString(1),
    `UTCEventTime` DateTime,
    `Age` UInt8,
    `Sex` UInt8,
    `Income` UInt8,
    `Interests` UInt16,
    `Robotness` UInt8,
    `GeneralInterests` Array(UInt16),
    `RemoteIP` UInt32,
    `RemoteIP6` FixedString(16),
    `WindowName` Int32,
    `OpenerName` Int32,
    `HistoryLength` Int16,
    `BrowserLanguage` FixedString(2),
    `BrowserCountry` FixedString(2),
    `SocialNetwork` String,
    `SocialAction` String,
    `HTTPError` UInt16,
    `SendTiming` Int32,
    `DNSTiming` Int32,
    `ConnectTiming` Int32,
    `ResponseStartTiming` Int32,
    `ResponseEndTiming` Int32,
    `FetchTiming` Int32,
    `RedirectTiming` Int32,
    `DOMInteractiveTiming` Int32,
    `DOMContentLoadedTiming` Int32,
    `DOMCompleteTiming` Int32,
    `LoadEventStartTiming` Int32,
    `LoadEventEndTiming` Int32,
    `NSToDOMContentLoadedTiming` Int32,
    `FirstPaintTiming` Int32,
    `RedirectCount` Int8,
    `SocialSourceNetworkID` UInt8,
    `SocialSourcePage` String,
    `ParamPrice` Int64,
    `ParamOrderID` String,
    `ParamCurrency` FixedString(3),
    `ParamCurrencyID` UInt16,
    `GoalsReached` Array(UInt32),
    `OpenstatServiceName` String,
    `OpenstatCampaignID` String,
    `OpenstatAdID` String,
    `OpenstatSourceID` String,
    `UTMSource` String,
    `UTMMedium` String,
    `UTMCampaign` String,
    `UTMContent` String,
    `UTMTerm` String,
    `FromTag` String,
    `HasGCLID` UInt8,
    `RefererHash` UInt64,
    `URLHash` UInt64,
    `CLID` UInt32,
    `YCLID` UInt64,
    `ShareService` String,
    `ShareURL` String,
    `ShareTitle` String,
    `ParsedParams` Nested(
        Key1 String,
        Key2 String,
        Key3 String,
        Key4 String,
        Key5 String,
        ValueDouble Float64),
    `IslandID` FixedString(16),
    `RequestNum` UInt32,
    `RequestTry` UInt8
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(EventDate)
ORDER BY (CounterID, EventDate, intHash32(UserID))
SAMPLE BY intHash32(UserID)
```

```clickhouse
CREATE TABLE tutorial.visits_v1
(
    `CounterID` UInt32,
    `StartDate` Date,
    `Sign` Int8,
    `IsNew` UInt8,
    `VisitID` UInt64,
    `UserID` UInt64,
    `StartTime` DateTime,
    `Duration` UInt32,
    `UTCStartTime` DateTime,
    `PageViews` Int32,
    `Hits` Int32,
    `IsBounce` UInt8,
    `Referer` String,
    `StartURL` String,
    `RefererDomain` String,
    `StartURLDomain` String,
    `EndURL` String,
    `LinkURL` String,
    `IsDownload` UInt8,
    `TraficSourceID` Int8,
    `SearchEngineID` UInt16,
    `SearchPhrase` String,
    `AdvEngineID` UInt8,
    `PlaceID` Int32,
    `RefererCategories` Array(UInt16),
    `URLCategories` Array(UInt16),
    `URLRegions` Array(UInt32),
    `RefererRegions` Array(UInt32),
    `IsYandex` UInt8,
    `GoalReachesDepth` Int32,
    `GoalReachesURL` Int32,
    `GoalReachesAny` Int32,
    `SocialSourceNetworkID` UInt8,
    `SocialSourcePage` String,
    `MobilePhoneModel` String,
    `ClientEventTime` DateTime,
    `RegionID` UInt32,
    `ClientIP` UInt32,
    `ClientIP6` FixedString(16),
    `RemoteIP` UInt32,
    `RemoteIP6` FixedString(16),
    `IPNetworkID` UInt32,
    `SilverlightVersion3` UInt32,
    `CodeVersion` UInt32,
    `ResolutionWidth` UInt16,
    `ResolutionHeight` UInt16,
    `UserAgentMajor` UInt16,
    `UserAgentMinor` UInt16,
    `WindowClientWidth` UInt16,
    `WindowClientHeight` UInt16,
    `SilverlightVersion2` UInt8,
    `SilverlightVersion4` UInt16,
    `FlashVersion3` UInt16,
    `FlashVersion4` UInt16,
    `ClientTimeZone` Int16,
    `OS` UInt8,
    `UserAgent` UInt8,
    `ResolutionDepth` UInt8,
    `FlashMajor` UInt8,
    `FlashMinor` UInt8,
    `NetMajor` UInt8,
    `NetMinor` UInt8,
    `MobilePhone` UInt8,
    `SilverlightVersion1` UInt8,
    `Age` UInt8,
    `Sex` UInt8,
    `Income` UInt8,
    `JavaEnable` UInt8,
    `CookieEnable` UInt8,
    `JavascriptEnable` UInt8,
    `IsMobile` UInt8,
    `BrowserLanguage` UInt16,
    `BrowserCountry` UInt16,
    `Interests` UInt16,
    `Robotness` UInt8,
    `GeneralInterests` Array(UInt16),
    `Params` Array(String),
    `Goals` Nested(
        ID UInt32,
        Serial UInt32,
        EventTime DateTime,
        Price Int64,
        OrderID String,
        CurrencyID UInt32),
    `WatchIDs` Array(UInt64),
    `ParamSumPrice` Int64,
    `ParamCurrency` FixedString(3),
    `ParamCurrencyID` UInt16,
    `ClickLogID` UInt64,
    `ClickEventID` Int32,
    `ClickGoodEvent` Int32,
    `ClickEventTime` DateTime,
    `ClickPriorityID` Int32,
    `ClickPhraseID` Int32,
    `ClickPageID` Int32,
    `ClickPlaceID` Int32,
    `ClickTypeID` Int32,
    `ClickResourceID` Int32,
    `ClickCost` UInt32,
    `ClickClientIP` UInt32,
    `ClickDomainID` UInt32,
    `ClickURL` String,
    `ClickAttempt` UInt8,
    `ClickOrderID` UInt32,
    `ClickBannerID` UInt32,
    `ClickMarketCategoryID` UInt32,
    `ClickMarketPP` UInt32,
    `ClickMarketCategoryName` String,
    `ClickMarketPPName` String,
    `ClickAWAPSCampaignName` String,
    `ClickPageName` String,
    `ClickTargetType` UInt16,
    `ClickTargetPhraseID` UInt64,
    `ClickContextType` UInt8,
    `ClickSelectType` Int8,
    `ClickOptions` String,
    `ClickGroupBannerID` Int32,
    `OpenstatServiceName` String,
    `OpenstatCampaignID` String,
    `OpenstatAdID` String,
    `OpenstatSourceID` String,
    `UTMSource` String,
    `UTMMedium` String,
    `UTMCampaign` String,
    `UTMContent` String,
    `UTMTerm` String,
    `FromTag` String,
    `HasGCLID` UInt8,
    `FirstVisit` DateTime,
    `PredLastVisit` Date,
    `LastVisit` Date,
    `TotalVisits` UInt32,
    `TraficSource` Nested(
        ID Int8,
        SearchEngineID UInt16,
        AdvEngineID UInt8,
        PlaceID UInt16,
        SocialSourceNetworkID UInt8,
        Domain String,
        SearchPhrase String,
        SocialSourcePage String),
    `Attendance` FixedString(16),
    `CLID` UInt32,
    `YCLID` UInt64,
    `NormalizedRefererHash` UInt64,
    `SearchPhraseHash` UInt64,
    `RefererDomainHash` UInt64,
    `NormalizedStartURLHash` UInt64,
    `StartURLDomainHash` UInt64,
    `NormalizedEndURLHash` UInt64,
    `TopLevelDomain` UInt64,
    `URLScheme` UInt64,
    `OpenstatServiceNameHash` UInt64,
    `OpenstatCampaignIDHash` UInt64,
    `OpenstatAdIDHash` UInt64,
    `OpenstatSourceIDHash` UInt64,
    `UTMSourceHash` UInt64,
    `UTMMediumHash` UInt64,
    `UTMCampaignHash` UInt64,
    `UTMContentHash` UInt64,
    `UTMTermHash` UInt64,
    `FromHash` UInt64,
    `WebVisorEnabled` UInt8,
    `WebVisorActivity` UInt32,
    `ParsedParams` Nested(
        Key1 String,
        Key2 String,
        Key3 String,
        Key4 String,
        Key5 String,
        ValueDouble Float64),
    `Market` Nested(
        Type UInt8,
        GoalID UInt32,
        OrderID String,
        OrderPrice Int64,
        PP UInt32,
        DirectPlaceID UInt32,
        DirectOrderID UInt32,
        DirectBannerID UInt32,
        GoodID String,
        GoodName String,
        GoodQuantity Int32,
        GoodPrice Int64),
    `IslandID` FixedString(16)
)
ENGINE = CollapsingMergeTree(Sign)
PARTITION BY toYYYYMM(StartDate)
ORDER BY (CounterID, StartDate, intHash32(UserID), VisitID)
SAMPLE BY intHash32(UserID)
```

Выйдём из клиента, чтобы заимпортить данные
```bash
exit
```

Заимпортим данные
```bash
clickhouse-client --host clickhouse-server --query "INSERT INTO tutorial.hits_v1 FORMAT TSV" --max_insert_block_size=100000 < /tmp/hits_v1.tsv
clickhouse-client --host clickhouse-server --query "INSERT INTO tutorial.visits_v1 FORMAT TSV" --max_insert_block_size=100000 < /tmp/visits_v1.tsv
```

Зайдём обратно в clickhouse-client для того, чтобы посмотреть что получилось
```bash
clickhouse-client --host clickhouse-server
```

Сначала оптимизируем получившиеся таблицы чтобы устранить возможне дубли и хранение их было оптимально
```clickhouse
OPTIMIZE TABLE tutorial.hits_v1 FINAL
```

```clickhouse
OPTIMIZE TABLE tutorial.visits_v1 FINAL
```

Посомтрим на количество записей
```clickhouse
SELECT COUNT(*) FROM tutorial.hits_v1
┌─count()─┐
│ 8873898 │
└─────────┘
```
```clickhouse
SELECT COUNT(*) FROM tutorial.visits_v1
┌─count()─┐
│ 1676861 │
└─────────┘
```

Выполним аггрегирующие запросы
```clickhouse
SELECT
    StartURL AS URL,
    AVG(Duration) AS AvgDuration
FROM tutorial.visits_v1
WHERE StartDate BETWEEN '2014-03-23' AND '2014-03-30'
GROUP BY URL
ORDER BY AvgDuration DESC
LIMIT 10
┌─URL─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬─AvgDuration─┐
│ http://itpalanija-pri-patrivative=0&ads_app_user                                                                                                                                                │       60127 │
│ http://renaul-myd-ukraine                                                                                                                                                                       │       58938 │
│ http://karta/Futbol/dynamo.kiev.ua/kawaica.su/648                                                                                                                                               │       56538 │
│ https://moda/vyikroforum1/top.ru/moscow/delo-product/trend_sms/multitryaset/news/2014/03/201000                                                                                                 │       55218 │
│ http://e.mail=on&default?abid=2061&scd=yes&option?r=city_inter.com/menu&site-zaferio.ru/c/m.ensor.net/ru/login=false&orderStage.php?Brandidatamalystyle/20Mar2014%2F007%2F94dc8d2e06e56ed56bbdd │       51378 │
│ http://karta/Futbol/dynas.com/haberler.ru/messages.yandsearchives/494503_lte_13800200319                                                                                                        │       49078 │
│ http://xmusic/vstreatings of speeds                                                                                                                                                             │       36925 │
│ http://news.ru/yandex.ru/api.php&api=http://toberria.ru/aphorizana                                                                                                                              │       36902 │
│ http://bashmelnykh-metode.net/video/#!/video/emberkas.ru/detskij-yazi.com/iframe/default.aspx?id=760928&noreask=1&source                                                                        │       34323 │
│ http://censonhaber/547-popalientLog=0&strizhki-petro%3D&comeback=search?lr=213&text                                                                                                             │       31773 │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────────┘

10 rows in set. Elapsed: 0.065 sec. Processed 1.44 million rows, 113.24 MB (22.16 million rows/s., 1.75 GB/s.)
```
В итоге было обработано 1.4 млн записей за 65 миллисекунд

Выполним ещё один запрос
```clickhouse
SELECT
    sum(Sign) AS visits,
    sumIf(Sign, has(Goals.ID, 1105530)) AS goal_visits,
    (100. * goal_visits) / visits AS goal_percent
FROM tutorial.visits_v1
WHERE (CounterID = 912887) AND (toYYYYMM(StartDate) = 201403) AND (domain(StartURL) = 'yandex.ru')
┌─visits─┬─goal_visits─┬──────goal_percent─┐
│  10543 │        8553 │ 81.12491700654462 │
└────────┴─────────────┴───────────────────┘

1 rows in set. Elapsed: 0.016 sec. Processed 19.72 thousand rows, 3.44 MB (1.23 million rows/s., 214.12 MB/s.)
```
В итоге было обработано 1.23 млн записей за 16 миллисекунд


## Дополнительная тестовая база
Создадим таблицу покупок в UK
```clickhouse
CREATE TABLE uk_price_paid
(
    price UInt32,
    date Date,
    postcode1 LowCardinality(String),
    postcode2 LowCardinality(String),
    type Enum8('terraced' = 1, 'semi-detached' = 2, 'detached' = 3, 'flat' = 4, 'other' = 0),
    is_new UInt8,
    duration Enum8('freehold' = 1, 'leasehold' = 2, 'unknown' = 0),
    addr1 String,
    addr2 String,
    street LowCardinality(String),
    locality LowCardinality(String),
    town LowCardinality(String),
    district LowCardinality(String),
    county LowCardinality(String)
)
    ENGINE = MergeTree
        ORDER BY (postcode1, postcode2, addr1, addr2);
```

и заполним её с интернет ресурса
```clickhouse
INSERT INTO uk_price_paid
WITH
    splitByChar(' ', postcode) AS p
SELECT
    toUInt32(price_string) AS price,
    parseDateTimeBestEffortUS(time) AS date,
    p[1] AS postcode1,
    p[2] AS postcode2,
    transform(a, ['T', 'S', 'D', 'F', 'O'], ['terraced', 'semi-detached', 'detached', 'flat', 'other']) AS type,
    b = 'Y' AS is_new,
    transform(c, ['F', 'L', 'U'], ['freehold', 'leasehold', 'unknown']) AS duration,
    addr1,
    addr2,
    street,
    locality,
    town,
    district,
    county
FROM url(
        'http://prod.publicdata.landregistry.gov.uk.s3-website-eu-west-1.amazonaws.com/pp-complete.csv',
        'CSV',
        'uuid_string String,
        price_string String,
        time String,
        postcode String,
        a String,
        b String,
        c String,
        addr1 String,
        addr2 String,
        street String,
        locality String,
        town String,
        district String,
        county String,
        d String,
        e String'
    ) SETTINGS max_http_get_redirects=10;
```

Проверим сколько записей создалось
```clickhouse
SELECT count(*) FROM uk_price_paid
┌──count()─┐
│ 27734966 │
└──────────┘
```

Проверим сколько места на диске она занимают
```clickhouse
SELECT formatReadableSize(total_bytes)
FROM system.tables
WHERE name = 'uk_price_paid'
┌─formatReadableSize(total_bytes)─┐
│ 298.17 MiB                      │
└─────────────────────────────────┘

```

Выполним аналитические запросы.
Средняя цена за год
```clickhouse
SELECT
   toYear(date) AS year,
   round(avg(price)) AS price,
   bar(price, 0, 1000000, 80)
FROM uk_price_paid
GROUP BY year
ORDER BY year

┌─year─┬──price─┬─bar(round(avg(price)), 0, 1000000, 80)─┐
│ 1995 │  67935 │ █████▍                                 │
│ 1996 │  71509 │ █████▋                                 │
│ 1997 │  78537 │ ██████▎                                │
│ 1998 │  85442 │ ██████▋                                │
│ 1999 │  96038 │ ███████▋                               │
│ 2000 │ 107487 │ ████████▌                              │
│ 2001 │ 118890 │ █████████▌                             │
│ 2002 │ 137955 │ ███████████                            │
│ 2003 │ 155895 │ ████████████▍                          │
│ 2004 │ 178889 │ ██████████████▎                        │
│ 2005 │ 189361 │ ███████████████▏                       │
│ 2006 │ 203532 │ ████████████████▎                      │
│ 2007 │ 219375 │ █████████████████▌                     │
│ 2008 │ 217042 │ █████████████████▎                     │
│ 2009 │ 213418 │ █████████████████                      │
│ 2010 │ 236111 │ ██████████████████▊                    │
│ 2011 │ 232807 │ ██████████████████▌                    │
│ 2012 │ 238383 │ ███████████████████                    │
│ 2013 │ 256927 │ ████████████████████▌                  │
│ 2014 │ 280012 │ ██████████████████████▍                │
│ 2015 │ 297280 │ ███████████████████████▋               │
│ 2016 │ 313533 │ █████████████████████████              │
│ 2017 │ 346430 │ ███████████████████████████▋           │
│ 2018 │ 350609 │ ████████████████████████████           │
│ 2019 │ 352401 │ ████████████████████████████▏          │
│ 2020 │ 376482 │ ██████████████████████████████         │
│ 2021 │ 382080 │ ██████████████████████████████▌        │
│ 2022 │ 378166 │ ██████████████████████████████▎        │
└──────┴────────┴────────────────────────────────────────┘

28 rows in set. Elapsed: 0.069 sec. Processed 27.73 million rows, 166.41 MB (403.27 million rows/s., 2.42 GB/s.) 
```

Самые дорогие места
```clickhouse
SELECT
    town,
    district,
    count() AS c,
    round(avg(price)) AS price,
    bar(price, 0, 5000000, 100)
FROM uk_price_paid
WHERE date >= '2020-01-01'
GROUP BY
    town,
    district
HAVING c >= 100
ORDER BY price DESC
LIMIT 50

┌─town─────────────────┬─district───────────────┬─────c─┬───price─┬─bar(round(avg(price)), 0, 5000000, 100)─────────────────────────┐
│ LONDON               │ CITY OF LONDON         │   648 │ 3122204 │ ██████████████████████████████████████████████████████████████▍ │
│ LONDON               │ CITY OF WESTMINSTER    │  8027 │ 2892541 │ █████████████████████████████████████████████████████████▋      │
│ LONDON               │ KENSINGTON AND CHELSEA │  5638 │ 2379209 │ ███████████████████████████████████████████████▌                │
│ LEATHERHEAD          │ ELMBRIDGE              │   231 │ 2111393 │ ██████████████████████████████████████████▏                     │
│ VIRGINIA WATER       │ RUNNYMEDE              │   340 │ 1952105 │ ███████████████████████████████████████                         │
│ LONDON               │ CAMDEN                 │  6412 │ 1656120 │ █████████████████████████████████                               │
│ NORTHWOOD            │ THREE RIVERS           │   133 │ 1440970 │ ████████████████████████████▋                                   │
│ WINDLESHAM           │ SURREY HEATH           │   206 │ 1381865 │ ███████████████████████████▋                                    │
│ LONDON               │ RICHMOND UPON THAMES   │  1533 │ 1276968 │ █████████████████████████▌                                      │
│ BARNET               │ ENFIELD                │   321 │ 1272131 │ █████████████████████████▍                                      │
│ COBHAM               │ ELMBRIDGE              │   805 │ 1253521 │ █████████████████████████                                       │
│ LONDON               │ ISLINGTON              │  6205 │ 1246671 │ ████████████████████████▊                                       │
│ BEACONSFIELD         │ BUCKINGHAMSHIRE        │   756 │ 1221156 │ ████████████████████████▍                                       │
│ LONDON               │ HOUNSLOW               │  1454 │ 1142654 │ ██████████████████████▋                                         │
│ RICHMOND             │ RICHMOND UPON THAMES   │  1866 │ 1127853 │ ██████████████████████▌                                         │
│ BURFORD              │ WEST OXFORDSHIRE       │   201 │ 1114455 │ ██████████████████████▎                                         │
│ LONDON               │ TOWER HAMLETS          │ 11060 │ 1114239 │ ██████████████████████▎                                         │
│ ASCOT                │ WINDSOR AND MAIDENHEAD │   878 │ 1097645 │ █████████████████████▊                                          │
│ KINGSTON UPON THAMES │ RICHMOND UPON THAMES   │   171 │ 1091598 │ █████████████████████▋                                          │
│ LONDON               │ HAMMERSMITH AND FULHAM │  6934 │ 1071048 │ █████████████████████▍                                          │
│ RADLETT              │ HERTSMERE              │   571 │ 1068349 │ █████████████████████▎                                          │
│ FARNHAM              │ HART                   │   110 │ 1059406 │ █████████████████████▏                                          │
│ LEATHERHEAD          │ GUILDFORD              │   399 │ 1046311 │ ████████████████████▊                                           │
│ WEYBRIDGE            │ ELMBRIDGE              │  1443 │ 1042370 │ ████████████████████▋                                           │
│ SURBITON             │ ELMBRIDGE              │   210 │ 1020546 │ ████████████████████▍                                           │
│ SALCOMBE             │ SOUTH HAMS             │   233 │ 1018012 │ ████████████████████▎                                           │
│ ESHER                │ ELMBRIDGE              │  1079 │ 1015143 │ ████████████████████▎                                           │
│ CHALFONT ST GILES    │ BUCKINGHAMSHIRE        │   326 │ 1004095 │ ████████████████████                                            │
│ FARNHAM              │ EAST HAMPSHIRE         │   118 │ 1001203 │ ████████████████████                                            │
│ GERRARDS CROSS       │ BUCKINGHAMSHIRE        │   932 │  998314 │ ███████████████████▊                                            │
│ BROCKENHURST         │ NEW FOREST             │   251 │  957445 │ ███████████████████▏                                            │
│ EAST MOLESEY         │ ELMBRIDGE              │   415 │  952550 │ ███████████████████                                             │
│ PETERSFIELD          │ CHICHESTER             │   100 │  950453 │ ███████████████████                                             │
│ GUILDFORD            │ WAVERLEY               │   287 │  945821 │ ██████████████████▊                                             │
│ LONDON               │ MERTON                 │  4919 │  936864 │ ██████████████████▋                                             │
│ SUTTON COLDFIELD     │ LICHFIELD              │   120 │  907375 │ ██████████████████▏                                             │
│ LONDON               │ WANDSWORTH             │ 14925 │  898653 │ █████████████████▊                                              │
│ HARPENDEN            │ ST ALBANS              │  1430 │  895874 │ █████████████████▊                                              │
│ HENLEY-ON-THAMES     │ SOUTH OXFORDSHIRE      │  1177 │  887385 │ █████████████████▋                                              │
│ LONDON               │ SOUTHWARK              │  8704 │  883531 │ █████████████████▋                                              │
│ OXFORD               │ SOUTH OXFORDSHIRE      │   738 │  877209 │ █████████████████▌                                              │
│ BILLINGSHURST        │ CHICHESTER             │   285 │  871811 │ █████████████████▍                                              │
│ POTTERS BAR          │ WELWYN HATFIELD        │   348 │  858915 │ █████████████████▏                                              │
│ INGATESTONE          │ CHELMSFORD             │   133 │  848976 │ ████████████████▊                                               │
│ KINGSTON UPON THAMES │ KINGSTON UPON THAMES   │  2043 │  847725 │ ████████████████▊                                               │
│ EAST GRINSTEAD       │ TANDRIDGE              │   124 │  834471 │ ████████████████▋                                               │
│ BELVEDERE            │ BEXLEY                 │   743 │  830758 │ ████████████████▌                                               │
│ LLANGOLLEN           │ WREXHAM                │   148 │  830385 │ ████████████████▌                                               │
│ LONDON               │ HACKNEY                │  7552 │  827224 │ ████████████████▌                                               │
│ LONDON               │ EALING                 │  6266 │  815602 │ ████████████████▎                                               │
└──────────────────────┴────────────────────────┴───────┴─────────┴─────────────────────────────────────────────────────────────────┘

50 rows in set. Elapsed: 0.041 sec. Processed 27.73 million rows, 77.44 MB (669.97 million rows/s., 1.87 GB/s.)
```

Построим следующий аналитический запрос
```clickhouse
SELECT
    toYear(date) AS year,
    round(avg(price)),
    bar(sum(price), 0, 1000000, 80)
FROM uk_price_paid
GROUP BY
    year,
    district,
    town
ORDER BY year ASC

77419 rows in set. Elapsed: 0.532 sec. Processed 27.73 million rows, 293.55 MB (52.15 million rows/s., 551.92 MB/s.) 
```

Создадим материализованное представление для ускорения запросов
```clickhouse
CREATE MATERIALIZED VIEW year_district_town_mv
ENGINE = SummingMergeTree
ORDER BY (date, district, town)
POPULATE
AS SELECT
    toYear(date) AS date,
    district,
    town,
    avg(price) AS avg_price,
    sum(price) AS sum_price,
    count() AS count
FROM uk_price_paid
GROUP BY
    date,
    district,
    town
```

Проверим корость запроса после оптимизации
```clickhouse
SELECT
    date,
    round(avg_price),
    bar(sum_price, 0, 1000000, 80)
FROM year_district_town_mv
ORDER BY date ASC

77419 rows in set. Elapsed: 0.095 sec. Processed 77.42 thousand rows, 1.39 MB (819.04 thousand rows/s., 14.74 MB/s.)
```

Результат тот же, но стал в 5 раз быстрее работать


## Настройка кластера
Поднимем контейнеры с кластером
```bash
docker-compose --file docker-compose-cluster.yml up -d
```

Аналогичныйм образом как и в случае одной базы создадим базу и таблицы

Зайдём в контейнер клиента clickhouse
```bash
docker exec -it clickhouse-cluster-client bash
```

Создадим тестовую базу на каджом сервере
```bash
clickhouse-client --host clickhouse-01 --query "CREATE DATABASE IF NOT EXISTS tutorial"
clickhouse-client --host clickhouse-02 --query "CREATE DATABASE IF NOT EXISTS tutorial"
clickhouse-client --host clickhouse-03 --query "CREATE DATABASE IF NOT EXISTS tutorial"
clickhouse-client --host clickhouse-04 --query "CREATE DATABASE IF NOT EXISTS tutorial"
clickhouse-client --host clickhouse-05 --query "CREATE DATABASE IF NOT EXISTS tutorial"
clickhouse-client --host clickhouse-06 --query "CREATE DATABASE IF NOT EXISTS tutorial"
```

Зайдём в clickhouse-cluster-client
```bash
clickhouse-client --host clickhouse-01
```

И создадим таблицу
```clickhouse
CREATE TABLE tutorial.hits_v1 ON CLUSTER cluster_1
(
    `WatchID` UInt64,
    `JavaEnable` UInt8,
    `Title` String,
    `GoodEvent` Int16,
    `EventTime` DateTime,
    `EventDate` Date,
    `CounterID` UInt32,
    `ClientIP` UInt32,
    `ClientIP6` FixedString(16),
    `RegionID` UInt32,
    `UserID` UInt64,
    `CounterClass` Int8,
    `OS` UInt8,
    `UserAgent` UInt8,
    `URL` String,
    `Referer` String,
    `URLDomain` String,
    `RefererDomain` String,
    `Refresh` UInt8,
    `IsRobot` UInt8,
    `RefererCategories` Array(UInt16),
    `URLCategories` Array(UInt16),
    `URLRegions` Array(UInt32),
    `RefererRegions` Array(UInt32),
    `ResolutionWidth` UInt16,
    `ResolutionHeight` UInt16,
    `ResolutionDepth` UInt8,
    `FlashMajor` UInt8,
    `FlashMinor` UInt8,
    `FlashMinor2` String,
    `NetMajor` UInt8,
    `NetMinor` UInt8,
    `UserAgentMajor` UInt16,
    `UserAgentMinor` FixedString(2),
    `CookieEnable` UInt8,
    `JavascriptEnable` UInt8,
    `IsMobile` UInt8,
    `MobilePhone` UInt8,
    `MobilePhoneModel` String,
    `Params` String,
    `IPNetworkID` UInt32,
    `TraficSourceID` Int8,
    `SearchEngineID` UInt16,
    `SearchPhrase` String,
    `AdvEngineID` UInt8,
    `IsArtifical` UInt8,
    `WindowClientWidth` UInt16,
    `WindowClientHeight` UInt16,
    `ClientTimeZone` Int16,
    `ClientEventTime` DateTime,
    `SilverlightVersion1` UInt8,
    `SilverlightVersion2` UInt8,
    `SilverlightVersion3` UInt32,
    `SilverlightVersion4` UInt16,
    `PageCharset` String,
    `CodeVersion` UInt32,
    `IsLink` UInt8,
    `IsDownload` UInt8,
    `IsNotBounce` UInt8,
    `FUniqID` UInt64,
    `HID` UInt32,
    `IsOldCounter` UInt8,
    `IsEvent` UInt8,
    `IsParameter` UInt8,
    `DontCountHits` UInt8,
    `WithHash` UInt8,
    `HitColor` FixedString(1),
    `UTCEventTime` DateTime,
    `Age` UInt8,
    `Sex` UInt8,
    `Income` UInt8,
    `Interests` UInt16,
    `Robotness` UInt8,
    `GeneralInterests` Array(UInt16),
    `RemoteIP` UInt32,
    `RemoteIP6` FixedString(16),
    `WindowName` Int32,
    `OpenerName` Int32,
    `HistoryLength` Int16,
    `BrowserLanguage` FixedString(2),
    `BrowserCountry` FixedString(2),
    `SocialNetwork` String,
    `SocialAction` String,
    `HTTPError` UInt16,
    `SendTiming` Int32,
    `DNSTiming` Int32,
    `ConnectTiming` Int32,
    `ResponseStartTiming` Int32,
    `ResponseEndTiming` Int32,
    `FetchTiming` Int32,
    `RedirectTiming` Int32,
    `DOMInteractiveTiming` Int32,
    `DOMContentLoadedTiming` Int32,
    `DOMCompleteTiming` Int32,
    `LoadEventStartTiming` Int32,
    `LoadEventEndTiming` Int32,
    `NSToDOMContentLoadedTiming` Int32,
    `FirstPaintTiming` Int32,
    `RedirectCount` Int8,
    `SocialSourceNetworkID` UInt8,
    `SocialSourcePage` String,
    `ParamPrice` Int64,
    `ParamOrderID` String,
    `ParamCurrency` FixedString(3),
    `ParamCurrencyID` UInt16,
    `GoalsReached` Array(UInt32),
    `OpenstatServiceName` String,
    `OpenstatCampaignID` String,
    `OpenstatAdID` String,
    `OpenstatSourceID` String,
    `UTMSource` String,
    `UTMMedium` String,
    `UTMCampaign` String,
    `UTMContent` String,
    `UTMTerm` String,
    `FromTag` String,
    `HasGCLID` UInt8,
    `RefererHash` UInt64,
    `URLHash` UInt64,
    `CLID` UInt32,
    `YCLID` UInt64,
    `ShareService` String,
    `ShareURL` String,
    `ShareTitle` String,
    `ParsedParams` Nested(
        Key1 String,
        Key2 String,
        Key3 String,
        Key4 String,
        Key5 String,
        ValueDouble Float64),
    `IslandID` FixedString(16),
    `RequestNum` UInt32,
    `RequestTry` UInt8
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(EventDate)
ORDER BY (CounterID, EventDate, intHash32(UserID))
SAMPLE BY intHash32(UserID)
SETTINGS index_granularity = 8192
```

Создадим распределённую таблицу
```clickhouse
CREATE TABLE tutorial.hits_v1_distributed ON CLUSTER cluster_1 AS tutorial.hits_v1
ENGINE = Distributed(cluster_1, tutorial, hits_v1, rand());
```

Заимпортим данные
```bash
clickhouse-client --host clickhouse-01 --query "INSERT INTO tutorial.hits_v1 FORMAT TSV" --max_insert_block_size=100000 < /tmp/hits_v1.tsv
```

Зайдём обратно в clickhouse-client для того, чтобы посмотреть что получилось
```bash
clickhouse-client --host clickhouse-01
```

Проверим сколько записей в исходной таблице
```clickhouse
SELECT count(*) FROM tutorial.hits_v1
┌─count()─┐
│ 8873898 │
└─────────┘

1 rows in set. Elapsed: 0.006 sec. 
```

Проверим сколько записей в распределённой таблице
```clickhouse
SELECT count(*) FROM tutorial.hits_v1_distributed
┌─count()─┐
│ 8873898 │
└─────────┘

1 rows in set. Elapsed: 0.008 sec.
```

Количество записей одинаковое

Выполним запрос сложнее
```clickhouse
SELECT count(*) FROM tutorial.hits_v1 GROUP BY BrowserCountry, EventDate, URLDomain ORDER BY URLDomain, BrowserCountry, EventDate

246635 rows in set. Elapsed: 0.270 sec. Processed 8.87 million rows, 257.25 MB (32.86 million rows/s., 952.63 MB/s.) 
```

Выполним его оп распределённой таблице
```clickhouse
SELECT count(*) FROM tutorial.hits_v1_distributed GROUP BY BrowserCountry, EventDate, URLDomain ORDER BY URLDomain, BrowserCountry, EventDate

246635 rows in set. Elapsed: 0.346 sec. Processed 8.87 million rows, 257.25 MB (25.62 million rows/s., 742.74 MB/s.)
```

Ответ одинаковый, однако время выполнения по распределённой таблице больше.
Похоже, что это связано с тем, что все сервера кластера находятся на одной ноде, а для того, чтобы ответ сагрегировать используются сетевые пакеты, что
в случае размещения всех серверов на одной ноде не даёт выйгрыш, так как ресурсы у всех серверов общие.