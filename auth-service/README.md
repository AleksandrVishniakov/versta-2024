# Orders service
Orders service is a microservice for storing and managing user orders

## Environment variables
### Public (initialized in ```docker-compose.yml```):
* ```HTTP_PORT``` - http port which service listens (default port is 8000)
* ```LOG_LEVEL``` - level of logging (default level is ```PRODUCTION```). All levels:
    + ```DEBUG```
    + ```INFO```
    + ```PRODUCTION```
    + ```WARNING```
    + ```ERROR```
* ```DB_HOST```
* ```DB_PORT```
* ```DB_USERNAME```
* ```DB_NAME```
* ```SESSION_EXPIRATION_TIME_MS``` - every user session will be expired after this time (default time is 900000 ms or 15 min)

### Private (initialized in ```.env```):
* ```DB_PASSWORD``` - password for PostgreSQL database

## API docs
### API overview
| Path                             | Method | Overview                                               |
|----------------------------------|--------|--------------------------------------------------------|
| [/ping]()                        | GET    | pings server, returns 200 OK if server is healthy      |