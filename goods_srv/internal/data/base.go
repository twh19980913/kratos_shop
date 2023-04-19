package data

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
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

type BaseModel struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}
