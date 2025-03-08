# Tiny REST API

Tiny REST API is a simple yet powerful tool to simulate a backend API with customizable JSON fields.
Perfect for:

- Testing frontends before the real backend is ready.
- Adapting JSON fields without modifying your project.
- Running a lightweight mock server without a database.


 **A brief tutorial is available [here](TUTORIAL.md)**.

 **Breve tutorial en ESPAÑOL [aquí](TUTORIAL-ES.md)**.

## **Download**
You can either **clone the project and compile it yourself**, or **download the precompiled and config files from `app/`**.

## **Features**
- **Simple configuration** using `config.ini`.
- **Multiple virtual servers**, allowing different paths and aliases.
- **Loads a test JSON file** and behaves like a database.
- **Dynamic field aliasing**, so you can adapt JSON fields to your project instead of adapting your project to the fields.
- **Supports RESTful methods:** `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`.
- **Concurrency-safe** using `sync.Mutex` to prevent race conditions in multiple requests.

## How to test 
You can test the API with:
```sh
curl -X <METHOD> [-H {header}] [-d {JSON data}]
```
### Console mode
Here are some examples using the dummy book shop from the tutorial:
```sh
# Get all books
curl -X GET http://localhost:8081/books

# Get the book with id 32
curl -X GET http://localhost:8081/books/32

# Add a new book (POST)
curl -X POST http://localhost:8081/books -H "Content-Type: application/json" -d '{"id": "34", "title": "The Last Novel", "author": "Mr. Gopher"}'

# Update a book (PUT - full replacement)
curl -X PUT http://localhost:8081/books/33 -H "Content-Type: application/json" -d '{"title": "The Penultimate Novel"}'

# Update a book (PATCH - only change specific fields)
curl -X PATCH http://localhost:8081/books/33 -H "Content-Type: application/json" -d '{"title": "The Penultimate Novel"}'

# Delete a book
curl -X DELETE http://localhost:8081/books/34
```

### Browser mode
You can test the API visually by using the `RESTer extension` for Firefox or Chrome.

## The "dummy" Data

Don't worry about modifying or deleting records!. All data is stored in RAM, so it is not written to disk.

## JSON Structure
A JSON file `data.json` is provided with some fields of general purpose:
|Field|Type|Max Length|Comments|
|:---|:---|---:|:---|
|id|Num|x||
|name|Text|25||
|surname|Text|25||
|age|Num|2| 1 to 99||
|phone|Num|11|Country(2) + Number(9)|
|country_code_2|Text|2|GB, US, JP, FR,...|
|country_code_3|Text|3|BGR, USA, JPN, FRA,...|
|country_name|Text|25||
|address|Text|40||
|zipcode4|Num|4|4 digits Postal code|
|zipcode5|Num|5|5 digits...|
|city|Text|25||
|province|Text|25||
|email|Text|25||
|url|Text|50||
|check1|Bool|x||
|check2|Bool|x||
|ean|Num|13|Bar code|
|isnb|Num|13|International Standard Book Number|
|price99|Num|5| 0.00 to 99.99|
|price999|Num|6| 0.00 to 999.99|
|text60|Text|60|Max 60 chars|
|text256|Text|256|Max 256 chars|
|comment|Text|2048|Between 256-2048|
