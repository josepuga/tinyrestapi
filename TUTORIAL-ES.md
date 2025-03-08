#Tiny REST API

## Example of use.

Imagina que estás desarrollando un frontend para una tienda de libros. Ya tienes los formularios y páginas listos, pero aún no tienes la base de datos del cliente.

El equipo del backend podría no haber terminado la API, el servidor podría estar caído, o simplemente necesitas una solución rápida para pruebas.

Una opción sería encontrar un JSON de ejemplo o una base de datos online y montar un backend improvisado. Sin embargo, esto tiene varios problemas:

- Los nombres de los campos pueden no coincidir con los de tu proyecto.
- Tendrás que modificar tu código o la estructura de datos (lo que genera más trabajo y posibles errores futuros).

`Tiny REST API`  soluciona esto al proporcionar un JSON flexible que se adapta automáticamente a tu proyecto sin necesidad de modificar el código.

## Configuration in `config.ini`

### 1. Crear un "servidor virtual"

Tiny REST API permite definir múltiples configuraciones. En config.ini, cada sección representa un entorno de datos diferente.

Por ejemplo, si queremos configurar un servidor para una librería llamada ACME, agregamos una sección:

```ini
[Book Shop ACME]
```
*(puedes usar cualquier nombre para la sección)*

### 2. Definir las rutas de la API

Si nuestra aplicación hace peticiones a `/books`, pero también queremos soportar `/libros` para usuarios en español, podemos definir ambas rutas::

```ini
paths=books,libros
```
Esto significa que `/books` y  `/libros` devolverán los mismos datos, permitiendo alias para las rutas.

### 3. Indicar los métodos HTTP permitidos
Podemos especificar qué métodos están habilitados para este endpoint. En este caso, activamos todos:

```ini
methods=GET,POST,PUT,PATCH,DELETE,HEAD
```

### 4. Habilitar o deshabilitar modificaciones

Si queremos que los datos sean solo de lectura y evitar que sean modificados accidentalmente, podemos activar `readonly`:

Si intentamos modificar o eliminar un registro cuando `readonly=true`, el servidor responderá con un error de acceso denegado. Esto es útil, por ejemplo, para pruebas de permisos.

```ini
readonly=false
```
*(Por defecto, readonly=false, lo que significa que las modificaciones están permitidas.)*

### 5. Asignar alias a los campos del JSON

Aquí viene la **característica más importante**: adaptar el JSON a los nombres de campos de nuestra app.

El JSON de prueba tiene un campo llamado `text60`, pero en nuestro frontend lo llamamos `title` (como texto dummy para el título del libro), y `phone` lo usamos como `telephone`. Podemos hacer esta conversión automáticamente:

```ini
field aliases=title,test60|telephone,phone
```
*(Tienes una relación de los campos del JSON en el README.md)*
- **Cuando la API devuelva datos, los nombres de los campos se transformarán automáticamente.**

- **Cuando envíes datos en una petición, también se convertirán a los nombres internos.**
  
#### Ejemplo:
JSON original almacenado en memoria
```json
{
  "text60": "The Go Programming Language and the Gin Gonic with Ice",
  "phone": "25123456789"
}
```

Respues de la API tras aplicar los alias:

```json
{
  "title": "The Go Programming Language and the Gin Gonic with Ice",
  "telephone": "25123456789"
}
```
Esto evita que tengas que modificar tu código para adaptarlo al backend.

## ¿Listo? ¡Ejecuta Tiny REST API!
Una vez configurado config.ini, simplemente ejecuta el servidor.
```bash
./tinyrestapi
```

Ahora puedes hacer peticiones a tu API sin preocuparte por nombres de campos incorrectos o por modificar accidentalmente el JSON original.

***Todas las modificaciones se hacen en RAM, el archivo original en disco no se altera.***