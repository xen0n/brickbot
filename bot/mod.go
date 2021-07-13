// SPDX-License-Identifier: GPL-3.0-or-later

package bot

import (
	"errors"
	"plugin"
	"reflect"

	"github.com/BurntSushi/toml"

	"github.com/xen0n/brickbot/bot/v1alpha1"
)

type LoadedPlugin struct {
	configFactoryFn v1alpha1.IPluginConfigFactoryFunc
	factoryFn       v1alpha1.IPluginFactoryFunc
}

func LoadPlugin(pluginPath string) (*LoadedPlugin, error) {
	pl, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	apiVersionSym, err := pl.Lookup("BrickbotPluginAPIVersion")
	if err != nil {
		return nil, err
	}

	apiVersion, ok := apiVersionSym.(*int)
	if !ok {
		return nil, errors.New("wrong type of BrickbotPluginAPIVersion symbol")
	}

	if apiVersion == nil || *apiVersion != v1alpha1.PluginAPIVersion {
		return nil, errors.New("plugin API version mismatch")
	}

	pluginConfigFactoryFnSym, err := pl.Lookup("BrickbotPluginConfigFactory")
	if err != nil {
		return nil, err
	}

	pluginConfigFactoryFn, ok := pluginConfigFactoryFnSym.(v1alpha1.IPluginConfigFactoryFunc)
	if !ok {
		return nil, errors.New("wrong type of BrickbotPluginConfigFactory symbol")
	}

	pluginFactoryFnSym, err := pl.Lookup("BrickbotPluginFactory")
	if err != nil {
		return nil, err
	}

	pluginFactoryFn, ok := pluginFactoryFnSym.(v1alpha1.IPluginFactoryFunc)
	if !ok {
		return nil, errors.New("wrong type of BrickbotPluginFactory symbol")
	}

	return &LoadedPlugin{
		configFactoryFn: pluginConfigFactoryFn,
		factoryFn:       pluginFactoryFn,
	}, nil
}

func (p *LoadedPlugin) InitWithConfigTOML(configPath string) (v1alpha1.IPlugin, error) {
	// get concrete type for unmarshaling
	configTypeTemplate := p.configFactoryFn()
	rv := reflect.New(reflect.TypeOf(configTypeTemplate))

	_, err := toml.DecodeFile(configPath, rv.Interface())
	if err != nil {
		return nil, err
	}

	return p.factoryFn(rv.Elem().Interface())
}
