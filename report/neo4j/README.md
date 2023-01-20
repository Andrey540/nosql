# Отчёт по домашнему заданию №11

## Сравнение построения модели в neo4j и SQL базах

За основу взял модель из лекции.
У нас есть узлы, которые могут быть типов - Actor, Director, Movie
И рёбра, которые могут быть типов - CREATED, PLAYED_IN
Кроме того у узлов и рёбер могут быть атрибуты

Для хранения данных в MySQL можно воспользоваться примерной схемой

```sql
CREATE TABLE node
(
    id INTEGER NOT NULL AUTO_INCREMENT,
    type VARCHAR(128) NOT NULL,
    data JSON,
    PRIMARY KEY (`id`)
)
```
Поле тип может принимать значения (Actor, Director, Movie), поле data - json, содержит атрибуты

```sql
CREATE TABLE edge
(
    id INTEGER NOT NULL AUTO_INCREMENT,
    type VARCHAR(128) NOT NULL,
    data JSON,
    src_node_id INTEGER NOT NULL,
    dsc_node_id INTEGER NOT NULL,
    PRIMARY KEY (`id`)
)
```
Поле тип может принимать значения (CREATED, PLAYED_IN), поле data - json, содержит атрибуты.
Поля src_node_id и dsc_node_id - узлы, которые ребро соединяет.

Заполним узлы
```sql
INSERT INTO node (type, data) VALUES
('Director', '{"name": "Joel Coen"}'),
('Movie', '{"title":"Blood Simple", "year":"1983"}'),
('Actor', '{"name": "Frances McDormand"}'),
('Director', '{"name": "Ethan Coen", "born":"1957"}'),
('Director', '{"name": "Martin McDonagh"}'),
('Movie', '{"title": "Three Billboards Outside Ebbing, Missouri"}'),
('Movie', '{"title": "Venom"}'),
('Actor', '{"name": "Woody Harrelson"}'),
('Actor', '{"name": "Tom Hardy"}'),
('Movie', '{"title": "Inception"}'),
('Actor', '{"name": "Leonardo DiCaprio"}'),
('Actor', '{"name": "Marion Cotillard"}'),
('Movie', '{"title": "The Dark Knight Rises"}'),
('Movie', '{"title": "Batman"}'),
('Director', '{"name": "Christopher Nolan"}'),
('Director', '{"name": "Ruben Fleischer"}')
;
```

И связи между ними
```sql
INSERT INTO edge (type, src_node_id, dsc_node_id) VALUES
('CREATED', 1, 2),
('PLAYED_IN', 3, 2),
('CREATED', 4, 2),
('CREATED', 5, 6),
('PLAYED_IN', 3, 6),
('PLAYED_IN', 8, 7),
('PLAYED_IN', 8, 6),
('PLAYED_IN', 9, 7),
('PLAYED_IN', 11, 10),
('PLAYED_IN', 9, 10),
('PLAYED_IN', 12, 10),
('PLAYED_IN', 12, 13),
('CREATED', 15, 14),
('CREATED', 15, 10),
('CREATED', 16, 7)
;
```

Попробуем построить запрос как для neo4j, аналогичный для sql
```sql
match (venom:Movie {title:'Venom'}) -[*1..3]- (d:Director) return d
```

По сути нужно получить все вершины, которые связаны с Movie {title:'Venom'} не более чем тремя рёбрами с типом Director.
Можно разбить на 3 подзапроса - первый получается вершины на расстоярии 1, второй - 2, третий - 3.

