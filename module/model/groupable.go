package model

import "strings"

const GROUP_KEY_DELIMITER = "|-|"

type Groupable interface {
	GetField(string) string
}

// Generic GroupByKey function
func GroupByKey[T Groupable](c T, groupBy []string) string {
	pieces := []string{}
	for _, key := range groupBy {
		v := c.GetField(key)
		pieces = append(pieces, v)
	}
	return strings.Join(pieces, GROUP_KEY_DELIMITER)
}
