# Simple Open Library

This repository contains the code of a simple backend application for scheduling book pickups from [Open Library](https://openlibrary.org)

Due to the nature of the pandemic, a local library is planning to reduce gatherings and crowds in its area by creating an online book pick up service. In order to borrow books from the library, users have to make an appointment at a specific time of day to come and pick up the books.

This backend application provides the necessary APIs to do so. Users can browse books from Open Library before making a schedule. Users can then submit a pick up schedule of the book they want to borrow. Librarians can get all appointed book pick up schedules with the book information. Librarians can also fetch a single appointment schedule. They can also update or delete an appointment.

## API Access

Users do not have to login to use the backend application to browse books from Open Library or previouly submitted books. But an API Key (owned by a librarian) is needed to submit a pick up schedule, list, update, or delete schedules. This is done to prevent users from adding unnecessary schedules or changing submitted schedules without a librarian's verification.

## API Specifications

For now, please copy the API Spec file ([apispec.json](apispec.json)) into a [Swagger Editor](https://editor.swagger.io/) to view the full API specifications.

# Setup local development

## Install tools

-   [Golang](https://go.dev/)
-   [MySQL](https://www.mysql.com/downloads/)
-   [Google Wire](https://github.com/google/wire)
    ```bash
    go install github.com/google/wire/cmd/wire@latest
    ```
-   [Golang Migrate](https://github.com/golang-migrate/migrate)
    ```bash
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

## Setup infrastructure

-   Install the necessary [tools](##install-tools) to run / build the application
-   Setup the .env file. The [.env.example](.env.example) file can be found in the root folder.
-   Setup a MySQL database, and migrate the [schema](db/migrations/20221218140111_init_library_schema.up.sql) using the command below.
    ```bash
    make migrateup
    ```
-   A separate database should be made for unit tests. Migrate the schema to the test database using the command below
    ```bash
    make migrateup_test
    ```
-   Optionally install an API testing tool such as [Postman](https://www.postman.com/downloads) to manually test APIs with ease

## How to run

-   Run server:

    ```bash
    make run
    ```

-   Build application:

    ```bash
    make build
    ```

-   Build and run application:

    ```bash
    make build_run
    ```

-   Run tests:

    ```bash
    make run_test
    ```

### Database Migration Commands

-   Run db migration up all versions
    ```bash
    make migrateup
    ```
-   Run db migration up 1 version
    ```bash
    make migrateup1
    ```
-   Run db migration down all versions
    ```bash
    make migratedown
    ```
-   Run db migration down 1 version
    ```bash
    make migratedown
    ```
-   Run test db migration up all versions
    ```bash
    make migrateup_test
    ```
-   Run test db migration up 1 version
    ```bash
    make migrateup1_test
    ```
-   Run test db migration down all versions
    ```bash
    make migratedown_test
    ```
-   Run test db migration down 1 version
    ```bash
    make migratedown_test
    ```

## Further Improvements

This software still has room for improvements:

-   A virtual container tool like [Docker](https://www.docker.com/) can improve development by automating deployment of applications in isolated lightweight containers so applications can work efficiently in different environments (machines)
-   API documentation should be more easily accessible

## Libraries

-   MySQL Driver ([github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
-   Chi ([github.com/go-chi/chi](github.com/go-chi/chi)) - Golang Router library
-   Validator ([github.com/go-playground/validator](https://github.com/go-playground/validator)) - Request / input validation
-   Google Wire ([github.com/google/wire](https://github.com/google/wire)) - Dependency Injection library
-   Testify ([github.com/stretchr/testify](https://github.com/stretchr/testify)) - Unit testing library
-   GoDotEnv ([github.com/joho/godotenv](https://github.com/joho/godotenv)) - Env loader library
-   sqlmock ([github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)) - SQL driver mock to simulate sql driver behavior in tests
-   Golang Migrate ([github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)) - Library for database migration

## Contributions

-   Gregorius Jovan Kresnadi - [BaconLover307](https://github.com/BaconLover307)
