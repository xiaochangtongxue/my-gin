package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
)

// Init 初始化验证器
func Init() {
	once.Do(func() {
		// 创建翻译器
		uni = ut.New(zh.New())

		// 获取中文翻译器
		trans, _ = uni.GetTranslator("zh")

		// 创建验证器
		validate = validator.New()

		validate.SetTagName("binding")

		// 注册字段名获取函数，优先使用 label tag，其次使用 json tag
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			if name == "" {
				name = fld.Tag.Get("json")
				// 去掉 json tag 中的 omitempty 等选项
				if idx := strings.Index(name, ","); idx != -1 {
					name = name[:idx]
				}
			}
			return name
		})

		// 注册中文翻译
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)

		// 注册自定义验证器及其中文翻译
		registerCustomValidators()

		// 替换gin的默认验证器
		binding.Validator = new(defaultValidator)
	})
}

// Get 获取验证器实例
func Get() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

// GetTranslator 获取翻译器
func GetTranslator() ut.Translator {
	if trans == nil {
		Init()
	}
	return trans
}

// TranslateError 翻译验证错误
func TranslateError(err error) string {
	if err == nil {
		return ""
	}

	// 验证错误
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var builder strings.Builder
		for i, e := range validationErrors {
			if i > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(e.Translate(trans))
		}
		return builder.String()
	}

	return err.Error()
}

// registerCustomValidators 注册自定义验证器及其中文翻译
func registerCustomValidators() {
	// 注册手机号验证器
	_ = validate.RegisterValidation("phone", validatePhone)
	// 注册身份证号验证器
	_ = validate.RegisterValidation("idcard", validateIDCard)
	// 注册密码强度验证器
	_ = validate.RegisterValidation("password", validatePassword)

	// 2. 注册翻译 (关键步骤)
	registerTagTranslation("phone", "{0}格式不合法")
	registerTagTranslation("idcard", "{0}格式不正确")
	registerTagTranslation("password", "{0}必须包含字母和数字，且长度不小于8位")
}

// validatePhone 手机号验证
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}
	// 简单验证：1开头，第二位为3-9
	if phone[0] != '1' {
		return false
	}
	second := phone[1]
	return second >= '3' && second <= '9'
}

// validateIDCard 身份证号验证
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if len(idCard) != 18 {
		return false
	}
	// 简单验证：前17位为数字，最后一位为数字或X
	for i := 0; i < 17; i++ {
		if idCard[i] < '0' || idCard[i] > '9' {
			return false
		}
	}
	last := idCard[17]
	return (last >= '0' && last <= '9') || last == 'X' || last == 'x'
}

// validatePassword 密码强度验证（至少8位，包含字母和数字）
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	var hasLetter, hasDigit bool
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

// 辅助函数：快速注册标签翻译
func registerTagTranslation(tag string, msg string) {
	_ = validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// defaultValidator gin默认验证器实现
type defaultValidator struct {
	v *validator.Validate
}

// ValidateStruct 验证结构体
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kind := reflect.ValueOf(obj).Kind(); kind == reflect.Struct || kind == reflect.Ptr {
		return Get().Struct(obj)
	}
	return nil
}

// Engine 获取验证引擎
func (v *defaultValidator) Engine() interface{} {
	return Get()
}

// GetErrorMsg 获取字段错误消息
func GetErrorMsg(field, tag string) string {
	switch tag {
	case "required":
		return field + "为必填项"
	case "email":
		return field + "格式不正确"
	case "phone":
		return field + "格式不正确"
	case "idcard":
		return field + "格式不正确"
	case "password":
		return field + "强度不足（至少8位，包含字母和数字）"
	case "min":
		return field + "长度不足"
	case "max":
		return field + "长度超出限制"
	case "len":
		return field + "长度不符合要求"
	default:
		return field + "验证失败"
	}
}
