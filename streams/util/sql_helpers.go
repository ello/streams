package util

import "strings"

// ToISO8601Sql constructs the sql to pull a timestamp w/o tz in ISO8601
func ToISO8601Sql(fieldName string) string {
	return strings.Join([]string{"to_char(", fieldName, "at time zone 'UTC', 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as", fieldName}, " ")
}
