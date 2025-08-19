# casino transaction manager


## Project layout

### main dirs

- `cmd/` - app entrypoints
    - `app/` - initialisation and configuration of app
    - `main.go` - main file to run app

- `internal/` - inner logic
    - `config/` - app configuration
    - `handler/` - handlers
    - `middleware/` - middlewares for server
    - `model/` - business models and data structures
      - `/dto` - data transfer objects for requests
      - `/mapper` - structure mapper
    - `storage/` - in memory storage realisation
    - `usecase/` - usecases for tasks
    - `server/` - http server realisation and setup 

- `pkg/logger` - async logger realisation
- 
### Docker files

- `docker-compose.yaml` - main file for running app

## API
Common ports:
- rest api on 8080

### Endpoints
Get all tasks:
```curl
    curl -X GET http://localhost:8080/tasks
```

Get task by id:
```curl
    curl -X GET http://localhost:8080/tasks/{task_id}
```

Get tasks filtered:
```curl
    curl -X GET http://localhost:8080/tasks?status=done # or inProgress or created
```

Post a new task:
```curl
    curl -X POST -H "Content-Type: application/json" -d '{"status": "created", "name": "test name", "description": "test description"}' \
    http://localhost:8080/tasks
```

## App starting

You can change app config in .env file, but for safety reasons don't do like me and dont push them in production repositories

### Main environment

```bash
make run
```
or
```bash
docker-compose up --build
```

### Unit tests
```bash
make test
```

*Note:* if you have issues with running go commands in make file, try: 
```bash
sudo -E env "PATH=$PATH" make *make command here*
```