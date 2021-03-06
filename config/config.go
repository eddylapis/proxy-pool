// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"os"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

var (
	box = packr.New("config", "../configs")
	env = os.Getenv("GO_ENV")
)

const (
	// Dev development env
	Dev = "dev"
	// Test test env
	Test = "test"
	// Production production env
	Production = "production"
)

type (
	Crawler struct {
		Name     string
		Interval time.Duration
		MaxPage  int
	}
	// Detect detect config
	Detect struct {
		URL      string
		Interval time.Duration
		Timeout  time.Duration
		MaxTimes int
	}
)

func init() {
	configType := "yml"
	configExt := "." + configType
	data, err := box.Find("default" + configExt)
	if err != nil {
		panic(err)
	}
	viper.SetConfigType(configType)
	v := viper.New()
	v.SetConfigType(configType)
	// 读取默认配置中的所有配置
	err = v.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		viper.SetDefault(k, v)
	}

	// 根据当前运行环境配置读取
	// 可根据不同的环境仅调整与default不一致的相关配置
	if env != "" {
		envConfigFile := env + configExt
		data, err = box.Find(envConfigFile)
		if err != nil {
			panic(err)
		}
		// 读取当前运行环境对应的配置
		err = viper.ReadConfig(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
	}
}

// GetCrawlers get crawlers config
func GetCrawlers() []*Crawler {
	crawlers := make([]*Crawler, 0)
	data := viper.GetStringSlice("crawler")
	for _, name := range data {
		interval := viper.GetDuration(name + ".interval")
		maxPage := viper.GetInt(name + ".maxPage")
		// 如果未配置抓取间隔时间，则设置为2分钟
		if interval == 0 {
			interval = 2 * time.Minute
		}
		crawlers = append(crawlers, &Crawler{
			Name:     name,
			Interval: interval,
			MaxPage:  maxPage,
		})
	}
	return crawlers
}

// GetDetect get detect config
func GetDetect() *Detect {
	prefix := "detect."
	conf := &Detect{
		Timeout:  viper.GetDuration(prefix + "timeout"),
		URL:      viper.GetString(prefix + "url"),
		Interval: viper.GetDuration(prefix + "interval"),
		MaxTimes: viper.GetInt(prefix + "maxTimes"),
	}
	if conf.Timeout == 0 {
		conf.Timeout = 3 * time.Second
	}
	if conf.Interval == 0 {
		conf.Interval = 30 * time.Minute
	}
	if conf.URL == "" {
		conf.URL = "https://www.baidu.com/"
	}
	if conf.MaxTimes <= 0 {
		conf.MaxTimes = 3
	}
	return conf
}

// GetListenAddr get listen address
func GetListenAddr() string {
	addr := viper.GetString("listen")
	if addr == "" {
		return ":4000"
	}
	return addr
}
