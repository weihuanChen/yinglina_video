package base

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

// BasedModel 定义
type BasedModel struct {
	ID          int64   `gorm:"primaryKey" json:"id"`
	GmtCreate   int64   `gorm:"type:bigint" json:"gmt_create"`
	GmtModified int64   `gorm:"type:bigint" json:"gmt_modified"`
	Ext         JSONMap `gorm:"type:json" json:"ext"`
	Deleted     bool    `gorm:"type:boolean" json:"deleted"`
}

// JSONMap 用于处理 json 类型字段
type JSONMap map[string]interface{}

// Scan 用于从数据库读取 json 数据
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// Value 用于将 map 转换为 json 数据存储到数据库中
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// BeforeCreate hook
func (m *BasedModel) BeforeCreate(db *gorm.DB) (err error) {
	currentTime := time.Now().Unix() // 当前时间戳（秒级）
	m.GmtCreate = currentTime
	m.GmtModified = currentTime
	m.Deleted = false // 默认值为 false (0)

	if m.Ext == nil {
		m.Ext = make(map[string]interface{})
	}
	return nil
}

// BeforeUpdate hook
func (m *BasedModel) BeforeUpdate(db *gorm.DB) (err error) {
	m.GmtModified = time.Now().Unix()
	return nil
}
