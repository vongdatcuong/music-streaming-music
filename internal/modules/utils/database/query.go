package database_utils

import "fmt"

func BuildSqlQueryStr(filter map[string]any, operator string) ([]string, []interface{}) {
	var wheres []string
	var values []interface{}
	for key, value := range filter {
		wheres = append(wheres, fmt.Sprintf("%s %s ?", key, operator))
		values = append(values, value)
	}

	return wheres, values
}

func GetSqlWhereClause(whereStr string) string {
	if whereStr != "" {
		return "WHERE " + whereStr
	}
	return ""
}
