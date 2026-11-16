# Tiny REST API

Tiny REST API is a lightweight, configurable mock REST server written in Go.

It is designed for development and testing scenarios where the real backend, database, or external service is not available yet. A JSON file is loaded into memory and exposed through configurable REST endpoints, allowing client applications to work against a predictable API without requiring a real database.

Typical use cases include:

- Testing frontends before the real backend is ready.
- Developing clients while another team is still implementing the API.
- Replacing an unavailable backend or external service during development.
- Adapting JSON field names to the contract expected by your application.
- Testing read-only or simulated write scenarios without modifying the source dataset.

**A brief tutorial is available [here](TUTORIAL.md).**

**Breve tutorial en ESPAÑOL [aquí](TUTORIAL-ES.md).**

## Download

Tiny REST API is fully functional, but the project is still adding some improvements and strengthening its test suite to consolidate a solid base before publishing formal releases.

For the moment, you can either:

- Clone the project and compile it yourself.
- Use the precompiled binary and configuration files available in `app/`.

Once the project reaches that stage, packaged releases can replace the temporary `app/` distribution method.

## Features

- **Simple configuration** using `config.ini`.
- **Configurable routes and HTTP methods**.
- **Multiple handlers ("virtual servers")**, each with its own paths and behavior.
- **Route aliases**, allowing several URLs to expose the same resource.
- **Dynamic field aliasing**, so JSON field names can be adapted to the client application without changing its code.
- **Supports RESTful methods:** `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`.
- **Read-only mode**, which rejects write operations.
- **Simulation mode**, which accepts write operations without modifying the in-memory dataset.
- **In-memory modifications only:** the original JSON file is never overwritten.
- **Concurrency-safe** using `sync.Mutex` to prevent race conditions during concurrent requests.
- **Configurable startup**, using a command-line option, an environment variable, or the default configuration file.
- **Unit tests** covering command-line options, configuration loading, aliases, handler configuration, fatal configuration errors, and recoverable warnings.

## Configuration file

Tiny REST API resolves the configuration file in this order:

1. `-c` or `--config`
2. `TINYRESTAPICONFIG`
3. `./config.ini`

Examples:

```sh
./tinyrestapi -c ./myconfig.ini
```

```sh
TINYRESTAPICONFIG=/path/to/config.ini ./tinyrestapi
```

If none of these are provided, Tiny REST API uses `config.ini` from the current working directory.

## Read-only and simulation modes

Tiny REST API provides two different ways to prevent changes to the dataset.

### Read-only mode

When `readonly=true`, write operations such as `POST`, `PUT`, `PATCH`, and `DELETE` are rejected with `405 Method Not Allowed`.

This is useful when testing client permissions or when an endpoint should behave as read-only.

### Simulation mode

When `simulation=true`, write requests are accepted and return their normal success status, but the in-memory dataset is left unchanged.

This is useful when a client needs to exercise write operations without changing the test data.

In short:

```text
readonly=true
    write operations are rejected

simulation=true
    write operations are accepted
    data remains unchanged

normal mode
    write operations modify the in-memory dataset
```

### Why simulation mode is not a traditional dry run

`SimulationMode` is intentionally not implemented as a traditional dry run.

A true dry run would require executing the full write operation without committing the result. Depending on the implementation, that could require cloning or temporarily duplicating part or all of the JSON dataset.

Tiny REST API cannot make assumptions about the size of the JSON file or the resources available on the machine where it runs. Duplicating a large dataset could introduce unnecessary memory usage, performance degradation, or even allocation failures.

Instead, simulation mode simply accepts write requests and returns the corresponding success status without modifying the in-memory dataset.

This keeps the feature lightweight and predictable. If state changes are required during testing, the user can run Tiny REST API against a disposable copy of the JSON dataset.

## How to test

You can test the API with `curl`:

```sh
curl -X <METHOD> [-H {header}] [-d {JSON data}]
```

Some examples using the dummy bookshop from the tutorial:

```sh
# Get all books
curl -X GET http://localhost:8081/books

# Get the book with id 32
curl -X GET http://localhost:8081/books/32

# Add a new book
curl -X POST http://localhost:8081/books \
  -H "Content-Type: application/json" \
  -d '{"id": "34", "title": "The Last Novel", "author": "Mr. Gopher"}'

# Update selected fields
curl -X PATCH http://localhost:8081/books/33 \
  -H "Content-Type: application/json" \
  -d '{"title": "The Penultimate Novel"}'

# Delete a book
curl -X DELETE http://localhost:8081/books/34
```

For a complete example, see [TUTORIAL.md](TUTORIAL.md).

### Browser mode

You can also test the API visually using a REST client extension such as RESTer for Firefox or Chrome.

## The test dataset

Tiny REST API loads its JSON dataset into memory. Changes made through the API affect only the in-memory representation and are never written back to the original file.

The supplied `data.json` contains general-purpose dummy fields that can be reused through field aliases:

| Field | Type | Max Length | Comments |
|:---|:---|---:|:---|
| id | Num | x | |
| name | Text | 25 | |
| surname | Text | 25 | |
| age | Num | 2 | 1 to 99 |
| phone | Num | 11 | Country(2) + Number(9) |
| country_code_2 | Text | 2 | GB, US, JP, FR,... |
| country_code_3 | Text | 3 | BGR, USA, JPN, FRA,... |
| country_name | Text | 25 | |
| address | Text | 40 | |
| zipcode4 | Num | 4 | 4-digit postal code |
| zipcode5 | Num | 5 | 5-digit postal code |
| city | Text | 25 | |
| province | Text | 25 | |
| email | Text | 25 | |
| url | Text | 50 | |
| check1 | Bool | x | |
| check2 | Bool | x | |
| ean | Num | 13 | Barcode |
| isnb | Num | 13 | International Standard Book Number |
| price99 | Num | 5 | 0.00 to 99.99 |
| price999 | Num | 6 | 0.00 to 999.99 |
| text60 | Text | 60 | Max 60 chars |
| text256 | Text | 256 | Max 256 chars |
| comment | Text | 2048 | Between 256-2048 |

## Scope

Tiny REST API is intentionally designed as a lightweight mock server for development and testing.

It is not intended to provide:

- High availability.
- Database persistence.
- Horizontal scalability.
- Production-grade authentication.
- High-performance serving of very large workloads.

The JSON dataset is considered disposable test data. Tiny REST API never writes changes back to the source JSON file.
