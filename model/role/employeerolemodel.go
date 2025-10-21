package role

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmployeeRoleModel = (*customEmployeeRoleModel)(nil)

type (
	// EmployeeRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmployeeRoleModel.
	EmployeeRoleModel interface {
		employeeRoleModel
		withSession(session sqlx.Session) EmployeeRoleModel
	}

	customEmployeeRoleModel struct {
		*defaultEmployeeRoleModel
	}
)

// NewEmployeeRoleModel returns a model for the database table.
func NewEmployeeRoleModel(conn sqlx.SqlConn) EmployeeRoleModel {
	return &customEmployeeRoleModel{
		defaultEmployeeRoleModel: newEmployeeRoleModel(conn),
	}
}

func (m *customEmployeeRoleModel) withSession(session sqlx.Session) EmployeeRoleModel {
	return NewEmployeeRoleModel(sqlx.NewSqlConnFromSession(session))
}
