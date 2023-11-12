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