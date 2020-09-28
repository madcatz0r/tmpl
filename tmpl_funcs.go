package templates

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

const (
	update = iota
	insert
	upsert
)

func (t Tmpl) values(m interface{}, op int) (result []interface{}) {
	var fields []string
	switch op {
	case insert:
		fields = t.InsertFields
		result = make([]interface{}, 0, len(fields))
	case update:
		fields = t.UpdateFields
		result = make([]interface{}, 0, len(fields)+1)
	case upsert:
		fields = make([]string, len(t.UpsertFields), len(t.UpsertFields))
		for key, value := range t.UpsertFields {
			if value != 0 {
				fields[value-1] = key
			} else {
				fields = fields[:len(fields)-1]
			}
		}
		result = make([]interface{}, 0, len(fields))
	}

	var str reflect.Value
	if t := reflect.TypeOf(m); t.Kind() == reflect.Ptr {
		str = reflect.ValueOf(m).Elem()
	} else {
		str = reflect.ValueOf(m)
	}
	for _, item := range fields {
		field := str.FieldByName(t.Mapping[item])
		switch field.Type().String() {
		case "time.Time", "*time.Time":
			def, ok := t.Default[item]
			if !ok {
				result = append(result, field.Interface())
				continue
			}
			// prevents panic on value.Interface call
			if field.Interface() == nil {
				t.Default[item] = parseTimeDiff(time.Time{}, def)
				continue
			}
			t.Default[item] = parseTimeDiff(reflect.Indirect(field).Interface().(time.Time), def)
		default:
			result = append(result, field.Interface())
		}
	}
	if op == update {
		field := str.FieldByName(t.Mapping[t.Primary])
		result = append(result, field.Interface())
	}
	return result
}

func (t *Tmpl) upsertMap() {
	counter := 1
	for _, item := range t.InsertFields {
		if _, ok := t.Default[item]; !ok {
			t.UpsertFields[item] = counter
			counter++
		} else {
			t.UpsertFields[item] = 0
		}
	}

	for _, item := range t.UpdateFields {
		if _, found := t.UpsertFields[item]; !found {
			if _, ok := t.Default[item]; !ok {
				t.UpsertFields[item] = counter
				counter++
			} else {
				t.UpsertFields[item] = 0
			}
		}
	}
}

// date.Time fields //
func SetTimeDiff(dur time.Duration) time.Time {
	zero := &time.Time{}
	return zero.Add(dur)
}

func parseTimeDiff(t time.Time, defaultString string) string {
	if t.IsZero() {
		return defaultString
	}
	zero := &time.Time{}
	diff := t.Sub(*zero).Seconds()
	if math.Signbit(diff) {
		return fmt.Sprintf("%s - %.0f * interval '1 second'", defaultString, math.Ceil(math.Abs(diff)))
	}
	return fmt.Sprintf("%s + %.0f * interval '1 second'", defaultString, math.Ceil(diff))
}

/////////////////////

func structName(str interface{}) string {
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
