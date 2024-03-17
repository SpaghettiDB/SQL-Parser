package sqlParser 

type TokenType string

const (
	ColumnToken    TokenType = "column"
	TableToken     TokenType = "table"
	ConditionToken TokenType = "condition"
	KeywordToken   TokenType = "keyword"
	ValueToken     TokenType = "value"
)

type QueryType string

const (
    SelectQuery QueryType = "SELECT"
    InsertQuery QueryType = "INSERT"
    UpdateQuery QueryType = "UPDATE"
    DeleteQuery QueryType = "DELETE"
    DropQuery   QueryType = "DROP"
)

func ContainsAll(arr []string, elements []string) (bool, string) {
	for _, element := range elements {
		found := false
		for _, value := range arr {
			if value == element {
				found = true
				break
			}
		}
		if !found {
			return false, element
		}
	}
	return true, ""
}

