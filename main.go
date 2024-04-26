package main

import (
	"fmt"
	sqlParser "github.com/SpaghettiDB/SQL-Parser/Parser"
)

func main() {
	// create a schema instance from sqlParser package
	schema := sqlParser.Schema{}
	parser := sqlParser.NewSQLParser(schema)
	tokens, err := parser.Tokenize("select name,age , ali FROM users WHERE id = 1")
	if err != nil {
		panic(err)
	}

	fmt.Println("Tokens: ", tokens)

}
