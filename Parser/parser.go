package sqlParser

import (
	"errors"
	"fmt"
	"strings"
)

// some struct and interface definitions can be separated into another file later on
type Token struct {
	Value string
	Type  TokenType
}

type ParsedStmt struct {
	QueryType  QueryType     // Type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	Tables     []string      // Tables involved in the statement
	Columns    []string      // Columns referenced in the statement
	Values     []string      // Values specified in the statement
	Conditions []interface{} // Conditions and operators specified in the statement
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

	tokens, err := parser.Tokenize(sql)
	if err != nil {
		return ParsedStmt{}, err
	}

	tokenError, err := parser.validateSyntax(tokens)
	if err != nil {
		return ParsedStmt{}, err
	}

	fmt.Println("Tokens: ", tokenError)

	parsedStmtPtr, err := parser.parse(tokens)
	if err != nil {
		return ParsedStmt{}, err
	}

	parsedStmt, err = parser.semanticAnalysis(*parsedStmtPtr)
	if err != nil {
		return ParsedStmt{}, err
	}

	return parsedStmt, nil
}

func (parser *SQLParser) Tokenize(sql string) ([]Token, error) {
	sql = strings.ToUpper(sql)
	fmt.Println("sql : ", sql)
	var tokens []Token

	// get the first token from the input string reading from start to the first space
	var firstStatement string = strings.Split(sql, " ")[0]

	queryType, err := paresQueryType(firstStatement)

	if err != nil {
		return nil, errors.New("invalid query type")
	}

	//add the first token to the tokens list
	tokens = append(tokens, Token{Type: "keyword", Value: firstStatement})
	//remove the first word from the string
	sql = strings.Replace(sql, firstStatement, "", 1)

	switch queryType {
	case SelectQuery:
		// SELECT col1, col2, ... FROM table_name WHERE condition;
		//get the collection of column names
		//get the substring that is between the first space and the first "FROM" keyword
		var colsString string = strings.Split(sql, "FROM")[0]
		//remove the spaces from the string
		colsString = strings.Replace(colsString, " ", "", -1)
		//split the string by the comma
		columns := strings.Split(colsString, ",")
		//add the columns to the tokens list
		for _, col := range columns {
			tokens = append(tokens, Token{Type: "column", Value: col})
		}
		//get the substring after the first "FROM" keyword
		var tablesString string = strings.Split(sql, "FROM")[1]
		//split at where keyword
		tablesString = strings.Split(tablesString, "WHERE")[0]
		//trim the spaces
		tablesString = strings.Replace(tablesString, " ", "", -1)
		//add the table to the tokens list
		tokens = append(tokens, Token{Type: "table", Value: tablesString})

		tokens = parseCondition(tokens, sql)

	case DeleteQuery:
		// DELETE FROM table_name WHERE condition;
		var tablesString string = strings.Split(sql, "FROM")[1]
		tablesString = strings.Split(tablesString, "WHERE")[0]
		tablesString = strings.Replace(tablesString, " ", "", -1)
		tokens = append(tokens, Token{Type: "table", Value: tablesString})

		tokens = parseCondition(tokens, sql)

	case UpdateQuery:
		// UPDATE table_name SET col1 = value1, col2 = value2, ... WHERE condition;
		var tablesString string = strings.Split(sql, "SET")[0]
		tablesString = strings.Replace(tablesString, " ", "", -1)
		tokens = append(tokens, Token{Type: "table", Value: tablesString})

		var colsString string = strings.Split(sql, "SET")[1]
		colsString = strings.Split(colsString, "WHERE")[0]
		colsString = strings.Replace(colsString, " ", "", -1)
		columns := strings.Split(colsString, ",")

		for _, col := range columns {
			val := strings.Split(col, "=")[1]
			col = strings.Split(col, "=")[0]
			tokens = append(tokens, Token{Type: "column", Value: col})
			tokens = append(tokens, Token{Type: "value", Value: val})
		}

		tokens = parseCondition(tokens, sql)

	case InsertQuery:
		// INSERT INTO table_name (col1, col2, ...) VALUES (value1, value2, ...);
		var tablesString string = strings.Split(sql, "INTO")[1]
		tablesString = strings.Split(tablesString, "(")[0]
		tablesString = strings.Replace(tablesString, " ", "", -1)
		tokens = append(tokens, Token{Type: "table", Value: tablesString})

		var colsString string = strings.Split(sql, "(")[1]
		colsString = strings.Split(colsString, ")")[0]
		colsString = strings.Replace(colsString, " ", "", -1)
		columns := strings.Split(colsString, ",")

		var valuesString string = strings.Split(sql, "VALUES")[1]
		valuesString = strings.Split(valuesString, "(")[1]
		valuesString = strings.Split(valuesString, ")")[0]
		valuesString = strings.Replace(valuesString, " ", "", -1)
		values := strings.Split(valuesString, ",")

		if len(columns) != len(values) {
			return nil, errors.New("invalid number of columns and values")
		}

		for i := range columns {
			tokens = append(tokens, Token{Type: "column", Value: columns[i]})
			tokens = append(tokens, Token{Type: "value", Value: values[i]})
		}

	case DropQuery:
		words := strings.Fields(sql)
		if len(words) > 2 {
			return nil, errors.New("invalid DROP statement")
		}
		dropType := words[0]
		switch dropType {
		case "TABLE":
			tableName := words[1]
			tokens = append(tokens, Token{Type: "table", Value: refineFieldName(tableName)})
		case "INDEX":
			indexName := words[1]
			tokens = append(tokens, Token{Type: "index", Value: refineFieldName(indexName)})
		default:
			return nil, errors.New("invalid DROP statement")
		}

	case CreateQuery:
		words := strings.Fields(sql)
		CreateType := words[0]
		switch CreateType {
		case "TABLE":
			tableName := words[1]
			tokens = append(tokens, Token{Type: "table", Value: refineFieldName(tableName)})

			var colsString string = strings.Split(sql, "(")[1]
			colsString = strings.Split(colsString, ")")[0]
			colsString = strings.Replace(colsString, " ", "", -1)
			columns := strings.Split(colsString, ",")

			for _, col := range columns {
				tokens = append(tokens, Token{Type: "column", Value: refineFieldName(col)})
			}

		case "INDEX":
			indexName := words[1]
			tokens = append(tokens, Token{Type: "index", Value: refineFieldName(indexName)})

			var tableName string = strings.Split(sql, "ON")[1]
			tableName = strings.Split(tableName, "(")[0]
			tokens = append(tokens, Token{Type: "table", Value: refineFieldName(tableName)})

			var colsString string = strings.Split(sql, "(")[1]
			colsString = strings.Split(colsString, ")")[0]
			colsString = strings.Replace(colsString, " ", "", -1)
			columns := strings.Split(colsString, ",")

			for _, col := range columns {
				tokens = append(tokens, Token{Type: "column", Value: refineFieldName(col)})
			}
		default:
			return nil, errors.New("invalid CREATE statement :(")

		}

	default:
		return nil, errors.New("unsupported query type yet :(")
	}

	return tokens, nil
}

