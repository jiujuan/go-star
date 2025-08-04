package validator

import (
	"reflect"
	"strings"
)

// ValidateStruct 验证结构体的便捷函数
func ValidateStruct(s interface{}) error {
	validator := New()
	return validator.Validate(s)
}

// ValidateField 验证单个字段
func ValidateField(field interface{}, tag string) error {
	validator := New()
	
	// 创建一个临时结构体来验证单个字段
	tempStruct := struct {
		Field interface{} `validate:""`
	}{
		Field: field,
	}
	
	// 动态设置验证标签
	structType := reflect.TypeOf(tempStruct)
	structField := structType.Field(0)
	
	// 创建新的结构体类型，设置验证标签
	newStructType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Field",
			Type: structField.Type,
			Tag:  reflect.StructTag(`validate:"` + tag + `"`),
		},
	})
	
	// 创建新结构体实例
	newStruct := reflect.New(newStructType).Elem()
	newStruct.Field(0).Set(reflect.ValueOf(field))
	
	return validator.Validate(newStruct.Interface())
}

// GetStructTags 获取结构体的验证标签
func GetStructTags(s interface{}) map[string]string {
	tags := make(map[string]string)
	
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	if t.Kind() != reflect.Struct {
		return tags
	}
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("validate"); tag != "" {
			fieldName := getFieldName(field)
			tags[fieldName] = tag
		}
	}
	
	return tags
}

// getFieldName 获取字段名（优先使用json标签）
func getFieldName(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}
	}
	return field.Name
}

// IsRequired 检查字段是否为必填
func IsRequired(s interface{}, fieldName string) bool {
	tags := GetStructTags(s)
	if tag, exists := tags[fieldName]; exists {
		return strings.Contains(tag, "required")
	}
	return false
}

// GetValidationRules 获取字段的验证规则
func GetValidationRules(s interface{}, fieldName string) []string {
	tags := GetStructTags(s)
	if tag, exists := tags[fieldName]; exists {
		return strings.Split(tag, ",")
	}
	return nil
}

// HasValidationRule 检查字段是否有指定的验证规则
func HasValidationRule(s interface{}, fieldName, rule string) bool {
	rules := GetValidationRules(s, fieldName)
	for _, r := range rules {
		if strings.TrimSpace(r) == rule {
			return true
		}
	}
	return false
}

// ExtractValidationParam 提取验证规则的参数
func ExtractValidationParam(rule string) (string, string) {
	parts := strings.SplitN(rule, "=", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(rule), ""
}

// BuildValidationTag 构建验证标签
func BuildValidationTag(rules []string) string {
	return strings.Join(rules, ",")
}

// MergeValidationTags 合并验证标签
func MergeValidationTags(tags ...string) string {
	var allRules []string
	for _, tag := range tags {
		if tag != "" {
			rules := strings.Split(tag, ",")
			for _, rule := range rules {
				rule = strings.TrimSpace(rule)
				if rule != "" {
					allRules = append(allRules, rule)
				}
			}
		}
	}
	return BuildValidationTag(allRules)
}