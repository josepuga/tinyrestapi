# **Tiny REST API**

## **Example of use**

Imagine you are developing a frontend for a bookshop. You have already set up your forms and pages, but you **don't have the client's database yet**.

The backend team might not have finished the API, the server might be down, or you simply need a **quick testing solution**.

One option would be to find a **sample JSON** or an **online database** and set up a temporary backend. However, this comes with several issues:

- **The field names might not match your project’s structure.**
- **You would have to modify your code or the data structure** (which adds extra work and potential future errors).

**Tiny REST API solves this problem by providing a flexible JSON that automatically adapts to your project** without needing code modifications.

---

## **🛠 Configuration in `config.ini`**

### **1️. Creating a "virtual server"**

Tiny REST API allows you to define **multiple configurations**. In `config.ini`, each section represents a different data environment.

For example, if we want to configure a server for a bookstore called **ACME**, we add a section:

```ini
[Book Shop ACME]
```

_(You can use any name for the section.)_

---

### **2️. Defining API routes**

If our application makes requests to `/books`, but we also want to support `/libros` for Spanish-speaking users, we can define both routes:

```ini
paths=books,libros
```

This means that **both `/books` and `/libros` will return the same data**, allowing **route aliases**.

---

### **3️. Specifying allowed HTTP methods**

We can define which HTTP methods are enabled for this endpoint. In this case, we allow **all methods**:

```ini
methods=GET,POST,PUT,PATCH,DELETE,HEAD
```

---

### **4️. Enabling or disabling data modifications**

If we want to keep the data **read-only** and prevent accidental modifications, we can enable `readonly` mode:

```ini
readonly=false
```

**If we try to modify or delete a record while `readonly=true`, the server will return an access error.**  
_(By default, `readonly=false`, meaning modifications are allowed.)_

This is useful for **permission testing** and **ensuring data integrity**.

---

### **5️. Assigning aliases to JSON fields**

Now comes the **most important feature**: **adapting the JSON fields to match your app’s structure**.

For example, in our **test JSON**, the field `text60` might be **used as a book title** (`title`) in our frontend, and `phone` might be used as `telephone`. We can automatically apply these aliases:

```ini
field aliases=title,text60|telephone,phone
```

_(You can find the full list of JSON fields in the `README.md`.)_

**How does it work?**

- **When the API returns data, it automatically transforms the field names to match your project.**
- **When you send data in a request, it is automatically converted to the original internal names.**

#### **Example**

**Original JSON stored in memory**

```json
{
  "text60": "The Go Programming Language and the Gin Gonic with Ice",
  "phone": "25123456789"
}
```

**API response after applying the aliases**

```json
{
  "title": "The Go Programming Language and the Gin Gonic with Ice",
  "telephone": "25123456789"
}
```

This **eliminates the need to modify your code** to match the backend's field names.

---

## **Ready? Run Tiny REST API!**

Once `config.ini` is configured, simply start the server:

```bash
./tinyrestapi
```

Now you can make requests to your API **without worrying about incorrect field names or accidentally modifying the original JSON**.

**All modifications are done in RAM, so the original file on disk is never altered.**
