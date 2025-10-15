package core

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AuthAccountModel = (*customAuthAccountModel)(nil)

type (
	// AuthAccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthAccountModel.
	AuthAccountModel interface {
		authAccountModel
		withSession(session sqlx.Session) AuthAccountModel
	}

	customAuthAccountModel struct {
		*defaultAuthAccountModel
	}
)

// NewAuthAccountModel returns a model for the database table.
func NewAuthAccountModel(conn sqlx.SqlConn) AuthAccountModel {
	return &customAuthAccountModel{
		defaultAuthAccountModel: newAuthAccountModel(conn),
	}
}

func (m *customAuthAccountModel) withSession(session sqlx.Session) AuthAccountModel {
	return NewAuthAccountModel(sqlx.NewSqlConnFromSession(session))
}
