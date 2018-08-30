package util

import (
	"reflect"
	"log"
	"fmt"
	"regexp"
)

// 检查结构体的属性是否为空
// 当前只判断string类型的空值
// 用于检查参数上报
// 所有属性不为空返回nil，否则返回对应错误
func CheckAttrExist(structName interface{}, attrs []string) error {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	s := reflect.ValueOf(structName).Elem()
	fieldNum := t.NumField() // 获取结构体的变量个数
	//for i := 0; i < fieldNum; i++ {
	//	fmt.Println(t.Field(i).Name, s.Field(i).Type().String(), s.Field(i).Interface())
	//}
	for _, attr := range attrs {
		for i := 0; i < fieldNum; i++ {
			if t.Field(i).Name == attr {
				switch s.Field(i).Type().String() {
				case "string":
					if s.Field(i).Interface() == "" {
						return fmt.Errorf("%s is not valid", t.Field(i).Name)
					}
				default:
					continue
				}
			}
		}
	}
	return nil
}

func PrintStruct(structName interface{}) error {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	s := reflect.ValueOf(structName).Elem()
	fieldNum := t.NumField() // 获取结构体的变量个数
	for i := 0; i < fieldNum; i++ {
		fmt.Println(t.Field(i).Name, s.Field(i).Type().String(), s.Field(i).Interface())
	}
	return nil
}

// 验证邮箱格式
func CheckEmailFormat(email string) error {
	matched, err := regexp.MatchString("^.+\\@(\\[?)[a-zA-Z0-9\\-\\.]+\\.([a-zA-Z]{2,3}|[0-9]{1,3})(\\]?)$", email)
	if err != nil {
		return fmt.Errorf("check email err: %v", err)
	}
	if matched {
		return nil
	} else {
		return fmt.Errorf("email format is not right")
	}
}
