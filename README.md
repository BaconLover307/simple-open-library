# Simple Open Library
This repository contains the code of a simple backend application for scheduling book pickups from [Open Library](https://openlibrary.org)

Due to the nature of the pandemic, Open Library is planning to reduce gatherings and crowds in its area by creating an online book pick up service. In order to borrow books from the library, users have to make an appointment at a specific time of day to come and pick up the books.

This backend application provides the necessary APIs to do so. Users can browse books from Open Library before making a schedule. Users can then submit a pick up schedule of the book they want to borrow. Librarians can get all appointed schedules and their respective books. Librarians can also fetch a single appointment schedule. They can also update or delete an appointment.

## API Access
Users do not have to login to use the backend application to browse books from Open Library or previouly submitted books. But an API Key (owned by a librarian) is needed to submit a pick up schedule, list schedules, or edit pick up schedules. This is done to prevent users from adding unnecessary schedules or changing submitted schedules without a librarian's verification.

## API Specifications
For now, please copy the API Spec file ([apispec.json](apispec.json)) into a [Swagger Editor](https://editor.swagger.io/) to view the full API specifications.

# Setup local development
## Install tools

- [Golang](https://go.dev/)
- [MySQL](https://www.mysql.com/downloads/)
- [Google Wire](https://github.com/google/wire)
    ```bash
    go install github.com/google/wire/cmd/wire@latest
    ```
## Prerequisites
* Install the necessary [tools](##install-tools) to run / build the application
* Setup the .env file. The [.env.example](.env.example) file can be found in the root folder.
* Migrate the [schema](##sql-schema) into MySQL database. A separate database should be made for unit tests.

## SQL Schema
```sql
CREATE TABLE book(
    bookId varchar(20) PRIMARY KEY,
    title varchar(200),
    edition integer
);
 
CREATE TABLE pickup(
    pickupId integer PRIMARY KEY AUTO_INCREMENT,
    bookId varchar(20),
    schedule datetime,
    FOREIGN KEY (bookId) REFERENCES book(bookId)
);
 
CREATE TABLE author(
   authorId varchar(20) PRIMARY KEY,
   name varchar(100)
);
 
CREATE TABLE authored(
   authorId varchar(20),
   bookId varchar(20),
   PRIMARY KEY (authorId, bookId),
   FOREIGN KEY (authorId) REFERENCES author(authorId),
   FOREIGN KEY (bookId) REFERENCES book(bookId)
);

```
## How to run
* Run server:
```bash
make run
```
* Build application:
```bash
make build
```
* Run tests:
```bash
make test
```
#
## Further Improvements
This software still has room for improvements:
* A proper database migration can improve developer experience when setting up / building the project
* Unit tests can be improved using mocking libraries. As of now, the unit tests still use a mock database with MySQL. Mocking libraries can make unit tests easier to setup and execute. 


## Libraries
* MySQL Driver ([github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
* HTTP Router ([github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter))
* Validator ([github.com/go-playground/validator](https://github.com/go-playground/validator)) - Request / input validation
* Google Wire ([github.com/google/wire](https://github.com/google/wire)) - Dependency Injection library
* Testify ([github.com/stretchr/testify](https://github.com/stretchr/testify)) - Unit testing library
* GoDotEnv ([github.com/joho/godotenv](https://github.com/joho/godotenv)) - Env loader library

## Contributions
* Gregorius Jovan Kresnadi - [BaconLover307](https://github.com/BaconLover307)