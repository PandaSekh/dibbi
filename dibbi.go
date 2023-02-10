package dibbi

var mb = newMemoryBackend()

// Query the database with the given input.
func Query(query string) (result *QueryResults, isResultPresent bool, error error) {
	ast, err := parse(query)
	if err != nil {
		return nil, false, err
	}

	for _, stmt := range ast.Statements {
		switch stmt.Type {
		case CreateTableType:
			err = mb.CreateTable(ast.Statements[0].createTableStatement)
			if err != nil {
				return nil, false, err
			}
		case InsertType:
			err = mb.Insert(stmt.InsertStatement)
			if err != nil {
				return nil, false, err
			}
		case SelectType:
			results, err := mb.Select(stmt.selectStatement)
			if err != nil {
				return nil, false, err
			}
			return results, true, nil
		}
	}
	return nil, false, nil
}
