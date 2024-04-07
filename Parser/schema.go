/* this file contains the schema information of the database will help you
in the semantic analysis of the SQL statement.

*/

package sqlParser

const schemaMetaFile = "schema.db"

// Schema represents the database schema information.
type ISchema interface {
	// Methods to access tables, columns, data types, constraints, etc.
	LoadSchema() error                                     // LoadSchema loads the database schema
	GetTableColumns(tableName string) []string             // GetTableColumns returns columns of a given table
	GetColumnDataType(tableName, columnName string) string // GetColumnDataType returns data type of a column
}

// SampleSchema represents a sample database schema.
type Schema struct {
	Tables map[string][]string // Maps table names to column names
}

func (Schema *Schema) LoadSchema(SchemaName string) error {
	// Implementation omitted for brevity, you would typically load schema metadata from the disk

	// open the dir with the schema name and load the schema
	// schemaDir := filepath.Join(schemaDir, schemaName)
	// schemaFile := filepath.Join(schemaDir, schemaMetaFile)
	// schemaData, err := ioutil.ReadFile(schemaFile)
	// if err != nil {
	//     return err
	// }
	return nil
}

func (Schema *Schema) GetSchemaTables() []string {
	tableNames := make([]string, 0, len(Schema.Tables))
	for tableName := range Schema.Tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
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
