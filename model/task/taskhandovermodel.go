package task

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskHandoverModel = (*customTaskHandoverModel)(nil)

type (
	// TaskHandoverModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskHandoverModel.
	TaskHandoverModel interface {
		taskHandoverModel
		withSession(session sqlx.Session) TaskHandoverModel
	}

	customTaskHandoverModel struct {
		*defaultTaskHandoverModel
	}
)

// NewTaskHandoverModel returns a model for the database table.
func NewTaskHandoverModel(conn sqlx.SqlConn) TaskHandoverModel {
	return &customTaskHandoverModel{
		defaultTaskHandoverModel: newTaskHandoverModel(conn),
	}
}

func (m *customTaskHandoverModel) withSession(session sqlx.Session) TaskHandoverModel {
	return NewTaskHandoverModel(sqlx.NewSqlConnFromSession(session))
}
