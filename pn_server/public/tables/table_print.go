package tables

import (
	"fmt"
	"reflect"

	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/table"
)

func ParseMultiple(v interface{}) []map[string]string {
	// 确认输入的类型是一个数组
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Slice {
		panic("ToMap() expects a slice")
	}

	// 创建 []map[string]string
	result := make([]map[string]string, value.Len())

	// 遍历数组，解析每个对象
	for i := 0; i < value.Len(); i++ {
		result[i] = ParseOne(value.Index(i).Interface())
	}

	return result
}

// 将单个对象解析为 map[string]string
func ParseOne(i interface{}) map[string]string {
	objectValue := reflect.ValueOf(i)

	// 创建 map[string]string
	ret := make(map[string]string)

	objectType := objectValue.Type()

	// 如果是指针，需要先解引用
	if objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
		objectValue = objectValue.Elem()
	}

	if objectType.Kind() == reflect.Map {
		// 如果是 map, 直接转化为 map[string]string
		m := objectValue.Interface().(map[string]interface{})
		for _, key := range objectValue.MapKeys() {
			ret[key.String()] = fmt.Sprintf("%v", m[key.String()])
		}
	} else if objectType.Kind() == reflect.Struct {
		// 如果是结构体, 遍历对象的字段，将每个字段解析为 string 并添加到 map 中
		for j := 0; j < objectType.NumField(); j++ {
			// 如果是匿名字段, 或者方法，跳过, 只处理公开字段
			if objectType.Field(j).Anonymous || objectType.Field(j).PkgPath != "" {
				continue
			}
			field := objectType.Field(j)
			fieldValue := objectValue.Field(j)
			ret[field.Name] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	} else if objectType.Kind() != reflect.Struct {
		// 如果不是结构体, 转化为字符串
		ret["Value"] = fmt.Sprintf("%v", objectValue.Interface())
	}

	return ret
}

type TableStruct struct {
	Value interface{}
}

func TableString(i interface{}) string {
	// 如果切片长度为 0, 直接返回提示
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Slice && value.Len() == 0 {
		return "No data"
	}

	// 如果不是切片类型, 转为原类型的切片
	if value.Kind() != reflect.Slice {
		oldType := reflect.TypeOf(i)
		newType := reflect.SliceOf(oldType)
		newValue := reflect.MakeSlice(newType, 1, 1)
		newValue.Index(0).Set(value)
		i = newValue.Interface()
	}

	// ToMap
	m := ParseMultiple(i)

	columns := make([]string, 0)
	if len(m) > 0 {
		for k := range m[0] {
			columns = append(columns, k)
		}
	}

	t, err := gotable.Create(columns...)
	if err != nil {
		panic(err)
	}

	t.AddRows(m)

	return t.String()
}
func Table(i interface{}) *table.Table {
	// 如果不是切片类型, 转为原类型的切片
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Slice {
		oldType := reflect.TypeOf(i)
		newType := reflect.SliceOf(oldType)
		newValue := reflect.MakeSlice(newType, 1, 1)
		newValue.Index(0).Set(value)
		i = newValue.Interface()
	}

	// ToMap
	m := ParseMultiple(i)

	columns := make([]string, 0)
	if len(m) > 0 {
		for k := range m[0] {
			columns = append(columns, k)
		}
	}

	t, err := gotable.Create(columns...)
	if err != nil {
		panic(err)
	}

	t.AddRows(m)

	return t
}
