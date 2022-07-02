package global

import "github.com/jinzhu/gorm"

var (
	// DBEngine 变量包含了当前数据库连接的信息
	DBEngine *gorm.DB
)
