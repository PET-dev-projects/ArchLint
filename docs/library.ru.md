# Руководство по использованию библиотеки

Документ ориентирован на разработчиков, которые встраивают модуль Go напрямую (без CLI). Здесь описано, как загружать YAML-модели, запускать структурную валидацию, выполнять правила и настраивать поведение.

## 1. Установка модуля

```
go get github.com/NovokshanovE/archlint
```

Импортируйте необходимые пакеты:

```go
import (
    "os"

    "github.com/NovokshanovE/archlint/pkg/archlint"
    "github.com/NovokshanovE/archlint/pkg/engine"
    "github.com/NovokshanovE/archlint/pkg/types"
)
```

## 2. Загрузка и валидация модели

```go
f, err := os.Open("path/to/arch.yaml")
if err != nil {
    return err
}
defer f.Close()

model, err := archlint.LoadModelFromYAML(f)
if err != nil {
    return err
}

structuralFindings := archlint.ValidateModel(model)
```

`LoadModelFromYAML` возвращает типизированную `*model.Architecture`. `ValidateModel` проводит проверку схемы (версия, обязательные поля, дубли, неизвестные ссылки). Эти находки обрабатываются в первую очередь: ошибка уровня `SeverityError` обычно означает, что последующие правила работать не смогут.

## 3. Запуск движка правил

```go
opts := engine.Options{} // по умолчанию включены все встроенные правила
ruleFindings := archlint.RunAll(model, opts)

allFindings := append(structuralFindings, ruleFindings...)
```

Доступные встроенные правила и их назначение:

| Rule ID | Назначение |
|---------|------------|
| `ARCH-ACYCLIC` | Поиск циклов зависимостей между контейнерами. |
| `ARCH-CRUD` | Контроль контрактов CRUD/БД (только помеченные сервисы ходят в БД, ограничения repo-only). |
| `ARCH-ACL` | Внешние интеграции разрешены только через ACL-контейнеры. |
| `ARCH-BOUNDARIES` | Отчёт по сплочённости/сцеплению границ (порог настраивается). |
| `ARCH-EXTERNAL-PROTOCOL` | Требует протоколы с разрешёнными префиксами для внешних вызовов. |
| `ARCH-DB-ISOLATION` | Гарантирует пассивность баз данных и предупреждает о неиспользуемых БД. |

Все правила возвращают `types.Finding` в детерминированном порядке.

## 4. Программная настройка правил

Заполняйте `engine.Options`, чтобы кастомизировать запуск:

```go
opts := engine.Options{
    EnabledRules: []string{"ARCH-ACYCLIC", "ARCH-ACL"},
    RuleConfig: map[string]map[string]any{
        "ARCH-BOUNDARIES": {
            "minInternalToCrossRatio": 2.0,
            "maxCrossRelations":       10,
        },
        "ARCH-EXTERNAL-PROTOCOL": {
            "allowedPrefixes": []string{"https://gateway.", "amqp://"},
        },
    },
}
```

- `EnabledRules` работает как allowlist. Оставьте `nil`, чтобы выполнить все зарегистрированные правила.
- `RuleConfig` пробрасывает произвольные JSON-подобные объекты в декодер конкретного правила (см. `pkg/checks/*` для списка полей).

## 5. Загрузка конфигурации правил из YAML

Вместо ручной сборки `engine.Options`, загрузите их из YAML через `pkg/config`:

```go
import "github.com/NovokshanovE/archlint/pkg/config"

opts, err := config.LoadOptionsFromFile("configs/rules.yaml")
if err != nil {
    return err
}
findings := archlint.RunAll(model, opts)
```

Схема YAML:

```yaml
rules:
  - id: ARCH-BOUNDARIES
    enabled: true        # опустить или false, чтобы пропустить
    config:
      minInternalToCrossRatio: 1.5
      maxCrossRelations: 5
  - id: ARCH-CRUD
    enabled: false       # отключено
```

## 6. Обработка находок

Структура находки:

```go
type Finding struct {
    RuleID   string
    Severity types.Severity // "error", "warn", "info"
    Message  string
    Path     string         // например, boundaries[0].relations[2]
    Meta     map[string]any // необязательно, правило-специфичный контекст
}
```

Типичный pipeline:

1. Сортируйте/фильтруйте по `Severity`, чтобы решать, падает ли CI или отправляется уведомление.
2. Выводите `Path`, чтобы пользователи могли перейти к нужному месту в YAML.
3. Если нужен готовый текст/JSON, используйте `pkg/report` (`report.WriteText` / `WriteJSON`).

## 7. Добавление собственных правил

Правила реализуют интерфейс из `pkg/checks/checks.go`:

```go
type Rule interface {
    ID() string
    Run(*model.Architecture, map[string]any) []types.Finding
}
```

Чтобы добавить правило:

1. Создайте файл в `pkg/checks` и реализуйте интерфейс.
2. Зарегистрируйте его в `pkg/checks/registry.go` (добавьте в слайс `rules`).
3. При необходимости опишите параметры через struct + `decodeConfig`.
4. Добавьте тесты и фикстуры в `pkg/checks` / `testdata`.

После регистрации правило доступно как программно, так и через CLI/конфигурацию.

## 8. Пример: полное встраивание

```go
func Evaluate(path string, rulesConfig string) ([]types.Finding, error) {
    fh, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer fh.Close()

    model, err := archlint.LoadModelFromYAML(fh)
    if err != nil {
        return nil, err
    }

    findings := archlint.ValidateModel(model)

    var opts engine.Options
    if rulesConfig != "" {
        opts, err = config.LoadOptionsFromFile(rulesConfig)
        if err != nil {
            return nil, err
        }
    }

    findings = append(findings, archlint.RunAll(model, opts)...)
    return findings, nil
}
```

Функция повторяет типичный сценарий сервиса: загрузить YAML, опционально прочитать профиль правил и вернуть все находки вызывающей стороне.
