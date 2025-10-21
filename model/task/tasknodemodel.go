package task

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskNodeModel = (*customTaskNodeModel)(nil)

type (
	// TaskNodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskNodeModel.
	TaskNodeModel interface {
		taskNodeModel
		withSession(session sqlx.Session) TaskNodeModel
	}

	customTaskNodeModel struct {
		*defaultTaskNodeModel
	}
)

// NewTaskNodeModel returns a model for the database table.
func NewTaskNodeModel(conn sqlx.SqlConn) TaskNodeModel {
	return &customTaskNodeModel{
		defaultTaskNodeModel: newTaskNodeModel(conn),
	}
}

func (m *customTaskNodeModel) withSession(session sqlx.Session) TaskNodeModel {
	return NewTaskNodeModel(sqlx.NewSqlConnFromSession(session))
}
