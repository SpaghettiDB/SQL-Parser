package sqlParser


type ParsedStmtInterface interface {
	GetQueryType() QueryType                  // GetQueryType returns the type of the query (e.g., SELECT, INSERT, UPDATE, DELETE)
	GetTables() []string                      // GetTables returns the tables involved in the statement
	GetColumns() []string                     // GetColumns returns the columns referenced in the statement
	GetConditions() []interface{} // GetConditionsAndOperators returns the conditions and operators specified in the statement
	GetValues() []string                      // GetValues returns the values specified in the statement
}



func (stmt *ParsedStmt ) GetQueryType() QueryType {
	return stmt.QueryType
}

func (stmt *ParsedStmt ) GetTables() []string {
	return stmt.Tables
}

func (stmt *ParsedStmt ) GetColumns() []string {
	return stmt.Columns
}

func (stmt *ParsedStmt ) GetConditions() []interface{} {
	return stmt.Conditions
}

func (stmt *ParsedStmt ) GetValues() []string {
	return stmt.Values
}