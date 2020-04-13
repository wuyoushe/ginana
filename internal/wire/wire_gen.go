// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server"
	"ginana/internal/server/controller/api"
	"ginana/internal/server/router"
	"ginana/internal/service"
	"ginana/internal/service/i_user"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() (*App, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := db.NewDB(configConfig)
	if err != nil {
		return nil, nil, err
	}
	iUser, err := i_user.New(gormDB, configConfig)
	if err != nil {
		return nil, nil, err
	}
	syncedEnforcer, err := db.NewCasbin(iUser, configConfig)
	if err != nil {
		return nil, nil, err
	}
	memcache, err := db.NewMC(configConfig)
	if err != nil {
		return nil, nil, err
	}
	serviceService, err := service.New(gormDB, syncedEnforcer, memcache, iUser)
	if err != nil {
		return nil, nil, err
	}
	cApi := api.New(serviceService)
	application := router.InitRouter(cApi, configConfig)
	httpServer, err := server.NewHttpServer(application, configConfig)
	if err != nil {
		return nil, nil, err
	}
	app, cleanup, err := NewApp(serviceService, httpServer)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewMC, db.NewCasbin)

var iProvider = wire.NewSet(i_user.New)

var cProvider = wire.NewSet(api.New)

var httpProvider = wire.NewSet(router.InitRouter, server.NewHttpServer)
