package conditions

import (
	"fmt"
	"strings"
)

type ConditionBuilder struct {
	strings.Builder
	err error
}

func String(s string) string {
	return fmt.Sprintf("'%s'", s)
}

func (c *ConditionBuilder) And() *ConditionBuilder {
	_, _ = c.WriteString(" AND")
	return c
}

func (c *ConditionBuilder) Or() *ConditionBuilder {
	_, _ = c.WriteString(" OR")
	return c
}

func (c *ConditionBuilder) IsNull(left string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "IS NULL"}, " "))
	return c
}

func (c *ConditionBuilder) IsNotNull(left string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "IS NOT NULL"}, " "))
	return c
}

func (c *ConditionBuilder) Eq(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "=", right}, " "))
	return c
}

func Eq(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, "=", right}, " "))
	return c
}

func (c *ConditionBuilder) Ne(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "!=", right}, " "))
	return c
}

func Ne(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, "!=", right}, " "))
	return c
}

func (c *ConditionBuilder) Gt(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, ">", right}, " "))
	return c
}

func Gt(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, ">", right}, " "))
	return c
}

func (c *ConditionBuilder) Ge(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, ">=", right}, " "))
	return c
}

func Ge(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, ">=", right}, " "))
	return c
}

func (c *ConditionBuilder) Lt(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "<", right}, " "))
	return c
}

func Lt(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, "<", right}, " "))
	return c
}

func (c *ConditionBuilder) Le(left, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, "<=", right}, " "))
	return c
}

func Le(left, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, "<=", right}, " "))
	return c
}

func (c *ConditionBuilder) Cond(left, operation, right string) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{"", left, operation, right}, " "))
	return c
}

func Cond(left, operation, right string) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{left, operation, right}, " "))
	return c
}

func (c *ConditionBuilder) Par(cond *ConditionBuilder) *ConditionBuilder {
	_, _ = c.WriteString(strings.Join([]string{" (", cond.String(), ")"}, ""))
	return c
}

func Par(cond *ConditionBuilder) *ConditionBuilder {
	c := &ConditionBuilder{}
	_, _ = c.WriteString(strings.Join([]string{"(", cond.String(), ")"}, ""))
	return c
}
