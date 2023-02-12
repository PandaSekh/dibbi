// Package dibbi contains functions to manipulate a dibbi instance.
package dibbi

// Query the database with the given input.
func Query(query string, db *Database) (result *Results, isResultPresent bool, error error) {
	ast, err := parse(query)
	if err != nil {
		return nil, false, err
	}

	for _, stmt := range ast.Statements {
		switch stmt.statementType {
		case CreateTableType:
			err = (*db).CreateTable(ast.Statements[0].createTableStatement)
			if err != nil {
				return nil, false, err
			}
		case InsertType:
			err = (*db).Insert(stmt.insertStatement)
			if err != nil {
				return nil, false, err
			}
		case SelectType:
			results, err := (*db).Select(stmt.selectStatement)
			if err != nil {
				return nil, false, err
			}
			return results, true, nil
		}
	}
	return nil, false, nil
}
