package go_tool

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func Md5(str string) string {
	var m = md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Uint8sToBytes(src []uint8) []byte {
	var dst []byte
	for _, b := range src {
		dst = append(dst, byte(b))
	}
	return dst
}

func AssetsMapOrS(kind reflect.Kind, msg string) {
	if kind != reflect.Map && kind != reflect.Struct {
		AssetsError(errors.New(msg))
	}
}

func AssetsSlice(kind reflect.Kind, msg string) {
	if kind != reflect.Slice {
		AssetsError(errors.New(msg))
	}
}

func AssetsMap(kind reflect.Kind, msg string) {
	if kind != reflect.Map {
		AssetsError(errors.New(msg))
	}
}

func StringJoin(ss ...string) string {
	buffer := new(strings.Builder)
	for _, s := range ss {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func ArrayToInterface(data interface{}) []interface{} {
	v := reflect.ValueOf(data)

	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make([]interface{}, vl)

	for i := 0; i < vl; i++ {
		dst[i] = v.Index(i).Interface()
	}

	return dst
}

func MapToInterface(data interface{}) map[interface{}]interface{} {
	v := reflect.ValueOf(data)

	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")

	keys := v.MapKeys()
	dst := make(map[interface{}]interface{}, len(keys))

	for _, kV := range keys {
		dst[kV.Interface()] = v.MapIndex(kV).Interface()
	}

	return dst
}

func MapMap(data interface{}, cb func(interface{}) interface{}) interface{} {
	dataTransform := MapToInterface(data)
	dst := make(map[interface{}]interface{}, len(dataTransform))
	for index, d := range dataTransform {
		dst[index] = cb(d)
	}
	return dst
}

func ArrayMap(data interface{}, cb func(interface{}) interface{}) interface{} {
	dataTransform := ArrayToInterface(data)
	dst := make([]interface{}, len(dataTransform))
	for index, d := range dataTransform {
		dst[index] = cb(d)
	}
	return dst
}

func ArrayFilter(data interface{}, cb func(interface{}) bool) interface{} {
	dataTransform := ArrayToInterface(data)
	dst := make([]interface{}, 0, len(dataTransform))
	for _, d := range dataTransform {
		if cb(d) {
			dst = append(dst, d)
		}
	}
	return dst
}

func ArrayReduce(data interface{}, cb func(float64, interface{}) float64, start float64) float64 {
	dataTransform := ArrayToInterface(data)
	for _, d := range dataTransform {
		start = cb(start, d)
	}
	return start
}

type ArrayKeyByS map[interface{}][]interface{}

func ArrayKeyBy(data interface{}, key string) ArrayKeyByS {
	dataTransform := ArrayToInterface(data)
	dst := make(ArrayKeyByS, len(dataTransform))
	for _, d := range dataTransform {
		k, _ := Get(d, key)
		dst[k] = append(dst[k], d)
	}
	return dst
}

func ArrayMakeKey(data interface{}, key string) map[interface{}]interface{} {
	dataTransform := ArrayToInterface(data)
	dst := make(map[interface{}]interface{}, len(dataTransform))
	for _, d := range dataTransform {
		k, _ := Get(d, key)
		dst[k] = d
	}
	return dst
}

func ArrayKeyByFunc(data interface{}, key string, cb func(interface{}, interface{}) interface{}) map[interface{}]interface{} {
	dataTransform := ArrayToInterface(data)
	dst := make(map[interface{}]interface{}, len(dataTransform))
	for _, d := range dataTransform {
		k, _ := Get(d, key)
		dst[k] = cb(dst[k], d)
	}
	return dst
}

func ArrayFirst(data interface{}) interface{} {
	v := reflect.ValueOf(data)

	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	if v.Len() == 0 {
		AssetsError(errors.New("数组为空, 无法识别子元素类型"))
	}

	return v.Index(0).Interface()
}

func GetArrayType(data interface{}) reflect.Type {
	return reflect.TypeOf(ArrayFirst(data))
}

func GetTypeFieldNum(t reflect.Type) int {
	kind := t.Kind()
	AssetsMapOrS(kind, "key提取接口错误, 非可提取接口")

	if kind == reflect.Map {
		return t.Len()
	} else {
		return t.NumField()
	}
}

func GetTypeFieldBySet(data interface{}, keys []string) []string {
	v := reflect.ValueOf(data)
	kind := v.Kind()
	AssetsMapOrS(kind, "无法提取字段类型")

	fieldNum := GetTypeFieldNum(v.Type())
	fields := make([]string, 0, fieldNum)
	if kind == reflect.Map {
		fields = ArrayMap(v.MapKeys(), func(d interface{}) interface{} {
			return d.(reflect.Value).Interface()
		}).([]string)
	} else {
		t := reflect.TypeOf(data)
		for i := 0; i < fieldNum; i++ {
			fields = append(fields, t.Field(i).Name)
		}
	}

	dstKeys := make([]string, 0, fieldNum)
	for _, k := range fields {
		if !InArray(keys, k) {
			dstKeys = append(dstKeys, k)
		}
	}
	return dstKeys
}

func InArray(data interface{}, find interface{}) bool {
	v := reflect.ValueOf(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	dNum := v.Len()
	exist := false
	for i := 0; i < dNum; i++ {
		exist = v.Index(i).Interface() == find
		if exist {
			break
		}
	}
	return exist
}

type ArrayPickS []map[string]interface{}

func ArrayPick(data interface{}, keys []string) ArrayPickS {
	dataTransform := ArrayToInterface(data)
	dst := make(ArrayPickS, 0, len(dataTransform))
	dstT := GetArrayType(dataTransform).Kind()
	fieldNum := len(keys)

	for _, d := range dataTransform {
		tmp := make(map[string]interface{}, fieldNum)
		v := reflect.ValueOf(d)
		for _, k := range keys {
			switch dstT {
			case reflect.Map:
				{
					kV := reflect.ValueOf(k)
					if v.MapIndex(kV).IsValid() {
						tmp[k] = v.MapIndex(kV).Interface()
					} else {
						tmp[k] = nil
					}
				}
			case reflect.Struct:
				{
					if v.FieldByName(k).IsValid() {
						tmp[k] = v.FieldByName(k).Interface()
					} else {
						tmp[k] = nil
					}
				}
			}
		}
		dst = append(dst, tmp)
	}
	return dst
}

func ArrayOmit(data interface{}, keys []string) ArrayPickS {
	return ArrayPick(data, GetTypeFieldBySet(ArrayFirst(data), keys))
}

func Set() bool {
	AssetsError(errors.New("不要尝试这样, 宝贝"))
	return true
}

func Has(data interface{}, path string) bool {
	_, exist := Get(data, path)
	return exist
}

func Get(data interface{}, path string) (interface{}, bool) {
	paths := strings.SplitN(path, ".", 2)
	v := reflect.ValueOf(data)
	kind := v.Kind()
	k := paths[0]
	shouldNext := len(paths) > 1

	if !InArray([]reflect.Kind{reflect.Slice, reflect.Map, reflect.Struct}, kind) {
		return nil, false
	} else {
		item := reflect.Value{}
		switch kind {
		case reflect.Slice:
			{
				len := v.Len()
				if index, err := strconv.Atoi(k); err != nil || index > len-1 || index < 0 {
					return nil, false
				} else {
					item = v.Index(index)
				}
			}
		case reflect.Map:
			{
				item = v.MapIndex(reflect.ValueOf(k))
			}
		case reflect.Struct:
			{
				item = v.FieldByName(k)
			}
		}

		if item.IsValid() {
			if shouldNext {
				return Get(item.Interface(), paths[1])
			}

			return item.Interface(), true
		} else {
			return nil, false
		}
	}
}

func AssetsReturn(assets bool, s interface{}, t interface{}) interface{} {
	if assets {
		return s
	} else {
		return t
	}
}

func AssetsError(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		buffer := new(strings.Builder)
		if ok {
			buffer.WriteString("file: ")
			buffer.WriteString(file)
			buffer.WriteString(" line: ")
			buffer.WriteString(strconv.Itoa(line))
			buffer.WriteString(" err: ")
		}
		buffer.WriteString(err.Error())
		panic(errors.New(buffer.String()))
	}
}
