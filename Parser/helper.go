package sqlParser

import (
	"errors"
	"strings"
)

// func Contains[T any](arr []any, element int) bool {
func ContainsAll(arr []string, elements []string) bool {
	for _, element := range elements {
		found := false
		for _, value := range arr {
			if value == element {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func paresQueryType(queryType string) (QueryType, error) {
	//this should work with upper and lower case
	queryType = strings.ToUpper(queryType)
	switch queryType {
	case "SELECT", "INSERT", "UPDATE", "DELETE", "DROP":
		return QueryType(queryType), nil
	default:
		return "", errors.New("Invalid query type")
	}
}

func isKeyword(word string) bool {
	//check if the word is a keyword
	word = strings.ToUpper(word)
	return word == "SELECT" || word == "FROM" || word == "WHERE" || word == "INSERT" || word == "UPDATE" || word == "DELETE" || word == "DROP"
}

func keyWordsToUpperCase(query string) string {
	//split the query and check each word if it is a keyword convert it to upper case
	words := strings.Fields(query)
	for i, word := range words {
		if isKeyword(word) {
			words[i] = strings.ToUpper(word)
		}
	}
	return strings.Join(words, " ")
}
