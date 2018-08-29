package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/controller"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/server"
	"github.com/golang/glog"
)

var options struct {
	Port int
}

func init() {
	flag.IntVar(&options.Port, "port", 8005, "use '--port' option to specify the port for broker to listen on")
	flag.Parse()
}

func main() {
	if flag.Arg(0) == "version" {
		fmt.Printf("%s/%s\n", path.Base(os.Args[0]), "UNKNOWN")
		return
	}

	controller := controller.NewBaseController()

	err := server.Run(context.Background(), ":"+strconv.Itoa(options.Port), controller)
	if err != nil {
		glog.Errorf("Alibaba Cloud Service Broker Server failed to start. Error: %v", err.Error())
	}
}
