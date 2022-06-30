package db

import (
	"fmt"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Charset   string
	Collation string
	Host      string
	Name      string
	Password  string
	Port      string
	Username  string
}

var Conn *gorm.DB

func Setup(config *Config) error {
	var err error
	if Conn, err = gorm.Open(mysql.Open(config.mysqlDSN()), &gorm.Config{}); err != nil {
		return err
	}
	return nil
}

func (cf *Config) mysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&collation=%s&parseTime=true&loc=Local",
		url.QueryEscape(cf.Username),
		url.QueryEscape(cf.Password),
		url.QueryEscape(cf.Host),
		url.QueryEscape(cf.Port),
		url.QueryEscape(cf.Name),
		url.QueryEscape(cf.Charset),
		url.QueryEscape(cf.Collation),
	)
}