func (parser *SQLParser) validateSyntax(tokens []Token) ([]Token, error) {
	// Implement syntax checking logic, such as checking for unmatched parentheses, missing semicolons, etc.
	// ...

	// Example of syntax error detection:
	for _, token := range tokens {
		if token.Type == "invalid" {
			return []Token{{Value: token.Value}}, errors.New("syntax error: Invalid token found")
		}
	}

	return nil, nil
}

// parse validates the token order and structure according to the chosen grammar.
func (parser *SQLParser) parse(tokens []Token) (parsedStmt *ParsedStmt, err error) {
	// Implement parsing logic to convert tokens into a ParsedStmt object.
	var queryType QueryType
	var tables []string
	var columns []string
	var values []string
	var conditionsAndOperators []interface{}

	for _, token := range tokens {
		switch token.Type {

		case "keyword":
			queryType, err = paresQueryType(token.Value)
			if err != nil {
				return nil, err
			}

		case "table":
			tables = append(tables, token.Value)

		case "column":
			columns = append(columns, token.Value)

		case "condition":

			words := strings.Fields(token.Value)
			if len(words) != 3 {
				return nil, errors.New("syntax Error: Invalid condition")
			}
			condition := Condition{Column: words[0], Operator: words[1], Value: words[2]}
			conditionsAndOperators = append(conditionsAndOperators, condition)

		case "operator":
			op := operator{operator: token.Value}
			conditionsAndOperators = append(conditionsAndOperators, op)

		case "value":
			values = append(values, token.Value)
		default:
			return nil, errors.New("syntax Error: Invalid token type")
		}

	}

	parsedStmt = &ParsedStmt{
		QueryType:  queryType,
		Tables:     tables,
		Columns:    columns,
		Values:     values,
		Conditions: conditionsAndOperators,
	}
	return parsedStmt, nil
}

