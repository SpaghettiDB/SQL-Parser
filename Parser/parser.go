package sqlParser

import (
	"errors"
	"fmt"
)

// some struct and interface definitions can be separated into another file later on
type Token struct {
	Value string
	Type  TokenType
}

// ParsedStmt represents the parsed SQL statement.
type ParsedStmt interface {
	GetQueryType() QueryType       // GetQueryType returns the type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	GetTables() []string        // GetTables returns the tables involved in the statement
	GetColumns() []string       // GetColumns returns the columns referenced in the statement
	GetConditions() []Condition // GetConditions returns the conditions specified in the statement
}

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
func (parser *SQLParser) ParseSQL(sql string) (parsedStmt ParsedStmt, err error) {
	tokens, err := parser.tokenize(sql)
	if err != nil {
		return nil, err
	}

	tokenError, err := parser.syntaxCheck(tokens)

	if err != nil {
		return nil, err
	}
	// remove it
	fmt.Println("Tokens: ", tokenError)

	parsedStmt, err = parser.parse(tokens)
	if err != nil {
		return nil, err
	}

	return parser.semanticAnalysis(parsedStmt)
}

// tokenize breaks down the SQL string into individual tokens.
func (parser *SQLParser) tokenize(sql string) ([]Token, error) {
	// Implement tokenization logic using regular expressions or a state machine
	// ...
	return nil, nil
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
func (parser *SQLParser) parse(tokens []Token) (parsedStmt ParsedStmt, err error) {

	return nil, nil
}

// semanticAnalysis interprets the parsed structure and assigns meaning based on the schema.
func (parser *SQLParser) semanticAnalysis(parsedStmt ParsedStmt) (res ParsedStmt, err error) {

	// then we need to check if the column is of the correct type if not return error
	// check on the type of the query and if the query is valid or not

	tables := parsedStmt.GetTables()
	Conditions := parsedStmt.GetConditions()

	// check if the table exists in the schema since it is common between all the queries
	if len(tables) == 0 {
		return nil, errors.New("No tables specified")
	}
	valid, table := ContainsAll(parser.Schema.GetSchemaTables(), tables)
	if !valid {
		return nil, errors.New(fmt.Sprintf("Table %s does not exist", table))
	}


	// check if the conditions are valid or not



	// check if the parsedStmt meets the requirements of the query type
	switch parsedStmt.GetQueryType() {
	case SelectQuery:

	case InsertQuery:
			// check if the column exists in the table
		for _, table := range tables {
			tableColumns := parser.Schema.GetTableColumns(table)
			// check if the table has columns
			if len(tableColumns) == 0 {
				return nil, errors.New(fmt.Sprintf("Table %s has no columns", table))
			}
			valid, column := ContainsAll(tableColumns, parsedStmt.GetTableColumns(table))
			if !valid {
				return nil, errors.New(fmt.Sprintf("Column %s does not exist in table", column))
			}
		}

		for _, condition := range Conditions {
			// check if the column exists in the table
			// [TO BE REFACTORED LATER ON]
			valid, column := ContainsAll(parsedStmt.GetTableColumns(parsedStmt.GetTables()[0]), []string{condition.Column})
			if !valid {
				return nil, errors.New(fmt.Sprintf("Column %s in Condition does not exist in table", column))
			}
			// check if the operator is valid
			valid, operator := ContainsAll([]string{"=", ">", "<", ">=", "<=", "!="}, []string{condition.Operator})
			if !valid {
				return nil, errors.New(fmt.Sprintf("Invalid operator %s in Condition", operator))
			}

		}
		// check if the values count is equal to the columns count
		if len(parsedStmt.GetTableColumns(parsedStmt.GetTables()[0])) != len(parsedStmt.GetTableColumns(parsedStmt.GetTables()[0])) {
			return nil, errors.New("Values count does not match the columns count")
		}
		// check if the value is of th e correct type
		// [TO BE REFACTORED LATER ON]

		//check if 




	case UpdateQuery:
		// check if the columns exist in the table
		// check if the columns are of the correct type
		// check if the conditions are valid
		// check if the conditions are of the correct type
		// check if the conditions are of the correct type
		// check if the conditions are of the correct type
	case DeleteQuery:
		// check if the columns exist in the table
		// check if the columns are of the correct type
		// check if the conditions are valid
		// check if the conditions are of the correct type
		// check if the conditions are of the correct type
		// check if the conditions are of the correct type
	case DropQuery:


	default:
		return nil, errors.New("Invalid query type")
	}
	return parsedStmt, nil
}

// funciton to ccheck if column exists in the table
// func (parser *SQLParser)  checkColumnExists(column string, table []string) bool {
// 	// check if the column exists in the table
// 	for _, table := range tables {
// 		tableColumns := parser.Schema.GetTableColumns(table)
// 		// check if the table has columns
// 		if len(tableColumns) == 0 {
// 			return nil, nil, errors.New(fmt.Sprintf("Table %s has no columns", table))
// 		}
// 		valid, column := ContainsAll(tableColumns, parsedStmt.GetTableColumns(table))
// 		if !valid {
// 			return nil, nil, errors.New(fmt.Sprintf("Column %s does not exist in table", column))
// 		}
// 	}
// }
// // 	return true
// // }


/*


sample query with conditions

SELECT * FROM table WHERE column = 1 AND column2 = 2




*/
