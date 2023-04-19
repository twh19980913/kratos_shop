package data

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type GormList []string

//Value 当数据传过来我们如何拿给数据库
func (g GormList) Value() (driver.Value,error) {
	return json.Marshal(g)
}
//Scan 将数据查询出来后怎么处理
func (g *GormList)Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte),&g)
}

type BaseModel struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}

func Paginate(page,pageSize int) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB {
		if page == 0{
			page = 1
		}
		switch{
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}