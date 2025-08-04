package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NestedValidator 嵌套验证器
type NestedValidator struct {
	*ValidatorImpl
}

// NewNestedValidator 创建嵌套验证器
func NewNestedValidator() *NestedValidator {
	return &NestedValidator{
		ValidatorImpl: New().(*ValidatorImpl),
	}
}

// ValidateNested 验证嵌套结构体
func (nv *NestedValidator) ValidateNested(s interface{}) error {
	return nv.validateWithPath(s, "")
}

// validateWithPath 带路径的验证
func (nv *NestedValidator) validateWithPath(s interface{}, path string) error {
	if err := nv.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nv.formatNestedValidationErrors(validationErrors, path)
		}
		return err
	}
	return nil
}

// formatNestedValidationErrors 格式化嵌套验证错误
func (nv *NestedValidator) formatNestedValidationErrors(validationErrors validator.ValidationErrors, basePath string) error {
	var errors ValidationErrors
	
	for _, err := range validationErrors {
		field := nv.buildFieldPath(err, basePath)
		tag := err.Tag()
		value := err.Value()
		param := err.Param()
		
		// 获取自定义错误消息
		message := nv.getErrorMessage(tag, field, param, value)
		
		validationError := ValidationError{
			Field:     field,
			Tag:       tag,
			Value:     nv.formatValue(value),
			Param:     param,
			Message:   message,
			Kind:      err.Kind().String(),
			Type:      err.Type().String(),
			Namespace: nv.buildNamespace(err, basePath),
		}
		
		errors = append(errors, validationError)
	}
	
	return errors
}

// buildFieldPath 构建字段路径
func (nv *NestedValidator) buildFieldPath(err validator.FieldError, basePath string) string {
	field := nv.getFieldPath(err)
	if basePath != "" {
		return basePath + "." + field
	}
	return field
}

// buildNamespace 构建命名空间
func (nv *NestedValidator) buildNamespace(err validator.FieldError, basePath string) string {
	namespace := err.Namespace()
	if basePath != "" {
		return basePath + "." + namespace
	}
	return namespace
}

// ValidatePartial 部分验证（只验证指定字段）
func (nv *NestedValidator) ValidatePartial(s interface{}, fields ...string) error {
	if err := nv.validate.StructPartial(s, fields...); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nv.formatNestedValidationErrors(validationErrors, "")
		}
		return err
	}
	return nil
}

// ValidateExcept 排除验证（排除指定字段）
func (nv *NestedValidator) ValidateExcept(s interface{}, fields ...string) error {
	if err := nv.validate.StructExcept(s, fields...); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nv.formatNestedValidationErrors(validationErrors, "")
		}
		return err
	}
	return nil
}

// ValidateSlice 验证切片中的每个元素
func (nv *NestedValidator) ValidateSlice(slice interface{}) error {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return &ValidationError{
			Field:   "slice",
			Tag:     "slice",
			Message: "input must be a slice or array",
		}
	}
	
	var allErrors ValidationErrors
	
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if err := nv.ValidateNested(item); err != nil {
			if ve, ok := AsValidationErrors(err); ok {
				// 为每个错误添加索引前缀
				for _, e := range ve {
					indexedError := e
					indexedError.Field = "[" + string(rune(i+'0')) + "]." + e.Field
					indexedError.Namespace = "[" + string(rune(i+'0')) + "]." + e.Namespace
					allErrors = append(allErrors, indexedError)
				}
			}
		}
	}
	
	if len(allErrors) > 0 {
		return allErrors
	}
	
	return nil
}

