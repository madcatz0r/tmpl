package order

import "strings"

func Desc(fields ...string) string {
	return strings.Join(fields, " DESC,") + " DESC"
}

// not necessary func, default sort is ASC
func Asc(fields ...string) string {
	return strings.Join(fields, " ASC,") + " ASC"
}
