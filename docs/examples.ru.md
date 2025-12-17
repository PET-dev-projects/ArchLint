# CLI над встроенными фикстурами

Ниже собраны сценарии, чтобы быстро проверить бинарник end-to-end. Предполагается, что вы находитесь в корне репозитория и уже собрали `./archlint` (`go build ./cmd/archlint`).

## 1. Позитивная архитектура

```
./archlint check -f examples/payments.yaml --format text --fail-on error
```

Ожидаемый вывод:

```
No findings
```

Так вы проверяете загрузку, валидацию и все проверки по умолчанию на YAML-образце из репозитория.

## 2. Циклы зависимостей (ARCH-ACYCLIC)

```
./archlint check -f testdata/arch_cycle.yaml --format text --fail-on error
```

Пример строки:

```
ARCH-ACYCLIC	error	boundaries[0].relations[2]	cycle detected: [repo api repo]
```

## 3. CRUD-правила (ARCH-CRUD)

```
./archlint check -f testdata/arch_crud_violation.yaml --format text --fail-on error
```

Образец находок:

```
ARCH-CRUD	error	boundaries[0].relations[0]	container api must declare one of [crud repo relay] to access databases
ARCH-CRUD	error	boundaries[0].relations[2]	container cache is restricted to database relations
```

## 4. ACL (ARCH-ACL)

```
./archlint check -f testdata/arch_acl_violation.yaml --format text --fail-on error
```

Ищите строки вида:

```
ARCH-ACL	error	boundaries[0].relations[2]	container api must declare one of [acl] to talk to external audit
```

## 5. Сплочённость границ (ARCH-BOUNDARIES)

```
./archlint check -f testdata/arch_boundary_weak.yaml --format text --fail-on warn
```

Правило выдаёт предупреждения, поэтому установите `--fail-on warn`, чтобы получить ненулевой код возврата. Пример:

```
ARCH-BOUNDARIES	warn	boundaries[0]	boundary Reporting cohesion/coupling ratio 0.50 below minimum 1.00	{"internal":1,"cross":2,"ratio":0.5}
```

## 6. Протоколы для внешних интеграций (ARCH-EXTERNAL-PROTOCOL)

```
./archlint check -f testdata/arch_external_protocol.yaml --format text --fail-on error
```

Ожидайте:

```
ARCH-EXTERNAL-PROTOCOL	error	boundaries[0].relations[0].protocol	protocol "http://public-gateway/antifraud" for external antifraud is not allowed	{"allowedPrefixes":["https://gateway.","kafka://"]}
ARCH-EXTERNAL-PROTOCOL	error	boundaries[0].relations[1]	relation from payments-repo to external antifraud must define protocol
```

## 7. Изоляция баз данных (ARCH-DB-ISOLATION)

```
./archlint check -f testdata/arch_db_isolation.yaml --format text --fail-on warn
```

Первая находка запрещает инициировать связи из базы данных, вторая — предупреждает о неиспользуемых БД.

## 8. Регрессионный прогон

```
GOCACHE=$(pwd)/.cache go test ./...
```

## 9. Пользовательские профили правил

```
./archlint check -f testdata/arch_valid.yaml --config configs/rules.yaml --format text
```

Команда читает `configs/rules.yaml`, включает только перечисленные правила и пробрасывает указанную конфигурацию.

## 10. Сценарий потокового сервиса масштаба Spotify

```
./archlint check -f examples/music_streaming.yaml --format text --fail-on error
```

Ожидаемый вывод:

```
No findings
```

YAML покрывает пять границ (playback, catalog, personalization, monetization, platform) и несколько внешних систем. Запустив этот сценарий, вы убедитесь, что движок справляется с крупными продакшн-топологиями и что все правила работают вместе.

Команда `GOCACHE=$(pwd)/.cache go test ./...` выше запускает юнит-тесты Go, которые повторяют PlantUML-сценарии. У каждого правила есть покрытие + проверка движка/валидации, поэтому сбои почти всегда сигнализируют о регрессии.
