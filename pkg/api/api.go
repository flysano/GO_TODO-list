package api

import (
	"fmt"
	"regexp"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	dstartParse, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parsed date: %w", err)
	}

	//определяем правило
	var rule byte
	if len(repeat) > 0 {
		rule = repeat[0]
	} else {
		return "", fmt.Errorf("the rule is incorrect")
	}

	//определяем параметр для правила
	ruleArg := regexp.MustCompile(`-?\d+`).FindAllString(repeat, -1)
	if len(ruleArg) == 0 {
		return "", fmt.Errorf("the rule is incorrect")
	}

	//
	switch rule {
	case 'd':
		fmt.Print("dd")
	case 'y':
		fmt.Print("yy")

	}
}
