# TODOLIST

TODOLIST - приложение для ведения списка дел, состоящее из веб-интерфейса и сервера на языке go. Сервер состоит из нескольких логических субмодулей:
* API модуль для обработки HTTP-запросов
* модуль для работы с базой данных
* логический модуль для работы с данными
* некоторые другие вспомогательные субмодули

# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Настройка сервера

На текущий момент реализована поддержка двух параметров приложения:
* TODO_PORT - порт, на котором осуществляет работу сервер
* TODO_DBFILE - путь к файлу базы данных sqlite

Эти переменные можно указать в файле **.env**. Также эти переменные можно указать в качестве флагов **--env** в команде `docker build` при сборке docker-образа приложения.

Адрес веб-интерфейса приложения в веб-браузере по адресу: `localhost:<port>`, где **port** соответствует переменной **TODO_PORT**.

# Docker

Для сборки докер-образа приложения используется команда `docker build -t <name>:<tag>`, где **name** и **tag** выбираются произвольно. 
Для запуска сервиса в контейнере используется команда `docker run -p <port>:<port> <name>:<tag>`, где **port** - соответствует порту, указанному в файле .env при сборке docker-образа.

# Выполненные задания

Помимо основного блока заданий были реализованы несколько дополнительных, со "звёздочкой":

* реализована функция поиска по содержимому заголовков, комментариев, а также дате задач
* реализована возможность параметризации рабочего порта сервера, а также файла базы данных
* добавлен Dockerfile для запуска приложения в контейнере