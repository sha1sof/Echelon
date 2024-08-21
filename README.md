# Echelon

Это проект, включающий несколько приложений для работы с видео-превью и миграциями базы данных. Проект включает сервер для предоставления превью изображений с YouTube, клиент для загрузки изображений и утилиту базы данных.

## Структура проекта

- **cmd/server**: Главный исполняемый файл для запуска gRPC сервера.
- **cmd/migrator**: Утилита для выполнения миграций базы данных.
- **cmd/client**: Клиент для загрузки изображений превью.
- **internal**: Внутренние пакеты проекта.
- **test**: Тесты для компонентов проекта.

## Установка

1. Убедитесь, что у вас установлен Go версии 1.23.0.
2. Клонируйте репозиторий и перейдите в директорию проекта:


    git clone https://github.com/sha1sof/Echelon
    
    cd Echelon


## Запуск мигратора

Для функционирования сервера, нужно сначала создать базу данных sqlite. Это нужно сделать следующей командой:

    go run ./cmd/migrator --storage=./youtube.db --migrations=./migrations


- **--storage**: Путь к базе данных.
- **--migrations**: Путь к папке с файлами миграций.

## Запуск сервера

Для запуска gRPC сервера используйте следующую команду:

    go run cmd/server/main.go

Сервер будет слушать на порту, указанном в конфигурационном файле ./config/prod.yaml.

## Запуск клиента

Для загрузки превью изображений видео используйте следующие команды:

### Асинхронный режим

Для асинхронной загрузки превью изображений видео:

    go run cmd/client/main.go --async https://www.youtube.com/watch?v=e_pY0btswmk https://www.youtube.com/watch?v=bcwpkiXlpno

### Синхранный режим

Для синхронной загрузки превью изображений видео:

    go run cmd/client/main.go https://www.youtube.com/watch?v=e_pY0btswmk

## Конфигурация проекта

Конфигурация проекта хранится в файле ./config/prod.yaml. Пример конфигурации:

    env: prod
    storage:
        type: memory
        storage_path: ./storage/preview.db
    grpc_server:
        port: 8080
    clients:
        preview:
            address: localhost:50051
            timeout: 5s
            retriesCount: 3
            output_dir: ./thumbnails

- **env**: Среда выполнения, отличия в типе логироваия (prod или local).
- **storage**: Параметры для хранения данных.
- **grpc_server**: Параметры для gRPC сервера.
- **clients.preview**: Параметры для клиента.

## Тестирование

Для запуска тестов используйте команду:

    go test ./test/...

### 2 раза подряд нельзя делать тесты, нужно менять параметр ***url*** в файле ***getThumbnail_test.go*** на другую ссылку youtube.
