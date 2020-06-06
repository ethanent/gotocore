package gotocore

import (
	"reflect"
	"strconv"
	"strings"
)

var intKinds []reflect.Kind = []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}

func getGotocoreFields(d interface{}) []reflect.StructField {
	// Get all Gotocore fields

	t := reflect.TypeOf(d)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	gotocoreFields := []reflect.StructField{}

	for i := 0; i < t.NumField(); i++ {
		curField := t.Field(i)

		if curField.Tag.Get("g") == "" {
			continue
		}

		gotocoreFields = append(gotocoreFields, curField)
	}

	return gotocoreFields
}

// Marshal will marshal data into a Gotocore-supported []byte
func Marshal(d interface{}) []byte {
	// Get all Gotocore fields

	gotocoreFields := getGotocoreFields(d)

	// Introspect and start relaying values from d

	v := reflect.Indirect(reflect.ValueOf(d))

	built := []byte{}

	for _, gField := range gotocoreFields {
		f := v.FieldByName(gField.Name)

		valueKind := f.Kind()
		protoDesc := strings.Split(gField.Tag.Get("g"), ",")

		protoType := protoDesc[1]

		switch protoType {
		case "varint":
			built = append(built, buildVarint(int(f.Int()))...)
		case "string":
			built = append(built, buildString(f.String())...)
		case "buffer":
			built = append(built, buildBuffer(f.Bytes())...)
		case "uint":
			isInt := false

			for _, k := range intKinds {
				if valueKind == k {
					isInt = true
				}
			}

			size, err := strconv.Atoi(protoDesc[2])

			if err != nil {
				panic(err)
			}

			if isInt {
				built = append(built, buildUInt(uint(f.Int()), size)...)
			} else {
				built = append(built, buildUInt(uint(f.Uint()), size)...)
			}
		default:
			panic("unknown kind " + protoType)
		}
	}

	return built
}
