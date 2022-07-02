package word

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

// 全部转换成大写/小写，使用标准库中的原生方法进行转换
func ToUpper(s string) string {
	return strings.ToUpper(s)
}
func ToLower(s string) string {
	return strings.ToLower(s)
}

// 下划线转大写驼峰
func UnderscoreToUpperCamelCase(s string) string {
	// 将下划线替换为空格字符
	s = strings.Replace(s, "_", " ", -1)
	// 将所有字符修改为其对应的首字母大写的格式,strings.Title()被抛弃了
	caser := cases.Title(language.English)
	s = caser.String(s)
	//s = strings.Title(s)
	// 将先前的空格字符替换为空
	return strings.Replace(s, " ", "", -1)
}

// 下划线转小写驼峰
func UnderscoreToLowerCamelCase(s string) string {
	// 主体逻辑可以直接复用大写驼峰的转换方法
	s = UnderscoreToUpperCamelCase(s)
	// 只需要对首字母进行处理即可，将字符串中的第一位取出，使用unicode.ToLower()转换即可
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// 驼峰转下划线
func CamelCaseToUnderscore(s string) string {
	var output []rune
	// 将字符串转换为小写的同时添加下划线
	for i, r := range s {
		// 首字母较为特殊，在其前不需要添加下划线
		if i == 0 {
			// 拼接
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
