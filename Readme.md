# Распределенный вычислитель арифметических выражений

## Обзор

Проект "Распределенный вычислитель арифметических выражений" представляет собой систему, разработанную для асинхронных вычислений сложных арифметических выражений. Система разделена на две основные части: оркестратор и агенты (вычислители). Оркестратор принимает выражения от пользователей, распределяет задачи между вычислителями, и отслеживает статус выполнения. Агенты отвечают за непосредственное выполнение арифметических операций.

## Начало работы

### Требования

- Go 1.22 или выше
- Доступ к СУБД для хранения состояния выражений

### Установка

1. Клонируйте репозиторий в вашу локальную среду разработки:
    ```bash
    git clone https://github.com/AZEN-SGG/distributed-arithmetic-expression-evaluator.git
    ```
2. Перейдите в каталог проекта:
    ```bash
    cd distributed-arithmetic-expression-evaluator
    ```

### Запуск

Для запуска сервера (оркестратора) выполните следующую команду:
```bash
go run main.go
```

## Архитектура системы

Система состоит из двух основных компонентов:

1. **Оркестратор** - сервер, который принимает запросы на вычисление арифметических выражений, управляет их распределением и отслеживает выполнение.

2. **Агенты (Вычислители)** - независимые вычислительные единицы, которые выполняют арифметические операции и возвращают результаты на оркестратор.

## Описание функций сервера и их задачи

### Оркестратор

Оркестратор является центральным узлом системы, который координирует процесс вычисления арифметических выражений. Вот ключевые функции оркестратора:

- **Добавление вычисления арифметического выражения (`/expression`)**: Принимает HTTP POST запросы с арифметическим выражением и уникальным идентификатором. Ответственна за парсинг выражения, его валидацию и добавление в очередь на вычисление.

- **Получение списка выражений со статусами (`/list`)**: Обрабатывает HTTP GET запросы для предоставления списка всех выражений вместе со статусами их выполнения, датой создания и предполагаемым временем вычисления.

- **Получение результата вычисления выражения (`/get`)**: Принимает HTTP POST запросы с идентификатором выражения для получения его результата, если вычисление было завершено.

- **Управление временем выполнения операций (`/math`)**: Позволяет получать текущее время выполнения арифметических операций и обновлять его значения через HTTP GET и POST запросы соответственно.

- **Управление вычислительными ресурсами (`/processes`)**: Отображает текущие вычислительные ресурсы и операции, выполняемые на них, через HTTP GET запрос.

### Агент (Вычислитель)

Агенты выполняют непосредственные арифметические операции, получая задачи от оркестратора. Каждый агент может работать параллельно в несколько потоков (горутин), количество которых регулируется переменной окружения.

## Возможности использования сервера

Сервер предоставляет широкие возможности для асинхронных вычислений арифметических выражений, что особенно полезно в условиях, когда каждая операция требует значительных вычислительных ресурсов. Возможности сервера включают в себя:

- **Масштабирование**: Система может масштабироваться путем добавления дополнительных вычислительных агентов без необходимости изменения оркестратора.

- **Надежность**: Оркестратор способен перезапускаться без потери состояния, благодаря хранению всех выражений в СУБД.

- **Контроль за выполнением задач**: Оркестратор отслеживает задачи, выполняющиеся слишком долго или в случае потери связи с агентом, делает их повторно доступными для вычислений.

## Взаимодействие с сервером

Взаимодействие с сервером возможно через следующие HTTP интерфейсы:

- **Отправка выражений на вычисление**: Пользователи могут отправлять выражения через HTTP POST запрос на `/expression`, получая в ответ уникальный идентификатор выражения.

- **Проверка статуса и получение результата**: Пользователи могут проверять статус выражения и получать его результат, отправляя запрос на `/get` с идентификатором выражения.

- **Настройка времени выполнения операций**: Позволяет пользователям получать и настраивать время выполнения арифметических операций через `/math`.

- **Просмотр списка выражений и их статусов**: Предоставляет обзор всех выражений, их статусов, времени создания и завершения через `/list`.

- **Просмотр и управление вычислительными ресурсами**: Пользователи могут видеть текущие вычислительные ресурсы и задачи, выполняемые на них, через `/processes`.

### Добавление вычисления арифметического выражения

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/expression -d "id=123&expression=2+2*2"
```
**Описание:** Этот запрос отправляет на сервер выражение "2 + 2 * 2" с уникальным идентификатором "123". В случае успешного приема, сервер возвращает статус `200 OK` и уникальный идентификатор выражения, который можно использовать для дальнейшего получения результата.

### Получение списка выражений со статусами

**Пример запроса:**
```bash
curl http://localhost:8080/list
```
**Описание:** Запрос возвращает список всех добавленных выражений со статусами их выполнения, датой создания и предполагаемым временем вычисления. Это позволяет пользователю отслеживать состояние своих выражений.

### Получение результата вычисления выражения

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/get -d "id=123"
```
**Описание:** Этот запрос используется для получения результата выражения по его уникальному идентификатору "123". Если выражение было успешно вычислено, сервер возвращает результат вычисления.

### Управление временем выполнения операций

**Пример получения времени выполнения операций:**
```bash
curl http://localhost:8080/math
```
**Описание:** Возвращает текущие значения времени выполнения для арифметических операций (+, -, *, /).

**Пример обновления времени выполнения операции сложения:**
```bash
curl -X POST http://localhost:8080/math -d "addition=1000"
```
**Описание:** Обновляет время выполнения операции сложения (+) до 1000 миллисекунд. Подобным образом можно обновить время выполнения и для других операций, добавляя параметры к запросу

**Параметры:**
- **addition** - сложение
- **subtraction** - вычитание
- **multiplication** - умножение
- **division** - деление

### Просмотр списка вычислительных ресурсов

**Пример запроса:**
```bash
curl http://localhost:8080/processes
```
**Описание:** Запрос возвращает список текущих вычислительных ресурсов и задач, выполняемых на них. Это позволяет пользователям видеть распределение вычислительных задач и доступные ресурсы.
