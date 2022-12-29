package commands

import (
	"BRGS/conf"
	"BRGS/management"
	"BRGS/routers"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ExitCommand struct {
	*management.ShareData
}

func (e *ExitCommand) Execute() bool {
	os.Exit(0)
	return true
}

func (e *ExitCommand) String() string {
	return conf.CommandNames.ExitCommand
}

// 启动web端
type StartServerCommand struct {
	*management.ShareData
}

func (e *StartServerCommand) Execute() bool {
	server := &http.Server{
		Addr:    ":" + conf.ServerConf.Port,
		Handler: routers.InitRouter(),
	}
	fmt.Println(conf.ServerConf.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// log.Println("shutdown server")
	// 	ctx, chanel :=context.WithTiout(context.Background(), 5*time.Second)
	// 	defer chanel()
	// 	if err := serv.Shutdown(ctx); err =nil {
	// 		g.Fatal("Server Shutdown:", err)
	// 	}
	// 	t.Println("Server exiting")
	return true
}

func (e *StartServerCommand) String() string {
	return conf.CommandNames.StartServerCommand
}

// 停止web端
type StopServerCommand struct {
	*management.ShareData
}

func (e *StopServerCommand) Execute() bool {

	log.Println("Server exiting")

	return true
}

func (e *StopServerCommand) String() string {
	return conf.CommandNames.StopServerCommand
}

func CloseServer(server *http.Server) {
	log.Println("shutdown server")
	// ctx, chanel := context.WithTimeout(coext.Background(), 5*time.Secod)
	// defer chane()
	// if err := serv.Shutdown(ctx); err != ni {
	// 	g.Fatal("Server Shutdown:", er)
	// }
	// t.Println("Server exitin")
}
