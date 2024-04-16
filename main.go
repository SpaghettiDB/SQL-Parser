package main

import (
	"fmt"
	"os"
	sqlParser "sqlParser/Parser"
	"strings"
)

func main() {


	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <query_type>")
		return
	}

	// Extract the query type from command-line arguments
	queryType := strings.ToUpper(os.Args[1])


	// create a schema instance from sqlParser package
	schema := sqlParser.Schema{}
	parser := sqlParser.NewSQLParser(schema)


	// Define example queries for each query type
	var query string
	switch queryType {
	case "SELECT":
		query = "SELECT a , b, c, d,e  FROM Customers WHERE Country = 'USA';"
	case "INSERT":
		query = "INSERT INTO Customers (CustomerName, ContactName, Country) VALUES ('John Doe', 'John Smith', 'UK');"
	case "UPDATE":
		query = "UPDATE Customers SET Country = 'Germany' WHERE CustomerID = 1;"
	case "DELETE":
		query = "DELETE FROM Customers WHERE CustomerID = 1;"
	case "DROP":
		query = "DROP TABLE Customers;"
	case "DROPINDEX":
		query = "DROP index Customers;"
	case "CREATE":
		query = "CREATE TABLE Customers (CustomerID, CustomerName , Country );"	
	case "CREATEINDEX":
		query = "CREATE index index_anme on Customers (CustomerID, CustomerName , Country );"
	default:
		fmt.Println("Unsupported query type.")
		return
	}

	tokens, err := parser.Tokenize(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tokens: ", tokens)

}
