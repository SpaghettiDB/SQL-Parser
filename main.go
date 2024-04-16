package main

import (
	"fmt"
	sqlParser "sqlParser/Parser"
)

func main() {
	// create a schema instance from sqlParser package
	schema := sqlParser.Schema{}
	parser := sqlParser.NewSQLParser(schema)
	tokens, err := parser.Tokenize("INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country) VALUES (Cardinal, Tom B. Erichsen, Skagen 21, Stavanger, 4006, Norway);")
	if err != nil {
		panic(err)
	}

	fmt.Println("Tokens: ", tokens)

}
