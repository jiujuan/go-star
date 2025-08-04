package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidator(t *testing.T) {
	Convey("验证器测试", t, func() {
		v := New()
		
		Convey("基本验证规则测试", func() {
			Convey("必填字段验证", func() {
				type TestStruct struct {
					Name string `validate:"required"`
				}
				
				// 测试空值
				err := v.Validate(&TestStruct{Name: ""})
				So(err, ShouldNotBeNil)
				
				ve, ok := AsValidationErrors(err)
				So(ok, ShouldBeTrue)
				So(len(ve), ShouldEqual, 1)
				So(ve[0].Field, ShouldEqual, "Name")
				So(ve[0].Tag, ShouldEqual, "required")
				
				// 测试有值
				err = v.Validate(&TestStruct{Name: "test"})
				So(err, ShouldBeNil)
			})
			
			Convey("邮箱验证", func() {
				type TestStruct struct {
					Email string `validate:"required,email"`
				}
				
				// 测试无效邮箱
				err := v.Validate(&TestStruct{Email: "invalid-email"})
				So(err, ShouldNotBeNil)
				
				ve, ok := AsValidationErrors(err)
				So(ok, ShouldBeTrue)
				So(ve[0].Tag, ShouldEqual, "email")
				
				// 测试有效邮箱
				err = v.Validate(&TestStruct{Email: "test@example.com"})
				So(err, ShouldBeNil)
			})
			
			Convey("长度验证", func() {
				type TestStruct struct {
					Username string `validate:"required,min=3,max=20"`
				}
				
				// 测试太短
				err := v.Validate(&TestStruct{Username: "ab"})
				So(err, ShouldNotBeNil)
				
				ve, ok := AsValidationErrors(err)
				So(ok, ShouldBeTrue)
				So(ve[0].Tag, ShouldEqual, "min")
				
				// 测试太长
				err = v.Validate(&TestStruct{Username: "this-is-a-very-long-username"})
				So(err, ShouldNotBeNil)
				
				ve, ok = AsValidationErrors(err)
				So(ok, ShouldBeTrue)
				So(ve[0].Tag, ShouldEqual, "max")
				
				// 测试正确长度
				err = v.Validate(&TestStruct{Username: "testuser"})
				So(err, ShouldBeNil)
			})
			
			Convey("数值范围验证", func() {
				type TestStruct struct {
					Age int `validate:"required,gte=18,lte=120"`
				}
				
				// 测试太小
				err := v.Validate(&TestStruct{Age: 17})
				So(err, ShouldNotBeNil)
				
				// 测试太大
				err = v.Validate(&TestStruct{Age: 121})
				So(err, ShouldNotBeNil)
				
				// 测试正确范围
				err = v.Validate(&TestStruct{Age: 25})
				So(err, ShouldBeNil)
			})
			
			Convey("枚举值验证", func() {
				type TestStruct struct {
					Role string `validate:"required,oneof=admin user guest"`
				}
				
				// 测试无效值
				err := v.Validate(&TestStruct{Role: "invalid"})
				So(err, ShouldNotBeNil)
				
				// 测试有效值
				err = v.Validate(&TestStruct{Role: "admin"})
				So(err, ShouldBeNil)
			})
		})
		
		Convey("自定义验证规则测试", func() {
			Convey("手机号验证", func() {
				type TestStruct struct {
					Mobile string `validate:"required,mobile"`
				}
				
				// 测试无效手机号
				err := v.Validate(&TestStruct{Mobile: "123456"})
				So(err, ShouldNotBeNil)
				
				ve, ok := AsValidationErrors(err)
				So(ok, ShouldBeTrue)
				So(ve[0].Tag, ShouldEqual, "mobile")
				
				// 测试有效手机号
				err = v.Validate(&TestStruct{Mobile: "13812345678"})
				So(err, ShouldBeNil)
			})
			
			Convey("密码强度验证", func() {
				type TestStruct struct {
					Password string `validate:"required,password"`
				}
				
				// 测试弱密码
				err := v.Validate(&TestStruct{Password: "123456"})
				So(err, ShouldNotBeNil)
				
				// 测试强密码
				err = v.Validate(&TestStruct{Password: "abc123"})
				So(err, ShouldBeNil)
			})
			
			Convey("身份证号验证", func() {
				type TestStruct struct {
					IDCard string `validate:"required,idcard"`
				}
				
				// 测试无效身份证号
				err := v.Validate(&TestStruct{IDCard: "123"})
				So(err, ShouldNotBeNil)
				
				// 测试有效身份证号
				err = v.Validate(&TestStruct{IDCard: "123456789012345"})
				So(err, ShouldBeNil)
			})
		})
		
		Convey("注册自定义验证规则", func() {
			// 注册自定义验证规则
			err := v.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "custom"
			})
			So(err, ShouldBeNil)
			
			// 注册自定义错误消息
			err = v.RegisterTranslation("custom", "{field} must be 'custom'")
			So(err, ShouldBeNil)
			
			type TestStruct struct {
				Value string `validate:"required,custom"`
			}
			
			// 测试自定义验证规则
			err = v.Validate(&TestStruct{Value: "invalid"})
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(ve[0].Tag, ShouldEqual, "custom")
			So(ve[0].Message, ShouldContainSubstring, "must be 'custom'")
			
			// 测试通过验证
			err = v.Validate(&TestStruct{Value: "custom"})
			So(err, ShouldBeNil)
		})
	})
}

