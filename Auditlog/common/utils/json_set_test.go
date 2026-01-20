package utils

import "testing"

func TestAddToJSON(t *testing.T) {
	jsonStr := `{"name":"John", "age":30, "cars":{"car1":"Ford","car2":"BMW"}}`
	jsonPath := "cars.car3"
	value := "Audi"

	expectedJSON := `{"name":"John", "age":30, "cars":{"car1":"Ford","car2":"BMW","car3":"Audi"}}`

	updatedJSON, err := AddToJSON(jsonStr, jsonPath, value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if updatedJSON != expectedJSON {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedJSON, updatedJSON)
	}
}

func TestAddKeyToJSONArray(t *testing.T) {
	jsonArrayStr := `[{"a":1},{"a":1},{"a":1}]`
	key := "b"
	value := 2

	expectedJSON := `[{"a":1,"b":2},{"a":1,"b":2},{"a":1,"b":2}]`

	updatedJSON, err := AddKeyToJSONArray(jsonArrayStr, key, value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if updatedJSON != expectedJSON {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedJSON, updatedJSON)
	}
}
