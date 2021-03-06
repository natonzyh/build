package postgres

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改界面
func (this *UpdateAction) RunGet(params struct {
	Action string
}) {
	this.Data["action"] = params.Action

	// 读取配置
	config, err := db.LoadPostgresConfig()
	if err != nil {
		config = db.DefaultPostgresConfig()
	}

	this.Data["config"] = config
	this.Data["typeIsChanged"] = db.SharedDBConfig().Type != db.DBTypePostgres

	this.Show()
}

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	Addr     string
	Username string
	Password string
	DBName   string `alias:"dbName"`
	PoolSize int
	Must     *actions.Must
}) {
	// 是否已改变
	sharedConfig := db.SharedDBConfig()
	isChanged := sharedConfig.Type != db.DBTypePostgres

	config, err := db.LoadPostgresConfig()
	if err != nil {
		config = db.DefaultPostgresConfig()
	}
	config.Addr = params.Addr
	config.Username = params.Username
	config.Password = params.Password
	config.DBName = params.DBName
	config.PoolSize = params.PoolSize
	config.DSN = config.ComposeDSN()
	err = config.Save()
	if err != nil {
		this.Fail(err.Error())
	}

	if isChanged {
		shouldRestart = true
		sharedConfig.Type = db.DBTypePostgres
		err = sharedConfig.Save()
		if err != nil {
			this.Fail(err.Error())
		}
	}

	teadb.ChangeDB()

	this.Success()
}
