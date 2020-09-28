package templates

import (
	"bytes"
)

func (t Tmpl) UpsertQuery() (string, error) {
	if len(t.UpsertFields) == 0 {
		t.upsertMap()
	}
	var upsertQuery = new(bytes.Buffer)
	err := upsertTmpl.Execute(upsertQuery, t)
	return upsertQuery.String(), err
}

func Upsert(m interface{}, fields ...string) (string, []interface{}, error) {
	t, err := GetTmpl(m)
	if err != nil {
		return "", nil, err
	}
	if fields != nil {
		t.UpdateFields = stripFieldNames(fields)
		if t.MustUpdate != "" {
			missedUpdate := true
			for _, item := range fields {
				if item == t.MustUpdate {
					missedUpdate = false
					break
				}
			}
			if missedUpdate {
				t.UpdateFields = append(t.UpdateFields, t.MustUpdate)
			}
		}
	}
	t.upsertMap()
	values := t.values(m, upsert)
	qry, err := t.UpsertQuery()
	if err != nil {
		return "", nil, err
	}
	return qry, values, err
}
