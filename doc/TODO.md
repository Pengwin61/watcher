# TO DO
1. ~~Поменять структуру проекта убрать пакеты в папку **/internal**~~

2. Вынести обработки для разных  OS  в отдельные методы

```
https://github.com/Pengwin61/watcher/blob/dev/healthcheck/healthcheck.go
```

3. Сократить функцию, разбить на более мелкие

```
https://github.com/Pengwin61/watcher/blob/dev/watch/watch.go
```

4. Безконечный цикл, нужно подумать как плавно завершать процессы не kill -9, прочитать про graceful shutdown

```
https://blog.ildarkarymov.ru/posts/graceful-shutdown/

https://github.com/Pengwin61/watcher/blob/dev/watch/watch.go#L20
```
~~5. Отказаться от стандартной библиотеки http в пользу gin~~

~~6. Перейти на парсинг конфигов spf13/viper~~

7. Перейти для актеров с ssh на gRPC.
8. Реализовать архивирование пользовательских каталогов через время.
9. Реализовать выгрузку архивированных пользовательских каталогов на s3.
10. Доделать квоты