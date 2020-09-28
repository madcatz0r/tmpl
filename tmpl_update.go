package templates

import (
	"bytes"
)

func (t Tmpl) UpdateQuery() (string, error) {
	var updateAllQuery = new(bytes.Buffer)
	err := updateAllTmpl.Execute(updateAllQuery, t)
	return updateAllQuery.String(), err
}

func Update(m interface{}, fields ...string) (string, []interface{}, error) {
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
	values := t.values(m, update)
	qry, err := t.UpdateQuery()
	if err != nil {
		return "", nil, err
	}
	return qry, values, err
}
