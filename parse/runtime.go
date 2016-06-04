package parse

import (
	"fmt"
	"reflect"
	"strings"
)

// Object extracts struct information from given object.
func Object(obj interface{}, schema, table string) (res *StructInfo, err error) {
	// convert any panic to error
	defer func() {
		p := recover()
		switch p := p.(type) {
		case error:
			err = p
		case nil:
			// nothing
		default:
			err = fmt.Errorf("%s", p)
		}
	}()

	t := reflect.ValueOf(obj).Elem().Type()
	res = &StructInfo{
		Type:         t.Name(),
		SQLSchema:    schema,
		SQLName:      table,
		PKFieldIndex: -1,
	}

	var n int
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("reform")
		if len(tag) == 0 {
			continue
		}

		// check for anonymous fields
		if f.Anonymous {
			return nil, fmt.Errorf(`reform: %s has anonymous field %s with "reform:" tag, it is not allowed`, res.Type, f.Name)
		}

		// check for exported name
		if f.PkgPath != "" {
			return nil, fmt.Errorf(`reform: %s has non-exported field %s with "reform:" tag, it is not allowed`, res.Type, f.Name)
		}

		// parse tag and type
		column, isPK := parseStructFieldTag(tag)
		if column == "" {
			return nil, fmt.Errorf(`reform: %s has field %s with invalid "reform:" tag value, it is not allowed`, res.Type, f.Name)
		}
		typ := f.Type.String()
		if isPK && strings.HasPrefix(typ, "*") {
			return nil, fmt.Errorf(`reform: %s has pointer field %s with with "pk" label in "reform:" tag, it is not allowed`, res.Type, f.Name)
		}
		if isPK && res.PKFieldIndex >= 0 {
			return nil, fmt.Errorf(`reform: %s has field %s with with duplicate "pk" label in "reform:" tag (first used by %s), it is not allowed`, res.Type, f.Name, res.Fields[res.PKFieldIndex].Name)
		}
		// if isPKOrOmitEmpty && strings.HasPrefix(typ, "*") {
		// 	return nil, fmt.Errorf(`reform: %s has pointer field %s with with "omitempty" label in "reform:" tag, it is not allowed`, res.Type, f.Name)
		// }

		// drop package name from qualified identifier if type is defined in this package
		if strings.Contains(typ, ".") && t.PkgPath() == f.Type.PkgPath() {
			typ = strings.Join(strings.Split(typ, ".")[1:], ".")
		}

		res.Fields = append(res.Fields, FieldInfo{
			Name:   f.Name,
			Type:   typ,
			Column: column,
			// PKOrOmitEmpty: isPKOrOmitEmpty,
		})
		if isPK {
			res.PKFieldIndex = n
		}
		n++
	}

	err = checkFields(res)
	if err != nil {
		return nil, err
	}

	return
}
