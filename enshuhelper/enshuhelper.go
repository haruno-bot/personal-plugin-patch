package enshuhelper

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/haruno-bot/haruno/coolq"
	"github.com/haruno-bot/haruno/logger"
)

// EnshuHelper 转推插件
type EnshuHelper struct {
	coolq.Plugin
	name    string
	version string
	module  string
}

// Name 插件名称
func (_plugin *EnshuHelper) Name() string {
	return fmt.Sprintf("%s@%s", _plugin.name, _plugin.version)
}

// Module 插件名称
func (_plugin *EnshuHelper) Module() string {
	return _plugin.module
}

func (_plugin *EnshuHelper) loadConfig() error {
	cfg := new(Config)
	_, err := toml.DecodeFile("config.toml", cfg)
	if err != nil {
		return err
	}
	pcfg := cfg.EnshuHelper
	_plugin.name = pcfg.Name
	_plugin.module = pcfg.Module
	_plugin.version = pcfg.Version
	return nil
}

// Callback 回调接口
func (_plugin *EnshuHelper) Callback(data string) {
	logger.Field("Plugin EnshuHelper").Successf("%s", data)
}

// Load 插件加载
func (_plugin *EnshuHelper) Load() error {
	err := _plugin.loadConfig()
	if err != nil {
		return err
	}
	return nil
}

// Instance 转推插件实体
var Instance = &EnshuHelper{}
