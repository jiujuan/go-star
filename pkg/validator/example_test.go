package validator

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

// ExampleValidator_Validate 基本验证示例
func ExampleValidator_Validate() {
	v := New()
	
	type User struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"required,gte=18,lte=120"`
	}
	
	// 验证有效数据
	user := User{
		Username: "john_doe",
		Email:    "john@example.com",
		Age:      25,
	}
	
	err := v.Validate(&user)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed")
	}
	
	// Output: Validation passed
}

// ExampleValidator_Validate_withErrors 验证错误示例
func ExampleValidator_Validate_withErrors() {
	v := New()
	
	type User struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"required,gte=18,lte=120"`
	}
	
	// 验证无效数据
	user := User{
		Username: "jo",              // 太短
		Email:    "invalid-email",   // 无效邮箱
		Age:      17,                // 年龄太小
	}
	
	err := v.Validate(&user)
	if err != nil {
		if ve, ok := AsValidationErrors(err); ok {
			fmt.Printf("Found %d validation errors:\n", len(ve))
			for _, e := range ve {
				fmt.Printf("- %s: %s\n", e.Field, e.Message)
			}
		}
	}
	
	// Output:
	// Found 3 validation errors:
	// - Username: Username must be at least 3 characters long
	// - Email: Email must be a valid email address
	// - Age: Age must be greater than or equal to 18
}

// ExampleValidator_RegisterValidation 自定义验证规则示例
func ExampleValidator_RegisterValidation() {
	v := New()
	
	// 注册自定义验证规则：检查是否为工作日
	err := v.RegisterValidation("weekday", func(fl validator.FieldLevel) bool {
		if date, ok := fl.Field().Interface().(time.Time); ok {
			weekday := date.Weekday()
			return weekday >= time.Monday && weekday <= time.Friday
		}
		return false
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// 注册自定义错误消息
	v.RegisterTranslation("weekday", "{field} must be a weekday")
	
	type Meeting struct {
		Title string    `validate:"required"`
		Date  time.Time `validate:"required,weekday"`
	}
	
	// 测试周末日期
	meeting := Meeting{
		Title: "Team Meeting",
		Date:  time.Date(2024, 1, 6, 10, 0, 0, 0, time.UTC), // Saturday
	}
	
	err = v.Validate(&meeting)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}
	
	// Output: Validation failed: Date must be a weekday
}

// ExampleNestedValidator_ValidateNested 嵌套验证示例
func ExampleNestedValidator_ValidateNested() {
	nv := NewNestedValidator()
	
	type Address struct {
		Street  string `validate:"required,min=5"`
		City    string `validate:"required,min=2"`
		ZipCode string `validate:"required,len=6,numeric"`
	}
	
	type User struct {
		Name    string  `validate:"required,min=2"`
		Email   string  `validate:"required,email"`
		Address Address `validate:"required"`
	}
	
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Address: Address{
			Street:  "123 Main Street",
			City:    "New York",
			ZipCode: "123456",
		},
	}
	
	err := nv.ValidateNested(&user)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Nested validation passed")
	}
	
	// Output: Nested validation passed
}

// ExampleNestedValidator_ValidateSlice 切片验证示例
func ExampleNestedValidator_ValidateSlice() {
	nv := NewNestedValidator()
	
	type Product struct {
		Name  string  `validate:"required,min=3"`
		Price float64 `validate:"required,gt=0"`
	}
	
	products := []Product{
		{Name: "Laptop", Price: 999.99},
		{Name: "Mouse", Price: 29.99},
		{Name: "KB", Price: 0}, // 名称太短，价格无效
	}
	
	err := nv.ValidateSlice(products)
	if err != nil {
		if ve, ok := AsValidationErrors(err); ok {
			fmt.Printf("Found %d validation errors in slice:\n", len(ve))
			for _, e := range ve {
				fmt.Printf("- %s: %s\n", e.Field, e.Message)
			}
		}
	}
	
	// Output:
	// Found 2 validation errors in slice:
	// - [2].Name: Name must be at least 3 characters long
	// - [2].Price: Price is required
}

