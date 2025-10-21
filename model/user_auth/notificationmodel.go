package user_auth

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ NotificationModel = (*customNotificationModel)(nil)

type (
	// NotificationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotificationModel.
	NotificationModel interface {
		notificationModel
		withSession(session sqlx.Session) NotificationModel
	}

	customNotificationModel struct {
		*defaultNotificationModel
	}
)

// NewNotificationModel returns a model for the database table.
func NewNotificationModel(conn sqlx.SqlConn) NotificationModel {
	return &customNotificationModel{
		defaultNotificationModel: newNotificationModel(conn),
	}
}

func (m *customNotificationModel) withSession(session sqlx.Session) NotificationModel {
	return NewNotificationModel(sqlx.NewSqlConnFromSession(session))
}
