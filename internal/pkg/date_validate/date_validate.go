package date_validate

import (
	"regexp"
)

func IsValidDate(date string) bool {
	re := regexp.MustCompile("^(0[1-9]|[12][0-9]|3[01]).(0[1-9]|1[0-2]).(\\d{4})$")
	return re.MatchString(date)
}
