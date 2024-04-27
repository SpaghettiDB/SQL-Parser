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

// ParsedStmt represents the parsed SQL statement.
type ParsedStmt interface {
	GetQueryType() QueryType    // GetQueryType returns the type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	GetTables() []string        // GetTables returns the tables involved in the statement
	GetColumns() []string       // GetColumns returns the columns referenced in the statement
	GetConditions() []Condition // GetConditions returns the conditions specified in the statement
	GetValues() []string        // GetValues returns the values specified in the statement
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
// tokenize breaks down the SQL string into individual tokens.
// It returns a slice of tokens and an error if the tokenization fails.
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

			whereKeyword := strings.Split(tablesString, "WHERE")
			if len(whereKeyword) > 1 {
				//split at where keyword
				tablesString = strings.Split(tablesString, "WHERE")[0]
				tablesString = strings.Trim(tablesString, " ")
				if len(strings.Split(tablesString, " ")) > 1 {
					return nil, errors.New("invalid table name in SELECT statement")
				}
				
				if len(tablesString) == 0 {
					return nil, errors.New("missing table name in SELECT statement")
				}
				//add the table to the tokens list
				tokens = append(tokens, Token{Type: "table", Value: tablesString})
				//parse the conditions
				tokens, err = parseCondition(tokens, sql)
				if err != nil {
					return nil, err
				}
			} else {
				tablesString = strings.Trim(tablesString, " ")
				if len(strings.Split(tablesString, " ")) > 1 {
					return nil, errors.New("invalid table name in SELECT statement")
				}
				if len(tablesString) == 0 {
					return nil, errors.New("missing table name in SELECT statement")
				}
				tablesString = strings.Replace(tablesString, ";", "", -1)
				//add the table to the tokens list
				tokens = append(tokens, Token{Type: "table", Value: tablesString})
			}

		case DeleteQuery: 
			// DELETE FROM table_name WHERE condition;
			var tablesString string = strings.Split(sql, "FROM")[1]
			if len(tablesString) == 0 {
				return nil, errors.New("missing table name in DELETE statement")
			}
			tablesString = strings.Split(tablesString, "WHERE")[0]
			tablesString = strings.Trim(tablesString, " ")
			if len(strings.Split(tablesString, " ")) > 1 {
				return nil, errors.New("invalid table name in DELETE statement")
			}

			if len(tablesString) == 0 {
				return nil, errors.New("missing table name in DELETE statement")
			}
			tokens = append(tokens, Token{Type: "table", Value: tablesString})

			tokens, err = parseCondition(tokens, sql)
			if err != nil {
				return nil, err
			}

		case UpdateQuery:
			// UPDATE table_name SET col1 = value1, col2 = value2, ... WHERE condition;
			var tablesString string = strings.Split(sql, "SET")[0]

			tablesString = strings.Trim(tablesString, " ")
			if len(strings.Split(tablesString, " ")) > 1 {
				return nil, errors.New("invalid table name in UPDATE statement")
			}

			if len(tablesString) == 0 {
				return nil, errors.New("missing table name in UPDATE statement")
			}
			tokens = append(tokens, Token{Type: "table", Value: tablesString})

			var colsString string = strings.Split(sql, "SET")[1]
			colsString = strings.Split(colsString, "WHERE")[0]

			if len(colsString) == 0 {
				return nil, errors.New("missing columns in UPDATE statement")
			}
			colsString = strings.Replace(colsString, " ", "", -1)
			columns := strings.Split(colsString, ",")

			for _, col := range columns {
				operator := strings.Split(col, "=")
				if len(operator) != 2 {
					return nil, errors.New("invalid column value pair in UPDATE statement")
				}
				val := strings.Split(col, "=")[1]
				col = strings.Split(col, "=")[0]
				tokens = append(tokens, Token{Type: "column", Value: col})
				tokens = append(tokens, Token{Type: "value", Value: val})
			}

			tokens, err = parseCondition(tokens, sql)
			if err != nil {
				return nil, err
			}

		case InsertQuery:
			// INSERT INTO table_name (col1, col2, ...) VALUES (value1, value2, ...);
			parenthesis := strings.Split(sql, "(")
			if len(parenthesis) < 3 {
				return nil, errors.New("missing parenthesis in INSERT statement")
			}
			var tablesString string = strings.Split(sql, "INTO")[1]
			tablesString = strings.Split(tablesString, "(")[0]
			tablesString = strings.Trim(tablesString, " ")
			if len(strings.Split(tablesString, " ")) > 1 {
				return nil, errors.New("invalid table name in INSERT statement")
			}

			tablesString = strings.Replace(tablesString, " ", "", -1)
			if len(tablesString) == 0 {
				return nil, errors.New("missing table name in INSERT statement")
			}
			tokens = append(tokens, Token{Type: "table", Value: tablesString})

			parenthesis = strings.Split(sql, ")")
			if len(parenthesis) < 3 {
				return nil, errors.New("missing parenthesis in INSERT statement")
			}

			var colsString string = strings.Split(sql, "(")[1]
			colsString = strings.Split(colsString, ")")[0]
			colsString = strings.Replace(colsString, " ", "", -1)
			columns := strings.Split(colsString, ",")
			if len(columns) == 0 {
				return nil, errors.New("missing columns in INSERT statement")
			}

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
			words := strings.Fields(sql);
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
			words := strings.Fields(sql);
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
		err := parser.validateColumnExistence(tables, parsedStmt.GetColumns())
		if err != nil {
			return nil, err
		}
		// check if the conditions exist in the table
		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return nil, err
		}
		return parsedStmt, nil

	case InsertQuery:
		// check if the columns exist in the table
		err := parser.validateColumnExistence(tables, parsedStmt.GetColumns())
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

		err := parser.validateColumnExistence(tables, parsedStmt.GetColumns())
		if err != nil {
			return nil, err
		}

		err = parser.validateConditionExistence(tables, Conditions)
		if err != nil {
			return nil, err
		}
		// check if the values are of the correct type (TODO: implement this)

		if len(parsedStmt.GetValues()) == 0 || (len(parsedStmt.GetValues()) != len(parsedStmt.GetColumns())) {
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

func (parser *SQLParser) validateColumnExistence(tables []string, columns []string) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return errors.New(fmt.Sprintf("Table %s has no columns", table))
		}

		if len(columns) == 1 && columns[0] == "*" {
			continue
		}

		valid, column := ContainsAll(tableColumns, columns)

		if !valid {
			return errors.New(fmt.Sprintf("Column %s does not exist in table", column))
		}
	}
	return nil
}

