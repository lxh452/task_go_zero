package company

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PositionModel = (*customPositionModel)(nil)

type (
	// PositionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPositionModel.
	PositionModel interface {
		positionModel
		withSession(session sqlx.Session) PositionModel
	}

	customPositionModel struct {
		*defaultPositionModel
	}
)

// NewPositionModel returns a model for the database table.
func NewPositionModel(conn sqlx.SqlConn) PositionModel {
	return &customPositionModel{
		defaultPositionModel: newPositionModel(conn),
	}
}

func (m *customPositionModel) withSession(session sqlx.Session) PositionModel {
	return NewPositionModel(sqlx.NewSqlConnFromSession(session))
}
