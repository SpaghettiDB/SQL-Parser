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
	GetQueryType() QueryType    // GetQueryType returns the type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	GetTables() []string        // GetTables returns the tables involved in the statement
	GetColumns() []string       // GetColumns returns the columns referenced in the statement
	GetConditions() []Condition // GetConditions returns the conditions specified in the statement
	GetValues() []string   // GetValues returns the values specified in the statement
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

	tables := parsedStmt.GetTables()
	Conditions := parsedStmt.GetConditions()

	// check if the tables exist in the schema
	err = parser.validateTableExistence(tables)
	if err != nil {
		return nil, err
	}

	switch parsedStmt.GetQueryType() {
	case SelectQuery:
		// check if the columns exist in the table
		err := parser.validateColumnExistence(tables,parsedStmt.GetColumns())
		if err != nil {
			return nil, err
		}
		// check if the conditions exist in the table
		err = parser.validateConditionExistence(tables,Conditions)
		if err != nil {
			return nil, err
		}
		return parsedStmt, nil

	case InsertQuery:
		// check if the columns exist in the table
		err := parser.validateColumnExistence(tables,parsedStmt.GetColumns())
		if err != nil {
			return nil, err
		}
		// check if the values count is equal to the columns count
		if len(parsedStmt.GetColumns()) != len(parsedStmt.GetValues()) {
			return nil, errors.New("Values count does not match the columns count")
		}
		// check if the values are of the correct type (TODO: implement this)
		return parsedStmt, nil

	case UpdateQuery:
		
		err := parser.validateColumnExistence(tables,parsedStmt.GetColumns())
		if err != nil {
			return nil, err
		}

		err = parser.validateConditionExistence(tables,Conditions)
		if err != nil {
			return nil, err
		}
		// check if the values are of the correct type (TODO: implement this)

		if len(parsedStmt.GetValues()) == 0 || (len(parsedStmt.GetValues()) != len(parsedStmt.GetColumns())){
			return nil, errors.New("Values count does not match the columns count")
		}
		return parsedStmt, nil

	case DeleteQuery:
		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return nil, err
		}
		return parsedStmt, nil

	case DropQuery:
		return parsedStmt, nil

	default:
		return nil, errors.New("Invalid query type")
	}
}

func (parser *SQLParser) validateTableExistence(tables []string) error {

	if len(tables) == 0 {
		return errors.New("No tables specified")
	}

	valid, table := ContainsAll(parser.Schema.GetSchemaTables(), tables)

	if !valid {
		return errors.New(fmt.Sprintf("Table %s does not exist in the schema", table))
	}
	return nil
}

func (parser *SQLParser) validateColumnExistence(tables []string,columns []string) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return  errors.New(fmt.Sprintf("Table %s has no columns", table))
		}

		if len(columns) == 1 &&columns[0] == "*" {
			continue
		}
	
		valid, column := ContainsAll(tableColumns, columns)

		if !valid {
			return  errors.New(fmt.Sprintf("Column %s does not exist in table", column))
		}
	}
	return nil
}

func (parser *SQLParser) validateConditionExistence(tables []string,conditions []Condition) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return  errors.New(fmt.Sprintf("Table %s has no columns", table))
		}

		for _, condition := range conditions {
			valid, column := ContainsAll(tableColumns, []string{condition.Column})
			if !valid {
				return  errors.New(fmt.Sprintf("Column %s in Condition does not exist in table", column))
			}
			// check if the operator is valid
			valid, operator := ContainsAll([]string{"=", ">", "<", ">=", "<=", "!="}, []string{condition.Operator})
			if !valid {
				return  errors.New(fmt.Sprintf("Invalid operator %s in Condition", operator))
			}
		}
	}
	return nil
}