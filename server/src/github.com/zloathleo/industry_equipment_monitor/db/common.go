package db

import (
	"encoding/json"
	"reflect"
)

var floatType = reflect.TypeOf(float64(0))
var stringType = reflect.TypeOf("")

type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	} else {
		return []byte("null"), nil
	}
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	switch t := value.(type) {
	case []uint8:
		{
			vv := []uint8(t)
			if len(vv) == 0 {
				ns.String, ns.Valid = "", false
				return nil
			} else {
				//[]uint8 to string
				b := make([]byte, len(vv))
				for i, v := range vv {
					b[i] = byte(v)
				}
				ns.String, ns.Valid = string(b), true
				return nil
			}
			break
		}
	default:
		{
			ns.String, ns.Valid = "", false
			return nil
		}
	}
	return nil

}

//sql float64 的null处理
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

func (ni *NullFloat64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Float64)
	} else {
		return []byte("null"), nil
	}
}

func (ni *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		ni.Float64, ni.Valid = 0, false
		return nil
	}
	if reflect.TypeOf(value).String() == "[]uint8" {
		ni.Float64, ni.Valid = 0, false
		return nil
	}

	v := reflect.ValueOf(value)
	fv := v.Convert(floatType)
	ni.Float64, ni.Valid = fv.Float(), true
	return nil
}
