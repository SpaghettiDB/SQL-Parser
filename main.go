package main

import (
	"fmt"
	sqlParser "sqlParser/Parser"
)

func main() {
	// create a schema instance from sqlParser package
	schema := sqlParser.Schema{}
	parser := sqlParser.NewSQLParser(schema)
	tokens, err := parser.Tokenize("UPDATE Customers SET ContactName = 'Alfred Schmidt', age= 30 WHERE CustomerID = 1;")
	if err != nil {
		panic(err)
	}

	fmt.Println("Tokens: ", tokens)

}
