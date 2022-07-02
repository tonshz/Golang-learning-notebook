# Go 语言编程之旅(二)：HTTP 应用(四) 

## 六、模块开发：标签管理

在初步完成了业务接口的入参校验的逻辑处理后，接下来进入正式的业务模块的业务逻辑开发，在本章节将完成标签模块的接口代码编写，涉及的接口如下：

| 功能         | HTTP 方法 | 路径      |
| :----------- | :-------- | :-------- |
| 新增标签     | POST      | /tags     |
| 删除指定标签 | DELETE    | /tags/:id |
| 更新指定标签 | PUT       | /tags/:id |
| 获取标签列表 | GET       | /tags     |

### 1. 新建 model 方法

首先需要针对标签表进行处理，修改 `internal/model` 目录下的` tag.go `文件，针对标签模块的模型操作进行封装，并且只与实体产生关系，代码如下：

```go
package model

import (
   "demo/ch02/pkg/app"
   "github.com/jinzhu/gorm"
)

type Tag struct {
   *Model
   Name  string `json:"name"`
   State uint8  `json:"state"`
}

// tag.go
type TagSwagger struct {
   List  []*Tag
   Pager *app.Pager
}

func (t Tag) TableName() string {
   return "blog_tag"
}

// 使用 db *grom.DB 作为函数首参数传入
func (t Tag) Count(db *gorm.DB) (int, error) {
   var count int
   if t.Name != "" {
      // Where 设置筛选条件，接受 map、struct、string作为条件
      db = db.Where("name = ?", t.Name)
   }
   db = db.Where("state = ?", t.State)
   // Model 指定运行 DB 操作的模型实例，默认解析该结构体的名字为表名
   // Count 统计行为，用于统计模型的记录数
   if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
      return 0, err
   }
   return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
   var tags []*Tag
   var err error
   if pageOffset >= 0 && pageSize > 0 {
      // Offset 偏移量，用于指定开始返回记录之前要跳过的记录数
      // Limit 限制检索的记录数
      db = db.Offset(pageOffset).Limit(pageSize)
   }
   if t.Name != "" {
      db = db.Where("name = ?", t.Name)
   }
   db = db.Where("state = ?", t.State)
   // Find 有两个参数，out 是数据接收者，where 是查询条件，可以代替 Where 来传入条件
   // err = e.g. db.Find(&tags, "is_del = 0").Error
   if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
      return nil, err
   }
   return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
   return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB) error {
   // Update 更新所选字段
   return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error
}

func (t Tag) Delete(db *gorm.DB) error {
   // Delete 删除数据
   return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
```

需要注意的是，在上述代码中，采用的是将`db *gorm.DB`作为函数首参数传入的方式，在实际开发中也可以基于结构体传入。

### 2. 处理 model 回调

在编写 model 代码时，并没有针对公共字段  created_on、modified_on、deleted_on、is_del 进行处理。可以通过设置 `model callback`的方式实现公共字段的处理，本项目使用的 ORM 是 GORM，其本身是提供回调支持的，因此可以根据自己的需要自定义 GORM 的回调操作，而在 GORM 中，可以分别进行如下的回调相关行为。

- 注册一个新的回调。
- 删除现有的回调。
- 替换现有的回调。
- 注册回调的先后顺序。

在本项目中使用到的“替换现有的回调”这一行为，修改`internal/model` 目录下的 model.go 文件，准备开始编写 model 的回调代码，下述所新增的回调代码均写入在` NewDBEngine `方法后。

```go
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {}
func updateTimeStampForCreateCallback(scope *gorm.Scope) {}
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {}
func deleteCallback(scope *gorm.Scope) {}
func addExtraSpaceIfExist(str string) string {}
```

#### a. 新增行为的回调

```go
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
```

- 通过调用 `scope.FieldByName` 方法，获取当前是否包含所需的字段。
- 通过判断 `Field.IsBlank` 的值，可以得知该字段的值是否为空。
- 若为空，则会调用 `Field.Set` 方法给该字段设置值，入参类型为 interface{}，内部也就是通过反射进行一系列操作赋值。

#### b. 更新行为的回调

```go
// 更新时间戳的更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
   // scope.Get() 获取当前设置了标识 gorm:update_column 的字段属性
   if _, ok := scope.Get("gorm:update_column"); !ok {
      // 若不存在，即未自定义设置 update_column
      // 设置默认字段 ModifiedOn 的值为当前时间戳
      _ = scope.SetColumn("ModifiedOn", time.Now().Unix())
   }
}
```

