# DevFlow

DevFlow — локальная CLI-тулза на Go для запуска dev-окружения проекта одной командой.

## Главная боль

В каждом backend-проекте перед работой нужно вспомнить:

- какие Docker-сервисы поднять;
- какой `.env` нужен;
- какие порты заняты;
- какие миграции прогнать;
- какие команды запустить;
- почему сервис не стартует.

DevFlow должен превратить это в один понятный сценарий:

```bash
devflow up
```

## MVP

Первый MVP фокусируется только на dev-окружении:

```bash
devflow init
devflow up
devflow down
devflow status
devflow doctor
```

## Что делает `devflow up`

1. читает `devflow.yaml`;
2. проверяет наличие `.env`;
3. валидирует обязательные env-переменные;
4. проверяет внешние команды: `docker`, `docker compose`;
5. проверяет порты;
6. запускает Docker Compose;
7. ждёт readiness сервисов;
8. выполняет prepare-команды: migrations/seed/install;
9. печатает понятный результат.

## Пример `devflow.yaml`

```yaml
project:
  name: billing-service

env:
  file: .env
  required:
    - DATABASE_URL
    - REDIS_URL

compose:
  files:
    - docker-compose.yml
  project_name: billing-service

checks:
  commands:
    - docker
  ports:
    - name: api
      host: localhost
      port: 3000
      required_free_before_up: true
    - name: postgres
      host: localhost
      port: 5432
      required_free_before_up: false

services:
  - name: postgres
    type: tcp
    host: localhost
    port: 5432
    timeout: 30s
  - name: api
    type: http
    url: http://localhost:3000/health
    timeout: 30s

commands:
  after_up:
    - name: migrations
      run: npm run migration:run
    - name: seed
      run: npm run seed
      optional: true
```

## Почему это хороший Go-проект

Stage 1 закрывает реальные Go-темы:

- CLI;
- YAML config;
- filesystem;
- env validation;
- external commands;
- context timeout;
- TCP/HTTP healthchecks;
- structured errors;
- testing;
- Docker Compose orchestration.

## Будущие стадии

После рабочего MVP можно добавить:

1. tmux workspace;
2. project registry;
3. file indexing;
4. AI/RAG;
5. AI-agent orchestration.

Но Stage 1 не должен уходить в эти темы.
