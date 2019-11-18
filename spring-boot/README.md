# spring-boot

开箱即用的 Go-Spring 程序启动框架。

```
import (
	"fmt"
	"testing"
	"time"

	"github.com/go-spring/go-spring/boot-starter"
	"github.com/go-spring/go-spring/spring-boot"
	"github.com/go-spring/go-spring/spring-core"
)

func init() {
	SpringBoot.RegisterModule(func(ctx SpringCore.SpringContext) {
		ctx.RegisterBean(new(MyModule))
	})
}

type MyModule struct {
}

func (m *MyModule) OnStartApplication(ctx SpringBoot.ApplicationContext) {
	fmt.Println("MyModule start")

	ctx.SafeGoroutine(func() {

		defer fmt.Println("go stop")
		fmt.Println("go start")

		time.Sleep(200 * time.Millisecond)
		BootStarter.Exit()
	})
}

func (m *MyModule) OnStopApplication(ctx SpringBoot.ApplicationContext) {
	fmt.Println("MyModule stop")
}

func TestRunApplication(t *testing.T) {
	SpringBoot.RunApplication("config/")
}
```