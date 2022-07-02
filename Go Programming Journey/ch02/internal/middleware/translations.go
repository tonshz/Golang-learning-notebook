package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	// 多语言包
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	// 通用翻译器
	"github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	// validator 的翻译器
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		// 通过 GetHeader 方法获取约定的 header 参数 locale,用于辨别当前请求的语言类别是en还是zh
		// 如果有其他语言环境要求，也可以继续引入其他语言类别，go-playground/locales 基本上都支持。
		locale := c.GetHeader("locale")
		// 对应语言的 Translator
		trans, _ := uni.GetTranslator(locale)
		// 验证器
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			// 调用 RegisterDefaultTranslations 方法将验证器和对应语言类型的 Translator 注册进来
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			}
			// 将 Translator 存储到全局上下文中，便于后续翻译时使用
			c.Set("trans", trans)
		}
		c.Next()
	}
}
