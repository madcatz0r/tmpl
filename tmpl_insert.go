package templates

import (
	"bytes"
)

func (t Tmpl) InsertQuery() (string, error) {
	var insertSql = new(bytes.Buffer)
	err := insertTmpl.Execute(insertSql, t)
	return insertSql.String(), err
}

func Insert(m interface{}) (string, []interface{}, error) {
	t, err := GetTmpl(m)
	if err != nil {
		return "", nil, err
	}
	values := t.values(m, insert)
	qry, err := t.InsertQuery()
	if err != nil {
		return "", nil, err
	}
	return qry, values, err
}
