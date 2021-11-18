package apollo

import (
	"fmt"
	"strings"

	gc "github.com/jiurongzhao/bootstrap-global/config"
	"github.com/jiurongzhao/bootstrap-global/util"

	av4 "github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

type ApolloConfig struct{}

type ApolloContainer struct {
	client    *av4.Client
	namespace string
}

type Apollo struct {
	AppId          string `name:"appId"`
	IP             string `name:"ip"`
	Cluster        string `name:"cluster" default:"default"`
	NamespaceName  string `name:"namespaceName"`
	Secret         string `name:"secret"`
	IsBackupConfig bool   `name:"isBackupConfig" default:"true"`
}

func (c *ApolloConfig) Load(filename string) (gc.Configer, error) {
	appConfig := Apollo{}
	if err := gc.Resolve("apollo", &appConfig); err != nil {
		return nil, err
	}
	appConfig2 := config.AppConfig{
		AppID:          appConfig.AppId,
		Cluster:        appConfig.Cluster,
		IP:             appConfig.IP,
		NamespaceName:  appConfig.NamespaceName,
		Secret:         appConfig.Secret,
		IsBackupConfig: appConfig.IsBackupConfig,
	}
	client, err := av4.StartWithConfig(func() (*config.AppConfig, error) {
		return &appConfig2, nil
	})
	if err != nil {
		return nil, err
	}

	return &ApolloContainer{
		client:    client,
		namespace: appConfig.NamespaceName,
	}, nil
}

func (c *ApolloContainer) Get(key string) (interface{}, bool) {
	if value, err := c.client.GetConfigCache(c.namespace).Get(key); err != nil {
		return nil, false
	} else {
		return value, true
	}
}

func (c *ApolloContainer) Resolve(prefix string, p interface{}) error {
	dict := make(map[string]interface{})
	cache := c.client.GetConfigCache(c.namespace)
	if cache == nil {
		return fmt.Errorf("not found cache with %v", c.namespace)
	}
	prefixIsEmpty := prefix == ""
	cache.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if !ok {
			return true
		}
		if prefixIsEmpty || strings.HasPrefix(keyStr, prefix) {
			dict[keyStr] = value
		}
		return true
	})
	return util.ResolveStruct(&dict, prefix, p)
}

func init() {
	gc.Register("apollo", &ApolloConfig{})
}
