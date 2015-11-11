package util_test

import (
	"fmt"
	"testing"

	"github.com/ello/ello-go/common/util"
)

func ExampleValidateInt() {
	value, err := util.ValidateInt("15", 10) // validates, returns parsed val

	fmt.Printf("%v | %v | ", value, err)

	value, err = util.ValidateInt("1a5", 10) // fails to validate, returns default

	fmt.Printf("%v | %v", value, err)
	// Output: 15 | <nil> | 10 | strconv.ParseInt: parsing "1a5": invalid syntax
}

func TestValidateInt(t *testing.T) {
	testVal, err := util.ValidateInt("15", 10)

	if err != nil {
		t.Error("'15' should pass validation")
	}
	if 15 != testVal {
		t.Error("validation should return 15")
	}

	testVal, err = util.ValidateInt("1a5", 10)

	if err == nil {
		t.Error("'1a5' should fail validation")
	}
	if 10 != testVal {
		t.Error("validation should return 10 (default value)")
	}

	testVal, err = util.ValidateInt("", 10)

	if err != nil {
		t.Error("'' should pass validation (empty string)")
	}
	if 10 != testVal {
		t.Error("validation should return 10 (default value)")
	}
}