- 通过调用 `scope.Get("gorm:update_column")` 去获取当前设置了标识 `gorm:update_column` 的字段属性。
- 若不存在，也就是没有自定义设置 `update_column`，那么将会在更新回调内设置默认字段` ModifiedOn `的值为当前的时间戳。

#### c. 删除行为的回调

```go
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
            // 获取表明
            scope.QuotedTableName(),
            addExtraSpaceIfExist(scope.CombinedConditionSql()),
            addExtraSpaceIfExist(extraOption),
         )).Exec()
      }
   }
}
```

```go
// CombinedConditionSql return combined condition sql
func (scope *Scope) CombinedConditionSql() string {
   joinSQL := scope.joinsSQL()
   whereSQL := scope.whereSQL()
   if scope.Search.raw {
      whereSQL = strings.TrimSuffix(strings.TrimPrefix(whereSQL, "WHERE ("), ")")
   }
   return joinSQL + whereSQL + scope.groupSQL() +
      scope.havingSQL() + scope.orderSQL() + scope.limitAndOffsetSQL()
}
```

- 通过调用 `scope.Get("gorm:delete_option")` 去获取当前设置了标识 `gorm:delete_option` 的字段属性。
- 判断是否存在 `DeletedOn` 和 `IsDel` 字段，若存在则调整为执行 UPDATE 操作进行软删除（修改` DeletedOn` 和` IsDel` 的值），否则执行 DELETE 进行硬删除。
- 调用 `scope.QuotedTableName` 方法获取当前所引用的表名，并调用一系列方法针对 SQL 语句的组成部分进行处理和转移，最后在完成一些所需参数设置后调用 `scope.CombinedConditionSql` 方法完成 组合条件 SQL 语句的组装。

#### d. 注册回调行为

```go
package model

import (
	"demo/ch02/global"
	"demo/ch02/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	// 引入 MYSQL驱动库进行初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	return db, nil
}

// 回调注册方法实现
func updateTimeStampForCreateCallback(scope *gorm.Scope) {...}
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {...}
func deleteCallback(scope *gorm.Scope) {...}
func addExtraSpaceIfExist(str string) string {...}
```

在最后回到 `NewDBEngine` 方法中，针对上述写的三个 Callback 方法进行回调注册，才能够让应用程序真正的使用上，至此，公共字段处理就完成了。

### 3. 新建 dao 方法

在项目的 `internal/dao` 目录下新建 `dao.go` 文件。

##### dao.go

```go
package dao

import "github.com/jinzhu/gorm"

type Dao struct {
   engine *gorm.DB
}
func New(engine *gorm.DB) *Dao {
   return &Dao{engine: engine}
}
```

接下来在同层级下新建 `tag.go` 文件，用于处理标签模块的 dao 操作。

##### tag.go

```go
package dao

import (
   "demo/ch02/internal/model"
   "demo/ch02/pkg/app"
)

func (d *Dao) CountTag(name string, state uint8) (int, error) {
   tag := model.Tag{Name: name, State: state}
   return tag.Count(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
   tag := model.Tag{Name: name, State: state}
   pageOffset := app.GetPageOffset(page, pageSize)
   return tag.List(d.engine, pageOffset, pageSize)
}

func (d *Dao) CreateTag(name string, state uint8, createdBy string) error {
   tag := model.Tag{
      Name:  name,
      State: state,
      Model: &model.Model{CreatedBy: createdBy},
   }
   return tag.Create(d.engine)
}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
   tag := model.Tag{
      Name:  name,
      State: state,
      Model: &model.Model{ID: id, ModifiedBy: modifiedBy},
   }
   return tag.Update(d.engine)
}

func (d *Dao) DeleteTag(id uint32) error {
   tag := model.Tag{Model: &model.Model{ID: id}}
   return tag.Delete(d.engine)
}
```

在 dao 层进行了数据访问对象的封装，并对针对业务所需字段进行了处理。

### 4. 新建 service 方法

在项目的 `internal/service` 目录下新建 `service.go` 文件。

##### service.go

```go
package service

import (
   "context"
   "demo/ch02/global"
   "demo/ch02/internal/dao"
)

type Service struct {
   ctx context.Context
   dao *dao.Dao
}
func New(ctx context.Context) Service {
   svc := Service{ctx: ctx}
   svc.dao = dao.New(global.DBEngine)
   return svc
}
```

修改同层级下的`tag.go`，用于处理标签模块的业务逻辑。

##### tag.go

