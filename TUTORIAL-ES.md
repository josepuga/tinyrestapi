# Tiny REST API

## Ejemplo de uso

Imagina que estás desarrollando el frontend de una tienda de libros. Ya tienes preparados los formularios y las páginas, pero todavía no dispones de la base de datos del cliente o de la API del backend.

El equipo de backend podría seguir implementando el servicio, el servidor podría estar temporalmente caído o simplemente podrías necesitar un entorno rápido y predecible para realizar pruebas.

Una opción sería buscar un JSON de ejemplo o una base de datos online y montar un backend temporal alrededor de esos datos. Sin embargo, esto plantea varios problemas:

- Los nombres de los campos pueden no coincidir con el contrato de datos de tu proyecto.
- Las rutas pueden no coincidir con las URLs que espera el cliente.
- Podrías tener que modificar el código de la aplicación únicamente para adaptarlo a datos temporales.
- Montar un backend provisional supone trabajo adicional que después será descartado.

Tiny REST API permite exponer un dataset JSON de prueba mediante una API REST configurable y adaptar tanto sus rutas como los nombres de los campos a la aplicación cliente.

De esta forma, el frontend o cliente puede continuar su desarrollo sin depender del backend definitivo.

---

## Configuración en `config.ini`

### 1. Crear un "servidor virtual"

Tiny REST API permite definir múltiples configuraciones de handlers. Cada sección de `config.ini` representa un entorno de datos o una configuración de endpoints diferente.

Por ejemplo, para configurar una librería llamada ACME:

```ini
[Book Shop ACME]
```

Puedes utilizar cualquier nombre para la sección.

---

### 2. Definir las rutas de la API

Si nuestra aplicación realiza peticiones a `/books`, pero también queremos exponer el mismo recurso mediante `/libros`, podemos definir ambas rutas:

```ini
paths=books,libros
```

Tanto `/books` como `/libros` exponen los mismos datos.

Los alias de rutas no tienen por qué utilizarse únicamente para distintos idiomas. También pueden resultar útiles para URLs antiguas, distintas convenciones entre clientes o rutas de transición de una API.

---

### 3. Indicar los métodos HTTP permitidos

Podemos especificar qué métodos HTTP están habilitados para este handler.

Por ejemplo:

```ini
methods=GET,POST,PUT,PATCH,DELETE,HEAD
```

Únicamente se expondrán los métodos configurados.

---

### 4. Modo de solo lectura

Si queremos que el endpoint se comporte como de solo lectura, podemos activar:

```ini
readonly=true
```

Cuando `readonly=true`, los métodos de escritura como `POST`, `PUT`, `PATCH` y `DELETE` son rechazados con:

```text
405 Method Not Allowed
```

Esto puede resultar útil para probar permisos o para endpoints que no deberían aceptar modificaciones.

Por defecto:

```ini
readonly=false
```

---

### 5. Modo de simulación

En algunas pruebas queremos que el cliente pueda realizar operaciones de escritura sin modificar el dataset utilizado como referencia.

Para ello podemos utilizar:

```ini
simulation=true
```

En modo simulación, las peticiones de escritura son aceptadas y devuelven su código de éxito correspondiente, pero no se aplica ningún cambio sobre los datos cargados en memoria.

Por ejemplo:

```text
DELETE /books/32
```

puede devolver una respuesta satisfactoria y, sin embargo, el libro `32` seguirá apareciendo en posteriores peticiones `GET`.

La diferencia entre los modos es:

```text
readonly=true
    las operaciones de escritura son rechazadas

simulation=true
    las operaciones de escritura son aceptadas
    los datos permanecen sin cambios

modo normal
    las operaciones de escritura modifican los datos en memoria
```

#### Por qué el modo simulación no es un dry run

El modo simulación es deliberadamente más sencillo que un dry run tradicional.

Un dry run real requeriría procesar completamente la modificación garantizando al mismo tiempo que el resultado no se consolide. Dependiendo de la implementación, esto podría requerir clonar o duplicar temporalmente una parte o la totalidad del dataset JSON.

Tiny REST API no puede saber de antemano qué tamaño tendrá el fichero JSON utilizado por el usuario ni cuánta memoria tendrá disponible la máquina donde se ejecute. Duplicar un dataset grande únicamente para simular una operación podría provocar un consumo innecesario de memoria, pérdida de rendimiento o incluso errores de asignación.

Por ese motivo, el modo simulación simplemente acepta la petición y devuelve el código de éxito correspondiente sin modificar el dataset.

Si durante las pruebas se necesitan cambios reales de estado, el usuario puede ejecutar Tiny REST API utilizando una copia desechable del fichero JSON.

---

### 6. Asignar alias a los campos del JSON

Una de las principales características de Tiny REST API es poder adaptar los campos del JSON a los nombres esperados por la aplicación cliente.

Por ejemplo, el JSON de prueba contiene un campo llamado `text60`, pero nuestro frontend espera un campo llamado `title`. También contiene `phone`, mientras que el cliente espera `telephone`.

Podemos definir los alias:

```ini
field aliases=title,text60|telephone,phone
```

Cuando la API devuelve datos, los nombres de los campos se transforman a los alias esperados por el cliente.

Cuando el cliente envía datos, los alias se convierten de nuevo a los nombres internos del JSON.

#### Ejemplo

JSON original almacenado en memoria:

```json
{
  "text60": "The Go Programming Language and the Gin Gonic with Ice",
  "phone": "25123456789"
}
```

Respuesta de la API después de aplicar los alias:

```json
{
  "title": "The Go Programming Language and the Gin Gonic with Ice",
  "telephone": "25123456789"
}
```

Esto permite que la aplicación cliente trabaje con su propio contrato de datos sin modificar el dataset original.

Puedes consultar la lista completa de campos incluidos en el JSON de ejemplo en `README.md`.

---

## Ejecutar Tiny REST API

Una vez configurado `config.ini`, inicia el servidor:

```bash
./tinyrestapi
```

También puedes indicar otro fichero de configuración:

```bash
./tinyrestapi -c ./myconfig.ini
```

o utilizar la variable de entorno:

```bash
TINYRESTAPICONFIG=/ruta/al/config.ini ./tinyrestapi
```

El fichero de configuración se resuelve siguiendo este orden:

1. `-c` o `--config`
2. `TINYRESTAPICONFIG`
3. `./config.ini`

A partir de ese momento, la aplicación cliente puede realizar peticiones contra la API mock sin depender del backend definitivo.

Todas las modificaciones normales se realizan en RAM. El fichero JSON original del disco nunca se modifica.

---

## Estado del proyecto

Tiny REST API es plenamente funcional.

Actualmente se están añadiendo algunas mejoras y reforzando la base de tests para consolidar una base sólida antes de publicar releases empaquetadas formalmente.

Por el momento, el binario precompilado y los ficheros de configuración están disponibles en `app/`, o bien el proyecto puede compilarse directamente desde el código fuente.
