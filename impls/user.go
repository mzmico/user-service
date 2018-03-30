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
	"github.com/spf13/viper"
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
					AppId:  viper.GetString("wxapp.appid"),
					Secret: viper.GetString("wxapp.secret"),
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
	case pb.LoginType_LOGIN_TYPE_WECHAT_JSCODE:

		session, err := wxapp.JavascriptCodeToSession(m.wApp, ask.Certificate)

		if err != nil {
			return ack, state.Error(err)
		}

		if len(session.UnionID) != 0 {
			ack.Account = session.UnionID
			ack.Type = pb.AccountType_ACCOUNT_TYPE_WECHAT_APP_UNIONID
		} else {
			ack.Account = session.OpenID
			ack.Type = pb.AccountType_ACCOUNT_TYPE_WECHAT_APP_OPENID
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
			err = nil
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

func (m *ServiceUser) CreateUser(ctx context.Context, ask *pb.CreateUserRequest) (ack *pb.CreateUserResponse, err error) {

	state := state.NewRpcState(
		rpc_service.GetService(),
		ctx,
		ask.Session,
		&err)

	ack = new(pb.CreateUserResponse)

	ack.Uid = utils.NewShortUUID()

	dbUser := db.Use("db_user")

	_, err = dbUser.ExecContext(
		ctx,
		"INSERT INTO tb_user(app_id, uid, avatar, extend) VALUES (?,?,?,?)",
		ask.AppId,
		ack.Uid,
		ask.Avatar,
		ask.Extend,
	)

	if err != nil {
		state.Error(err)
	}

	return
}
