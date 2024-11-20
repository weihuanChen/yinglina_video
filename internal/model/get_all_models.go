package model

// GetAllModels 获取并注册所有数据库模型
func GetAllModels() []interface{} {
	return []interface{}{
		&Account{},
	}
}
