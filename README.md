# go_dibbi

Based on [this great blog post](https://notes.eatonphil.com/database-basics.html).  
I started by following the tutorial with the intention of adding new features down the line as a learning process.

Currently supported queries: 
- SELECT (basic)
- CREATE TABLE (int and text columns)
- INSERT (needs to specify columns)

Features:
- Automatic migrations on startup

TODO List:
- `SELECT *` to get all columns
- `WHERE` clause
- More column types (bool, uuid)
- Automatic uuid on insertion, if specified during `CREATE TABLE`
- Use a persistent backend
