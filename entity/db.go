package entity

import "gorm.io/gorm"

// DB instance. Placed here so that it is globally accessible across the module
var DB *gorm.DB
