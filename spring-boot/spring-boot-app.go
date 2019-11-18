/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//
// 开箱即用的 Go-Spring 程序启动框架。
//
package SpringBoot

import (
	"fmt"
	"os"
)

//
// 定义 SpringBoot 应用。
//
type Application struct {
	AppContext     ApplicationContext // 应用上下文
	ConfigLocation string             // 配置文件目录
}

//
// 应用运行过程中的事件。
//
type ApplicationEvent interface {
	OnStartApplication(ctx ApplicationContext) // 应用启动的事件
	OnStopApplication(ctx ApplicationContext)  // 应用停止的事件
}

//
// 工厂函数
//
func NewApplication(configLocation string) *Application {
	return &Application{
		AppContext:     NewDefaultApplicationContext(),
		ConfigLocation: configLocation,
	}
}

//
// BootStarter.AppRunner.$Start
//
func (app *Application) Start() {

	// 加载配置文件
	app.loadConfigFiles()

	// 注册 ApplicationContext Bean 对象
	app.AppContext.RegisterBean(app.AppContext)

	// 初始化所有的 SpringBoot 模块
	for _, fn := range modules {
		fn(app.AppContext)
	}

	// 依赖注入、属性绑定、Bean 初始化
	app.AppContext.AutoWireBeans()

	// 通知应用启动事件
	var eventBeans []ApplicationEvent
	app.AppContext.CollectBeans(&eventBeans)

	if eventBeans != nil && len(eventBeans) > 0 {
		for _, bean := range eventBeans {
			bean.OnStartApplication(app.AppContext)
		}
	}

	fmt.Println("spring boot started")
}

func (app *Application) loadConfigFiles() {

	// 加载默认的应用配置文件，如 application.properties
	filePath := app.ConfigLocation + "application"
	app.loadConfigFiles0(filePath)

	// 加载用户设置的配置文件，如 application-test.properties
	if env := os.Getenv("spring.profile"); len(env) > 0 {
		filePath = fmt.Sprintf(app.ConfigLocation+"application-%s", env)
		app.loadConfigFiles0(filePath)
	}
}

func (app *Application) loadConfigFiles0(filePath string) {
	for _, ext := range []string{".properties", ".yaml", ".toml"} {
		app.AppContext.LoadProperties(filePath + ext)
	}
}

//
// BootStarter.AppRunner.$ShutDown
//
func (app *Application) ShutDown() {

	// 通知应用停止事件
	var eventBeans []ApplicationEvent
	app.AppContext.CollectBeans(&eventBeans)

	if eventBeans != nil && len(eventBeans) > 0 {
		for _, bean := range eventBeans {
			bean.OnStopApplication(app.AppContext)
		}
	}

	// 等待所有 goroutine 退出
	app.AppContext.Wait()

	fmt.Println("spring boot exit")
}