// ExampleValidationErrors_ToNestedMap 嵌套错误映射示例
func ExampleValidationErrors_ToNestedMap() {
	v := New()
	
	type Address struct {
		Street string `validate:"required,min=5"`
		City   string `validate:"required"`
	}
	
	type Profile struct {
		FirstName string  `validate:"required,min=2"`
		Address   Address `validate:"required"`
	}
	
	type User struct {
		Username string  `validate:"required,min=3"`
		Profile  Profile `validate:"required"`
	}
	
	user := User{
		Username: "jo", // 太短
		Profile: Profile{
			FirstName: "J", // 太短
			Address: Address{
				Street: "123", // 太短
				City:   "",    // 必填
			},
		},
	}
	
	err := v.Validate(&user)
	if err != nil {
		if ve, ok := AsValidationErrors(err); ok {
			nestedMap := ve.ToNestedMap()
			
			// 打印嵌套错误结构
			fmt.Printf("Username error: %v\n", nestedMap["Username"])
			
			if profile, ok := nestedMap["Profile"].(map[string]interface{}); ok {
				fmt.Printf("Profile.FirstName error: %v\n", profile["FirstName"])
				
				if address, ok := profile["Address"].(map[string]interface{}); ok {
					fmt.Printf("Profile.Address.Street error: %v\n", address["Street"])
					fmt.Printf("Profile.Address.City error: %v\n", address["City"])
				}
			}
		}
	}
	
	// Output:
	// Username error: Username must be at least 3 characters long
	// Profile.FirstName error: Profile.FirstName must be at least 2 characters long
	// Profile.Address.Street error: Profile.Address.Street must be at least 5 characters long
	// Profile.Address.City error: Profile.Address.City is required
}

// ExampleValidatorConfig 配置示例
func ExampleValidatorConfig() {
	config := &ValidatorConfig{
		EnableCustomValidations: true,
		EnableTranslations:      true,
		Language:               "zh",
		CustomTranslations: map[string]string{
			"required": "{field}是必填字段",
			"email":    "{field}必须是有效的邮箱地址",
			"min":      "{field}长度不能少于{param}个字符",
		},
	}
	
	v := NewWithConfig(config)
	
	type User struct {
		Name  string `json:"姓名" validate:"required,min=2"`
		Email string `json:"邮箱" validate:"required,email"`
	}
	
	user := User{
		Name:  "李",              // 太短
		Email: "invalid-email", // 无效邮箱
	}
	
	err := v.Validate(&user)
	if err != nil {
		if ve, ok := AsValidationErrors(err); ok {
			for _, e := range ve {
				fmt.Printf("%s\n", e.Message)
			}
		}
	}
	
	// Output:
	// Name长度不能少于2个字符
	// Email必须是有效的邮箱地址
}

// ExampleValidateStruct 便捷函数示例
func ExampleValidateStruct() {
	type LoginRequest struct {
		Username string `validate:"required,min=3"`
		Password string `validate:"required,min=6"`
	}
	
	request := LoginRequest{
		Username: "user123",
		Password: "secret123",
	}
	
	err := ValidateStruct(&request)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Login request is valid")
	}
	
	// Output: Login request is valid
}

// ExampleGetStructTags 获取结构体标签示例
func ExampleGetStructTags() {
	type User struct {
		Username string `json:"username" validate:"required,min=3,max=20"`
		Email    string `json:"email" validate:"required,email"`
		Age      int    `json:"age" validate:"omitempty,gte=18"`
	}
	
	tags := GetStructTags(User{})
	
	// Sort the output for consistent results
	fields := []string{"username", "email", "age"}
	for _, field := range fields {
		if tag, exists := tags[field]; exists {
			fmt.Printf("%s: %s\n", field, tag)
		}
	}
	
	// Output:
	// username: required,min=3,max=20
	// email: required,email
	// age: omitempty,gte=18
}

// ExampleValidationErrors_GroupByStruct 按结构体分组示例
func ExampleValidationErrors_GroupByStruct() {
	v := New()
	
	type Address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
	}
	
	type User struct {
		Name    string  `validate:"required"`
		Email   string  `validate:"required,email"`
		Address Address `validate:"required"`
	}
	
	user := User{
		Name:  "",                // 必填
		Email: "invalid-email",   // 无效
		Address: Address{
			Street: "", // 必填
			City:   "", // 必填
		},
	}
	
	err := v.Validate(&user)
	if err != nil {
		if ve, ok := AsValidationErrors(err); ok {
			grouped := ve.GroupByStruct()
			
			fmt.Printf("Root level errors: %d\n", len(grouped["root"]))
			fmt.Printf("Address errors: %d\n", len(grouped["Address"]))
			
			// Print in consistent order
			if rootErrors, exists := grouped["root"]; exists {
				fmt.Printf("\nroot:\n")
				for _, e := range rootErrors {
					fmt.Printf("  - %s: %s\n", e.Field, e.Message)
				}
			}
			
			if addressErrors, exists := grouped["Address"]; exists {
				fmt.Printf("\nAddress:\n")
				for _, e := range addressErrors {
					fmt.Printf("  - %s: %s\n", e.Field, e.Message)
				}
			}
		}
	}
	
	// Output:
	// Root level errors: 2
	// Address errors: 2
	//
	// root:
	//   - Name: Name is required
	//   - Email: Email must be a valid email address
	//
	// Address:
	//   - Address.Street: Address.Street is required
	//   - Address.City: Address.City is required
}