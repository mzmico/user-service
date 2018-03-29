package impls

import (
	"context"
	"database/sql"

	"github.com/mzmico/mz"
	"github.com/mzmico/mz/rpc_service"
	"github.com/mzmico/toolkit/db"
	"github.com/mzmico/toolkit/state"
	"github.com/mzmico/toolkit/utils"
	"github.com/mzmico/toolkit/wechat/wxapp"
	pb "github.com/mzmico/user-service/protobuf"
)

type ServiceUser struct {
	pb.UserServer
	wApp *wxapp.Config
}

func init() {
	mz.AddRpcServer(func(server *mz.RPCServer) {
		pb.RegisterUserServer(
			server,
			&ServiceUser{
				wApp: &wxapp.Config{
					AppId:  "",
					Secret: "",
				},
			},
		)
	})
}

type AccountType int

const (
	ACCOUNT_TYPE_UNKNOWN                    = 0
	ACCOUNT_TYPE_WXAPP_UNION_ID AccountType = 1
	ACCOUNT_TYPE_WXAPP_OPEN_ID  AccountType = 2
)

func (m *ServiceUser) Login(ctx context.Context, ask *pb.LoginRequest) (ack *pb.LoginResponse, err error) {

	state := state.NewRpcState(
		rpc_service.GetService(),
		ctx,
		ask.Session,
		&err)

	ack = new(pb.LoginResponse)

	var (
		acType AccountType = ACCOUNT_TYPE_UNKNOWN
	)

	switch ask.Type {
	case pb.AccountType_ACCOUNT_TYPE_WECHAT_JSCODE:

		session, err := wxapp.JavascriptCodeToSession(m.wApp, ask.Certificate)

		if err != nil {
			return nil, err
		}

		if len(session.UnionID) != 0 {
			acType = ACCOUNT_TYPE_WXAPP_UNION_ID
		} else {
			acType = ACCOUNT_TYPE_WXAPP_OPEN_ID
		}
	default:
		return nil, state.Errorf(
			"account type %s not support", ask.Type,
		)
	}

	dbUser := db.Use("db_user")

	type account struct {
		Uid         string `db:"uid"`
		Certificate string `json:"certificate"`
	}

	a := account{}

	err = dbUser.Get(
		&a,
		`SELECT 
					uid,
					certificate 
				FROM tb_account 
				WHERE app_id=? AND account=? AND type=?`,
		ask.AppId,
		ask.Account,
		acType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			ack.Status = pb.LoginStatus_LOGIN_STATUS_NOT_EXISTS
			return
		}

		return nil, state.Error(err)
	}

	if a.Certificate != ask.Certificate {
		ack.Uid = a.Uid
		ack.Status = pb.LoginStatus_LOGIN_STATUS_PASSOWRD_ERROR

		return
	}

	ack.Token = utils.NewShortUUID()
	ack.Status = pb.LoginStatus_LOGIN_STATUS_OK
	return
}
