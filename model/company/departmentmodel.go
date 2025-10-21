package company

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DepartmentModel = (*customDepartmentModel)(nil)

type (
	// DepartmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDepartmentModel.
	DepartmentModel interface {
		departmentModel
		withSession(session sqlx.Session) DepartmentModel
	}

	customDepartmentModel struct {
		*defaultDepartmentModel
	}
)

// NewDepartmentModel returns a model for the database table.
func NewDepartmentModel(conn sqlx.SqlConn) DepartmentModel {
	return &customDepartmentModel{
		defaultDepartmentModel: newDepartmentModel(conn),
	}
}

func (m *customDepartmentModel) withSession(session sqlx.Session) DepartmentModel {
	return NewDepartmentModel(sqlx.NewSqlConnFromSession(session))
}
