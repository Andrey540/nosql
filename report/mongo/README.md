# Отчёт по домашнему заданию №2

Скачаем dataset вопросов и сохраним как QUESTIONS.json в папку mongo
https://www.reddit.com/r/datasets/comments/1uyd0t/200000_jeopardy_questions_in_a_json_file/

Поднимем контейнер
```bash
docker-compose --file mongo/docker-compose-mongo.yml up
```

Войдём в контейнер
```bash
docker exec -it mongo bash
```

Заимпортим коллекцию
```bash
mongoimport --file /tmp/QUESTIONS.json --jsonArray --host mongo --db questionsDB -c questions -u otus -p 1234 --authenticationDatabase admin
```

Подключимся к базе
```bash
mongosh mongodb://otus:1234@mongo
use questionsDB
```

Выполним несколко селектов

```bash
db.questions.countDocuments()
216930
```

```bash
db.questions.countDocuments({show_number:{$gt:'5957'}})
23258
```

```bash
db.questions.find({show_number:{$gt:'5957'}}).limit(3)
[
  {
    _id: ObjectId("636d3ea3e70b6fb6b3b415da"),
    category: 'TIMELESS TV',
    air_date: '2010-12-07',
    question: "'September 2010 brought the 45th edition of this comedian's telethon'",
    value: '$200',
    answer: 'Jerry Lewis',
    round: 'Jeopardy!',
    show_number: '6037'
  },
  {
    _id: ObjectId("636d3ea3e70b6fb6b3b415db"),
    category: "LET'S HIT IT",
    air_date: '2010-12-07',
    question: `'Hit this paper mache container, Spanish for "jug", if you want candy and small gifts'`,
    value: '$200',
    answer: 'pinata',
    round: 'Jeopardy!',
    show_number: '6037'
  },
  {
    _id: ObjectId("636d3ea3e70b6fb6b3b415dc"),
    category: 'LOST IN SPACE',
    air_date: '2010-12-07',
    question: "'While making repairs on the Intl. Space Station, Scott Parazynski lost a needle-nose pair of these'",
    value: '$200',
    answer: 'pliers',
    round: 'Jeopardy!',
    show_number: '6037'
  }
]
```

```bash
db.questions.aggregate({
  $group: { _id: "$value", count: { $sum: 1 } }
}, {
  $sort: { count: -1 }
})
[
  { _id: '$400', count: 42244 },
  { _id: '$800', count: 31860 },
  { _id: '$200', count: 30455 },
  { _id: '$600', count: 20377 },
  { _id: '$1000', count: 19539 },
  { _id: '$1200', count: 11331 },
  { _id: '$2000', count: 11243 },
  { _id: '$1600', count: 10801 },
  { _id: '$100', count: 9029 },
  { _id: '$500', count: 9016 },
  { _id: '$300', count: 8663 },
  { _id: null, count: 3634 },
  { _id: '$1,000', count: 2101 },
  { _id: '$2,000', count: 1586 },
  { _id: '$3,000', count: 769 },
  { _id: '$1,500', count: 546 },
  { _id: '$1,200', count: 441 },
  { _id: '$4,000', count: 349 },
  { _id: '$1,600', count: 239 },
  { _id: '$2,500', count: 232 }
```

