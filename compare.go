package compare

import (
	"reflect"
	"fmt"
	"github.com/fatih/structs"
	"errors"
	"strings"
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

func Struct(a, b interface{}, customMapping []CustomMapping, ignoreKeys []string) error {
	return compare(a, b, customMapping, ignoreKeys)
}

func compare(aInt, bInt interface{}, customMapping []CustomMapping, ignoreKeys []string) error {
	a := extractNameValues(aInt)
	b := extractNameValues(bInt)
	for _, i := range ignoreKeys {
		//TODO move down a little
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
		var aFound bool
		var bFound bool

		for aN := range a {
			if aN == mapping.A {
				aFound = true
			}
		}
		for bN := range b {
			if bN == mapping.A {
				aFound = true
			}
		}

		if !aFound || !bFound {
			matchErrors = append(matchErrors, fmt.Sprintf("No mapping found for %s : %s", mapping.A, mapping.B))
		}
	}

	for _, mapping := range customMapping {
		if mapping.Value == "" {
			if isEqual(a[mapping.A], b[mapping.B]) {
				delete(a, mapping.A)
				delete(b, mapping.B)
			}
		}
	}

	for _, mapping := range customMapping {
		if isEqual(mapping.Value, b[mapping.B]) {
			delete(a, mapping.A)
			delete(b, mapping.B)
		}
	}

	for n, v := range a {
		matchErrors = append(matchErrors, fmt.Sprintf("No Unique Match found for %s : %s", n, v))
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

func hasUniqueValueMatch(m map[string]interface{}, value interface{}) (bool, string) {
	var hasSingleMatch bool
	var matchKey string

	for k, v := range m {
		if isEqual(value, v) {
			if hasSingleMatch {
				// not unique
				return false, ""
			}
			hasSingleMatch = true
			matchKey = k

		}
	}

	return hasSingleMatch, matchKey
}

func isEqual(a, b interface{}) bool {

	if a == nil || b == nil {
		return a == nil && b == nil
	}

	var aVal = reflect.ValueOf(a)
	var bVal = reflect.ValueOf(b)

	if aVal.Kind() == reflect.Ptr {
		aVal = aVal.Elem()
	}
	if bVal.Kind() == reflect.Ptr {
		bVal = bVal.Elem()
	}

	return aVal.Interface() == bVal.Interface()
}
func extractNameValues(a interface{}) map[string]interface{} {
	var flattened = make(map[string]interface{})
	for n, v := range structs.Map(a) {
		flatten(flattened, n, v)
	}
	return flattened
}

func flatten(m map[string]interface{}, name string, value interface{}) {
	switch reflect.TypeOf(value).Kind() {
	default:
		m[name] = value
	case reflect.Map:
		for n, v := range value.(map[string]interface{}) {
			flatten(m, fmt.Sprint(name, ".", n), v)
		}
	}
}
