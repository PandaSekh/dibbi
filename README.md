# 🐳 dibbi

[![Build](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml)

In-memory non-persistent relational database.

## Examples

**Create Table**

```sql
CREATE TABLE users
(
  age  int,
  name text
);
```

**Insert**

```sql
INSERT INTO users
VALUES (24, 'Alessio');
```

**Select**

```sql
SELECT age
from users;
```

**Select All**

```sql
SELECT *
from users;
```

## Run

`go run main.go`

## Test

`go test -v ./...`

## Features

- `SELECT`
- CREATE TABLE
  - int
  - text
  - bool
- INSERT
- `SELECT *` to get all Columns

## TODO

- Insert with column specification
- `WHERE` clause
- More column types (uuid, date)
- Automatic uuid on insertion, if specified during `CREATE TABLE`
- default values
- Support float

## Resources

- Idea and tutorial followed for the beginning: [database basics](https://notes.eatonphil.com/database-basics.html).
- [Writing an interpreter in go](https://interpreterbook.com/)