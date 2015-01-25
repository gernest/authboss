package authboss

import (
	"errors"
	"testing"
)

func TestErrorList_Error(t *testing.T) {
	t.Parallel()

	errList := ErrorList{errors.New("one"), errors.New("two")}
	if e := errList.Error(); e != "one, two" {
		t.Error("Wrong value for error:", e)
	}
}

func TestErrorList_Map(t *testing.T) {
	t.Parallel()

	errNotLong := "not long enough"
	errEmail := "should be an email"
	errAsploded := "asploded"

	errList := ErrorList{
		FieldError{"username", errors.New(errNotLong)},
		FieldError{"username", errors.New(errEmail)},
		FieldError{"password", errors.New(errNotLong)},
		errors.New(errAsploded),
	}

	m := errList.Map()
	if len(m) != 3 {
		t.Error("Wrong number of fields:", len(m))
	}

	usernameErrs := m["username"]
	if len(usernameErrs) != 2 {
		t.Error("Wrong number of username errors:", len(usernameErrs))
	}
	if usernameErrs[0] != errNotLong {
		t.Error("Wrong username error at 0:", usernameErrs[0])
	}
	if usernameErrs[1] != errEmail {
		t.Error("Wrong username error at 1:", usernameErrs[1])
	}

	passwordErrs := m["password"]
	if len(passwordErrs) != 1 {
		t.Error("Wrong number of password errors:", len(passwordErrs))
	}
	if passwordErrs[0] != errNotLong {
		t.Error("Wrong password error at 0:", passwordErrs[0])
	}

	unknownErrs := m[""]
	if len(unknownErrs) != 1 {
		t.Error("Wrong number of unkown errors:", len(unknownErrs))
	}
	if unknownErrs[0] != errAsploded {
		t.Error("Wrong unkown error at 0:", unknownErrs[0])
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	ctx := mockRequestContext("username", "john", "email", "john@john.com")

	errList := ctx.Validate([]Validator{
		mockValidator{
			FieldName: "username", Errs: ErrorList{FieldError{"username", errors.New("must be longer than 4")}},
		},
		mockValidator{
			FieldName: "missing_field", Errs: ErrorList{FieldError{"missing_field", errors.New("Expected field to exist.")}},
		},
		mockValidator{
			FieldName: "email", Errs: nil,
		},
	})

	errs := errList.Map()
	if errs["username"][0] != "must be longer than 4" {
		t.Error("Expected a different error for username:", errs["username"][0])
	}
	if errs["missing_field"][0] != "Expected field to exist." {
		t.Error("Expected a different error for missing_field:", errs["missing_field"][0])
	}
	if _, ok := errs["email"]; ok {
		t.Error("Expected no errors for email.")
	}
}