# xm-platform

### Run & Setup Database
* In order to run application and database simply run `make compose`.
* In order to setup the database (database and table with some mock data) simply run `make setupdb` (use 'xm123' as password). This should be executed once.

### API Calls
Once the application is up and database is set up use the below routes to make API calls:

```
* GET /v1/health => Check app's health (connectivity with the database)
example curl:
curl --location --request GET 'localhost:8080/v1/health'

* POST /v1/login => Authenticate user. Returns a JWT token.
example curl:
curl --location --request POST 'localhost:8080/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "xm",
    "password": "xm123"
}'

* GET /v1/company/{id} => Gets a single company
example curl:
curl --location --request GET 'localhost:8080/v1/company/2c5165c5-6a66-45fd-a57e-d91726a8b32f'

* POST /v1/company => Creates a company
example curl:
curl --location --request POST 'localhost:8080/v1/company' \
--header 'Authorization: Bearer <token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Test2",
    "type": "Corporations",
    "employees": 1000,
    "registered": false
}'

* DELETE /v1/company/{id} => Deletes a company
example curl:
curl --location --request DELETE 'localhost:8080/v1/company/489f09c2-9849-482f-ffff-fb985c93a7a6' \
--header 'Authorization: Bearer <token>' \

* PATCH /v1/company/{id} => Patches a company
example curl:
curl --location --request PATCH 'localhost:8080/v1/company/489f09c2-9849-482f-90c6-fb985c93a7a6' \
--header 'Authorization: Bearer <token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "employees": 5,
    "type": "Cooperative"
}'
```

### Hints
* Run `make lint` to run linter on project
* `.env` file has every project's sensitive configuration (passwords, tokens, etc)