func TestNestedValidation(t *testing.T) {
	Convey("嵌套验证测试", t, func() {
		v := New()
		
		Convey("简单嵌套结构验证", func() {
			type Address struct {
				Street  string `validate:"required,min=5"`
				City    string `validate:"required,min=2"`
				ZipCode string `validate:"required,len=6,numeric"`
			}
			
			type User struct {
				Name    string  `validate:"required,min=2"`
				Address Address `validate:"required"`
			}
			
			user := User{
				Name: "John",
				Address: Address{
					Street:  "123",  // 太短
					City:    "A",    // 太短
					ZipCode: "abc",  // 不是数字且长度不对
				},
			}
			
			err := v.Validate(&user)
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldEqual, 3)
			
			// 检查嵌套字段路径
			fieldMap := ve.Map()
			So(fieldMap["Address.Street"], ShouldNotBeEmpty)
			So(fieldMap["Address.City"], ShouldNotBeEmpty)
			So(fieldMap["Address.ZipCode"], ShouldNotBeEmpty)
		})
		
		Convey("深度嵌套结构验证", func() {
			type Country struct {
				Name string `validate:"required,min=2"`
				Code string `validate:"required,len=2"`
			}
			
			type Address struct {
				Street  string  `validate:"required,min=5"`
				City    string  `validate:"required,min=2"`
				Country Country `validate:"required"`
			}
			
			type Profile struct {
				FirstName string  `validate:"required,min=2"`
				LastName  string  `validate:"required,min=2"`
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
					LastName:  "",  // 必填
					Address: Address{
						Street: "123", // 太短
						City:   "",    // 必填
						Country: Country{
							Name: "", // 必填
							Code: "USA", // 长度不对
						},
					},
				},
			}
			
			err := v.Validate(&user)
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldBeGreaterThan, 5)
			
			// 检查深度嵌套字段路径
			fieldMap := ve.Map()
			So(fieldMap["Profile.Address.Country.Name"], ShouldNotBeEmpty)
			So(fieldMap["Profile.Address.Country.Code"], ShouldNotBeEmpty)
		})
		
		Convey("指针嵌套结构验证", func() {
			type Address struct {
				Street string `validate:"required,min=5"`
				City   string `validate:"required,min=2"`
			}
			
			type User struct {
				Name    string   `validate:"required,min=2"`
				Address *Address `validate:"required"`
			}
			
			// 测试nil指针
			user := User{
				Name:    "John",
				Address: nil,
			}
			
			err := v.Validate(&user)
			So(err, ShouldNotBeNil)
			
			// 测试有效指针
			user.Address = &Address{
				Street: "123 Main St",
				City:   "New York",
			}
			
			err = v.Validate(&user)
			So(err, ShouldBeNil)
		})
		
		Convey("切片嵌套验证", func() {
			type Tag struct {
				Name string `validate:"required,min=2"`
			}
			
			type Post struct {
				Title string `validate:"required,min=5"`
				Tags  []Tag  `validate:"required,min=1,dive"`
			}
			
			post := Post{
				Title: "Test", // 太短
				Tags: []Tag{
					{Name: "A"}, // 太短
					{Name: ""},  // 必填
				},
			}
			
			err := v.Validate(&post)
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldBeGreaterThan, 2)
		})
	})
}

