//  Vikunja is a todo-list application to facilitate your life.
//  Copyright 2019 Vikunja and contributors. All rights reserved.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/log"
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // Because.
	_ "github.com/mattn/go-sqlite3"    // Because.
)

// CreateDBEngine initializes a db engine from the config
func CreateDBEngine() (engine *xorm.Engine, err error) {
	// If the database type is not set, this likely means we need to initialize the config first
	if config.DatabaseType.GetString() == "" {
		config.InitConfig()
	}

	// Use Mysql if set
	if config.DatabaseType.GetString() == "mysql" {
		engine, err = initMysqlEngine()
		if err != nil {
			return
		}
	} else {
		// Otherwise use sqlite
		engine, err = initSqliteEngine()
		if err != nil {
			return
		}
	}

	engine.SetMapper(core.GonicMapper{})
	engine.ShowSQL(config.LogDatabase.GetString() != "off")
	engine.SetLogger(xorm.NewSimpleLogger(log.GetLogWriter("database")))

	return
}

func initMysqlEngine() (engine *xorm.Engine, err error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		config.DatabaseUser.GetString(),
		config.DatabasePassword.GetString(),
		config.DatabaseHost.GetString(),
		config.DatabaseDatabase.GetString())
	engine, err = xorm.NewEngine("mysql", connStr)
	if err != nil {
		return
	}
	engine.SetMaxOpenConns(config.DatabaseMaxOpenConnections.GetInt())
	engine.SetMaxIdleConns(config.DatabaseMaxIdleConnections.GetInt())
	max, err := time.ParseDuration(strconv.Itoa(config.DatabaseMaxConnectionLifetime.GetInt()) + `ms`)
	if err != nil {
		return
	}
	engine.SetConnMaxLifetime(max)
	return
}

func initSqliteEngine() (engine *xorm.Engine, err error) {
	path := config.DatabasePath.GetString()
	if path == "" {
		path = "./db.db"
	}

	return xorm.NewEngine("sqlite3", path)
}