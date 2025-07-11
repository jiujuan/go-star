package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func Init() {
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v = engine
		// 注册自定义 tag name 函数，让错误提示使用 json 字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func Struct(s interface{}) error {
	return v.Struct(s)
}

var Module = fx.Invoke(func() { Init() })