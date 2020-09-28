package templates

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestTmpl_UpdateQuery(t *testing.T) {
	tests := []struct {
		name       string
		testStruct interface{}
		expected   string
	}{
		{
			name:       "exampleStruct",
			testStruct: exampleStruct{},
			expected: `UPDATE example_struct SET
amount = $1,another_amount = $2,some_useful_thing = $3,updated_at = now() AT TIME ZONE 'UTC'
 WHERE id=$4 returning id`,
		},
		{
			name:       "exampleStruct ref",
			testStruct: &exampleStruct{},
			expected: `UPDATE example_struct SET
amount = $1,another_amount = $2,some_useful_thing = $3,updated_at = now() AT TIME ZONE 'UTC'
 WHERE id=$4 returning id`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := GetTmpl(tt.testStruct)
			if err != nil {
				t.Fatal(err)
			}
			if got, _ := tmpl.UpdateQuery(); got != tt.expected {
				t.Errorf("query = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestTmpl_Update(t *testing.T) {
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
			name:       "exampleStruct",
			testStruct: exampleStruct{ID: id, Amount: 2, AnotherAmount: 3},
			fields:     []string{"example_struct.amount", "example_struct.another_amount"},
			expectedString: `UPDATE example_struct SET
amount = $1,another_amount = $2,updated_at = now() AT TIME ZONE 'UTC'
 WHERE id=$3 returning id`,
			expectedValues: []interface{}{int64(2), int64(3), id},
		},
		{
			name:       "exampleStruct ref",
			testStruct: &exampleStruct{ID: id, Amount: 2, AnotherAmount: 3, UpdatedAt: SetTimeDiff(time.Hour)},
			fields:     []string{"example_struct.amount", "example_struct.another_amount", "example_struct.updated_at"},
			expectedString: `UPDATE example_struct SET
amount = $1,another_amount = $2,updated_at = now() AT TIME ZONE 'UTC' + 3600 * interval '1 second'
 WHERE id=$3 returning id`,
			expectedValues: []interface{}{int64(2), int64(3), id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			query, values, err := Update(tt.testStruct, tt.fields...)
			if query != tt.expectedString {
				t1.Errorf("Update() got string = %v, expected %v", query, tt.expectedString)
			}
			if !reflect.DeepEqual(values, tt.expectedValues) {
				t1.Errorf("Update() got values = %v, expected %v", values, tt.expectedValues)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTmpl_Update_All(t *testing.T) {
	id, err := uuid.ParseBytes([]byte("342923c6-c073-42dd-b944-c7095b0576a0"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name           string
		testStruct     interface{}
		expectedString string
		expectedValues []interface{}
	}{
		{
			name:       "exampleStruct",
			testStruct: exampleStruct{ID: id, Amount: 2, SomeUsefulThing: "check"},
			expectedString: `UPDATE example_struct SET
amount = $1,another_amount = $2,some_useful_thing = $3,updated_at = now() AT TIME ZONE 'UTC'
 WHERE id=$4 returning id`,
			expectedValues: []interface{}{int64(2), int64(0), "check", id},
		},
		{
			name:       "exampleStruct ref",
			testStruct: &exampleStruct{ID: id, Amount: 2, AnotherAmount: 3, SomeUsefulThing: "check"},
			expectedString: `UPDATE example_struct SET
amount = $1,another_amount = $2,some_useful_thing = $3,updated_at = now() AT TIME ZONE 'UTC'
 WHERE id=$4 returning id`,
			expectedValues: []interface{}{int64(2), int64(3), "check", id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			query, values, err := Update(tt.testStruct)
			if query != tt.expectedString {
				t1.Errorf("Update() got string = %v, expected %v", query, tt.expectedString)
			}
			if !reflect.DeepEqual(values, tt.expectedValues) {
				t1.Errorf("Update() got values = %v, expected %v", values, tt.expectedValues)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