```bash
db.questions.aggregate({
  $match: { show_number: { $gt: '5957' } }
}, {
  $sort: { air_date: 1 }
}, {
  $limit: 3
})
[
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494c"),
    category: 'TRIVIA',
    air_date: '1984-09-20',
    question: "'They were formerly called the Sandwich Islands'",
    value: '$100',
    answer: 'Hawaii',
    round: 'Jeopardy!',
    show_number: '9'
  },
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494a"),
    category: 'RACY LADIES',
    air_date: '1984-09-20',
    question: "'Driving with a broken wrist, Janet Guthrie was the 1st woman to compete in this race'",
    value: '$100',
    answer: 'the Indianapolis 500',
    round: 'Jeopardy!',
    show_number: '9'
  },
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b74949"),
    category: 'FIRST LADIES',
    air_date: '1984-09-20',
    question: `'She talked to Arnold's class about drugs on TV's "Diff'rent Strokes"'`,
    value: '$100',
    answer: 'Nancy Reagan',
    round: 'Jeopardy!',
    show_number: '9'
  }
]
```

```bash
db.questions.aggregate({
  $limit: 3
}, {
  $sort: { air_date: 1 }
}, {
  $match: { show_number: { $gt: '5957' } }
})
```

```bash
db.questions.find({show_number:{$gt:'5957'}}).limit(3).sort({air_date:1})
[
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494b"),
    category: 'ADDRESSES',
    air_date: '1984-09-20',
    question: "'The prime minister of England lives here'",
    value: '$100',
    answer: '10 Downing Street',
    round: 'Jeopardy!',
    show_number: '9'
  },
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494a"),
    category: 'RACY LADIES',
    air_date: '1984-09-20',
    question: "'Driving with a broken wrist, Janet Guthrie was the 1st woman to compete in this race'",
    value: '$100',
    answer: 'the Indianapolis 500',
    round: 'Jeopardy!',
    show_number: '9'
  },
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b74949"),
    category: 'FIRST LADIES',
    air_date: '1984-09-20',
    question: `'She talked to Arnold's class about drugs on TV's "Diff'rent Strokes"'`,
    value: '$100',
    answer: 'Nancy Reagan',
    round: 'Jeopardy!',
    show_number: '9'
  }
]
```

```bash
db.questions.find({_id:ObjectId("636d3ea5e70b6fb6b3b7494b")})
[
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494b"),
    category: 'ADDRESSES',
    air_date: '1984-09-20',
    question: "'The prime minister of England lives here'",
    value: '$100',
    answer: '10 Downing Street',
    round: 'Jeopardy!',
    show_number: '9'
  }
]

db.questions.updateOne(
   { _id: ObjectId("636d3ea5e70b6fb6b3b7494b") },
   {
     $set: { air_date: '2022-09-23', value: '$1000' },
     $currentDate: { lastModified: true }
   }
)  
  
db.questions.findOne({_id:ObjectId("636d3ea5e70b6fb6b3b7494b")})
[
  {
    _id: ObjectId("636d3ea5e70b6fb6b3b7494b"),
    category: 'ADDRESSES',
    air_date: '2022-09-23',
    question: "'The prime minister of England lives here'",
    value: '$1000',
    answer: '10 Downing Street',
    round: 'Jeopardy!',
    show_number: '9',
    lastModified: ISODate("2022-11-10T19:47:21.434Z")
  }
]
```

И апдейтов
```bash
db.questions.countDocuments({category:'ADDRESSES'})
10
```

```bash
db.questions.updateMany(
   { category: 'ADDRESSES' },
   {
     $set: { category: 'HOUSES' },
     $currentDate: { lastModified: true }
   }
)
```

```bash
db.questions.countDocuments({category:'ADDRESSES'})
0
```

```bash
db.questions.countDocuments({category:'HOUSES'})
10
```

Проверим скоость работы запроса с фильтрацией и сортировкой
```bash
db.questions.aggregate({
  $match: { show_number: { $gt: '5957' } }
}, {
  $sort: { show_number: 1, air_date: 1 }
}, {
  $limit: 3
}).explain("executionStats").executionStats.executionTimeMillis
113
```

Создадим индекс по полям show_number и air_date
```bash
db.questions.createIndex({ show_number: 1, air_date: 1 })
```

Сравним полученное время
```bash
db.questions.aggregate({
  $match: { show_number: { $gt: '5957' } }
}, {
  $sort: { show_number: 1, air_date: 1 }
}, {
  $limit: 3
}).explain("executionStats").executionStats.executionTimeMillis
1
```

По данному запросу время сократилось в 100 раз
