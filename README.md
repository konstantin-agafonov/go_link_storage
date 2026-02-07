Это развитие бота, который получился у Николая Тузова в его курсе по созданию телеграм-бота на Golang.

[Ссылка на плейлист курса на Youtube](https://www.youtube.com/watch?v=PnOrFYtZJUI&list=PLFAQFisfyqlWDwouVTUztKX2wUjYQ4T3l)

[Ссылка на репозиторий с кодом получившегося бота на Github](https://github.com/GolangLessons/Read-Adviser-Bot/tree/lessons) 

Изменения:

- Обернул в Docker
- Air на лету подхватывает изменения в коде и перезапускает приложение
- Добавил реализацию Storage на PostgreSQL
- Выделил Fetcher отдельно от Processor
- Добавил реализацию Client и Fetcher через библиотеку [https://github.com/go-telegram/bot](https://github.com/go-telegram/bot)

Инструкция по запуску:

1) Создать файл `.env` и cкопировать туда содежимое `.env.example`
2) В файле `.env` поставить значение ключа своего телеграм-бота
3) `docker compose up --build`
