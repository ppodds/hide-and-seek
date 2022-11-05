package protobuf

import (
	"errors"
)

//type Marshaler interface {
//	MarshalProtoBuf() (any, error)
//}
//
//func Marshal(v any) (any, error) {
//	fmt.Println(v)
//	t := reflect.TypeOf(v)
//	fmt.Println(t.Kind())
//	switch t.Kind() {
//	case reflect.Pointer:
//		t2 := t.Elem()
//		//fmt.Println(t2.)
//		fmt.Println(t2.Kind())
//		fmt.Println(t2.String())
//		// get ptr in callee
//		data, err := Marshal(t2)
//		if err != nil {
//			return nil, err
//		}
//		return data, nil
//	case reflect.Slice:
//		marshaler, ok := v.([]Marshaler)
//		if ok {
//			l := make([]any, len(marshaler))
//			for _, item := range marshaler {
//				data, err := item.MarshalProtoBuf()
//				if err != nil {
//					return nil, err
//				}
//				l = append(l, data)
//			}
//			return l, nil
//		}
//		return nil, errors.New("type assertion failed")
//	case reflect.Map:
//		val := reflect.ValueOf(v)
//		m := make(map[any]any)
//		for _, key := range val.MapKeys() {
//			value := val.MapIndex(key)
//			marshaler, ok := reflect.ValueOf(value).Interface().(Marshaler)
//			if ok {
//				data, err := marshaler.MarshalProtoBuf()
//				if err != nil {
//					return nil, err
//				}
//				m[reflect.ValueOf(key)] = data
//			} else {
//				return nil, errors.New("type assertion failed")
//			}
//		}
//		return m, nil
//	case reflect.Struct:
//		marshaler, ok := v.(Marshaler)
//		if !ok {
//			return marshaler, errors.New("type assertion failed")
//		}
//		data, err := marshaler.MarshalProtoBuf()
//		if err != nil {
//			return nil, err
//		}
//		return data, nil
//	}
//	return nil, errors.New("unsupported format")
//}

func GetProtobuf(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return buf, nil
	}
	for i, data := range buf {
		if data == 0 {
			return buf[:i], nil
		}
	}
	return nil, errors.New("buffer is too short")
}
