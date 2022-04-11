package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pavelerokhin/user-microservice-go/controller"
	"github.com/pavelerokhin/user-microservice-go/repository"
	"github.com/pavelerokhin/user-microservice-go/router"
	"github.com/pavelerokhin/user-microservice-go/service"
)

var (
	userRouter     router.Router
	userRepository repository.UserRepository
	userService    service.UserService
	userController controller.UserController
)

func main() {
	var err error

	logger := log.New(os.Stdout, "faceit-test-commitment", log.LstdFlags|log.Lshortfile)

	userRepository, err = repository.NewSqliteRepo(logger)
	userService = service.New(userRepository, logger)
	userController = controller.New(userService, logger)

	userRouter = router.NewMuxRouter(logger)
	if err != nil {
		logger.Fatal(err)
	}

	portPtr := *flag.String("port", "8080", "Server port. Default: 8080")
	flag.Parse()

	if portPtr != "" {
		portPtr = fmt.Sprintf(":%s", portPtr)
	}

	userRouter.SERVE(fmt.Sprintf(":%v", portPtr))
}
