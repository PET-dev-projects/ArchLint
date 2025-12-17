# ArchLint – проверка архитектуры на Go

ArchLint поставляет готовый к продакшену модуль на Go для проверки архитектур сервисов, описанных в YAML, с переиспользуемой библиотекой и опциональным CLI.

## Возможности
- Загружайте описания архитектуры из YAML (`examples/payments.yaml`).
- Проверяйте движок на примере потокового сервиса масштаба Spotify (`examples/music_streaming.yaml`), который задействует каждое правило (диаграмма в `docs/music_streaming.md`).
- Структурная валидация с подробными находками (версия схемы, дубли контейнеров, неизвестные связи).
- Встроенные правила со стабильными идентификаторами:
  - `ARCH-ACYCLIC` – поиск циклов зависимостей.
  - `ARCH-CRUD` – контроль доступа к БД (CRUD/repo/relay паттерны).
  - `ARCH-ACL` – доступ к внешним системам только через ACL-контейнеры.
  - `ARCH-BOUNDARIES` – коэффициенты сплочённости/сцепления границ, настраиваемые пороги.
  - `ARCH-EXTERNAL-PROTOCOL` – допустимые протоколы/транспорты при интеграции с внешними системами.
  - `ARCH-DB-ISOLATION` – базы данных пассивны (нет исходящих вызовов, предупреждение о неиспользуемых БД).
- Настройка правил через YAML (`configs/rules.yaml`): включайте/отключайте проверки и задавайте параметры для каждого правила.
- Детерминированный формат находок для дальнейшей автоматизации.
- Тонкая CLI-обёртка (`cmd/archlint`) для CI.
- Покрытие `go test` с фикстурами в `testdata/` и текстовыми отчётами из `pkg/report`.

## YAML-модель
Исходником архитектуры служит YAML (с версией схемы). Минимальная структура:

```yaml
version: 1
meta:
  owner: platform-team
boundaries:
  - name: Payments
    description: Context information (optional)
    tags: [core]
    containers:
      - name: payments-api
        type: service            # service | database | external
        tags: [acl]
      - name: payments-repo
        type: service
        tags: [repo]
      - name: payments-db
        type: database
        technology: postgres
    boundaries: []               # допускаются вложенные контексты
    relations:
      - from: payments-api
        to: payments-repo
        kind: sync               # sync | async | db
        description: orchestrates domain logic
      - from: payments-repo
        to: payments-db
        kind: db
externals:
  - name: antifraud
    type: external
    description: third-party system
```

Все сущности живут внутри `boundaries`; `externals` — опциональные помощники с `type: external`. У каждой связи есть путь (`Path`), чтобы находки ссылались на `boundaries[0].relations[1]` и т.д.

## Контракт находок

```go
type Severity string // "error" | "warn" | "info"
type Finding struct {
    RuleID   string
    Severity Severity
    Message  string
    Path     string // JSONPointer-подобный путь внутри YAML
    Meta     map[string]any
}
```

Находки сортируются по `RuleID`, затем по `Path`, что гарантирует повторяемость.

## Использование библиотеки

```go
package main

import (
    "os"

    "github.com/PET-dev-projects/ArchLint/pkg/archlint"
    "github.com/PET-dev-projects/ArchLint/pkg/engine"
)

func main() {
    f, _ := os.Open("examples/payments.yaml")
    model, _ := archlint.LoadModelFromYAML(f)

    findings := append(
        archlint.ValidateModel(model),
        archlint.RunAll(model, engine.Options{}),
    )

    // findings — это []types.Finding, готовый к логированию/CI.
}
```

`engine.Options` позволяет включить подмножество правил и передать каждому JSON-подобную конфигурацию (например, ослабить пороги границ или добавить дополнительные CRUD-теги).

## CLI

```
go install ./cmd/archlint

archlint check -f examples/payments.yaml \
  --format text \   # text или json
  --fail-on error   # порог (error|warn|info|none)
```

Код возврата отличен от нуля, если найдена хотя бы одна находка с выбранной серьёзностью (`--fail-on`, по умолчанию `error`).

Посмотрите `docs/examples.md` для дополнительных сценариев, демонстрирующих каждое встроенное правило на готовых фикстурах.

### Файл конфигурации правил

Используйте `--config configs/rules.yaml` (или свой YAML), чтобы управлять набором правил и их параметрами:

```yaml
rules:
  - id: ARCH-ACYCLIC
    enabled: true
  - id: ARCH-BOUNDARIES
    enabled: true
    config:
      minInternalToCrossRatio: 1.5
      maxCrossRelations: 5
```

Каждая запись привязана к идентификатору правила. Уберите её или выставьте `enabled: false`, чтобы пропустить правило. Любой объект `config` передаётся декодеру соответствующего правила. Если конфигурация не указана, запускаются все встроенные проверки со значениями по умолчанию.

## Тесты и фикстуры
- `testdata/*.yaml` — наследие PlantUML-сценариев: циклы, CRUD-нарушения, ACL, слабые границы и т.д.
- `examples/music_streaming.yaml` описывает крупную потоковую платформу с несколькими границами, внешними системами и потоками данных — используйте её для проверки сложных ландшафтов.
- `go test ./...` покрывает валидацию моделей, каждое правило и оркестрацию движка. Если окружение ограничивает домашний каталог, задайте `GOCACHE=$(pwd)/.cache`.

## Расширение
- Добавляйте правила в `pkg/checks` и регистрируйте их в `pkg/checks/registry.go`.
- Используйте `pkg/report` для форматирования вывода в текст/JSON.
- `examples/payments.yaml` можно взять за основу при миграции со старых PlantUML-файлов.
- Подробности по внедрению библиотеки (API, настройка правил, расширение) — в `docs/library.md`.
