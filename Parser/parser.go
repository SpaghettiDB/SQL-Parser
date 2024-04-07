package sqlParser

import (
	"errors"
	"fmt"
	"strings"
)

// some struct and interface definitions can be separated into another file later on

// Token represents a parsed token from the SQL statement.
type Token struct {
	Type  string // Type of the token (e.g., keyword, identifier, operator, constant)
	Value string // Value of the token
}

// ParsedStmt represents the parsed SQL statement.
type ParsedStmt interface {
	GetQueryType() string // GetQueryType returns the type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	// Methods to access the parsed statement
	GetTables() []string        // GetTables returns the tables involved in the statement
	GetColumns() []string       // GetColumns returns the columns referenced in the statement
	GetConditions() []Condition // GetConditions returns the conditions specified in the statement
}

// inteface hold stmt

type baseOperation interface {
	GetQueryType() string
	GetTables() []string
	GetColumns() []string
	GetConditions() []Condition
}

// same for all other struct and interface definitions
// func (stmt *baseOperation) GetTables() []string {
// 	return stmt.GetTables()
// }

// func (stmt *baseOperation) GetColumns() []string {
// 	return stmt.GetColumns()
// }

// func (stmt *baseOperation) GetConditions() string {
// 	return stmt.GetConditions()
// }

// SQLParser represents an SQL parser instance.
type SQLParser struct {
	Schema Schema
}

// NewSQLParser creates a new SQLParser instance.
func NewSQLParser(schema Schema) *SQLParser {
	schema.LoadSchema("schema.db")
	return &SQLParser{Schema: schema}
}

// ParseSQL parses the given SQL statement and returns the parsed representation.
func (parser *SQLParser) ParseSQL(sql string) (parsedStmt ParsedStmt, warnings []string, err error) {
	tokens, err := parser.Tokenize(sql)
	if err != nil {
		return nil, nil, err
	}

	tokenError, err := parser.syntaxCheck(tokens)

	if err != nil {
		return nil, nil, err
	}
	// remove it
	fmt.Println("Tokens: ", tokenError)

	parsedStmt, warnings, err = parser.parse(tokens)
	if err != nil {
		return nil, warnings, err
	}

	return parser.semanticAnalysis(parsedStmt)
}

// tokenize breaks down the SQL string into individual tokens.
// tokenize breaks down the SQL string into individual tokens.
// It returns a slice of tokens and an error if the tokenization fails.
func (parser *SQLParser) Tokenize(sql string) ([]Token, error) {
	sql = keyWordsToUpperCase(sql)

	fmt.Println("sql : ", sql)
	var tokens []Token

	// get the first token from the input string reading from start to the first space
	var firstStatement string = strings.Split(sql, " ")[0]

	queryType, err := paresQueryType(firstStatement)

	if err != nil {
		return nil, errors.New("Invalid Query Type")
	}

	if queryType == SelectQuery {
		//add the first token to the tokens list
		tokens = append(tokens, Token{Type: "keyword", Value: string(queryType)})
		//get the collection of column names
		//get the substring that is between the first space and the first "FROM" keyword
		var colsString string = strings.Split(sql, "FROM")[0]
		//remove the first word from the string
		colsString = strings.Replace(colsString, firstStatement, "", 1)
		//remove the spaces from the string
		colsString = strings.Replace(colsString, " ", "", -1)
		//split the string by the comma
		columns := strings.Split(colsString, ",")
		//add the columns to the tokens list
		for _, col := range columns {
			tokens = append(tokens, Token{Type: "field", Value: col})
		}
		//get the substring after the first "FROM" keyword
		var tablesString string = strings.Split(sql, "FROM")[1]
		//split at where keyword
		tablesString = strings.Split(tablesString, "WHERE")[0]
		//trim the spaces
		tablesString = strings.Replace(tablesString, " ", "", -1)
		//add the table to the tokens list
		tokens = append(tokens, Token{Type: "table", Value: tablesString})
	}

	return tokens, nil
}

// syntaxCheck checks for syntax errors in the SQL statement
func (parser *SQLParser) syntaxCheck(tokens []Token) ([]Token, error) {
	// Implement syntax checking logic, such as checking for unmatched parentheses, missing semicolons, etc.
	// ...

	// Example of syntax error detection:
	for _, token := range tokens {
		if token.Type == "invalid" {
			return []Token{{Value: token.Value}}, errors.New("Syntax error: Invalid token found")
		}
	}

	return nil, nil
}

// parse validates the token order and structure according to the chosen grammar.
func (parser *SQLParser) parse(tokens []Token) (parsedStmt ParsedStmt, warnings []string, err error) {

	return nil, nil, nil
}

// semanticAnalysis interprets the parsed structure and assigns meaning based on the schema.
func (parser *SQLParser) semanticAnalysis(parsedStmt ParsedStmt) (res ParsedStmt, warnings []string, err error) {

	// first we need to check if the table exists in the schema if not exist return error
	// then we need to check if the column exists in the table if not exist return error
	// then we need to check if the column is of the correct type if not return error
	// check on the type of the query and if the query is valid or not

	tables := parsedStmt.GetTables()
	// columns := parsedStmt.GetColumns()
	// conditions := parsedStmt.GetConditions()

	// var warningsList []string

	found := ContainsAll(parser.Schema.GetSchemaTables(), tables)
	if !found {
		return nil, nil, errors.New("Table does not exist")
	}

	return parsedStmt, nil, nil
}
