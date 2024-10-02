package sex

import (
	"fmt"
	"strings"
)

type Sex uint8

const (
	Unknown Sex = iota
	Male
	Female
)

const (
	maleStr   = "male"
	femaleStr = "female"
	unknown   = "unknown"
	empty     = ""
)

func (s Sex) String() string {
	switch s {
	case Male:
		return maleStr
	case Female:
		return femaleStr
	default:
		return unknown
	}
}

func FromString(sex string) (Sex, error) {
	switch strings.ToLower(sex) {
	case empty, unknown:
		return Unknown, nil
	case maleStr:
		return Male, nil
	case femaleStr:
		return Female, nil
	default:
		return Unknown, fmt.Errorf("unknown sex '%s': only %v are available", sex, []string{maleStr, femaleStr})
	}
}
