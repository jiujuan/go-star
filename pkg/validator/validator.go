package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 验证器接口
type Validator interface {
	// Validate 验证结构体
	Validate(s interface{}) error
	// RegisterValidation 注册自定义验证规则
	RegisterValidation(tag string, fn validator.Func) error
	// RegisterTranslation 注册错误信息翻译
	RegisterTranslation(tag string, message string) error
}

// ValidatorImpl 验证器实现
type ValidatorImpl struct {
	validate     *validator.Validate
	translations map[string]string
}

// New 创建新的验证器实例
func New() Validator {
	v := &ValidatorImpl{
		validate:     validator.New(),
		translations: make(map[string]string),
	}
	
	// 注册默认的字段名获取函数
	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
	
	// 注册默认的自定义验证规则
	v.registerDefaultValidations()
	
	return v
}

// Validate 验证结构体
func (v *ValidatorImpl) Validate(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// RegisterValidation 注册自定义验证规则
func (v *ValidatorImpl) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// RegisterTranslation 注册错误信息翻译
func (v *ValidatorImpl) RegisterTranslation(tag string, message string) error {
	v.translations[tag] = message
	return nil
}

// formatValidationErrors 格式化验证错误
func (v *ValidatorImpl) formatValidationErrors(validationErrors validator.ValidationErrors) error {
	var errors ValidationErrors
	
	for _, err := range validationErrors {
		field := v.getFieldPath(err)
		tag := err.Tag()
		value := err.Value()
		param := err.Param()
		
		// 获取自定义错误消息
		message := v.getErrorMessage(tag, field, param, value)
		
		validationError := ValidationError{
			Field:     field,
			Tag:       tag,
			Value:     v.formatValue(value),
			Param:     param,
			Message:   message,
			Kind:      err.Kind().String(),
			Type:      err.Type().String(),
			Namespace: err.Namespace(),
		}
		
		errors = append(errors, validationError)
	}
	
	return errors
}

// getFieldPath 获取字段路径（支持嵌套结构）
func (v *ValidatorImpl) getFieldPath(err validator.FieldError) string {
	namespace := err.Namespace()
	structNamespace := err.StructNamespace()
	
	// 优先使用结构体命名空间
	if structNamespace != "" {
		// 移除根结构体名称，只保留字段路径
		parts := strings.Split(structNamespace, ".")
		if len(parts) > 1 {
			// 构建相对路径，去掉结构体前缀，只保留字段名
			var fieldParts []string
			for i := 1; i < len(parts); i++ {
				// 获取字段的JSON标签名或原始名
				fieldParts = append(fieldParts, v.getFieldDisplayName(parts[i]))
			}
			return strings.Join(fieldParts, ".")
		}
	}
	
	// 如果有命名空间，使用它
	if namespace != "" {
		parts := strings.Split(namespace, ".")
		if len(parts) > 1 {
			var fieldParts []string
			for i := 1; i < len(parts); i++ {
				fieldParts = append(fieldParts, v.getFieldDisplayName(parts[i]))
			}
			return strings.Join(fieldParts, ".")
		}
	}
	
	// 否则使用字段名
	return v.getFieldDisplayName(err.Field())
}

// getFieldDisplayName 获取字段显示名称（优先使用JSON标签）
func (v *ValidatorImpl) getFieldDisplayName(fieldName string) string {
	// 这里简化处理，直接返回字段名
	// 在实际应用中，可以通过反射获取JSON标签
	return fieldName
}

// getErrorMessage 获取错误消息
func (v *ValidatorImpl) getErrorMessage(tag, field, param string, value interface{}) string {
	// 检查是否有自定义翻译
	if message, exists := v.translations[tag]; exists {
		return strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(message, "{field}", field),
				"{param}", param),
			"{value}", v.formatValue(value))
	}
	
	// 返回默认错误消息
	return v.getDefaultErrorMessage(tag, field, param, value)
}

// getDefaultErrorMessage 获取默认错误消息
func (v *ValidatorImpl) getDefaultErrorMessage(tag, field, param string, value interface{}) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " must be at most " + param + " characters long"
	case "len":
		return field + " must be exactly " + param + " characters long"
	case "oneof":
		return field + " must be one of: " + param
	case "numeric":
		return field + " must be numeric"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "gt":
		return field + " must be greater than " + param
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lt":
		return field + " must be less than " + param
	case "lte":
		return field + " must be less than or equal to " + param
	case "eq":
		return field + " must be equal to " + param
	case "ne":
		return field + " must not be equal to " + param
	default:
		return field + " validation failed for tag '" + tag + "'"
	}
}

// formatValue 格式化值用于显示
func (v *ValidatorImpl) formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}
	
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return strings.Trim(strings.Replace(reflect.ValueOf(v).String(), " ", "", -1), "<>")
	case uint, uint8, uint16, uint32, uint64:
		return strings.Trim(strings.Replace(reflect.ValueOf(v).String(), " ", "", -1), "<>")
	case float32, float64:
		return strings.Trim(strings.Replace(reflect.ValueOf(v).String(), " ", "", -1), "<>")
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return reflect.ValueOf(v).String()
	}
}

// registerDefaultValidations 注册默认的自定义验证规则
func (v *ValidatorImpl) registerDefaultValidations() {
	// 注册手机号验证
	v.validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		mobile := fl.Field().String()
		if len(mobile) != 11 {
			return false
		}
		// 简单的手机号格式验证
		return strings.HasPrefix(mobile, "1") && isNumeric(mobile)
	})
	
	// 注册身份证号验证
	v.validate.RegisterValidation("idcard", func(fl validator.FieldLevel) bool {
		idcard := fl.Field().String()
		// 简单的身份证号长度验证
		return len(idcard) == 15 || len(idcard) == 18
	})
	
	// 注册密码强度验证
	v.validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 6 {
			return false
		}
		// 至少包含一个字母和一个数字
		hasLetter := false
		hasNumber := false
		for _, char := range password {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				hasLetter = true
			}
			if char >= '0' && char <= '9' {
				hasNumber = true
			}
		}
		return hasLetter && hasNumber
	})
	
	// 注册默认翻译
	v.translations["mobile"] = "{field} must be a valid mobile number"
	v.translations["idcard"] = "{field} must be a valid ID card number"
	v.translations["password"] = "{field} must contain at least one letter and one number"
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}