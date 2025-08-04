package validator

// ValidatorConfig 验证器配置
type ValidatorConfig struct {
	// EnableCustomValidations 是否启用自定义验证规则
	EnableCustomValidations bool `mapstructure:"enable_custom_validations"`
	// EnableTranslations 是否启用错误消息翻译
	EnableTranslations bool `mapstructure:"enable_translations"`
	// Language 错误消息语言
	Language string `mapstructure:"language"`
	// CustomValidations 自定义验证规则配置
	CustomValidations map[string]string `mapstructure:"custom_validations"`
	// CustomTranslations 自定义错误消息翻译
	CustomTranslations map[string]string `mapstructure:"custom_translations"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *ValidatorConfig {
	return &ValidatorConfig{
		EnableCustomValidations: true,
		EnableTranslations:      true,
		Language:               "en",
		CustomValidations:      make(map[string]string),
		CustomTranslations:     make(map[string]string),
	}
}

// NewWithConfig 使用配置创建验证器
func NewWithConfig(config *ValidatorConfig) Validator {
	if config == nil {
		config = DefaultConfig()
	}
	
	v := New().(*ValidatorImpl)
	
	// 应用自定义翻译
	if config.EnableTranslations {
		for tag, message := range config.CustomTranslations {
			v.RegisterTranslation(tag, message)
		}
	}
	
	return v
}