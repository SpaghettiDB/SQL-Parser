package main

import (
	"fmt"
	sqlParser "sqlParser/Parser"
)

func main() {
	// create a schema instance from sqlParser package
	schema := sqlParser.Schema{}
	parser := sqlParser.NewSQLParser(schema)
	tokens, err := parser.Tokenize("select name, id from Customers where ContactName = 'Alfred Schmidt' and age= 30;")
	if err != nil {
		panic(err)
	}

	fmt.Println("Tokens: ", tokens)

}
