package config

import (
	"fmt"
	"log"
	"time"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func InitDB(isProduction bool) *gorm.DB {
	// Read and Write Connection
	dsnRW := fmt.Sprintf(constant.STR_DSN, constant.DB_USER, constant.DB_PASS, constant.RW_HOST, constant.DB_PORT, constant.DB_NAME)
	db, err := gorm.Open(mysql.Open(dsnRW), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)

	// Read Only Connection
	dsnRO := fmt.Sprintf(constant.STR_DSN, constant.DB_USER, constant.DB_PASS, constant.RO_HOST, constant.DB_PORT, constant.DB_NAME)

	db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas:          []gorm.Dialector{mysql.Open(dsnRO)},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).SetMaxOpenConns(1000).SetMaxIdleConns(200),
	)

	if constant.DB_CACHING {
		var redisClient *redis.Client

		if isProduction {
			redisOptions := &redis.Options{
				Addr:     fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
				Password: constant.REDIS_PASS,
			}

			redisClient = redis.NewClient(redisOptions)
		} else {
			redisOptions := &redis.Options{
				Addr: fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
			}

			redisClient = redis.NewClient(redisOptions)
		}

		if _, err := redisClient.Ping().Result(); err != nil {
			if err.Error() == "ERR AUTH <password> called without any password configured for the default user. Are you sure your configuration is correct?" {
				redisOptions := &redis.Options{
					Addr: fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
				}

				redisClient = redis.NewClient(redisOptions)

				if _, err := redisClient.Ping().Result(); err != nil {
					log.Fatal("error connection to redis server, error: ", err.Error())
				}
			} else {
				log.Fatal("error connection to redis server, error: ", err.Error())
			}
		}

		gormCache, _ := cache.NewGorm2Cache(&config.CacheConfig{
			CacheLevel:           config.CacheLevelAll,
			CacheStorage:         config.CacheStorageRedis,
			RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
			InvalidateWhenUpdate: true,
			CacheTTL:             2500, // 2 1/2 second
			CacheMaxItemCnt:      0,    // 0 = cache all queries
		})

		db.Use(gormCache)
	}

	if err := db.Raw(helper.ConvertToInLineQuery("SET GLOBAL FOREIGN_KEY_CHECKS = 0;")).Error; err != nil {
		log.Fatal(err.Error())
	}

	return db
}
