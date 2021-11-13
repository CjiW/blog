package conf

import (
	"gorm.io/gorm"
)

var (
	DATABASE = "root:020103@tcp(127.0.0.1:3306)/blog?parseTime=true&loc=Local&charset=utf8mb4"
	GlobalDB gorm.DB
)
