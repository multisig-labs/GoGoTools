// Copyright (C) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// 2025 Multisig Labs, Inc.
// Copied from https://github.com/ava-labs/avalanche-cli/blob/5c0c29f660ceb19b3f332b2f148f82f65e4fd542/pkg/contract/contract.go
// and modified to accept field names in returned tuples
// i.e. "getValidator(bytes32 validationID)->((uint8 status,bytes nodeID,uint64 startingWeight,uint64 messageNonce,uint64 weight,uint64 startedAt,uint64 endedAt))"
package validatormanagercmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/multisig-labs/gogotools/pkg/utils"
)

func removeSurroundingParenthesis(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) > 0 {
		if string(s[0]) != "(" || string(s[len(s)-1]) != ")" {
			return "", fmt.Errorf("expected esp %q to be surrounded by parenthesis", s)
		}
		s = s[1 : len(s)-1]
	}
	return s, nil
}

func removeSurroundingBrackets(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) > 0 {
		if string(s[0]) != "[" || string(s[len(s)-1]) != "]" {
			return "", fmt.Errorf("expected esp %q to be surrounded by parenthesis", s)
		}
		s = s[1 : len(s)-1]
	}
	return s, nil
}

func getWords(s string) []string {
	words := []string{}
	word := ""
	parenthesisCount := 0
	insideBrackets := false
	for _, rune := range s {
		c := string(rune)
		if parenthesisCount > 0 {
			word += c
			if c == "(" {
				parenthesisCount++
			}
			if c == ")" {
				parenthesisCount--
				if parenthesisCount == 0 {
					words = append(words, word)
					word = ""
				}
			}
			continue
		}
		if insideBrackets {
			word += c
			if c == "]" {
				words = append(words, word)
				word = ""
				insideBrackets = false
			}
			continue
		}
		// Changed to not split on spaces to keep type and name together
		if c == "," || c == "(" || c == "[" {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		}
		if c == "," {
			continue
		}
		if c == "(" {
			parenthesisCount++
		}
		if c == "[" {
			insideBrackets = true
		}
		word += c
	}
	if word != "" {
		words = append(words, strings.TrimSpace(word))
	}
	return words
}

func getMap(types []string, params interface{}) ([]map[string]interface{}, error) {
	r := []map[string]interface{}{}
	for i, t := range types {
		var (
			param      interface{}
			name       string
			structName string
		)
		parts := strings.Fields(t)
		typeStr := parts[0]
		if len(parts) == 2 { // i.e. "uint8 status"
			name = parts[1]
			t = typeStr // Update t to only contain the type
		}
		rt := reflect.ValueOf(params)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Slice {
			if rt.Len() != len(types) {
				if rt.Len() == 1 {
					return getMap(types, rt.Index(0).Interface())
				} else {
					return nil, fmt.Errorf(
						"inconsistency in slice len between method esp %q and given params %#v: expected %d got %d",
						types,
						params,
						len(types),
						rt.Len(),
					)
				}
			}
			param = rt.Index(i).Interface()
		} else if rt.Kind() == reflect.Struct {
			if rt.NumField() < len(types) {
				return nil, fmt.Errorf(
					"inconsistency in struct len between method esp %q and given params %#v: expected %d got %d",
					types,
					params,
					len(types),
					rt.NumField(),
				)
			}
			name = rt.Type().Field(i).Name
			structName = rt.Type().Field(i).Type.Name()
			param = rt.Field(i).Interface()
		}
		m := map[string]interface{}{}
		switch {
		case string(t[0]) == "(":
			// struct type
			var err error
			t, err = removeSurroundingParenthesis(t)
			if err != nil {
				return nil, err
			}
			m["components"], err = getMap(getWords(t), param)
			if err != nil {
				return nil, err
			}
			if structName != "" {
				m["internalType"] = "struct " + structName
			} else {
				m["internalType"] = "tuple"
			}
			m["type"] = "tuple"
			m["name"] = name
		case string(t[0]) == "[":
			var err error
			t, err = removeSurroundingBrackets(t)
			if err != nil {
				return nil, err
			}
			if string(t[0]) == "(" {
				t, err = removeSurroundingParenthesis(t)
				if err != nil {
					return nil, err
				}
				rt := reflect.ValueOf(param)
				if rt.Kind() != reflect.Slice {
					return nil, fmt.Errorf("expected param for field %d of esp %q to be an slice", i, types)
				}
				param = reflect.Zero(rt.Type().Elem()).Interface()
				structName = rt.Type().Elem().Name()
				m["components"], err = getMap(getWords(t), param)
				if err != nil {
					return nil, err
				}
				if structName != "" {
					m["internalType"] = "struct " + structName + "[]"
				} else {
					m["internalType"] = "tuple[]"
				}
				m["type"] = "tuple[]"
				m["name"] = name
			} else {
				m["internalType"] = fmt.Sprintf("%s[]", t)
				m["type"] = fmt.Sprintf("%s[]", t)
				m["name"] = name
			}
		default:
			m["internalType"] = t
			m["type"] = t
			m["name"] = name
		}
		r = append(r, m)
	}
	return r, nil
}

func ParseSpec(
	esp string,
	indexedFields []int,
	constructor bool,
	event bool,
	paid bool,
	view bool,
	params ...interface{},
) (string, string, error) {
	index := strings.Index(esp, "(")
	if index == -1 {
		return esp, "", nil
	}
	name := esp[:index]
	types := esp[index:]
	inputs := ""
	outputs := ""
	index = strings.Index(types, "->")
	if index == -1 {
		inputs = types
	} else {
		inputs = types[:index]
		outputs = types[index+2:]
	}
	var err error
	inputs, err = removeSurroundingParenthesis(inputs)
	if err != nil {
		return "", "", err
	}
	outputs, err = removeSurroundingParenthesis(outputs)
	if err != nil {
		return "", "", err
	}
	inputTypes := getWords(inputs)
	outputTypes := getWords(outputs)
	inputsMaps, err := getMap(inputTypes, params)
	if err != nil {
		return "", "", err
	}
	outputsMaps, err := getMap(outputTypes, nil)
	if err != nil {
		return "", "", err
	}
	if event {
		for i := range inputsMaps {
			if utils.Belongs(indexedFields, i) {
				inputsMaps[i]["indexed"] = true
			}
		}
	}
	abiMap := []map[string]interface{}{
		{
			"inputs": inputsMaps,
		},
	}
	switch {
	case paid:
		abiMap[0]["stateMutability"] = "payable"
	case view:
		abiMap[0]["stateMutability"] = "view"
	default:
		abiMap[0]["stateMutability"] = "nonpayable"
	}
	switch {
	case constructor:
		abiMap[0]["type"] = "constructor"
	case event:
		abiMap[0]["type"] = "event"
		abiMap[0]["name"] = name
		delete(abiMap[0], "stateMutability")
	default:
		abiMap[0]["type"] = "function"
		abiMap[0]["outputs"] = outputsMaps
		abiMap[0]["name"] = name
	}
	abiBytes, err := json.MarshalIndent(abiMap, "", "  ")
	if err != nil {
		return "", "", err
	}
	return name, string(abiBytes), nil
}
