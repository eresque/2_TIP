# Практическое занятие №8
# Настройка GitHub Actions / GitLab CI для деплоя приложения

**Дисциплина:** Технологии индустриального программирования  
**Студент:** Гордеев Артём Ильич, ЭФМО-01-25

---

## Требования к проекту

- Go 1.23+
- Docker (Docker Desktop или Docker Engine)
- Аккаунт GitHub с доступом к GitHub Actions
- Для публикации образа: настроенные GitHub Secrets (`REGISTRY_USERNAME`, `REGISTRY_PASSWORD`)

---

## Краткое описание проекта

Реализован CI/CD pipeline на **GitHub Actions** для Go-сервиса **tasks**.

Pipeline состоит из трёх job:
- `test-and-build` — установка зависимостей, запуск тестов (`go test ./...`), сборка бинарника (`go build ./...`);
- `docker-build` — сборка Docker-образа, тегируемого хешем коммита (`github.sha`);
- `publish` — публикация образа в GitHub Container Registry (`ghcr.io`), выполняется только при push в `main`/`master` при наличии настроенных секретов.

Основой послужил HTTP-сервис `tasks` из практической работы №7. Для корректной работы `go test` в `internal/` добавлен unit-тест обработчика `/health` на основе `httptest`.

---

## Структура проекта

```
pr8/
├── .github/
│   └── workflows/
│       └── ci.yml
├── services/
│   └── tasks/
│       ├── cmd/
│       │   └── tasks/
│       │       └── main.go
│       ├── internal/
│       │   ├── handler.go
│       │   └── handler_test.go
│       ├── .dockerignore
│       ├── Dockerfile
│       ├── go.mod
│       └── go.sum
└── README.md
```

---

## CI/CD Pipeline

### Выбранная платформа

Использован **GitHub Actions** — файл `.github/workflows/ci.yml`.

### Полный YAML-файл pipeline

```yaml
name: CI Pipeline

on:
  push:
    branches: ["main", "master"]
  pull_request:
    branches: ["main", "master"]

jobs:
  test-and-build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Show Go version
        run: go version

      - name: Download dependencies
        run: go mod tidy
        working-directory: ./services/tasks

      - name: Run tests
        run: go test ./...
        working-directory: ./services/tasks

      - name: Build application
        run: go build ./...
        working-directory: ./services/tasks

  docker-build:
    runs-on: ubuntu-latest
    needs: test-and-build

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        run: docker build -t techip-tasks:${{ github.sha }} .
        working-directory: ./services/tasks

  publish:
    runs-on: ubuntu-latest
    needs: docker-build
    if: github.event_name == 'push' && (github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master')

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Login to GitHub Container Registry
        run: echo "${{ secrets.REGISTRY_PASSWORD }}" | docker login -u "${{ secrets.REGISTRY_USERNAME }}" --password-stdin ghcr.io

      - name: Build and tag image for registry
        run: |
          docker build \
            -t ghcr.io/${{ github.repository_owner }}/techip-tasks:${{ github.sha }} \
            -t ghcr.io/${{ github.repository_owner }}/techip-tasks:latest \
            .
        working-directory: ./services/tasks

      - name: Push image to registry
        run: |
          docker push ghcr.io/${{ github.repository_owner }}/techip-tasks:${{ github.sha }}
          docker push ghcr.io/${{ github.repository_owner }}/techip-tasks:latest
```

### Пояснение шагов pipeline

**Job `test-and-build`:**
| Шаг | Команда | Назначение |
|-----|---------|------------|
| Checkout repository | `actions/checkout@v4` | Клонирует репозиторий на runner |
| Setup Go | `actions/setup-go@v5` | Устанавливает Go 1.23 |
| Show Go version | `go version` | Верификация установленной версии |
| Download dependencies | `go mod tidy` | Синхронизирует go.mod и go.sum |
| Run tests | `go test ./...` | Запускает все тесты в проекте |
| Build application | `go build ./...` | Компилирует бинарник сервиса |

**Job `docker-build`** (запускается только после успешного `test-and-build`):
| Шаг | Назначение |
|-----|------------|
| Set up Docker Buildx | Настраивает расширенный сборщик образов |
| Build Docker image | Собирает образ с тегом из хеша коммита |

**Job `publish`** (только push в main/master, требует secrets):
| Шаг | Назначение |
|-----|------------|
| Login to registry | Авторизация в ghcr.io через secrets |
| Build and tag | Сборка образа с двумя тегами: SHA и latest |
| Push image | Публикация обоих тегов в registry |

