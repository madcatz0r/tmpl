package templates

import (
	"bytes"
	"fmt"
	"github.com/madcatz0r/tmpl/conditions"
	"strings"
)

const (
	inner string = "INNER JOIN"
	left         = "LEFT OUTER JOIN"
	right        = "RIGHT OUTER JOIN"
	full         = "FULL OUTER JOIN"
)

type sel struct {
	table   string
	fields  []string
	join    []*join
	where   string
	groupBy string
	orderBy string
	limit   int
	offset  int
	err     error
}

type selString struct {
	Table   string
	Fields  []string
	Join    []*join
	Where   string
	GroupBy string
	OrderBy string
	Limit   int
	Offset  int
	err     error
}

func Select(fields ...string) *sel {
	return &sel{fields: fields}
}

func (s *sel) From(model interface{}) *sel {
	t, err := GetTmpl(model)
	if err != nil {
		s.err = err
		return s
	}
	s.table = t.TableName
	return s
}

func (s *sel) Where(builder *conditions.ConditionBuilder) *sel {
	s.where = builder.String()
	return s
}

func (s *sel) Limit(limit int) *sel {
	s.limit = limit
	return s
}

func (s *sel) Offset(offset int) *sel {
	s.offset = offset
	return s
}

func (s *sel) GroupBy(fields ...string) *sel {
	s.groupBy = strings.Join(fields, ",")
	return s
}

func (s *sel) OrderBy(fields ...string) *sel {
	s.orderBy = strings.Join(fields, ",")
	return s
}

func (s *sel) InnerJoin(model interface{}) *join {
	j := &join{JoinType: inner, sel: s}
	return j.join(model)
}

func (s *sel) LeftJoin(model interface{}) *join {
	j := &join{JoinType: left, sel: s}
	return j.join(model)
}

func (s *sel) RightJoin(model interface{}) *join {
	j := &join{JoinType: right, sel: s}
	return j.join(model)
}

func (s *sel) FullJoin(model interface{}) *join {
	j := &join{JoinType: full, sel: s}
	return j.join(model)
}

type join struct {
	Table     string
	JoinType  string
	Condition string
	sel       *sel
}

func (j *join) join(model interface{}) *join {
	t, err := GetTmpl(model)
	if err != nil {
		j.sel.err = err
		return j
	}
	j.Table = t.TableName
	j.sel.join = append(j.sel.join, j)
	return j
}

func (j *join) On(builder *conditions.ConditionBuilder) *sel {
	j.Condition = builder.String()
	return j.sel
}

func Syn(syn string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}
	synAdd := fmt.Sprintf(",%s.", syn)
	return synAdd[1:] + strings.Join(stripFieldNames(fields), synAdd)
}

func (s *sel) String() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	selStr := &selString{
		Table:   s.table,
		Fields:  s.fields,
		Join:    s.join,
		Where:   s.where,
		GroupBy: s.groupBy,
		OrderBy: s.orderBy,
		Limit:   s.limit,
		Offset:  s.offset,
		err:     s.err,
	}
	var selectSql = new(bytes.Buffer)
	err := selectTmpl.Execute(selectSql, selStr)
	return selectSql.String(), err
}
