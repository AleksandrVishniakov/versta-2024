# Landing page parser
Landing page parser is a simple server that returns a landing page from ```./web/app/build``` which is available on [http://localhost:8080](http://localhost:8080)
## Environment variables
* ```HTTP_PORT``` - http port which service listens (default port is 8080)
* ```LOG_LEVEL``` - level of logging (default level is ```PRODUCTION```). All levels:
    + ```DEBUG```
    + ```INFO```
    + ```PRODUCTION```
    + ```WARNING```
    + ```ERROR```
* ```ORDERS_SERVICE_HOST``` - host of [order service]() which page uses to manage user orders 