# Накопительная система лояльности «Гофермарт»
![GoLand](https://img.shields.io/badge/GoLand-0f0f0f?&style=for-the-badge&logo=goland&logoColor=white)
![image](https://pictures.s3.yandex.net:443/resources/gophermart2x_1634502166.png)

Система представляет собой HTTP API со следующими требованиями к бизнес-логике:

* регистрация, аутентификация и авторизация пользователей;
* приём номеров заказов от зарегистрированных пользователей;
* учёт и ведение списка переданных номеров заказов зарегистрированного пользователя;
* учёт и ведение накопительного счёта зарегистрированного пользователя;
* проверка принятых номеров заказов через систему расчёта баллов лояльности;
* начисление за каждый подходящий номер заказа положенного вознаграждения на счёт лояльности пользователя.

### Конфигурирование сервиса накопительной системы лояльности

Сервис поддерживает конфигурирование следующими методами:

- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`


### Прогресс
- [X] Енд-поинты
  - [X] POST /api/user/register — регистрация пользователя;
  - [X] POST /api/user/login — аутентификация пользователя;
  - [X] POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
  - [X] GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
  - [X] GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
  - [X] POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
  - [X] GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
- [X] Хранилища
  - [X] MemStorage — в памяти приложения
  - [X] PgStorage — с использованием PostgresSQL