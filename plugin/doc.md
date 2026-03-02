Установка и использование плагина loglint
#Требования:

1.Go (версия соответствует go.mod)
2.golangci-lint v2.x
3.Доступ к исходному коду проекта (репозиторий с analyzer/ и plugin/)

Проверить версию:

golangci-lint --version

(Версия должна начинаться с v2)

Если установлен v1 — нужно обновить:

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

Плагин loglint реализован как module plugin для golangci-lint v2.

В корне проекта должен находиться файл .custom-gcl.yml.

Пример .custom-gcl.yml:

version: v2.10.1
name: custom-gcl
destination: .

plugins:
  - module: "linter.com/loglint"
    path: .
    import: "linter.com/loglint/plugin/loglint"

В корне проекта выполнить:

golangci-lint custom -v

После сборки появится файл:

Linux/macOS: custom-gcl
Windows: custom-gcl.exe

#Настройка golangci-lint

Создать файл .golangci.yml (формат v2):

version: "2"

linters:
  default: none
  enable:
    - loglint

  settings:
    custom:
      loglint:
        type: "module"
        description: "Checks slog/zap log messages for style & safety rules"
        settings: {}

Важно:

используется формат v2 (linters.settings.custom)
линтер называется ровно так же, как в register.Plugin("loglint", ...)

#Запуск

Использовать собранный бинарь, а не системный golangci-lint:

./custom-gcl run -c .golangci.yml ./...

Windows:

./custom-gcl.exe run -c .golangci.yml ./...