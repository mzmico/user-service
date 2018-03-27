package impls

import (
	"context"

	"github.com/mzmico/mz"
	"github.com/mzmico/toolkit/state"
	pb "github.com/mzmico/user-service/protobuf"
)

type ServiceUser struct {
	pb.UserServer
}

func init() {
	mz.AddRpcServer(func(server *mz.RPCServer) {
		pb.RegisterUserServer(
			server,
			&ServiceUser{},
		)
	})
}

func (m *ServiceUser) Login(ctx context.Context, ask *pb.LoginRequest) (ack *pb.LoginResponse, err error) {

	_ = state.NewRpcState(ctx, ask.Session, &err)

	ack = &pb.LoginResponse{
		Token: "xxxx",
	}
	return
}
