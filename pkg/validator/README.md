# Validator Package

go-star框架的数据验证组件，基于`go-playground/validator/v10`构建，提供了强大的数据验证功能。

## 特性

- 🔍 **基础验证规则**: 支持所有标准验证规则（required, email, min, max等）
- 🏗️ **嵌套结构验证**: 完整支持嵌套结构体验证
- 🎯 **自定义验证规则**: 支持注册自定义验证规则
- 🌐 **错误消息翻译**: 支持自定义错误消息和多语言
- 📊 **详细错误信息**: 提供结构化的错误信息
- 🔧 **灵活配置**: 支持配置驱动的验证器创建
- 🚀 **高性能**: 基于成熟的validator库，性能优异

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "go-star/pkg/validator"
)

type User struct {
    Username string `validate:"required,min=3,max=20"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"required,gte=18,lte=120"`
}

func main() {
    v := validator.New()
    
    user := User{
        Username: "john_doe",
        Email:    "john@example.com",
        Age:      25,
    }
    
    if err := v.Validate(&user); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation passed")
    }
}
```

### 嵌套结构验证

```go
type Address struct {
    Street  string `validate:"required,min=5"`
    City    string `validate:"required,min=2"`
    ZipCode string `validate:"required,len=6,numeric"`
}

type User struct {
    Name    string  `validate:"required,min=2"`
    Address Address `validate:"required"`
}

func main() {
    v := validator.New()
    
    user := User{
        Name: "John Doe",
        Address: Address{
            Street:  "123 Main Street",
            City:    "New York",
            ZipCode: "123456",
        },
    }
    
    if err := v.Validate(&user); err != nil {
        if ve, ok := validator.AsValidationErrors(err); ok {
            for _, e := range ve {
                fmt.Printf("Field: %s, Error: %s\n", e.Field, e.Message)
            }
        }
    }
}
```

### 自定义验证规则

```go
func main() {
    v := validator.New()
    
    // 注册自定义验证规则
    v.RegisterValidation("weekday", func(fl validator.FieldLevel) bool {
        if date, ok := fl.Field().Interface().(time.Time); ok {
            weekday := date.Weekday()
            return weekday >= time.Monday && weekday <= time.Friday
        }
        return false
    })
    
    // 注册自定义错误消息
    v.RegisterTranslation("weekday", "{field} must be a weekday")
    
    type Meeting struct {
        Title string    `validate:"required"`
        Date  time.Time `validate:"required,weekday"`
    }
    
    // 验证
    meeting := Meeting{
        Title: "Team Meeting",
        Date:  time.Now(),
    }
    
    if err := v.Validate(&meeting); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    }
}
```

## 内置自定义验证规则

### mobile - 手机号验证
```go
type User struct {
    Mobile string `validate:"required,mobile"`
}
```

### password - 密码强度验证
```go
type User struct {
    Password string `validate:"required,password"`
}
```

### idcard - 身份证号验证
```go
type User struct {
    IDCard string `validate:"required,idcard"`
}
```

## 错误处理

### ValidationError 结构

```go
type ValidationError struct {
    Field     string // 字段名（支持嵌套路径）
    Tag       string // 验证标签
    Value     string // 字段值
    Param     string // 验证参数
    Message   string // 错误消息
    Kind      string // 字段类型种类
    Type      string // 字段类型
    Namespace string // 完整命名空间
}
```

### ValidationErrors 方法

```go
// 基本方法
errors.Error()           // 返回错误字符串
errors.JSON()            // 返回JSON格式
errors.Map()             // 返回map[string]string格式
errors.Fields()          // 返回所有错误字段名

// 嵌套处理
errors.HasNestedErrors()         // 检查是否有嵌套错误
errors.GetNestedErrors("prefix") // 获取指定前缀的嵌套错误
errors.GetRootErrors()           // 获取根级别错误
errors.ToNestedMap()             // 转换为嵌套map结构

// 分组方法
errors.GroupByField()    // 按字段分组
errors.GroupByStruct()   // 按结构体分组
```

## 嵌套验证器

对于复杂的嵌套验证需求，可以使用`NestedValidator`：

```go
nv := validator.NewNestedValidator()

// 验证切片
items := []Item{{Name: "test", Value: 10}}
err := nv.ValidateSlice(items)

// 验证map
configs := map[string]Config{"db": {Host: "localhost"}}
err := nv.ValidateMap(configs)

// 部分验证
err := nv.ValidatePartial(&user, "Username", "Email")

// 排除验证
err := nv.ValidateExcept(&user, "Password")
```

## 工具函数

```go
// 便捷验证
err := validator.ValidateStruct(&user)

// 获取结构体标签
tags := validator.GetStructTags(User{})

// 检查字段是否必填
required := validator.IsRequired(User{}, "Username")

// 获取验证规则
rules := validator.GetValidationRules(User{}, "Username")

// 构建验证标签
tag := validator.BuildValidationTag([]string{"required", "min=3"})
```

## 配置

```go
config := &validator.ValidatorConfig{
    EnableCustomValidations: true,
    EnableTranslations:      true,
    Language:               "zh",
    CustomTranslations: map[string]string{
        "required": "{field}是必填字段",
        "email":    "{field}必须是有效的邮箱地址",
    },
}

v := validator.NewWithConfig(config)
```

## 常用验证标签

### 基础验证
- `required` - 必填
- `omitempty` - 空值时跳过验证

### 字符串验证
- `min=n` - 最小长度
- `max=n` - 最大长度
- `len=n` - 固定长度
- `email` - 邮箱格式
- `url` - URL格式
- `alpha` - 只包含字母
- `alphanum` - 只包含字母和数字
- `numeric` - 只包含数字

### 数值验证
- `gt=n` - 大于
- `gte=n` - 大于等于
- `lt=n` - 小于
- `lte=n` - 小于等于
- `eq=n` - 等于
- `ne=n` - 不等于

### 枚举验证
- `oneof=a b c` - 值必须是其中之一

### 嵌套验证
- `dive` - 深入验证切片/数组元素
- `required` - 嵌套结构体必须存在

### 字段比较
- `eqfield=Field` - 与其他字段相等
- `nefield=Field` - 与其他字段不相等

## 性能

该验证器基于高性能的`go-playground/validator/v10`库构建，在保持功能完整性的同时提供了优异的性能表现。

基准测试结果：
- 基本验证: ~1000ns/op
- 嵌套验证: ~2000ns/op

## 最佳实践

1. **使用结构体标签**: 在结构体定义时就指定验证规则
2. **合理使用omitempty**: 对于可选字段使用omitempty标签
3. **自定义错误消息**: 为用户友好的错误提示注册自定义翻译
4. **嵌套验证**: 对于复杂结构使用嵌套验证器
5. **性能考虑**: 对于高频验证场景，考虑复用验证器实例

## 示例项目

查看`example.go`和`example_test.go`文件获取更多使用示例。