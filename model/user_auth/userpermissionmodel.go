package user_auth

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserPermissionModel = (*customUserPermissionModel)(nil)

type (
	// UserPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserPermissionModel.
	UserPermissionModel interface {
		userPermissionModel
		withSession(session sqlx.Session) UserPermissionModel
	}

	customUserPermissionModel struct {
		*defaultUserPermissionModel
	}
)

// NewUserPermissionModel returns a model for the database table.
func NewUserPermissionModel(conn sqlx.SqlConn) UserPermissionModel {
	return &customUserPermissionModel{
		defaultUserPermissionModel: newUserPermissionModel(conn),
	}
}

func (m *customUserPermissionModel) withSession(session sqlx.Session) UserPermissionModel {
	return NewUserPermissionModel(sqlx.NewSqlConnFromSession(session))
}
