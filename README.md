# 🐳 dibbi

[![Build](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml)

In-memory non persistent database.

Based on [this great blog post](https://notes.eatonphil.com/database-basics.html).  
I started by following the tutorial with the intention of adding new features down the line as a learning process.

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
- INSERT
- `SELECT *` to get all Columns
- Automatic migrations on startup

## TODO

- Insert with column specification
- `WHERE` clause
- More column types (bool, uuid)
- Automatic uuid on insertion, if specified during `CREATE TABLE`
- Support float