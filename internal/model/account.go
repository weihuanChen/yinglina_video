package model

import "yunosphere.com/yun-fiber-scaffold/internal/model/base"

type Account struct {
	base.BasedModel
	UserId   int64  `gorm:"type:bigint;not null;uniqueIndex" json:"user_id"`        // 用户 ID
	Email    string `gorm:"type:varchar(64);default:null;uniqueIndex" json:"email"` // 邮箱
	Password string `gorm:"type:varchar(255);not null" json:"password"`             // 加密密码
	Nickname string `gorm:"type:varchar(255);default:null" json:"nickname"`         // 昵称
	RoleCode string `gorm:"type:varchar(32);not null" json:"role_code"`             // 角色编码
}

func (Account) TableName() string {
	return "account"
}
