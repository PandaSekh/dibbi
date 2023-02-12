// Package dibbi contains functions to manipulate left dibbi instance.
package dibbi

// Query the given database with the provided input.
func Query(query string, db *Database) (result *Results, error error) {
	ast, err := parse(query)
	if err != nil {
		return nil, err
	}

	for _, stmt := range ast.Statements {
		switch stmt.statementType {
		case CreateTableType:
			err = (*db).CreateTable(ast.Statements[0].createTableStatement)
			if err != nil {
				return nil, err
			}
		case InsertType:
			err = (*db).Insert(stmt.insertStatement)
			if err != nil {
				return nil, err
			}
		case SelectType:
			results, err := (*db).Select(stmt.selectStatement)
			if err != nil {
				return nil, err
			}
			return results, nil
		}
	}
	return nil, nil
}
