# Setting up

To start work with the project you nee to have installed some tools:

+ docker - to start the project
+ Go language - to run unit and integration tests

Optionaly some tool to run a makefile commands.

If you have installed above tools next thing is create `.env` file in a source folder. It should contains variables like

+ BINARY - name of created binary if you decide to build application
+ DB_CONTAINER_NAME - name of database conteiner
+ MONGO_DB - name of db
+ MONGO_DB_USERNAME - database username
+ MONGO_DB_PASSWORD - database password
+ SWIFT_APP - name of application image and container 

# Starting application

To start application you can use a `make up` command or a direct command from a up section in `Makefile`. To shut down containers you can use `make down` command or `docker-compose down`. Application works on `http://localhost:8080`, database is available on `http://localhost:27017`. To chceck database content you can use mongo compas, where you can connect to the database. To access API you can use e.g. postman here are list of provided endpoints:

+ GET `http://localhost:8080/v1/healthcheck` - chceck a API availibility
+ POST `http://localhost:8080/v1/swift-codes` - add a new swift code
+ GET `http://localhost:8080/v1/swift-codes` - get all swift codes
+ GET `http://localhost:8080/v1/swift-codes/{swift-code}` - get a swift code by swift code field
+ GET `http://localhost:8080/v1/swift-codes/country/{countryISO2code}` - get all swift codes with matching provided ISO2 code
+ DELETE `http://localhost:8080/v1/swift-codes/{swift-code}` - delete swift code witch matching swift code field

# Tests

To run a tests you need to download all dependencies by `go mod download`, the you can use make test command to run tests