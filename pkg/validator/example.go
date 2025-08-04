package validator

import "time"

// User 用户模型示例
type User struct {
	ID       uint   `json:"id" validate:"required,gt=0"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Age      int    `json:"age" validate:"required,gte=18,lte=120"`
	Mobile   string `json:"mobile" validate:"required,mobile"`
	Role     string `json:"role" validate:"required,oneof=admin user guest"`
	Status   int    `json:"status" validate:"required,oneof=0 1"`
	Profile  Profile `json:"profile" validate:"required"`
}

// Profile 用户资料示例（嵌套结构）
type Profile struct {
	FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
	Birthday  time.Time `json:"birthday" validate:"required"`
	Address   Address   `json:"address" validate:"required"`
}

// Address 地址示例（深度嵌套结构）
type Address struct {
	Country  string `json:"country" validate:"required,min=2,max=50"`
	Province string `json:"province" validate:"required,min=2,max=50"`
	City     string `json:"city" validate:"required,min=2,max=50"`
	Street   string `json:"street" validate:"required,min=5,max=200"`
	ZipCode  string `json:"zip_code" validate:"required,len=6,numeric"`
}

// CreateUserRequest 创建用户请求示例
type CreateUserRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Age             int    `json:"age" validate:"required,gte=18,lte=120"`
	Mobile          string `json:"mobile" validate:"required,mobile"`
	Terms           bool   `json:"terms" validate:"required,eq=true"`
}

// UpdateUserRequest 更新用户请求示例
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=20,alphanum"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Age      *int    `json:"age,omitempty" validate:"omitempty,gte=18,lte=120"`
	Mobile   *string `json:"mobile,omitempty" validate:"omitempty,mobile"`
}

// LoginRequest 登录请求示例
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// PaginationRequest 分页请求示例
type PaginationRequest struct {
	Page     int    `json:"page" validate:"required,gte=1"`
	PageSize int    `json:"page_size" validate:"required,gte=1,lte=100"`
	OrderBy  string `json:"order_by" validate:"omitempty,oneof=id username email created_at updated_at"`
	Order    string `json:"order" validate:"omitempty,oneof=asc desc"`
}

// SearchRequest 搜索请求示例
type SearchRequest struct {
	Keyword  string   `json:"keyword" validate:"required,min=1,max=100"`
	Fields   []string `json:"fields" validate:"required,min=1,dive,oneof=username email mobile"`
	Filters  map[string]interface{} `json:"filters" validate:"omitempty"`
	PaginationRequest
}

// FileUploadRequest 文件上传请求示例
type FileUploadRequest struct {
	FileName string `json:"file_name" validate:"required,min=1,max=255"`
	FileSize int64  `json:"file_size" validate:"required,gt=0,lte=10485760"` // 最大10MB
	FileType string `json:"file_type" validate:"required,oneof=image/jpeg image/png image/gif application/pdf"`
}

// ConfigRequest 配置请求示例
type ConfigRequest struct {
	AppName     string            `json:"app_name" validate:"required,min=3,max=50,alphanum"`
	Debug       bool              `json:"debug"`
	Port        int               `json:"port" validate:"required,gte=1024,lte=65535"`
	Database    DatabaseConfig    `json:"database" validate:"required"`
	Redis       RedisConfig       `json:"redis" validate:"required"`
	JWT         JWTConfig         `json:"jwt" validate:"required"`
	Features    []string          `json:"features" validate:"omitempty,dive,oneof=auth cache logging monitoring"`
	Environment string            `json:"environment" validate:"required,oneof=development staging production"`
	Metadata    map[string]string `json:"metadata" validate:"omitempty"`
}

// DatabaseConfig 数据库配置示例
type DatabaseConfig struct {
	Host     string `json:"host" validate:"required,hostname_rfc1123|ip"`
	Port     int    `json:"port" validate:"required,gte=1,lte=65535"`
	Username string `json:"username" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Database string `json:"database" validate:"required,min=1,max=50,alphanum"`
	Charset  string `json:"charset" validate:"omitempty,oneof=utf8 utf8mb4"`
}

// RedisConfig Redis配置示例
type RedisConfig struct {
	Host     string `json:"host" validate:"required,hostname_rfc1123|ip"`
	Port     int    `json:"port" validate:"required,gte=1,lte=65535"`
	Password string `json:"password" validate:"omitempty,min=6"`
	DB       int    `json:"db" validate:"gte=0,lte=15"`
}

// JWTConfig JWT配置示例
type JWTConfig struct {
	SecretKey       string `json:"secret_key" validate:"required,min=32"`
	TokenDuration   int    `json:"token_duration" validate:"required,gte=300,lte=86400"`   // 5分钟到24小时
	RefreshDuration int    `json:"refresh_duration" validate:"required,gte=3600,lte=604800"` // 1小时到7天
	Issuer          string `json:"issuer" validate:"required,min=3,max=50"`
}