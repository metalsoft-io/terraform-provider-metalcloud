package metalcloud

import (
	"fmt"
	"regexp"
)

//ID interface is used where the type of the ID can be either an int or a string
type ID interface{}

func checkLabel(s string) error {
	r := regexp.MustCompilePOSIX("^[a-zA-Z]{1,1}[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1,1}|[a-zA-Z]{1,1}$")

	if s != r.FindString(s) {
		return fmt.Errorf("ID must be a label format, which is leters and numbers, no underscore, only dashes. It was %s", s)
	}
	return nil
}

func checkID(i ID) error {

	switch v := i.(type) {
	case int:
		if v < 0 {
			return fmt.Errorf("id cannot be less than 0. It was %d", v)
		}
		if v == 0 {
			return fmt.Errorf("id cannot be 0")
		}
	case string:
		return checkLabel(v)
	default:
		return fmt.Errorf("ID must be an int or a string that matches the label format. It was %+v", v)
	}

	return nil
}
