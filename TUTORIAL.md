# Tiny REST API

## Example of use

Imagine you are developing a frontend for a bookshop. You already have the forms and pages ready, but the client's database or backend API is not available yet.

The backend team might still be implementing the service, the server might temporarily be unavailable, or you might simply need a quick and predictable environment for testing.

One option would be to find a sample JSON file or an online database and build a temporary backend around it. However, this creates several problems:

- The field names may not match your project's data contract.
- The routes may not match the URLs expected by your client.
- You may have to modify application code just to accommodate temporary test data.
- A temporary backend introduces additional work that will later be discarded.

Tiny REST API lets you expose a test JSON dataset through a configurable REST API and adapt its routes and field names to the client application.

This allows the frontend or client application to continue development without depending on the final backend.

---

## Configuration in `config.ini`

### 1. Creating a "virtual server"

Tiny REST API allows multiple handler configurations. Each section in `config.ini` represents a different data environment or endpoint configuration.

For example, to configure a bookshop called ACME:

```ini
[Book Shop ACME]
```

You can use any name for the section.

---

### 2. Defining API routes

If the application makes requests to `/books`, but we also want to expose the same resource as `/libros`, we can configure both routes:

```ini
paths=books,libros
```

Both `/books` and `/libros` expose the same data.

Route aliases are not limited to languages. They can also be useful for legacy URLs, different client conventions, or transitional API paths.

---

### 3. Specifying allowed HTTP methods

We can define which HTTP methods are enabled for this handler.

For example:

```ini
methods=GET,POST,PUT,PATCH,DELETE,HEAD
```

Only the configured methods are exposed.

---

### 4. Read-only mode

If we want the endpoint to behave as read-only, we can enable:

```ini
readonly=true
```

When `readonly=true`, write methods such as `POST`, `PUT`, `PATCH`, and `DELETE` are rejected with:

```text
405 Method Not Allowed
```

This can be useful for permission testing or for endpoints that should not accept modifications.

By default:

```ini
readonly=false
```

---

### 5. Simulation mode

Sometimes we want the client to exercise write operations without changing the test dataset.

For this case:

```ini
simulation=true
```

In simulation mode, write requests are accepted and return their normal success status, but no changes are applied to the in-memory dataset.

For example:

```text
DELETE /books/32
```

can return a successful response while book `32` is still present in later `GET` requests.

The difference between the modes is:

```text
readonly=true
    write operations are rejected

simulation=true
    write operations are accepted
    data remains unchanged

normal mode
    write operations modify the in-memory dataset
```

#### Why simulation mode is not a dry run

Simulation mode is intentionally simpler than a traditional dry run.

A true dry run would require processing the complete mutation while ensuring that its result is not committed. Depending on the implementation, this could require cloning or temporarily duplicating part or all of the JSON dataset.

Tiny REST API cannot know in advance how large the user's JSON file will be or how much memory is available on the machine where it runs. Duplicating a large dataset only to simulate an operation could create unnecessary memory usage, performance degradation, or allocation failures.

Simulation mode therefore accepts the request and returns the corresponding success status without modifying the dataset.

If real state changes are required during testing, the user can simply run Tiny REST API with a disposable copy of the JSON file.

---

### 6. Assigning aliases to JSON fields

One of Tiny REST API's main features is the ability to adapt JSON fields to the names expected by the client application.

For example, the test JSON contains a field named `text60`, but the frontend expects a field named `title`. It also contains `phone`, while the client expects `telephone`.

We can define aliases:

```ini
field aliases=title,text60|telephone,phone
```

When the API returns data, field names are transformed to the aliases expected by the client.

When the client sends data, aliases are converted back to the internal JSON field names.

#### Example

Original JSON stored in memory:

```json
{
  "text60": "The Go Programming Language and the Gin Gonic with Ice",
  "phone": "25123456789"
}
```

API response after applying aliases:

```json
{
  "title": "The Go Programming Language and the Gin Gonic with Ice",
  "telephone": "25123456789"
}
```

This allows the client application to work with its own data contract without modifying the original dataset.

You can find the full list of fields included in the sample JSON in the `README.md`.

---

## Running Tiny REST API

Once `config.ini` is configured, start the server:

```bash
./tinyrestapi
```

You can also specify a different configuration file:

```bash
./tinyrestapi -c ./myconfig.ini
```

or use the environment variable:

```bash
TINYRESTAPICONFIG=/path/to/config.ini ./tinyrestapi
```

The configuration file is resolved in this order:

1. `-c` or `--config`
2. `TINYRESTAPICONFIG`
3. `./config.ini`

Now the client application can make requests against the mock API without depending on the final backend.

All normal modifications are performed in RAM. The original JSON file on disk is never modified.

---

## Project status

Tiny REST API is fully functional.

The project is currently receiving additional improvements and strengthening its test suite to consolidate a solid foundation before formal packaged releases are published.

For the moment, the precompiled binary and configuration files are available under `app/`, or the project can be compiled directly from source.
