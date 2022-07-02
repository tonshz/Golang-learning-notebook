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

// 修改传入零值后，数据库未发生变化的问题
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	// Update 更新所选字段r
	return db.Model(&t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error
}

func (t Tag) Delete(db *gorm.DB) error {
	// Delete 删除数据
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}

// 文章管理新增代码
func (t Tag) Get(db *gorm.DB) (*Tag, error) {
	var tag *Tag
	err := db.Where("id = ? and is_del = ? and state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}
	return tag, nil
}

func (t Tag) ListByIDs(db *gorm.DB, ids []uint32) ([]*Tag, error) {
	var tags []*Tag
	db = db.Where("state = ? AND is_del = ?", t.State, 0)
	err := db.Where("id IN (?)", ids).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}
