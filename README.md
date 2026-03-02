# selectel-linter

# loglint

`loglint` — это линтер для Go, реализованный с использованием `go/analysis`.

Он проверяет сообщения логирования в:

`log/slog`
`go.uber.org/zap`

и выявляет потенциальные нарушения стиля и безопасности.

Необходимо обеспечить:

- единый стиль логирования
- отсутствие небезопасных лог-сообщений
- базовую защиту от утечки чувствительной информации

Анализ выполняется на этапе статической проверки кода.

#Структура проекта

analyzer/
  analyzer.go        — реализация анализатора
  rules_test.go      — unit-тесты правил
  analyzer_test.go   — analysistest

cmd/loglint/
  main.go            

testdata/
  src/a/a.go         — тестовые кейсы для analysistest

demo/
  main.go            — пример использования

#Сборка и запуск 

Из корня проекта:

go build -o loglint ./cmd/loglint

Windows:

go build -o loglint.exe ./cmd/loglint

1. Запуск без сборки (через go run)

Проверить конкретный пакет/директорию:

go run ./cmd/loglint ./demo

Проверить весь проект:

go run ./cmd/loglint ./...

2. Запуск собранного бинаря

Linux/macOS:

./loglint ./...

Windows:

./loglint.exe ./...


Линтер анализирует вызовы:

log/slog
go.uber.org/zap

и применяет правила к первому аргументу (сообщению):

1. Сообщение должно начинаться с маленькой буквы.
2. Сообщение должно быть на английском (без кириллицы).
3. Сообщение не должно содержать спецсимволы или emoji (разрешены только буквы/цифры/пробелы).
4. Сообщение не должно содержать приватную информацию (token/password/api_key и т.п.).

Особенность:

если сообщение формируется конкатенацией, извлекается константная часть ("token: " + token).

правила (1)-(3) применяются только если сообщение полностью константное.

правило (4) применяется и для частично-константных сообщений.

Примеры:

```go
package main

import "log/slog"

func main() {
	slog.Info("Server started") // ❌ должно начинаться с маленькой буквы
}
```

```go
package main

import "log/slog"

func main() {
	slog.Info("сервер запущен") // ❌ сообщение должно быть на английском
}
```

```go
package main

import "log/slog"

func main() {
	slog.Info("server started!") // ❌ спецсимвол
	slog.Info("server 🚀")        // ❌ emoji
}
```

```go
package main

import "log/slog"

func main() {
	token := "secret"
	slog.Info("token: " + token) // ❌ приватная информация
}
```

```go
package main

import "log/slog"

func main() {
	logger := slog.Default()
	logger.Info("Server started") // ❌ uppercase
}
```

```go
package main

import "go.uber.org/zap"

func main() {
	log, _ := zap.NewProduction()
	log.Info("Server started") // ❌ uppercase
}
```

В репозитории есть папка demo/ с примерами нарушений.

Запуск проверки demo:

go run ./cmd/loglint ./demo


#Тестирование

Проект содержит:

Unit-тесты для отдельных правил

Тесты через analysistest для проверки поведения анализатора

Запуск всех тестов:

go test ./…




