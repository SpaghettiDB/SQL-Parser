package parser


/* here we define the struct and interface definitions of the statments 
 that we gonna use .
*/

type QueryType string

const (
    SelectQuery QueryType = "SELECT"
    InsertQuery QueryType = "INSERT"
    UpdateQuery QueryType = "UPDATE"
    DeleteQuery QueryType = "DELETE"
    DropQuery   QueryType = "DROP"
)

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
    Tables     []string
    Columns    []string
    Conditions []Condition
    Limit      int
}

// UpdateStmt represents an UPDATE statement.
type UpdateStmt struct {
    Table      string
    SetValues  map[string]interface{}
    Conditions []Condition
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

type Drop struct {
    Table string
}

type Condition struct {
    Column string
    Operator string
    Value string
}