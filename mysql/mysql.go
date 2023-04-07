/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    mysql
 * @Date:    2022/7/26 23:35
 * @package: mysql
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package mysql

import (
	"errors"
	"fmt"
	"github.com/jageros/hawox/flags"
	"github.com/jageros/hawox/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var (
	_DBMap       = map[string]*gorm.DB{}
	_RwMx        sync.RWMutex
	NotInitDBErr = errors.New("db not init")
)

type OpFun func(opt *Option)

type Option struct {
	Addr             string
	User             string
	Password         string
	Database         string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnMaxLifetime  time.Duration
	LogSlowThreshold time.Duration
}

func defaultOption() *Option {
	return &Option{
		Addr:             "127.0.0.1:3306",
		User:             "root",
		Password:         "123456",
		Database:         "db_test",
		MaxIdleConns:     10,
		MaxOpenConns:     256,
		ConnMaxLifetime:  time.Hour,
		LogSlowThreshold: time.Millisecond * 500,
	}
}

func Conn(opfs ...OpFun) (*gorm.DB, error) {
	opt := defaultOption()
	for _, opf := range opfs {
		opf(opt)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", opt.User, opt.Password, opt.Addr, opt.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(logx.Logger(), logger.Config{
			SlowThreshold:             opt.LogSlowThreshold,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	sqlDB.SetMaxIdleConns(opt.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opt.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opt.ConnMaxLifetime)
	_RwMx.Lock()
	_DBMap[opt.Database] = db
	_RwMx.Unlock()
	return db, nil
}

func DB(dbName string) (*gorm.DB, error) {
	_RwMx.RLock()
	db, ok := _DBMap[dbName]
	_RwMx.RUnlock()
	if ok {
		return db, nil
	}
	user, pwd, addr, err := _Conf(dbName)
	if err == nil {
		db, err = Conn(func(opt *Option) {
			opt.Database = dbName
			opt.User = user
			opt.Password = pwd
			opt.Addr = addr
		})
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}

func _Conf(dbName string) (user string, pwd string, addr string, err error) {
	user = flags.GetString("mysql." + dbName + ".user")
	pwd = flags.GetString("mysql." + dbName + ".password")
	addr = flags.GetString("mysql." + dbName + ".addr")
	if len(user) <= 0 || len(pwd) <= 0 || len(addr) <= 0 {
		err = errors.New("can not find db config")
	}
	return
}
