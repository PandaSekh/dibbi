# 🐳 dibbi
[![Build](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/PandaSekh/dibbi/actions/workflows/build_and_test.yml)

In-memory non persistent database.

Based on [this great blog post](https://notes.eatonphil.com/database-basics.html).  
I started by following the tutorial with the intention of adding new features down the line as a learning process.

## Run

## Test

## Features
- `SELECT`
- CREATE TABLE (int and text columns)
- INSERT (needs to specify columns)
- `SELECT *` to get all columns
- Automatic migrations on startup

## TODO
- ~~`SELECT *` to get all columns~~
- `WHERE` clause
- More column types (bool, uuid)
- Automatic uuid on insertion, if specified during `CREATE TABLE`