```go
package service

import (
   "demo/ch02/internal/model"
   "demo/ch02/pkg/app"
)

// 设置方法的请求结构体和参数校验规则
type CountTagRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}
type TagListRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type CreateTagRequest struct {
   Name      string `form:"name" binding:"required,min=3,max=100"`
   CreatedBy string `form:"created_by" binding:"required,min=3,max=100"`
   State     uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type UpdateTagRequest struct {
   ID         uint32 `form:"id" binding:"required,gte=1"`
   Name       string `form:"name" binding:"min=3,max=100"`
   State      uint8  `form:"state" binding:"required,oneof=0 1"`
   ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}
type DeleteTagRequest struct {
   ID uint32 `form:"id" binding:"required,gte=1"`
}

func (svc *Service) CountTag(param *CountTagRequest) (int, error) {
   return svc.dao.CountTag(param.Name, param.State)
}
func (svc *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
   return svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}
func (svc *Service) CreateTag(param *CreateTagRequest) error {
   return svc.dao.CreateTag(param.Name, param.State, param.CreatedBy)
}
func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
   return svc.dao.UgopdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}
func (svc *Service) DeleteTag(param *DeleteTagRequest) error {
   return svc.dao.DeleteTag(param.ID)
}
```

在上述代码中，主要是定义了 Request 结构体作为接口入参的基准，而本项目由于并不会太复杂，所以直接放在了 service 层中便于使用，若后续业务不断增长，程序越来越复杂，service 也冗杂了，可以考虑将抽离一层接口校验层，便于解耦逻辑。

另外还在 service 中进行了一些简单的逻辑封装，在应用分层中，service 层主要是针对业务逻辑的封装，如果有一些业务聚合和处理可以在该层进行编码，同时也能较好的隔离上下两层的逻辑。

### 5. 新增业务错误码

在项目的 `pkg/errcode` 下新建 `module_code.go` 文件，针对标签模块，写入错误代码。

```go
package errcode

var (
   ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
   ErrorCreateTagFail  = NewError(20010002, "创建标签失败")
   ErrorUpdateTagFail  = NewError(20010003, "更新标签失败")
   ErrorDeleteTagFail  = NewError(20010004, "删除标签失败")
   ErrorCountTagFail   = NewError(20010005, "统计标签失败")
)
```

### 6. 新增路由方法

修改 `internal/routers/api/v1` 项目目录下的` tag.go `文件

```go
func (t Tag) List(c *gin.Context) {
   // 设置入参格式与参数校验规则
   param := service.TagListRequest{}
   // 初始化响应
   response := app.NewResponse(c)
   // 进行入参校验
   valid, errs := app.BindAndValid(c, &param)
   if !valid {
      global.Logger.Errorf("app.BindAndValid errs: %v", errs)
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
      return
   }

   svc := service.New(c.Request.Context())
   pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
   // 获取标签总数
   totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
   if err != nil {
      global.Logger.Errorf("svc.CountTag err: %v", err)
      response.ToErrorResponse(errcode.ErrorCountTagFail)
      return
   }
   // 获取标签列表
   tags, err := svc.GetTagList(&param, &pager)
   if err != nil {
      global.Logger.Errorf("svc.GetTagList err: %v", err)
      response.ToErrorResponse(errcode.ErrorGetTagListFail)
      return
   }
   
   // 序列化结果集
   response.ToResponseList(tags, totalRows)
   return
}
```

在上述代码中，完成了获取标签列表接口的处理方法，在方法中完成了入参校验和绑定、获取标签总数、获取标签列表、 序列化结果集等四大功能板块的逻辑串联和日志、错误处理。需要注意的是方法实现中的入参校验和绑定的处理代码基本都差不多，`tag.go`全部代码如下：

```go
package v1

import (
	"demo/ch02/global"
	"demo/ch02/internal/service"
	"demo/ch02/pkg/app"
	"demo/ch02/pkg/convert"
	"demo/ch02/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Tag struct {}
func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context) {}

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	// 设置入参格式与参数校验规则
	param := service.TagListRequest{}
	// 初始化响应
	response := app.NewResponse(c)
	// 进行入参校验
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	// 获取标签总数
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.Errorf("svc.CountTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}
	// 获取标签列表
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errorf("svc.GetTagList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	// 序列化结果集
	response.ToResponseList(tags, totalRows)
	return
}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string true "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	// 入参校验与绑定
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.CreateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	param := service.UpdateTagRequest{
		// 将 string 类型转换为 uint32
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.UpdateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	param := service.DeleteTagRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.DeleteTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
```