Получился примерный запрос
```sql
(
    SELECT DISTINCT JSON_EXTRACT(n1.data, "$.name") FROM node n
    INNER JOIN edge e1 ON (e1.src_node_id = n.id OR e1.dsc_node_id = n.id) AND e1.src_node_id <> e1.dsc_node_id
    INNER JOIN node n1 ON (n1.id = e1.dsc_node_id OR n1.id = e1.src_node_id) AND n1.type = 'Director'
    WHERE n.type = 'Movie' AND JSON_CONTAINS(n.data, JSON_QUOTE('Venom'),'$.title') = 1
)
UNION
(
    SELECT DISTINCT JSON_EXTRACT(n1.data, "$.name") FROM node n
    INNER JOIN edge e1 ON (e1.src_node_id = n.id OR e1.dsc_node_id = n.id) AND e1.src_node_id <> e1.dsc_node_id
    INNER JOIN edge e2 ON (e2.src_node_id = e1.dsc_node_id OR e2.src_node_id = e1.src_node_id OR e2.dsc_node_id = e1.src_node_id OR e2.dsc_node_id = e1.dsc_node_id) AND e2.src_node_id <> e2.dsc_node_id AND e2.id <> e1.id
    INNER JOIN node n1 ON (n1.id = e2.dsc_node_id OR n1.id = e2.src_node_id) AND n1.type = 'Director'
    WHERE n.type = 'Movie' AND JSON_CONTAINS(n.data, JSON_QUOTE('Venom'),'$.title') = 1
)
UNION
(
    SELECT DISTINCT JSON_EXTRACT(n1.data, "$.name") FROM node n
    INNER JOIN edge e1 ON (e1.src_node_id = n.id OR e1.dsc_node_id = n.id) AND e1.src_node_id <> e1.dsc_node_id
    INNER JOIN edge e2 ON (e2.src_node_id = e1.dsc_node_id OR e2.src_node_id = e1.src_node_id OR e2.dsc_node_id = e1.src_node_id OR e2.dsc_node_id = e1.dsc_node_id) AND e2.src_node_id <> e2.dsc_node_id AND e2.id <> e1.id
    INNER JOIN edge e3 ON (e3.src_node_id = e2.dsc_node_id OR e3.src_node_id = e2.src_node_id OR e3.dsc_node_id = e2.src_node_id OR e3.dsc_node_id = e2.dsc_node_id) AND e3.src_node_id <> e3.dsc_node_id AND e3.id <> e2.id
    INNER JOIN node n1 ON (n1.id = e3.dsc_node_id OR n1.id = e3.src_node_id) AND n1.type = 'Director'
    WHERE n.type = 'Movie' AND JSON_CONTAINS(n.data, JSON_QUOTE('Venom'),'$.title') = 1
)
```

В sql моделирование аналогичной модели графа, значительно сложнее, а также гораздо сложнее составление запроса по связям в глубину.
Простой запрос из neo4j в sql эквиваленте становится сложным, а также сложна его отладка и в нём легко ошибится. MySql не поддерживает рекурсивные выборки,
поэтому универсального механизма для поиска в глубину у него нет. Возможно подобные операции проще с использованием MySql делать следющим образом -
получать все вершины и рёбра и уже в приложении рекурсивно обходить граф, если же количество вершить будет много, то этот подход будет работать не оптимально.
Работа со всязами между объектами реализована значительно проще в neo4j.


## Возможная область применения neo4j

Карта местности с путями. При этом какие-то пееркрёстки или станции (например, метро) могут быть вершиными, а пути или дороги - рёбрами с определённым весом,
при этом рёбра могут иметь напрввление. Наиболее частая задача - нахождение оптимального пути.

Система рекомендаций. Можно хранить в графовой базе данных историю покупок пользоветля. Например, вершинами могут бытиь типы товаров, а рёбрами покупки пользователей.
Тогда при покупке товара категории, можно рекомендовать наиболее популярный товар из этой категории (С этим товаром часто покупают).
Система рекомендаций также может быть в социальной сети, например, рекомендации в знакоместве. Узлами могут быть пользователи, а факт нахождения в друзьях или подписка на человека - это рёбра.
Если несколько пользователей, с кем общается пользователь имеет в друзьях одного и того же человека, то можно рекомендовать его в друзья.

Возможно графовая база данных подойдёт для 3D проектирования и создания чертежей или электрических схем.
Вершины могу хранить координаты, рёбра тип линии и её параметры (прямая, дуга, толщина линии или ещё что-то).