/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    config
 * @Date:    2022/3/14 11:26 上午
 * @package: mgoattr
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package mgoattr

import "fmt"

var dbConfig *Option

type Option struct {
	Addr     string
	DB       string
	User     string
	Password string
}

func (o *Option) format() string {
	return fmt.Sprintf("Addr=%s DB=%s User=%s", o.Addr, o.DB, o.User)
}

func defaultOption() *Option {
	return &Option{
		Addr: "127.0.0.1:27017",
		DB:   "Attribute",
	}
}

func initDBConfig(opfs ...func(*Option)) {
	dbConfig = defaultOption()

	for _, opf := range opfs {
		opf(dbConfig)
	}
}
