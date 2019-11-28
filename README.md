## "Анти-брутфорс"
#### Usage
 Запуск проекта осуществляется командой `make up`

Запуск интеграционных тестов `make test` 

#### Сli
```Usage:
./abf [command]

Available Commands:

  add         Command for adding ip to the blacklist or whitelist
  bucket      Command resets bucket by login or ip
  delete      The command removes ip from the list
  grpc        Command to start grpc server
  help        Help about any command

Flags:
  -h, --help         help for ./abf
  -n, --net string   ip with mask example 127.0.0.0/24
  
Use "./abf [command] --help" for more information about a command.
```

*Комментарий:*

* Хранилище для bucket-ов реализовано в памяти. При создании bucket-а ему передается ссылка на [канал](https://github.com/ios116/antibruteforce/blob/c61903644a8402a02ac6cf1d876c3c237a05535d/antibruteforce/internal/domain/entities/bucket_entities.go#L57), по которому через timeout(время жизни bucket-а) посылает сообщение [коллектору](https://github.com/ios116/antibruteforce/blob/ed0bba00f85f29acebe7a658a1b8fe58b86146f7/antibruteforce/internal/usecase/bucketusecase/bucket_usecase.go#L102) на удаление из памяти. Коллектор это горутина которая [запускается](https://github.com/ios116/antibruteforce/blob/9605f140d3b5fad2792b0ae1679fee58c06e8c04/antibruteforce/internal/grpcserver/grpc.go#L39) при запуске grpc сервера. У каждого bucket-а есть счетчик (проще говоря это кол-во запросов в ед. времени) при каждом запросе счетчик уменьшается на 1 и при достижении 0 запрос блокируется.  

# ТЗ на сервис "Анти-брутфорс".

## Общее описание проекта

Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.
Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.
Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

## Алгоритм работы

Сервис ограничивает частоту попыток авторизации для различных комбинаций параметров, например:
* не более N = 10 попыток в минуту для данного логина.
* не более M = 100 попыток в минуту для данного пароля (защита от обратного brute-force).
* не более K = 1000 попыток в минуту для данного IP (число большое, т.к. NAT).

Для подсчета и ограничения частоты запросов, можно использовать например алгоритм leaky bucket.
Или иные аналогичные: https://en.wikipedia.org/wiki/Rate_limiting
Причем сервис будет поддерживать множество bucket-ов, по одному на каждый логин/пароль/ip.
Bucket-ы можно хранить 
* в памяти (в таком случае нужно продумать удаление неактивных bucket-ов, что бы избежать утечек памяти)
* во внешнем хранилище (например redis или СУБД, в таком случае нужно продумать производительность)

White/black листы содержат списки адресов сетей, которые обрабатываются более простым способом.
Если входящий ip в whitelist - сервис безусловно разрешает авторизацию (ok=true), если в blacklist - отклоняет (ok=false).

## Архитектура

Микросервис состоит из API, базы данных для хранения настроек и black/white списков.
Опционально - хранилище для bucket-ов.
Сервис должен предоставлять GRPC или REST API.

## Описание методов API

### Попытка авторизации
Запрос:
* login
* password
* ip

Ответ:
* ok (true/false) - сервис должен возвращать ok=true, если считает что запрос нормальный 
                    и ok=false, если считает что происходит bruteforce.

### Сброс bucket
* login
* ip
Должен очистить bucket-ы соответствующие переданным login и ip

### Добавление IP в blacklist
* subnet (ip + mask)

### Удаление IP из blacklist
* subnet (ip + mask)

### Добавление IP в whitelist
* subnet (ip + mask)

### Удаление IP из whitelist
* subnet (ip + mask)

## Конфигурация
Основные параметры конфигурации: N, M, K - лимиты по достижению которых, сервис считает попытку брутфорсом.

## Command-Line интерфейс
Необходимо так же разработать command-line интерфейс для ручного администрирования сервиса.
Через CLI должна быть возможность вызвать сброс бакета и управлять whitelist/blacklist-ом.
CLI работает через GRPC интерфейс.

## Развертывание

Развертывание микросервиса осуществляться командой `docker-compose up` в директории с проектом.

## Тестирование

Рекомендуется выделить модуль обработки одного bucket и протестировать его с помощью unit-тестов.
Так же необходимо написать интеграционные тесты, которые поднимаю сервис и все необходимые базы в docker-compose
и проверяют все вызовы API.
