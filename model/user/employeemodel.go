package user

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmployeeModel = (*customEmployeeModel)(nil)

type (
	// EmployeeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmployeeModel.
	EmployeeModel interface {
		employeeModel
		withSession(session sqlx.Session) EmployeeModel
	}

	customEmployeeModel struct {
		*defaultEmployeeModel
	}
)

// NewEmployeeModel returns a model for the database table.
func NewEmployeeModel(conn sqlx.SqlConn) EmployeeModel {
	return &customEmployeeModel{
		defaultEmployeeModel: newEmployeeModel(conn),
	}
}

func (m *customEmployeeModel) withSession(session sqlx.Session) EmployeeModel {
	return NewEmployeeModel(sqlx.NewSqlConnFromSession(session))
}
