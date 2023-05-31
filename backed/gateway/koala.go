package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"koala/gateway/internal/config"
	accountServer "koala/gateway/internal/server/account"
	apiServer "koala/gateway/internal/server/api"
	applicationServer "koala/gateway/internal/server/application"
	authorityServer "koala/gateway/internal/server/authority"
	deptServer "koala/gateway/internal/server/dept"
	menuServer "koala/gateway/internal/server/menu"
	postServer "koala/gateway/internal/server/post"
	roleServer "koala/gateway/internal/server/role"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/interceptor"
	"koala/gateway/pb"
)

var configFile = flag.String("f", "gateway/etc/koala.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	sGroup := service.NewServiceGroup()

	rpc := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAuthorityServer(grpcServer, authorityServer.NewAuthorityServer(ctx))
		pb.RegisterAccountServer(grpcServer, accountServer.NewAccountServer(ctx))
		pb.RegisterApplicationServer(grpcServer, applicationServer.NewApplicationServer(ctx))
		pb.RegisterMenuServer(grpcServer, menuServer.NewMenuServer(ctx))
		pb.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		pb.RegisterDeptServer(grpcServer, deptServer.NewDeptServer(ctx))
		pb.RegisterPostServer(grpcServer, postServer.NewPostServer(ctx))
		pb.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	//gateway log,grpc的全局拦截器
	rpc.AddUnaryInterceptors(interceptor.LoggerInterceptor)
	rpc.AddUnaryInterceptors()
	sGroup.Add(rpc)

	api := gateway.MustNewServer(c.Gateway)
	api.Use(interceptor.NewAuthorityMiddleware().Handle)
	sGroup.Add(api)

	defer sGroup.Stop()

	fmt.Printf("Starting gateway server at %s...\n", c.ListenOn)
	fmt.Printf("Starting api at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	sGroup.Start()
}
