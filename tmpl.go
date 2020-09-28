package templates

import (
	"fmt"
	"github.com/madcatz0r/tmpl/snake_case"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

var (
	tmplMap = make(map[string]*Tmpl)
	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	//	SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
	//[ * | expression [ [ AS ] output_name ] [, ...] ]
	//[ FROM from_item [, ...] ]
	//[ WHERE condition ]
	//[ GROUP BY grouping_element [, ...] ]
	//[ HAVING condition [, ...] ]
	//[ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
	//[ LIMIT { count | ALL } ]
	//[ OFFSET start [ ROW | ROWS ] ]
	selectTmpl = template.Must(template.New("selectTmpl").Funcs(funcMap).Parse(`SELECT {{range $index, $element := .Fields}}{{if $index}},{{end}}{{$element}}{{end}}{{if .Table }} FROM {{ .Table }}{{end}}{{ if .Join }}
 {{range $index, $element := .Join}}{{if $index}}
 {{end}}{{ $element.JoinType }} {{ $element.Table }} ON {{ $element.Condition }}{{end}}{{end}}{{if .Where }}
 WHERE {{ .Where }}{{end}}{{if .GroupBy}}
 GROUP BY {{ .GroupBy }}{{end}}{{if .OrderBy}}
 ORDER BY {{ .OrderBy }}{{end}}{{if .Limit }}
 LIMIT {{ .Limit }}{{end}}{{if .Offset}}
 OFFSET {{ .Offset }}{{end}}`))
	insertTmpl = template.Must(template.New("insertTmpl").Funcs(funcMap).Parse(`INSERT INTO {{.TableName}} 
({{ $counter := 0 }}{{range $index, $element := .InsertFields}}{{if $index}},{{end}}{{ $element }}{{end}})
 VALUES ({{$default_map := .Default}}{{range $index, $element := .InsertFields}}{{ $default := index $default_map $element}}{{if $index}},{{end}}{{if eq $default ""}}{{ $counter = inc $counter }}${{ $counter }}{{else}}{{ $default }}{{end}}{{end}}) returning {{ .Primary }}`))
	updateAllTmpl = template.Must(template.New("updateTmpl").Funcs(funcMap).Parse(`UPDATE {{.TableName}} SET
{{ $counter := 0 }}{{$default_map := .Default}}{{range $index, $element := .UpdateFields}}{{ $default := index $default_map $element}}{{if $index}},{{end}}{{ $element }} = {{if eq $default ""}}{{ $counter = inc $counter }}${{ $counter }}{{else}}{{ $default }}{{end}}{{end}}
 WHERE {{ .Primary }}=${{ inc $counter }} returning {{ .Primary }}`))
	upsertTmpl = template.Must(template.New("upsert").Parse(`INSERT INTO {{.TableName}}
 ({{$default_map := .Default}}{{range $index, $element := .InsertFields}}{{if $index}},{{end}}{{ $element }}{{end}})
 VALUES ({{$m := .UpsertFields}}{{range $index, $element := .InsertFields}}{{ $default := index $default_map $element}}{{if $index}},{{end}}{{if eq $default ""}}${{ index $m $element }}{{else}}{{ $default }}{{end}}{{end}})
 ON CONFLICT ({{.Outer}})
 DO UPDATE SET
 {{range $index, $element := .UpdateFields}}{{ $default := index $default_map $element}}{{if $index}},{{end}}{{ $element }} = {{if eq $default ""}}${{ index $m $element }}{{else}}{{ $default }}{{end}}{{end}}
 RETURNING {{ .Primary }}`))
)

var tagRegExp = regexp.MustCompile("^type=(?P<type>.[a-z,_]+)$|^type=(?P<type2>.[a-z,_]+),default=(?P<default>.*)$")

type Tmpl struct {
	TableName    string
	Primary      string
	Outer        string // insert only, use custom fields to update
	MustUpdate   string
	SelectFields []string
	InsertFields []string
	UpdateFields []string
	Mapping      map[string]string // k: db field v: struct field
	Default      map[string]string
	UpsertFields map[string]int // k: db field v: counter
}

func (t *Tmpl) getClone() *Tmpl {
	cloned := Tmpl{
		TableName:    t.TableName,
		Primary:      t.Primary,
		Outer:        t.Outer,
		MustUpdate:   t.MustUpdate,
		InsertFields: make([]string, len(t.InsertFields)),
		UpdateFields: make([]string, len(t.UpdateFields)),
		SelectFields: make([]string, len(t.SelectFields)),
		Mapping:      make(map[string]string),
		Default:      make(map[string]string),
		UpsertFields: make(map[string]int),
	}
	copy(cloned.InsertFields, t.InsertFields)
	copy(cloned.UpdateFields, t.UpdateFields)
	for key, value := range t.Mapping {
		cloned.Mapping[key] = value
	}
	for key, value := range t.Default {
		cloned.Default[key] = value
	}
	return &cloned
}

func GetTmpl(m interface{}) (*Tmpl, error) {
	name := structName(m)
	temp, ok := tmplMap[name]
	if !ok {
		return nil, fmt.Errorf("struct %s is not initialized", name)
	}
	return temp.getClone(), nil
}

type tag struct {
	Name    string
	ReqType string
	Default string
}

func ParseTags(m interface{}) error {
	str := reflect.TypeOf(m)
	tmpl := &Tmpl{
		TableName:    getTableName(str),
		Mapping:      make(map[string]string),
		Default:      make(map[string]string),
		UpsertFields: make(map[string]int),
	}
	for i := 0; i < str.NumField(); i++ {
		field := str.Field(i)
		tagString := field.Tag.Get("tmpl")
		tG := parseTag(tagString)
		if tG == nil {
			continue
		}
		tG.Name = snake_case.ToSnakeCase(field.Name)
		switch tG.ReqType {
		case "insert":
			tmpl.InsertFields = append(tmpl.InsertFields, tG.Name)
		case "update":
			tmpl.UpdateFields = append(tmpl.UpdateFields, tG.Name)
		case "upsert":
			tmpl.InsertFields = append(tmpl.InsertFields, tG.Name)
			tmpl.UpdateFields = append(tmpl.UpdateFields, tG.Name)
		case "primary":
			tmpl.Primary = tG.Name
		case "outer":
			tmpl.Outer = tG.Name
			tmpl.InsertFields = append(tmpl.InsertFields, tG.Name)
		case "must_upd":
			tmpl.MustUpdate = tG.Name
			tmpl.InsertFields = append(tmpl.InsertFields, tG.Name)
			tmpl.UpdateFields = append(tmpl.UpdateFields, tG.Name)
		default:
			err := fmt.Errorf("unknown tmpl %s field %s type tag: %s", tmpl.TableName, field.Name, tG.ReqType)
			return err
		}
		tmpl.Mapping[tG.Name] = field.Name
		tmpl.SelectFields = append(tmpl.SelectFields, fmt.Sprintf("%s.%s", tmpl.TableName, tG.Name))
		if tG.Default != "" {
			tmpl.Default[tG.Name] = tG.Default
		}
	}
	tmplMap[structName(m)] = tmpl
	return nil
}

func parseTag(tagString string) *tag {
	// []string{"", "type", "type2", "default"}
	groups := tagRegExp.FindStringSubmatch(tagString)
	if groups == nil {
		return nil
	}
	if groups[1] == "" {
		return &tag{ReqType: groups[2], Default: groups[3]}
	}
	return &tag{ReqType: groups[1]}
}

func getTableName(t reflect.Type) (name string) {
	for i := 0; i < t.NumField(); i++ {
		name = t.Field(i).Tag.Get("name")
		if name != "" {
			return name
		}
	}
	name = t.Name()
	if name == "" {
		name = t.Elem().Name()
	}
	return snake_case.ToSnakeCase(name)
}

func stripFieldNames(fields []string) []string {
	for i, field := range fields {
		dotPos := strings.IndexRune(field, '.')
		if dotPos > 0 {
			fields[i] = field[dotPos+1:]
		}
	}
	return fields
}
