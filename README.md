# Restaurant Search Service

Сервис для работы с данными о ресторанах Москвы, использующий ElasticSearch для поиска и рекомендаций.

## Особенности
- Загрузка данных из CSV в Elasticsearch
- Веб-интерфейс с пагинацией
- JSON API для получения данных
- Поиск ближайших ресторанов по координатам

## Требования
- Go 1.16+ 
- Elasticsearch 8.x
- Данные из materials/data.csv

## Установка и запуск

1. Запустите Elasticsearch:
```bash
path/to/elasticsearch/bin/elasticsearch
```

2. Настройте индекс:
```bash
curl -XPUT -H "Content-Type: application/json" "http://localhost:9200/places/_settings" -d '
{
  "index" : {
    "max_result_window" : 20000
  }'
```

3. Соберите и запустите сервис:
```bash
make build
./loadData # загрузка данных
./startServer # запуск сервера
```

## Makefile команды
```bash 
make build # сборка проекта
make clean # удаление бинарников
```

## API endpoints
### Веб-интерфейс
GET **/?page=**<номер_страницы>
### JSON API
GET **/api/places?page=**<номер_страницы>
### Поиск ближайших ресторанов
GET **/api/recommend?lat=**<широта>**&lon**=<долгота>

## Пример запросов
```bash
# Получение HTML-списка
curl http://localhost:8888/?page=1

# Получение JSON-данных
curl http://localhost:8888/api/places?page=3

# Поиск ближайших ресторанов
curl http://localhost:8888/api/recommend?lat=55.674&lon=37.666
```