package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/lib/pq"
	"github.com/lithammer/shortuuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateKey using in set keys
func GenCode() string {
	id := shortuuid.New()
	return strings.ToUpper(id[0:10])
}

// PatternGet using in get keys
func PatternGet(id uint) string {
	return strconv.Itoa(int(id)) + "-:--*"
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		var field interface{}
		if reflectValue.Field(i).Kind() == reflect.Ptr && !reflectValue.Field(i).IsNil() {
			field = reflectValue.Field(i).Elem()
		} else {
			field = reflectValue.Field(i).Interface()
		}
		if tag != "" && tag != "-" {
			tag = strcase.ToSnake(tag)
			res[tag] = field
		}

	}
	return res
}

func StructToMapType(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			tag = strcase.ToSnake(tag)

			res[tag] = reflect.TypeOf(field)
		}
	}
	return res
}

func InterfaceToFilterQuery(typeMap, valueMap map[string]interface{}, isExactSearch bool) string {
	query := "id > 0"
	for key, val := range typeMap {
		tmpVal := fmt.Sprintf("%v", val)
		if tmpVal == "int" || tmpVal == "float32" || tmpVal == "float64" || tmpVal == "uint" || tmpVal == "*int" || tmpVal == "*float32" || tmpVal == "*float64" || tmpVal == "*uint" {
			tmp := fmt.Sprintf("%v", valueMap[key])
			if tmp != "<nil>" && tmp != "99999999" {
				query += " AND \"" + key + "\" = " + tmp
			}

		}
		if tmpVal == "bool" || tmpVal == "*bool" {
			tmp := fmt.Sprintf("%v", valueMap[key])
			if tmp != "<nil>" {
				query += " AND \"" + key + "\" = " + tmp
			}
		}
		if tmpVal == "string" || tmpVal == "*string" {
			tmp := fmt.Sprintf("%v", valueMap[key])
			switch tmp {
			case "":
				break
			case "<nil>":
				break
			default:
				if isExactSearch {
					query += " AND \"" + key + "\" = '" + tmp + "'"
					break
				}
				tmp = strings.Replace(tmp, "đ", "d", -1)
				tmp = strings.Replace(tmp, "Đ", "d", -1)
				t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
				result, _, _ := transform.String(t, tmp)
				query += " AND unaccent(" + key + ") ILIKE '%" + result + "%'"
				break
			}
		}
		if tmpVal == "pq.StringArray" {
			var tmp []string
			for index := range valueMap[key].(pq.StringArray) {
				tmp = append(tmp, valueMap[key].(pq.StringArray)[index])
			}

			if len(tmp) > 0 {
				queryValue := strings.Join(tmp, ",")
				query += " AND \"" + key + "\" @> '{" + queryValue + "}'"
			}
		}

		if tmpVal == "pq.IntArray" {
			tmp := fmt.Sprintf("%v", valueMap[key])

			if tmp != "[]" {
				query += " AND \"" + key + "\" @> '" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(tmp, "[", "{"), "]", "}"), " ", ",") + "'"
			}
		}
		if strings.Contains(tmpVal, "time.Time") {
			tmp := fmt.Sprintf("%v", valueMap[key])
			if tmp != "<nil>" && StringTimeToString(tmp) != "0001-01-01 00:00:00" {
				query += " AND \"" + key + "\"::date = '" + StringTimeToString(tmp) + "'"
			}
		}
	}
	return query
}

// StringTimeToString converts a string in time.Time format to a string in format of 2006-01-02 15:04:05.
func StringTimeToString(s string) string {
	tmpTime, err := time.Parse("2006-01-02 15:05:05 +0000 UTC", s)
	if err != nil {
		fmt.Println(err)
	}
	return tmpTime.Format("2006-01-02 15:04:05")
}
