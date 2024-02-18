**Распределенный вычислитель арифметических выражений**

**Обзор**

Проект "Распределенный вычислитель арифметических выражений" представляет собой систему, разработанную для асинхронных вычислений сложных арифметических выражений. 

Система разделена на две основные части: оркестратор и агенты (вычислители). 

Оркестратор принимает выражения от пользователей, распределяет задачи между вычислителями, и отслеживает статус выполнения. 

Агенты отвечают за непосредственное выполнение арифметических операций.

**Начало работы**

**Требования**

Go 1.22 или выше

Доступ к СУБД для хранения состояния выражений

**Установка**

Клонируйте репозиторий в вашу локальную среду разработки:

``` bash
git clone https://github.com/AZEN-SGG/distributed-arithmetic-expression-evaluator.git
```
Перейдите в каталог проекта:

``` bash
cd distributed-arithmetic-expression-evaluator
```
Убедитесь, что ваша СУБД настроена и доступна для подключения.

**Запуск**

Для запуска сервера (оркестратора) выполните следующую команду:

``` bash
go run main.go
```
**Использование**

Добавление вычисления арифметического выражения
Для отправки выражения на вычисление используйте HTTP POST запрос на /expression с телом запроса, содержащим арифметическое выражение и уникальный идентификатор. Пример запроса:

``` bash
curl -X POST http://localhost:8080/expression -d "id=unique_id&expression=2+2*2"
```
**Получение списка выражений со статусами**

Отправьте HTTP GET запрос на /list для получения списка всех выражений со статусами их выполнения. Пример запроса:

``` bash
curl http://localhost:8080/list
````
**Получение результата вычисления выражения**

Чтобы получить результат вычисления, отправьте HTTP POST запрос на /get с идентификатором выражения. Пример запроса:

``` bash
curl -X POST http://localhost:8080/get -d "id=unique_id"
```
**Архитектура системы**

Система состоит из двух основных компонентов:

**Оркестратор** - сервер, который принимает запросы на вычисление арифметических выражений, управляет их распределением и отслеживает выполнение.

**Агенты** (Вычислители) - независимые вычислительные единицы, которые выполняют арифметические операции и возвращают результаты на оркестратор.

Примеры использования:

После запуска оркестратора и настройки вычислителей пользователи могут отправлять арифметические выражения через предоставленный GUI или напрямую через API, как описано выше. 
В зависимости от сложности выражения и загрузки системы, результат будет доступен через определенный интервал времени.

