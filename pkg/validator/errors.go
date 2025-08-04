package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationError 单个验证错误
type ValidationError struct {
	Field     string `json:"field"`               // 字段名（支持嵌套路径，如 profile.address.city）
	Tag       string `json:"tag"`                 // 验证标签
	Value     string `json:"value"`               // 字段值
	Param     string `json:"param"`               // 验证参数
	Message   string `json:"message"`             // 错误消息
	Kind      string `json:"kind,omitempty"`      // 字段类型种类
	Type      string `json:"type,omitempty"`      // 字段类型
	Namespace string `json:"namespace,omitempty"` // 完整命名空间
}

// Error 实现error接口
func (ve ValidationError) Error() string {
	return ve.Message
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现error接口
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	
	if len(ve) == 1 {
		return ve[0].Error()
	}
	
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	
	return strings.Join(messages, "; ")
}

// JSON 返回JSON格式的错误信息
func (ve ValidationErrors) JSON() string {
	data, _ := json.Marshal(ve)
	return string(data)
}

// Map 返回map格式的错误信息，key为字段名，value为错误消息
func (ve ValidationErrors) Map() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = err.Message
	}
	return result
}

// Fields 返回所有出错的字段名
func (ve ValidationErrors) Fields() []string {
	var fields []string
	for _, err := range ve {
		fields = append(fields, err.Field)
	}
	return fields
}

// HasField 检查是否包含指定字段的错误
func (ve ValidationErrors) HasField(field string) bool {
	for _, err := range ve {
		if err.Field == field {
			return true
		}
	}
	return false
}

// GetFieldError 获取指定字段的错误
func (ve ValidationErrors) GetFieldError(field string) *ValidationError {
	for _, err := range ve {
		if err.Field == field {
			return &err
		}
	}
	return nil
}

// String 返回字符串格式的错误信息
func (ve ValidationErrors) String() string {
	return ve.Error()
}

// Format 格式化输出错误信息
func (ve ValidationErrors) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			// 详细格式
			for i, err := range ve {
				if i > 0 {
					fmt.Fprint(f, "\n")
				}
				fmt.Fprintf(f, "Field: %s, Tag: %s, Value: %s, Message: %s", 
					err.Field, err.Tag, err.Value, err.Message)
			}
		} else {
			// 简单格式
			fmt.Fprint(f, ve.Error())
		}
	case 's':
		fmt.Fprint(f, ve.Error())
	case 'q':
		fmt.Fprintf(f, "%q", ve.Error())
	}
}

// IsValidationError 检查错误是否为验证错误
func IsValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	if ok {
		return true
	}
	_, ok = err.(ValidationError)
	return ok
}

// AsValidationErrors 将错误转换为ValidationErrors
func AsValidationErrors(err error) (ValidationErrors, bool) {
	if ve, ok := err.(ValidationErrors); ok {
		return ve, true
	}
	if ve, ok := err.(ValidationError); ok {
		return ValidationErrors{ve}, true
	}
	return nil, false
}

// GroupByField 按字段分组错误
func (ve ValidationErrors) GroupByField() map[string][]ValidationError {
	result := make(map[string][]ValidationError)
	for _, err := range ve {
		result[err.Field] = append(result[err.Field], err)
	}
	return result
}

// GroupByStruct 按结构体分组错误（用于嵌套结构）
func (ve ValidationErrors) GroupByStruct() map[string][]ValidationError {
	result := make(map[string][]ValidationError)
	for _, err := range ve {
		// 提取结构体路径
		parts := strings.Split(err.Field, ".")
		if len(parts) > 1 {
			structPath := strings.Join(parts[:len(parts)-1], ".")
			result[structPath] = append(result[structPath], err)
		} else {
			result["root"] = append(result["root"], err)
		}
	}
	return result
}

// GetNestedErrors 获取嵌套结构的错误
func (ve ValidationErrors) GetNestedErrors(prefix string) ValidationErrors {
	var nested ValidationErrors
	for _, err := range ve {
		if strings.HasPrefix(err.Field, prefix+".") {
			// 创建相对路径的错误副本
			nestedErr := err
			nestedErr.Field = strings.TrimPrefix(err.Field, prefix+".")
			nested = append(nested, nestedErr)
		}
	}
	return nested
}

// HasNestedErrors 检查是否有嵌套错误
func (ve ValidationErrors) HasNestedErrors() bool {
	for _, err := range ve {
		if strings.Contains(err.Field, ".") {
			return true
		}
	}
	return false
}

// GetRootErrors 获取根级别的错误（非嵌套）
func (ve ValidationErrors) GetRootErrors() ValidationErrors {
	var root ValidationErrors
	for _, err := range ve {
		if !strings.Contains(err.Field, ".") {
			root = append(root, err)
		}
	}
	return root
}

// ToNestedMap 转换为嵌套的map结构
func (ve ValidationErrors) ToNestedMap() map[string]interface{} {
	result := make(map[string]interface{})
	
	for _, err := range ve {
		parts := strings.Split(err.Field, ".")
		current := result
		
		// 构建嵌套结构
		for i, part := range parts {
			if i == len(parts)-1 {
				// 最后一个部分，设置错误消息
				current[part] = err.Message
			} else {
				// 中间部分，创建嵌套map
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				if nested, ok := current[part].(map[string]interface{}); ok {
					current = nested
				}
			}
		}
	}
	
	return result
}