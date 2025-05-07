package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"miaosha-system/config"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Log    *logrus.Logger
	Redis  *redis.Client
)
