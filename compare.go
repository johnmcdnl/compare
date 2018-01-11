package compare

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/bxcodec/faker"
)

type nameMapping struct {
	A string
	B string
}

func newNameMapping(a, b string) nameMapping {
	return nameMapping{A: a, B: b}
}

type CustomMapping struct {
	nameMapping
	Value string //should be interface{}
}

func NewCustomMapping(a, b string, value string) CustomMapping {
	return CustomMapping{nameMapping: newNameMapping(a, b), Value: value}
}

func Fill(i interface{}) {
	faker.FakeData(&i)
}

func Struct(a, b interface{}, customMapping []CustomMapping, ignoreKeys []string) error {
	return compare(a, b, customMapping, ignoreKeys)
}

func compare(aInt, bInt interface{}, customMapping []CustomMapping, ignoreKeys []string) error {
	a := Flatten(structs.Map(aInt))
	b := Flatten(structs.Map(bInt))
	for _, i := range ignoreKeys {
		delete(a, i)
		delete(b, i)
	}
	var matchErrors []string

	for aName, aValue := range a {
		hasMatch, matchKey := hasUniqueValueMatch(b, aValue)
		if hasMatch {
			delete(a, aName)
			delete(b, matchKey)
		}
	}

	for _, mapping := range customMapping {
		if !mappingIsValid(a, mapping) {
			matchErrors = append(matchErrors, fmt.Sprintf("No mapping found for %s : %s", mapping.A, mapping.B))
		}
		if !mappingIsValid(b, mapping) {
			matchErrors = append(matchErrors, fmt.Sprintf("No mapping found for %s : %s", mapping.A, mapping.B))
		}
	}

	for _, mapping := range customMapping {
		if mapping.Value == "" {
			if a[mapping.A] == b[mapping.B] {
				delete(a, mapping.A)
				delete(b, mapping.B)
			}
		}
	}

	for _, mapping := range customMapping {
		if mapping.Value == b[mapping.B] {
			delete(a, mapping.A)
			delete(b, mapping.B)
		}
	}

	for k, v := range a {
		matchErrors = append(matchErrors, fmt.Sprintf("No Unique Match found for %s : %s", k, v))
	}

	if len(matchErrors) > 0 {
		return errors.New(fmt.Sprint(
			"\n",
			strings.Join(matchErrors, "\n"),
			"\n\n\n",
			fmt.Sprintln("Unmatched Data"),
			fmt.Sprintln(a),
			fmt.Sprintln(b),
		))
	}

	return nil
}

func mappingIsValid(m map[string]string, mapping CustomMapping) bool {
	for k := range m {
		if k == mapping.B {
			return true
		}
	}
	return false
}

func hasUniqueValueMatch(m map[string]string, value string) (bool, string) {
	var hasSingleMatch bool
	var matchKey string

	for k, v := range m {
		if v == value {
			if hasSingleMatch {
				return false, ""
			}
			hasSingleMatch = true
			matchKey = k
		}
	}

	return hasSingleMatch, matchKey
}
