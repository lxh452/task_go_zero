package core

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserTaskLogModel = (*customUserTaskLogModel)(nil)

type (
	// UserTaskLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTaskLogModel.
	UserTaskLogModel interface {
		userTaskLogModel
		withSession(session sqlx.Session) UserTaskLogModel
	}

	customUserTaskLogModel struct {
		*defaultUserTaskLogModel
	}
)

// NewUserTaskLogModel returns a model for the database table.
func NewUserTaskLogModel(conn sqlx.SqlConn) UserTaskLogModel {
	return &customUserTaskLogModel{
		defaultUserTaskLogModel: newUserTaskLogModel(conn),
	}
}

func (m *customUserTaskLogModel) withSession(session sqlx.Session) UserTaskLogModel {
	return NewUserTaskLogModel(sqlx.NewSqlConnFromSession(session))
}