// semanticAnalysis interprets the parsed structure and assigns meaning based on the schema.
func (parser *SQLParser) semanticAnalysis(parsedStmt ParsedStmt) (res ParsedStmt, err error) {

	tables := parsedStmt.Tables
	Conditions := parsedStmt.Conditions

	// check if the tables exist in the schema
	err = parser.validateTableExistence(tables)
	if err != nil {
		return ParsedStmt{}, err
	}

	switch parsedStmt.QueryType {
	case SelectQuery:
		// check if the columns exist in the table
		err := parser.validateColumnExistence(tables, parsedStmt.Columns)
		if err != nil {
			return ParsedStmt{}, err
		}
		// check if the conditions exist in the table
		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return ParsedStmt{}, err
		}
		return parsedStmt, nil

	case InsertQuery:
		// check if the columns exist in the table
		err := parser.validateColumnExistence(tables, parsedStmt.Columns)
		if err != nil {
			return ParsedStmt{}, err
		}
		// check if the values count is equal to the columns count
		if len(parsedStmt.GetColumns()) != len(parsedStmt.Values) {
			return ParsedStmt{}, errors.New("values count does not match the columns count")
		}
		// check if the values are of the correct type (TODO: implement this)
		return parsedStmt, nil

	case UpdateQuery:

		err := parser.validateColumnExistence(tables, parsedStmt.Columns)
		if err != nil {
			return ParsedStmt{}, err
		}

		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return ParsedStmt{}, err
		}
		// check if the values are of the correct type (TODO: implement this)

		if len(parsedStmt.GetValues()) == 0 || (len(parsedStmt.GetValues()) != len(parsedStmt.Columns)) {
			return ParsedStmt{}, errors.New("values count does not match the columns count")
		}
		return parsedStmt, nil

	case DeleteQuery:
		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return ParsedStmt{}, err
		}
		return parsedStmt, nil

	case DropQuery:
		return parsedStmt, nil

	default:
		return ParsedStmt{}, errors.New("invalid query type")
	}
}

func (parser *SQLParser) validateTableExistence(tables []string) error {

	if len(tables) == 0 {
		return fmt.Errorf("no tables specified")
	}

	valid, table := ContainsAll(parser.Schema.GetSchemaTables(), tables)

	if !valid {
		return fmt.Errorf("table %s does not exist in the schema", table)
	}
	return nil
}

func (parser *SQLParser) validateColumnExistence(tables []string, columns []string) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return fmt.Errorf("table %s has no columns", table)
		}

		if len(columns) == 1 && columns[0] == "*" {
			continue
		}

		valid, column := ContainsAll(tableColumns, columns)

		if !valid {
			return fmt.Errorf("column %s does not exist in table", column)
		}
	}
	return nil
}

func (parser *SQLParser) validateConditionExistence(tables []string, conditions []interface{}) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return fmt.Errorf("table %s has no columns", table)
		}

		for _, item := range conditions {
			switch elem := item.(type) {
				case Condition:
					valid, column := ContainsAll(tableColumns, []string{elem.Column})
					if !valid {
						return fmt.Errorf("column %s in Condition does not exist in table", column)
					}
					// check if the operator is valid
					valid, operator := ContainsAll([]string{"=", ">", "<", ">=", "<=", "!="}, []string{elem.Operator})
					if !valid {
						return fmt.Errorf("invalid operator %s in Condition", operator)
					}
				case operator:
					valid, operator := ContainsAll([]string{"AND", "OR"}, []string{elem.operator})
					if !valid {
						return fmt.Errorf("invalid operator %s", operator)
					}
				default: 
					return fmt.Errorf("invalid condition type")
			}
		}
	}
	return nil
}

func parseCondition(token []Token, sql string) []Token {
	// get the substring after the first "WHERE" keyword
	// WHERE condition1 AND condition2 AND condition3 ...;
	var conditionString string = strings.Split(sql, "WHERE")[1]
	conditionString = strings.Replace(conditionString, ";", "", -1)

	// split the string by the AND or OR keyword
	andConditions := strings.Split(strings.ToLower(conditionString), "and")
	orConditions := strings.Split(strings.ToLower(conditionString), "or")

	if len(andConditions) > 1 {
		for _, cond := range andConditions {
			token = append(token, Token{Type: "condition", Value: strings.Replace(cond, " ", "", -1)})
			token = append(token, Token{Type: "operator", Value: "AND"})
		}
		token = token[:len(token)-1]
	} else if len(orConditions) > 1 {
		for _, cond := range orConditions {
			token = append(token, Token{Type: "condition", Value: strings.Replace(cond, " ", "", -1)})
			token = append(token, Token{Type: "operator", Value: "OR"})
		}
		token = token[:len(token)-1]
	} else {
		token = append(token, Token{Type: "condition", Value: strings.Replace(conditionString, " ", "", -1)})
	}

	return token
}
