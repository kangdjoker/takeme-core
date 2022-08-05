package utils

import "strings"

func IsContainSpecialCharacter(text string) bool {

	if strings.Contains(text, "!") {
		return true
	} else if strings.Contains(text, "@") {
		return true
	} else if strings.Contains(text, "#") {
		return true
	} else if strings.Contains(text, "$") {
		return true
	} else if strings.Contains(text, "%") {
		return true
	} else if strings.Contains(text, "^") {
		return true
	} else if strings.Contains(text, "&") {
		return true
	} else if strings.Contains(text, "*") {
		return true
	} else if strings.Contains(text, "(") {
		return true
	} else if strings.Contains(text, ")") {
		return true
	} else if strings.Contains(text, "-") {
		return true
	} else if strings.Contains(text, "_") {
		return true
	} else if strings.Contains(text, "+") {
		return true
	} else if strings.Contains(text, "=") {
		return true
	} else if strings.Contains(text, ".") {
		return true
	} else if strings.Contains(text, ",") {
		return true
	} else if strings.Contains(text, "?") {
		return true
	} else if strings.Contains(text, "/") {
		return true
	} else if strings.Contains(text, "<") {
		return true
	} else if strings.Contains(text, ">") {
		return true
	} else if strings.Contains(text, "1") {
		return true
	} else if strings.Contains(text, "2") {
		return true
	} else if strings.Contains(text, "3") {
		return true
	} else if strings.Contains(text, "4") {
		return true
	} else if strings.Contains(text, "5") {
		return true
	} else if strings.Contains(text, "6") {
		return true
	} else if strings.Contains(text, "7") {
		return true
	} else if strings.Contains(text, "8") {
		return true
	} else if strings.Contains(text, "9") {
		return true
	} else {
		return false
	}
}
