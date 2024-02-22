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
* ```AUTH_SERVICE_HOST``` - host of [auth service]() for authenticating

### Private (initialized in ```.env```):
* ```DB_PASSWORD``` - password for PostgreSQL database 
* ```COOKIE_KEY``` - key for encrypting cookies

## API docs
### API overview
| Path                              | Method | Overview                                                |
|-----------------------------------|--------|---------------------------------------------------------|
| [/ping]()                         | GET    | pings server, returns 200 OK if server is healthy       |
| [/auth]()                         | POST   | user authentication                                     |
| [/api/orders]()                   | GET    | returns list of all user orders                         |
| [/api/order]()                    | POST   | creates new order                                       |
| [/api/order/{order_id}]()         | GET    | returns all information about order with ```order_id``` |
| [/api/order/{order_id}/verify]()  | PUT    | marks order with ```order_id``` as verified             |
| [/api/order/{order_id}/compete]() | PUT    | marks order with ```order_id``` as completed            |
| [/api/order/{order_id}]()         | DELETE | deletes order with ```order_id```                       |