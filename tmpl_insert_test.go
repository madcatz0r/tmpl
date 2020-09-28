package templates

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

type exampleStruct struct {
	ID              uuid.UUID `tmpl:"type=primary"`
	OuterID         int64     `tmpl:"type=outer"`
	Amount          int64     `tmpl:"type=upsert"`
	AnotherAmount   int64     `tmpl:"type=upsert"`
	SomeUnusedThing int64     `tmpl:"-"`
	SomeUsefulThing string    `tmpl:"type=upsert"`
	CreatedAt       time.Time `tmpl:"type=insert,default=now() AT TIME ZONE 'UTC'"`
	UpdatedAt       time.Time `tmpl:"type=must_upd,default=now() AT TIME ZONE 'UTC'"`
}

func init() {
	_ = ParseTags(exampleStruct{})
}

func TestTmpl_InsertQuery(t *testing.T) {
	tests := []struct {
		name       string
		testStruct interface{}
		expected   string
	}{
		{
			name:       "exampleStruct",
			testStruct: exampleStruct{},
			expected: `INSERT INTO example_struct 
(outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC') returning id`,
		},
		{
			name:       "exampleStruct ref",
			testStruct: &exampleStruct{},
			expected: `INSERT INTO example_struct 
(outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC') returning id`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := GetTmpl(tt.testStruct)
			if err != nil {
				t.Fatal(err)
			}
			if got, _ := tmpl.InsertQuery(); got != tt.expected {
				t.Errorf("InsertQuery() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_Insert(t *testing.T) {
	var tests = []struct {
		name           string
		testStruct     interface{}
		expectedString string
		expected       []interface{}
	}{
		{
			name: "exampleStruct",
			testStruct: exampleStruct{
				ID:            uuid.UUID{},
				OuterID:       123,
				Amount:        1,
				AnotherAmount: 2,
				CreatedAt:     time.Time{},
				UpdatedAt:     time.Time{},
			},
			expectedString: `INSERT INTO example_struct 
(outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC') returning id`,
			expected: []interface{}{int64(123), int64(1), int64(2), ""},
		},
		{
			name: "exampleStruct ref",
			testStruct: &exampleStruct{
				ID:            uuid.UUID{},
				OuterID:       123,
				Amount:        2,
				AnotherAmount: 3,
				CreatedAt:     time.Time{},
				UpdatedAt:     SetTimeDiff(time.Hour),
			},
			expectedString: `INSERT INTO example_struct 
(outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC' + 3600 * interval '1 second') returning id`,
			expected: []interface{}{int64(123), int64(2), int64(3), ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, values, err := Insert(tt.testStruct)
			if query != tt.expectedString {
				t.Errorf("query = %s, expected %s", query, tt.expectedString)
			}
			if !reflect.DeepEqual(values, tt.expected) {
				t.Errorf("Insert() got values = %v, expected %v", values, tt.expected)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
