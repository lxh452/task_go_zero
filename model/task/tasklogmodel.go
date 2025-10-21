package task

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskLogModel = (*customTaskLogModel)(nil)

type (
	// TaskLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskLogModel.
	TaskLogModel interface {
		taskLogModel
		withSession(session sqlx.Session) TaskLogModel
	}

	customTaskLogModel struct {
		*defaultTaskLogModel
	}
)

// NewTaskLogModel returns a model for the database table.
func NewTaskLogModel(conn sqlx.SqlConn) TaskLogModel {
	return &customTaskLogModel{
		defaultTaskLogModel: newTaskLogModel(conn),
	}
}

func (m *customTaskLogModel) withSession(session sqlx.Session) TaskLogModel {
	return NewTaskLogModel(sqlx.NewSqlConnFromSession(session))
}