### 7. 验证接口

启动服务，对标签模块的接口进行验证，请注意，验证示例中的 `{id}`，代指占位符，也就是填写实际调用中希望处理的标签 ID 即可。

#### a. 新增标签

使用 postman 进行接口测试，使用 `json` 作为入参无法成功创建标签，会报错入参错误，原因不明。

![image-20220514181031590](https://raw.githubusercontent.com/tonshz/test/master/img/202205141810834.png)

若不满足参数校验规则则会报错，例如设置 name 值为 Go。

```json
{
    "code": 10000001,
    "details": [
        "Name长度必须至少为3个字符"
    ],
    "msg": "入参错误"
}
```

#### b. 获取标签列表

![image-20220514181652124](https://raw.githubusercontent.com/tonshz/test/master/img/202205141816180.png)

```json
{
    "list": [
        {
            "id": 1,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652522978,
            "modified_on": 1652522978,
            "deleted_on": 0,
            "is_del": 0,
            "name": "create_tag_test",
            "state": 1
        },
        {
            "id": 2,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652523217,
            "modified_on": 1652523217,
            "deleted_on": 0,
            "is_del": 0,
            "name": "Java",
            "state": 1
        }
    ],
    "pager": {
        "page": 1, 
        "page_size": 2,
        "total_rows": 3
    }
}
```

修改 page 参数为 2。

```json
{
    "list": [
        {
            "id": 3,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652523235,
            "modified_on": 1652523235,
            "deleted_on": 0,
            "is_del": 0,
            "name": "Golang",
            "state": 1
        }
    ],
    "pager": {
        "page": 2,
        "page_size": 2,
        "total_rows": 3
    }
}
```

#### c. 修改标签

此处 postman 失败，不清楚原因，与参数校验相关，始终报错 `"Name长度必须至少为3个字符"`。

```bash
$ curl -X PUT http://127.0.0.1:8000/api/v1/tags/{id} -F state=0 -F modified_by=eddycjy
{}
```

#### d. 删除标签

![image-20220514185320956](https://raw.githubusercontent.com/tonshz/test/master/img/202205141853990.png)

删除标签后，数据库内容更新。

![image-20220514185303772](https://raw.githubusercontent.com/tonshz/test/master/img/202205141853815.png)

### 8. 发现问题：零值未更新

在完成了接口的检验后，还需要确定一下数据库内的数据变更是否正确。在经过一系列的对比后，发现在调用修改标签的接口时，通过接口入参，是希望将 id 为 1 的标签状态修改为 0，但是在对比后发现数据库内它的状态值仍然是 1，而且 SQL 语句内也没有出现 state 字段的设置，控制台输出的 SQL 语句如下：

```bash
UPDATE `blog_tag` SET `id` = 1, `modified_by` = 'eddycjy', `modified_on` = xxxxx  WHERE `blog_tag`.`id` = 1
```

**原因是只要字段是零值的情况下，GORM 就不会对该字段进行变更。**实际上，这有一个概念上的问题，先入为主的认为它一定会变更，其实是不对的，因为在程序中使用的是 struct 的方式进行更新操作，而在 GORM 中使用 struct 类型传入进行更新时，**GORM 是不会对值为零值的字段进行变更**。更根本的原因是因为在识别这个结构体中的这个字段值时，**很难判定是真的是零值，还是外部传入恰好是该类型的零值**，GORM 在这块并没有过多的去做特殊识别。

### 9. 解决问题

修改项目的 `internal/model` 目录下的` tag.go `文件里的` Update `方法。

```go
// 修改传入零值后，数据库未发生变化的问题
func (t Tag) Update(db *gorm.DB, values interface{}) error {
   // Update 更新所选字段
   if err := db.Model(t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error; err != nil {
      return err
   }
   return nil
}
```

修改项目的 `internal/dao` 目录下的`tag.go `文件里的 `UpdateTag `方法。

```go
func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
   tag := model.Tag{
      Model: &model.Model{ID: id},
   }
   values := map[string]interface{}{
      "state":       state,
      "modified_by": modifiedBy,
   }
   if name != "" {
      values["name"] = name
   }
   return tag.Update(d.engine, values)
}
```

重新运行程序，请求修改标签接口，检查数据是否正常修改，在正确的情况下，该 id 为 1 的标签，modified_by 为 test，modified_on 应修改为当前时间戳，state 为 0。

### 10. 文章管理模块

参见 Github 代码。

-------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



