package validator

import "testing"

func TestValidator(t *testing.T) {
	type testStruct struct {
		Value string `json:"value" validate:"email,required"`
	}

	data := testStruct{
		Value: "",
	}

	errs := ValidateStruct(data)
	if len(errs) == 0 {
		t.Fatal("ValidateStruct method couldn't validate.")
	}
	t.Log(errs)
}
