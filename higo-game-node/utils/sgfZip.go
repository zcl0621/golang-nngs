package utils

import (
	"strings"
)

func ZIPSgf(sgf string) string {
	if sgf == "()" {
		return ""
	}
	info := ParseSingle(sgf)
	return info.ZIPSgf
}

type Info struct {
	ZIPSgf string
}

func ParseSingle(sgfStr string) Info {
	var infoOutput Info
	dataString := sgfStr

	dataString = strings.ReplaceAll(dataString, "\n", "")
	dataString = strings.ReplaceAll(dataString, "(", "")
	dataString = strings.ReplaceAll(dataString, ")", "")
	dataString = strings.ReplaceAll(dataString, "tt", "")
loop:
	for i := 0; i < len(dataString); i++ {
		str := string(dataString[i])
		if str == ";" {
			for j := i; j >= 0; j-- {
				str1 := string(dataString[j])
				if str1 != "]" {
					dataString = dataString[:j] + ";" + dataString[j+1:]
				} else {
					continue loop
				}
			}
		}
	}
	dataString = strings.ReplaceAll(dataString, ";", "")
	dataStringSplit := strings.Split(dataString, "]")
	lastKey := ""
	for i, s := range dataStringSplit {
		if len(s) == 0 {
			continue
		}
		if string(s[0]) == "[" {
			dataStringSplit[i] = lastKey + s
		}
		lastKey = strings.Split(dataStringSplit[i], "[")[0]
	}
	prevKey := ""
	for _, keyValueString := range dataStringSplit {
		prevKey += infoOutput.parseKeyValue(keyValueString, prevKey)
	}
	infoOutput.ZIPSgf = prevKey
	return infoOutput
}

func (info *Info) parseKeyValue(input string, prevKey string) string {
	input = strings.TrimSpace(input)
	kvPair := strings.Split(input, "[")

	thisKeyWithMultipleValues := ""
	if len(kvPair) != 2 {
		return thisKeyWithMultipleValues
	}

	if kvPair[0] == "" {
		kvPair[0] = prevKey
	}

	switch kvPair[0] {
	case "B":
		if kvPair[1] == "" || kvPair[1] == "tt" { // if yield
			//pass
			thisKeyWithMultipleValues += "TT"
		} else {
			thisKeyWithMultipleValues += strings.ToUpper(kvPair[1])
		}
	case "W":
		if kvPair[1] == "" || kvPair[1] == "tt" { // if yield
			// pass
			thisKeyWithMultipleValues += "TT"
		} else {
			thisKeyWithMultipleValues += strings.ToUpper(kvPair[1])
		}

	case "AB":
		thisKeyWithMultipleValues += strings.ToUpper(kvPair[1])
		thisKeyWithMultipleValues += "TT"
	case "AW":
		thisKeyWithMultipleValues += strings.ToUpper(kvPair[1])
		thisKeyWithMultipleValues += "TT"
	default:
	}

	return strings.ToLower(thisKeyWithMultipleValues)
}
