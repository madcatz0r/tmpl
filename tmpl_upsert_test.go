package templates

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestTmpl_UpsertQuery(t *testing.T) {
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
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC')
 ON CONFLICT (outer_id)
 DO UPDATE SET
 amount = $2,another_amount = $3,some_useful_thing = $4,updated_at = now() AT TIME ZONE 'UTC'
 RETURNING id`,
		},
		{
			name:       "exampleStruct ref",
			testStruct: &exampleStruct{},
			expected: `INSERT INTO example_struct
 (outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC')
 ON CONFLICT (outer_id)
 DO UPDATE SET
 amount = $2,another_amount = $3,some_useful_thing = $4,updated_at = now() AT TIME ZONE 'UTC'
 RETURNING id`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := GetTmpl(tt.testStruct)
			if err != nil {
				t.Fatal(err)
			}
			if got, _ := tmpl.UpsertQuery(); got != tt.expected {
				t.Errorf("UpsertQuery() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_Upsert(t *testing.T) {
	id, err := uuid.ParseBytes([]byte("342923c6-c073-42dd-b944-c7095b0576a0"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name           string
		testStruct     interface{}
		fields         []string
		expectedString string
		expectedValues []interface{}
	}{
		{
			name:       "exampleStruct full",
			testStruct: exampleStruct{ID: id, OuterID: 1, Amount: 2, AnotherAmount: 3},
			expectedString: `INSERT INTO example_struct
 (outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC')
 ON CONFLICT (outer_id)
 DO UPDATE SET
 amount = $2,another_amount = $3,some_useful_thing = $4,updated_at = now() AT TIME ZONE 'UTC'
 RETURNING id`,
			expectedValues: []interface{}{int64(1), int64(2), int64(3), ""},
		},
		{
			name:       "exampleStruct ref + partial",
			testStruct: &exampleStruct{ID: id, Amount: 3, AnotherAmount: 4},
			fields:     []string{"example_struct.amount", "example_struct.another_amount"},
			expectedString: `INSERT INTO example_struct
 (outer_id,amount,another_amount,some_useful_thing,created_at,updated_at)
 VALUES ($1,$2,$3,$4,now() AT TIME ZONE 'UTC',now() AT TIME ZONE 'UTC')
 ON CONFLICT (outer_id)
 DO UPDATE SET
 amount = $2,another_amount = $3,updated_at = now() AT TIME ZONE 'UTC'
 RETURNING id`,
			expectedValues: []interface{}{int64(0), int64(3), int64(4), ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var query string
			var values []interface{}
			var err error
			if tt.fields == nil {
				query, values, err = Upsert(tt.testStruct)
			} else {
				query, values, err = Upsert(tt.testStruct, tt.fields...)
			}

			if query != tt.expectedString {
				t1.Errorf("upsert got string = %v, expected %v", query, tt.expectedString)
			}
			if !reflect.DeepEqual(values, tt.expectedValues) {
				t1.Errorf("upsert got values = %v, expected %v", values, tt.expectedValues)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
