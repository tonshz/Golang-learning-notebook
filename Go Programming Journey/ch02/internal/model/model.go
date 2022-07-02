package model

import (
	"demo/ch02/global"
	"demo/ch02/pkg/setting"
	"fmt"
	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/jinzhu/gorm"
	"time"

	// 引入 MYSQL驱动库进行初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

// 公共字段
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// 新增 NewDBEngine()
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	// gorm.Open() 初始化一个 MySQL 连接，首先需要导入驱动
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		// 显示日志输出
		db.LogMode(true)
	}
	// 使用单表操作
	db.SingularTable(true)

	// 注册回调行为
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	// 设置空闲连接池中的最大连接数
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	// 设置数据库的最大打开连接数。
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	// 添加 OpenTracing 注册回调
	otgorm.AddGormCallbacks(db)
	return db, nil
}

// 更新时间戳的新增行为的回调
// 当对数据库执行任何操作时，Scope 包含当前操作的信息
// Scope 允许复用通用的逻辑
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		// 获取当前时间
		nowTime := time.Now().Unix()
		// scope.FiledByName() 获取当前是否包含所需字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			// 若创建时间为空，则设置创建时间为当前时间
			// 通过判断 Filed.IsBlank 的值，可得知该字段值是否为空
			if createTimeField.IsBlank {
				// 通过 Filed.Set() 为字段赋值
				_ = createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			// 若修改时间为空，则设置修改时间为当前时间
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新时间戳的更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	// scope.Get() 获取当前设置了标识 gorm:update_column 的字段属性
	if _, ok := scope.Get("gorm:update_column"); !ok {
		// 若不存在，即未自定义设置 update_column
		// 设置默认字段 ModifiedOn 的值为当前时间戳
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		// 通过 scope.Get() 获取当前设置了标识 gorm:delete_option 的字段属性
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		// 判断是否存在 DeletedOn 与 IsDel 字段
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")

		// 若存在执行 UPDATE 进行软删除
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			// scope.Raw() 设置原始 sql
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				// scope.QuotedTableName() 获取当前所引用的表名
				scope.QuotedTableName(),
				// Quote() 使用引用字符串对数据库进行转义
				scope.Quote(deletedOnField.DBName),
				// AddToVars() 添加 value 作为 sql 的 vars，用于防止 SQL 注入
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				// scope.CombinedConditionSql() 返回组合条件的 sql
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec() // scope.Exec() 执行生成的sql
		} else {
			// 否则执行 DELETE 进行硬删除
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				// 获取表名
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// 若 str 不为空添加额外的空格
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
