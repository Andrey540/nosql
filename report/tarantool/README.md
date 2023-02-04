# Отчёт по домашнему заданию №13

Перейдём в директорию tarantool
```bash
cd tarantool
```

Поднимем контейнеры
```bash
docker-compose up -d
```

Зайдём в контейнер
```bash
docker exec -it tarantool console
```

Установим необходимые модули
```bash
uuid = require('uuid')
fiber = require('fiber')
http_client = require('http.client')
expirationd = require('expirationd')
```

Создадим спейс
```bash
space = box.schema.space.create('billing', {if_not_exists = true})
```

Опишем схему
```bash
space:format({
{name = 'id', type = 'uuid'},
{name = 'amount', type = 'unsigned'}
})
```

Создадим индекс
```bash
space:create_index('primary', {
type = 'tree',
parts = {'id'},
if_not_exists = true
})
```

Создадим пользователя
```bash
userId = uuid.new()
```

Создадим функцию создания счёта
```bash
function createAccount(userId, amount)
    box.space.billing:insert{userId, amount}
end
```

Создадим функцию пополенеия счёта
```bash
function topUpAccount(userId, amount)
    box.space.billing:update(userId, { {"+", 2, amount} })
end
```

Создадим функцию списания денег со счёта
```bash
function spendMoney(userId, amount)
    while true do
        tuple = box.space.billing:get(userId)
        if ((tuple ~= nil) and (tuple.amount > 0)) then
            value = math.min(tuple.amount, amount)
            box.space.billing:update(userId, { {"-", 2, amount} })
            fiber.sleep(1)
        else
            return
        end
    end
end
```

Создадим счёт
```bash
createAccount(userId, 100)
```

Пополним счёт
```bash
topUpAccount(userId, 100)
```

Проверим состояние счёта
```bash
box.space.billing:get(userId)
---
- [d449cdf7-4950-43f8-8fb7-161bd74d0eb2, 200]
...
```

Спишем деньги с шагом в 10
```bash
spendMoney(userId, 10)
```

Проверим состояние счёта
```bash
box.space.billing:get(userId)
---
- [d449cdf7-4950-43f8-8fb7-161bd74d0eb2, 0]
...
```

Пополним счёт
```bash
topUpAccount(userId, 100)
```

Создадим задачу отслеживания состояния счёта
```bash
job_name = "checkBilling"

function isUserBillingEmpty(args, tuple)
    return tuple[2] == 0
end

function alertUserEmpty(space, args, tuple) 
    response = http_client.get('https://www.tarantool.io/en/?userId=' .. tostring(tuple[1]))
    
    if ((response ~= nil) and (response.status == 200)) then
        print('user ' .. tostring(tuple[1]) .. ' processed success!!')
    else
        print('user ' .. tostring(tuple[1]) .. ' processed with error!!')
    end 
    io.flush()
    box.space[space]:delete{tuple[1]}
end

expirationd.start(job_name, space.id, isUserBillingEmpty, {
    process_expired_tuple = alertUserEmpty,
    args = nil,
    tuples_per_iteration = 1,
    full_scan_time = 3600
})
```

Запустим функцию списания денег
```bash
spendMoney(userId, 10)

user d449cdf7-4950-43f8-8fb7-161bd74d0eb2 processed success!!
```
Получили сообщение о том, что пользователь обработан