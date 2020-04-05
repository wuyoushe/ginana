// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"ginana/internal/config"
	"ginana/internal/db"
	"ginana/internal/server/http"
	"ginana/internal/server/http/h_user"
	"ginana/internal/server/http/router"
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
	serviceService, err := service.New(gormDB, syncedEnforcer, iUser)
	if err != nil {
		return nil, nil, err
	}
	hUser := h_user.New(serviceService)
	engine := router.InitRouter(hUser)
	server, err := http.NewHttpServer(engine, configConfig)
	if err != nil {
		return nil, nil, err
	}
	app, cleanup, err := NewApp(serviceService, server)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewCasbin)

var iProvider = wire.NewSet(i_user.New)

var hProvider = wire.NewSet(h_user.New)

var httpProvider = wire.NewSet(router.InitRouter, http.NewHttpServer)
