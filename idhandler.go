package toolkit

import (
	"errors"
	"reflect"
)

func IdInfo(i interface{}) (idfield string, id interface{}) {
	//_ = "breakpoint"
	idFields := []interface{}{"_id", "ID", "Id", "id"}
	rv := reflect.ValueOf(i)

	//-- get key
	//found := false
	if rv.Kind() == reflect.Map {
		mapkeys := rv.MapKeys()
		for _, mapkey := range mapkeys {
			idkey := mapkey.String()
			if HasMember(idFields, idkey) {
				idValue := rv.MapIndex(mapkey)
				if idValue.IsValid() {
					idfield = idkey
					id = idValue.Interface()
					return
				}
			}
		}
	} else if rv.Kind() == reflect.Struct {
		for _, idkey := range idFields {
			idValue := rv.FieldByName(idkey.(string))
			if idValue.IsValid() {
				idfield = idkey.(string)
				id = idValue.Interface()
				return
			}
		}
	} else if rv.Kind() == reflect.Ptr {
		elem := rv.Elem()
		for _, idkey := range idFields {
			idValue := elem.FieldByName(idkey.(string))
			if idValue.IsValid() {
				idfield = idkey.(string)
				id = idValue.Interface()
				return
			}
		}
	} else {
		//_ = "breakpoint"
		//fmt.Printf("Kind: %s \n", rv.Kind().String())
	}
	return
}

func Id(i interface{}) interface{} {
	f, i := IdInfo(i)
	if f == "" {
		return nil
	}
	return i
}

func IdField(i interface{}) string {
	f, _ := IdInfo(i)
	return f
}

func SetValue(rv *reflect.Value, value interface{}) error {
	v := reflect.ValueOf(value)
	rv.Set(v)
	return nil
}

func SetId(i interface{}, id interface{}) error {
	idfield := IdField(i)
	if idfield == "" {
		return errors.New("toolkit.SetId: No ID field")
	}
	rv := reflect.ValueOf(i)
	//-- get key
	//found := false
	if rv.Kind() == reflect.Map {
		mapkeys := rv.MapKeys()
		for _, mapkey := range mapkeys {
			idkey := mapkey.String()
			if idkey == idfield {
				mapvalue := rv.MapIndex(mapkey)
				return SetValue(&mapvalue, id)
			}
		}
	} else if rv.Kind() == reflect.Struct {
		idValue := rv.FieldByName(idfield)
		return SetValue(&idValue, id)
	} else if rv.Kind() == reflect.Ptr {
		elem := rv.Elem()
		idValue := elem.FieldByName(idfield)
		return SetValue(&idValue, id)
	}
	return errors.New("toolkit.SetID: Invalid type " + rv.Type().String())
}
