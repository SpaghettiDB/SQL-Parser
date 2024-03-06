/* this file contains the schema information of the database will help you 
in the semantic analysis of the SQL statement.

*/

package parser

// Schema represents the database schema information.
type ISchema interface {
    // Methods to access tables, columns, data types, constraints, etc.
    GetTableColumns(tableName string) []string        // GetTableColumns returns columns of a given table
    GetColumnDataType(tableName, columnName string) string // GetColumnDataType returns data type of a column
}

// SampleSchema represents a sample database schema.
type Schema struct {
    Tables map[string][]string // Maps table names to column names
    // You can add more fields to represent constraints, data types, etc.
}

// GetTableColumns returns the columns of a given table.
func (schema *Schema) GetTableColumns(tableName string) []string {
    if columns, ok := schema.Tables[tableName]; ok {
        return columns
    }
    return nil
}

// GetColumnDataType returns the data type of a column in a given table.
func (schema *Schema) GetColumnDataType(tableName, columnName string) string {
    // Implementation omitted for brevity, you would typically retrieve data type from the schema
    return "varchar(255)" // Sample return type
}