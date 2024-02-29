package parser


/* here we define the struct and interface definitions of the statments 
 that we gonna use .
*/

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
    Tables     []string
    Columns    []string
    Conditions string
}

// UpdateStmt represents an UPDATE statement.
type UpdateStmt struct {
    Table      string
    SetValues  map[string]interface{}
    Conditions string
}


// InsertStmt represents an INSERT statement.
type InsertStmt struct {
    Table   string
    Columns []string
    Values  []interface{}
}

// DeleteStmt represents a DELETE statement.
type DeleteStmt struct {
	Table      string
	Conditions string
}
