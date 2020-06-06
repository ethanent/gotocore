package gotocore

import (
	"errors"
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
func Marshal(d interface{}) ([]byte, error) {
	// Get all Gotocore fields

	gotocoreFields := getGotocoreFields(d)

	// Introspect and start relaying values from d

	v := reflect.Indirect(reflect.ValueOf(d))

	built := []byte{}

	for _, gField := range gotocoreFields {
		f := v.FieldByName(gField.Name)

		valueKind := f.Kind()

		// Handle struct if nested struct / pointer

		if valueKind == reflect.Struct {
			marshalledNested, err := Marshal(f.Interface())

			if err != nil {
				return nil, err
			}

			built = append(built, marshalledNested...)
			continue
		}

		if valueKind == reflect.Ptr {
			elem := f.Elem()

			if elem.Kind() == reflect.Struct {
				marshalledNested, err := Marshal(f.Interface())

				if err != nil {
					return nil, err
				}

				built = append(built, marshalledNested...)
			} else {
				panic("protostruct does not support marshalling pointers to non-struct values")
			}
			continue
		}

		// Otherwise handle based on tag-specified kind

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

			var uintValue uint

			if isInt {
				uintValue = uint(f.Int())
			} else {
				uintValue = uint(f.Uint())
			}

			builtUint, err := buildUInt(uintValue, size)

			if err != nil {
				return nil, err
			}

			built = append(built, builtUint...)
		default:
			panic("unknown kind " + protoType)
		}
	}

	return built, nil
}

// Unmarshal parses d into dest struct
func Unmarshal(data []byte, d interface{}) (n int, err error) {
	// Get all Gotocore fields

	gotocoreFields := getGotocoreFields(d)

	// Extract data from gotocore fields

	v := reflect.Indirect(reflect.ValueOf(d))

	curLoc := 0

	for _, gField := range gotocoreFields {
		f := v.FieldByName(gField.Name)

		valueKind := f.Kind()

		// Handle struct or pointer to struct

		if gField.Type.Kind() == reflect.Struct {
			fv := reflect.New(gField.Type)

			end, err := Unmarshal(data[curLoc:], fv.Interface())

			if err != nil {
				return 0, err
			}

			curLoc += end

			if !f.CanSet() {
				return 0, errors.New("cannot set field " + gField.Name)
			}

			f.Set(fv.Elem())

			continue
		}

		if gField.Type.Kind() == reflect.Ptr {
			fv := reflect.New(gField.Type.Elem())

			end, err := Unmarshal(data[curLoc:], fv.Interface())

			if err != nil {
				return 0, err
			}

			curLoc += end

			if !f.CanSet() {
				return 0, errors.New("cannot set field " + gField.Name)
			}

			f.Set(reflect.ValueOf(fv.Interface()))

			continue
		}

		// Otherwise handle based on tag-specified kind

		protoDesc := strings.Split(gField.Tag.Get("g"), ",")

		protoType := protoDesc[1]

		switch protoType {
		case "varint":
			parsedFieldValue, readCount, err := parseVarint(data, curLoc, &Component{
				Name: gField.Name,
				Kind: Varint,
				Size: -1,
			})

			if err != nil {
				return 0, err
			}

			curLoc += readCount

			f.SetInt(int64(parsedFieldValue))
		case "string":
			parsedFieldValue, readCount, err := parseString(data, curLoc, &Component{
				Name: gField.Name,
				Kind: String,
				Size: -1,
			})

			if err != nil {
				return 0, err
			}

			curLoc += readCount

			f.SetString(parsedFieldValue)
		case "buffer":
			parsedFieldValue, readCount, err := parseBuffer(data, curLoc, &Component{
				Name: gField.Name,
				Kind: Buffer,
				Size: -1,
			})

			if err != nil {
				return 0, err
			}

			curLoc += readCount

			f.SetBytes(parsedFieldValue)
		case "uint":
			size, err := strconv.Atoi(protoDesc[2])

			if err != nil {
				panic(err)
			}

			parsedFieldValue, readCount, err := parseUInt(data, curLoc, &Component{
				Name: gField.Name,
				Kind: Uint,
				Size: size,
			})

			if err != nil {
				return 0, err
			}

			curLoc += readCount

			isInt := false

			for _, k := range intKinds {
				if valueKind == k {
					isInt = true
				}
			}

			if isInt {
				f.SetInt(int64(parsedFieldValue))
			} else {
				f.SetUint(uint64(parsedFieldValue))
			}
		default:
			panic("unknown kind " + protoType)
		}
	}

	return curLoc, nil
}