### Формирование тега Docker-образа

Образ тегируется двумя способами:
- `techip-tasks:${{ github.sha }}` — полный SHA коммита (локальная сборка в CI);
- `ghcr.io/.../techip-tasks:${{ github.sha }}` + `latest` — при публикации в registry.

Это позволяет точно отследить, из какого коммита собран конкретный образ.

### Хранение секретов

Секреты (`REGISTRY_USERNAME`, `REGISTRY_PASSWORD`) хранятся в **GitHub Secrets** (`Settings → Secrets and variables → Actions`). Они недоступны из кода репозитория, не попадают в логи и не коммитятся в YAML-файл. В pipeline они передаются через синтаксис `${{ secrets.NAME }}`.

---

## Локальная проверка перед push

```bash
cd services/tasks

# Тесты
go test ./...

# Сборка бинарника
go build ./...

# Сборка Docker-образа
docker build -t techip-tasks:0.1 .
```

---

## Результаты выполнения (скриншоты)

### Успешный запуск pipeline — вкладка Actions
![pipeline overview](readme-images/test-1.png)

### Job test-and-build — прохождение тестов и сборки
![test and build job](readme-images/test-2.png)

### Job docker-build — успешная сборка образа
![docker build job](readme-images/test-3.png)

### Локальный запуск тестов
```
go test ./...
```
![local tests](readme-images/test-4.png)

### Локальная сборка Docker-образа
```
docker build -t techip-tasks:0.1 .
```
![local docker build](readme-images/test-5.png)

---

## Ответы на контрольные вопросы

**1. Чем CI отличается от CD?**  
CI (Continuous Integration) — автоматическая проверка и сборка кода после каждого изменения: запуск тестов, компиляция. CD (Continuous Delivery / Deployment) — следующий шаг: упаковка результата в артефакт (Docker-образ), доставка в registry и/или автоматическое развёртывание на сервер.

**2. Почему pipeline должен запускать тесты?**  
Тесты в CI гарантируют, что каждый push не ломает рабочую функциональность. Это даёт раннее обнаружение регрессий: разработчик видит ошибку сразу после отправки кода, а не на этапе деплоя.

**3. Зачем нужен автоматический build?**  
Автоматическая сборка подтверждает, что код компилируется в чистой среде (не только на машине разработчика). Это исключает ошибки «у меня работает, а на сервере нет» из-за различий в окружении.

**4. Почему важно собирать Docker-образ в CI, а не только локально?**  
Сборка в CI обеспечивает воспроизводимость: образ собирается из конкретного коммита в чистой среде без локальных изменений. Это гарантирует, что образ точно соответствует коду в репозитории.

**5. Что такое CI secrets?**  
Secrets — зашифрованные переменные, хранящиеся на стороне CI-платформы (GitHub Secrets, GitLab CI/CD Variables). Они передаются в pipeline только во время выполнения и никогда не попадают в репозиторий или логи.

**6. Почему нельзя хранить токены и SSH-ключи в репозитории?**  
Репозиторий доступен всем участникам, а история коммитов — навсегда. Если токен однажды попал в Git, он скомпрометирован, даже если его позже удалить. Secrets-хранилища CI решают эту проблему.

**7. Для чего нужен тег Docker-образа?**  
Тег позволяет однозначно идентифицировать версию образа. Тег из SHA коммита связывает образ с конкретным состоянием кода, что упрощает откат и отладку. Без тега все образы перезаписывали бы друг друга.

**8. Что делает job docker-build?**  
Job `docker-build` запускается после успешного `test-and-build`, получает код репозитория, настраивает Docker Buildx и выполняет `docker build` — собирает образ с тегом из хеша коммита. Это подтверждает, что Dockerfile рабочий.

**9. Почему в multi-service проекте важен working-directory?**  
В репозитории с несколькими сервисами каждый сервис находится в своей директории со своим `go.mod`. Без явного `working-directory` команды `go test` и `go build` выполняются в корне репозитория и не найдут нужный модуль.

**10. Какие риски возникают при полностью автоматическом деплое?**  
Если тест не покрывает какую-то ошибку, она автоматически попадает в продакшн. Также возможны проблемы с конкурентными деплоями, откатом при сбое, несовместимостью БД-миграций. Поэтому на практике между CI и CD часто добавляют ручное подтверждение (manual approval gate).
