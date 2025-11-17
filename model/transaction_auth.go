package model

type TransactionAuth struct {
	Model
	Action         string      `gorm:"column:action"`
	StatusRegister string      `gorm:"column:status_register"` // pending, success, failed
	Remark         string      `gorm:"column:payment_method"`
	PreRegisterUid string      `gorm:"column:pre_register_uid"`
	PreRegister    PreRegister `gorm:"foreignKey:PreRegisterUid;references:Id"`
}