func TestNestedValidator(t *testing.T) {
	Convey("嵌套验证器测试", t, func() {
		nv := NewNestedValidator()
		
		Convey("验证切片", func() {
			type Item struct {
				Name  string `validate:"required,min=3"`
				Value int    `validate:"required,gt=0"`
			}
			
			items := []Item{
				{Name: "ab", Value: 0},    // 两个错误
				{Name: "valid", Value: 5}, // 正确
				{Name: "", Value: -1},     // 两个错误
			}
			
			err := nv.ValidateSlice(items)
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldEqual, 4) // 总共4个错误
		})
		
		Convey("验证map", func() {
			type Config struct {
				Host string `validate:"required,min=3"`
				Port int    `validate:"required,gt=0"`
			}
			
			configs := map[string]Config{
				"db": {Host: "ab", Port: 0},       // 两个错误
				"redis": {Host: "localhost", Port: 6379}, // 正确
				"cache": {Host: "", Port: -1},     // 两个错误
			}
			
			err := nv.ValidateMap(configs)
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldEqual, 4) // 总共4个错误
		})
		
		Convey("部分验证", func() {
			type User struct {
				Username string `validate:"required,min=3"`
				Email    string `validate:"required,email"`
				Age      int    `validate:"required,gte=18"`
			}
			
			user := User{
				Username: "ab",              // 错误
				Email:    "invalid-email",   // 错误
				Age:      17,                // 错误
			}
			
			// 只验证Username字段
			err := nv.ValidatePartial(&user, "Username")
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(len(ve), ShouldEqual, 1)
			So(ve[0].Field, ShouldContainSubstring, "Username")
		})
		
		Convey("排除验证", func() {
			type User struct {
				Username string `validate:"required,min=3"`
				Email    string `validate:"required,email"`
				Age      int    `validate:"required,gte=18"`
			}
			
			user := User{
				Username: "validuser",
				Email:    "invalid-email",   // 错误
				Age:      17,                // 错误
			}
			
			// 排除Email和Age字段
			err := nv.ValidateExcept(&user, "Email", "Age")
			So(err, ShouldBeNil) // Username是正确的
		})
		
		Convey("获取验证深度", func() {
			type Level3 struct {
				Value string `validate:"required"`
			}
			
			type Level2 struct {
				Level3 Level3 `validate:"required"`
			}
			
			type Level1 struct {
				Level2 Level2 `validate:"required"`
			}
			
			depth := nv.GetValidationDepth(Level1{})
			So(depth, ShouldEqual, 2) // 0-based depth
		})
	})
}

func TestValidationErrors(t *testing.T) {
	Convey("验证错误测试", t, func() {
		Convey("错误分组", func() {
			errors := ValidationErrors{
				{Field: "name", Message: "name is required"},
				{Field: "email", Message: "email is invalid"},
				{Field: "profile.age", Message: "age must be greater than 18"},
				{Field: "profile.address.city", Message: "city is required"},
			}
			
			Convey("按字段分组", func() {
				grouped := errors.GroupByField()
				So(len(grouped), ShouldEqual, 4)
				So(len(grouped["name"]), ShouldEqual, 1)
			})
			
			Convey("按结构体分组", func() {
				grouped := errors.GroupByStruct()
				So(len(grouped), ShouldEqual, 3)
				So(len(grouped["root"]), ShouldEqual, 2)
				So(len(grouped["profile"]), ShouldEqual, 1)
				So(len(grouped["profile.address"]), ShouldEqual, 1)
			})
			
			Convey("获取嵌套错误", func() {
				nested := errors.GetNestedErrors("profile")
				So(len(nested), ShouldEqual, 2)
				So(nested[0].Field, ShouldEqual, "age")
				So(nested[1].Field, ShouldEqual, "address.city")
			})
			
			Convey("检查嵌套错误", func() {
				So(errors.HasNestedErrors(), ShouldBeTrue)
				
				rootErrors := errors.GetRootErrors()
				So(len(rootErrors), ShouldEqual, 2)
			})
			
			Convey("转换为嵌套map", func() {
				nestedMap := errors.ToNestedMap()
				So(nestedMap["name"], ShouldEqual, "name is required")
				So(nestedMap["email"], ShouldEqual, "email is invalid")
				
				profile, ok := nestedMap["profile"].(map[string]interface{})
				So(ok, ShouldBeTrue)
				So(profile["age"], ShouldEqual, "age must be greater than 18")
				
				address, ok := profile["address"].(map[string]interface{})
				So(ok, ShouldBeTrue)
				So(address["city"], ShouldEqual, "city is required")
			})
		})
		
		Convey("错误格式化", func() {
			errors := ValidationErrors{
				{Field: "name", Tag: "required", Message: "name is required"},
				{Field: "email", Tag: "email", Message: "email is invalid"},
			}
			
			Convey("字符串格式", func() {
				str := errors.String()
				So(str, ShouldContainSubstring, "name is required")
				So(str, ShouldContainSubstring, "email is invalid")
			})
			
			Convey("JSON格式", func() {
				json := errors.JSON()
				So(json, ShouldContainSubstring, "name")
				So(json, ShouldContainSubstring, "email")
			})
			
			Convey("Map格式", func() {
				m := errors.Map()
				So(m["name"], ShouldEqual, "name is required")
				So(m["email"], ShouldEqual, "email is invalid")
			})
		})
	})
}

