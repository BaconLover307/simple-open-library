# simple-open-library
Simple backend application for scheduling book pickups from Open Library

## Unit Test
Unit tests in the same package must be run altogether due to dependencies of records on the database.
To run unit tests, navigate to the respective package (`cd package-name`), and enter the command `go test -v`


## Dependencies
* MySQL Driver ([github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
* HTTP Router ([github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter))
* Validator ([github.com/go-playground/validator](https://github.com/go-playground/validator)) - Request / input validation
* Google Wire: [github.com/google/wire](https://github.com/google/wire) - Dependency Injection library
* Testify ([github.com/stretchr/testify](https://github.com/stretchr/testify)) - Unit testing library
* GoDotEnv ([github.com/joho/godotenv](https://github.com/joho/godotenv)) - Env loader library