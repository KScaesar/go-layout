package utility

import (
	"time"
)

type Hack string

func (hack Hack) IsOk(value string) bool {
	if hack != "" {
		return string(hack) == value
	}
	const YYYYMMDD = "20060102"
	return time.Now().Format(YYYYMMDD) == value
}