func (parser *SQLParser) validateConditionExistence(tables []string, conditions []Condition) error {
	for _, table := range tables {
		tableColumns := parser.Schema.GetTableColumns(table)

		if len(tableColumns) == 0 {
			return errors.New(fmt.Sprintf("Table %s has no columns", table))
		}

		for _, condition := range conditions {
			valid, column := ContainsAll(tableColumns, []string{condition.Column})
			if !valid {
				return errors.New(fmt.Sprintf("Column %s in Condition does not exist in table", column))
			}
			// check if the operator is valid
			valid, operator := ContainsAll([]string{"=", ">", "<", ">=", "<=", "!="}, []string{condition.Operator})
			if !valid {
				return errors.New(fmt.Sprintf("Invalid operator %s in Condition", operator))
			}
		}
	}
	return nil
}

func parseCondition(token []Token, sql string) ([]Token, error) {
	// get the substring after the first "WHERE" keyword
	// WHERE condition1 AND condition2 AND condition3 ...;
	var conditionString string = strings.Split(sql, "WHERE")[1]
	conditionString = strings.Replace(conditionString, " ", "", -1)
	conditionString = strings.Replace(conditionString, ";", "", -1)
	if len(conditionString) == 0 {
		return nil, errors.New("missing condition in WHERE clause")
	}

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

	return token, nil
}
