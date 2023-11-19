# ports-app

## Quickstart

1. Spin-up server on `localhost:3000`:  
   `docker-compose up`


2. From a second terminal, import the ports in the file `testdata/ports.json`

```shell
curl --location 'http://localhost:3000/ports/import' \
--header 'Content-Type: application/json' \
--data "@testdata/ports.json"
```

3. Server should have logged: `imported 1632 ports`


