/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    mysql
 * @Date:    2021/8/24 6:56 下午
 * @package: mysql
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package mysql

import (
	"fmt"
	"github.com/jageros/hawox/redis"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

var opt_ = &Option{
	Addr:     "127.0.0.1:3306",
	Username: "root",
	Password: "QianYin@66",
	Database: "loverstore",
}
var conn = sqlx.NewMysql(opt_.DataSource())

type Option struct {
	Addr     string
	Username string
	Password string
	Database string
}

func (op Option) DataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", op.Username, op.Password, op.Addr, op.Database)
}

func Initialize(opfs ...func(opt *Option)) {
	for _, opf := range opfs {
		opf(opt_)
	}
	conn = sqlx.NewMysql(opt_.DataSource())
}

func Conn() sqlx.SqlConn {
	if conn == nil {
		conn = sqlx.NewMysql(opt_.DataSource())
	}
	return conn
}

func ConnWithCache() (sqlx.SqlConn, cache.CacheConf) {
	return Conn(), redis.CacheConf()
}