func TestValidatorUtils(t *testing.T) {
	Convey("验证器工具函数测试", t, func() {
		Convey("验证结构体", func() {
			type TestStruct struct {
				Name string `validate:"required,min=3"`
			}
			
			err := ValidateStruct(&TestStruct{Name: "ab"})
			So(err, ShouldNotBeNil)
			
			err = ValidateStruct(&TestStruct{Name: "valid"})
			So(err, ShouldBeNil)
		})
		
		Convey("获取结构体标签", func() {
			type TestStruct struct {
				Name  string `json:"name" validate:"required,min=3"`
				Email string `json:"email" validate:"required,email"`
				Age   int    `json:"age" validate:"gte=18"`
			}
			
			tags := GetStructTags(TestStruct{})
			So(len(tags), ShouldEqual, 3)
			So(tags["name"], ShouldEqual, "required,min=3")
			So(tags["email"], ShouldEqual, "required,email")
			So(tags["age"], ShouldEqual, "gte=18")
		})
		
		Convey("检查必填字段", func() {
			type TestStruct struct {
				Name     string `validate:"required"`
				Optional string `validate:"omitempty,min=3"`
			}
			
			So(IsRequired(TestStruct{}, "Name"), ShouldBeTrue)
			So(IsRequired(TestStruct{}, "Optional"), ShouldBeFalse)
		})
		
		Convey("获取验证规则", func() {
			type TestStruct struct {
				Name string `validate:"required,min=3,max=20"`
			}
			
			rules := GetValidationRules(TestStruct{}, "Name")
			So(len(rules), ShouldEqual, 3)
			So(rules[0], ShouldEqual, "required")
			So(rules[1], ShouldEqual, "min=3")
			So(rules[2], ShouldEqual, "max=20")
		})
		
		Convey("检查验证规则", func() {
			type TestStruct struct {
				Name string `validate:"required,min=3,max=20"`
			}
			
			So(HasValidationRule(TestStruct{}, "Name", "required"), ShouldBeTrue)
			So(HasValidationRule(TestStruct{}, "Name", "email"), ShouldBeFalse)
		})
		
		Convey("提取验证参数", func() {
			rule, param := ExtractValidationParam("min=3")
			So(rule, ShouldEqual, "min")
			So(param, ShouldEqual, "3")
			
			rule, param = ExtractValidationParam("required")
			So(rule, ShouldEqual, "required")
			So(param, ShouldEqual, "")
		})
		
		Convey("构建验证标签", func() {
			rules := []string{"required", "min=3", "max=20"}
			tag := BuildValidationTag(rules)
			So(tag, ShouldEqual, "required,min=3,max=20")
		})
		
		Convey("合并验证标签", func() {
			tag := MergeValidationTags("required,min=3", "max=20,email")
			So(tag, ShouldContainSubstring, "required")
			So(tag, ShouldContainSubstring, "min=3")
			So(tag, ShouldContainSubstring, "max=20")
			So(tag, ShouldContainSubstring, "email")
		})
	})
}

func TestValidatorConfig(t *testing.T) {
	Convey("验证器配置测试", t, func() {
		Convey("默认配置", func() {
			config := DefaultConfig()
			So(config.EnableCustomValidations, ShouldBeTrue)
			So(config.EnableTranslations, ShouldBeTrue)
			So(config.Language, ShouldEqual, "en")
		})
		
		Convey("使用配置创建验证器", func() {
			config := &ValidatorConfig{
				EnableCustomValidations: true,
				EnableTranslations:      true,
				Language:               "zh",
				CustomTranslations: map[string]string{
					"required": "{field}是必填字段",
				},
			}
			
			v := NewWithConfig(config)
			So(v, ShouldNotBeNil)
			
			type TestStruct struct {
				Name string `validate:"required"`
			}
			
			err := v.Validate(&TestStruct{Name: ""})
			So(err, ShouldNotBeNil)
			
			ve, ok := AsValidationErrors(err)
			So(ok, ShouldBeTrue)
			So(ve[0].Message, ShouldContainSubstring, "是必填字段")
		})
	})
}

// 基准测试
func BenchmarkValidator(b *testing.B) {
	v := New()
	
	type User struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"required,gte=18,lte=120"`
	}
	
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Age:      25,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Validate(&user)
	}
}

func BenchmarkNestedValidator(b *testing.B) {
	nv := NewNestedValidator()
	
	type Address struct {
		Street  string `validate:"required,min=5"`
		City    string `validate:"required,min=2"`
		ZipCode string `validate:"required,len=6"`
	}
	
	type User struct {
		Username string  `validate:"required,min=3,max=20"`
		Email    string  `validate:"required,email"`
		Address  Address `validate:"required"`
	}
	
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			ZipCode: "123456",
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nv.ValidateNested(&user)
	}
}