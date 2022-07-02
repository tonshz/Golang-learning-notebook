package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"strings"
)

// 实现 error 接口
type ValidError struct {
	Key     string
	Message string
}
type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}
func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}
func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	// 通过 ShouldBind 进行参数绑定和入参校验
	// ShouldBind 能够基于请求的不同，自动提取JSON、form表单和QueryString类型的数据，并把值绑定到指定的结构体对象
	// 此方法从上下文中获取传入的方法入参并进行绑定
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}
		// 通过中间件 Translations 设置的 Translator 来对错误消息体进行翻译
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		return false, errs
	}
	return true, nil
}