// ValidateMap 验证map中的每个值
func (nv *NestedValidator) ValidateMap(m interface{}) error {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return &ValidationError{
			Field:   "map",
			Tag:     "map",
			Message: "input must be a map",
		}
	}
	
	var allErrors ValidationErrors
	
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key).Interface()
		keyStr := nv.formatValue(key.Interface())
		
		if err := nv.ValidateNested(value); err != nil {
			if ve, ok := AsValidationErrors(err); ok {
				// 为每个错误添加键前缀
				for _, e := range ve {
					keyedError := e
					keyedError.Field = "[" + keyStr + "]." + e.Field
					keyedError.Namespace = "[" + keyStr + "]." + e.Namespace
					allErrors = append(allErrors, keyedError)
				}
			}
		}
	}
	
	if len(allErrors) > 0 {
		return allErrors
	}
	
	return nil
}

// ValidateConditional 条件验证
func (nv *NestedValidator) ValidateConditional(s interface{}, condition func(interface{}) bool) error {
	if !condition(s) {
		return nil
	}
	return nv.ValidateNested(s)
}

// ValidateWithContext 带上下文的验证
func (nv *NestedValidator) ValidateWithContext(s interface{}, context map[string]interface{}) error {
	// 可以根据上下文信息调整验证逻辑
	// 这里先实现基本功能
	return nv.ValidateNested(s)
}

// GetValidationDepth 获取验证深度
func (nv *NestedValidator) GetValidationDepth(s interface{}) int {
	return nv.getStructDepth(reflect.TypeOf(s), 0)
}

// getStructDepth 递归获取结构体深度
func (nv *NestedValidator) getStructDepth(t reflect.Type, currentDepth int) int {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	if t.Kind() != reflect.Struct {
		return currentDepth
	}
	
	maxDepth := currentDepth
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		
		if fieldType.Kind() == reflect.Struct {
			depth := nv.getStructDepth(fieldType, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		} else if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if elemType.Kind() == reflect.Struct {
				depth := nv.getStructDepth(elemType, currentDepth+1)
				if depth > maxDepth {
					maxDepth = depth
				}
			}
		}
	}
	
	return maxDepth
}

// ValidateFieldPath 验证指定路径的字段
func (nv *NestedValidator) ValidateFieldPath(s interface{}, fieldPath string) error {
	parts := strings.Split(fieldPath, ".")
	current := reflect.ValueOf(s)
	
	// 遍历路径找到目标字段
	for i, part := range parts {
		if current.Kind() == reflect.Ptr {
			if current.IsNil() {
				return &ValidationError{
					Field:   strings.Join(parts[:i+1], "."),
					Tag:     "required",
					Message: "field is nil",
				}
			}
			current = current.Elem()
		}
		
		if current.Kind() != reflect.Struct {
			return &ValidationError{
				Field:   strings.Join(parts[:i], "."),
				Tag:     "struct",
				Message: "field is not a struct",
			}
		}
		
		field := current.FieldByName(part)
		if !field.IsValid() {
			return &ValidationError{
				Field:   strings.Join(parts[:i+1], "."),
				Tag:     "exists",
				Message: "field does not exist",
			}
		}
		
		current = field
	}
	
	// 验证目标字段
	if current.Kind() == reflect.Struct {
		return nv.validateWithPath(current.Interface(), strings.Join(parts, "."))
	}
	
	// 对于非结构体字段，创建临时结构体进行验证
	tempStruct := reflect.New(reflect.StructOf([]reflect.StructField{
		{
			Name: "Field",
			Type: current.Type(),
			Tag:  nv.getFieldValidationTag(s, fieldPath),
		},
	})).Elem()
	
	tempStruct.Field(0).Set(current)
	
	return nv.validateWithPath(tempStruct.Interface(), fieldPath)
}

// getFieldValidationTag 获取字段的验证标签
func (nv *NestedValidator) getFieldValidationTag(s interface{}, fieldPath string) reflect.StructTag {
	parts := strings.Split(fieldPath, ".")
	current := reflect.TypeOf(s)
	
	for _, part := range parts {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}
		
		if current.Kind() != reflect.Struct {
			break
		}
		
		field, found := current.FieldByName(part)
		if !found {
			break
		}
		
		if len(parts) == 1 {
			return field.Tag
		}
		
		current = field.Type
	}
	
	return ""
}