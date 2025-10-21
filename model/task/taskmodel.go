package task

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskModel = (*customTaskModel)(nil)

type (
	// TaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskModel.
	TaskModel interface {
		taskModel
		withSession(session sqlx.Session) TaskModel
	}

	customTaskModel struct {
		*defaultTaskModel
	}
)

// NewTaskModel returns a model for the database table.
func NewTaskModel(conn sqlx.SqlConn) TaskModel {
	return &customTaskModel{
		defaultTaskModel: newTaskModel(conn),
	}
}

func (m *customTaskModel) withSession(session sqlx.Session) TaskModel {
	return NewTaskModel(sqlx.NewSqlConnFromSession(session))
}
