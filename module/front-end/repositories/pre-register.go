package repositories

import (
	"context"

	"github.com/tomatosAt/reskill-go-react/model"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (r *Repository) GetPreRegisterByEmailTelNoRepo(ctx context.Context, tx *gorm.DB, email, telNo string) (*model.PreRegister, error) {
	ctx, span := r.Trace(ctx, "svc.repo.GetPreRegisterByEmailTelNoRepo", oteltrace.WithAttributes())
	defer span.End()
	// * tx
	if tx == nil {
		tx = r.DB().Ctx()
	}
	var modelPreRegister *model.PreRegister
	db := r.DB().Ctx().WithContext(ctx)
	err := db.Where("email = ? AND tel = ?", email, telNo).First(&modelPreRegister).Error
	if err != nil {
		return nil, err
	}
	return modelPreRegister, nil
}

func (r *Repository) InsertPreRegisterRepo(ctx context.Context, tx *gorm.DB, data model.PreRegister) (*model.PreRegister, error) {
	ctx, span := r.Trace(ctx, "svc.repo.InsertPreRegisterRepo", oteltrace.WithAttributes())
	defer span.End()
	// * tx
	if tx == nil {
		tx = r.DB().Ctx()
	}
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		util.RecordSpanError(span, err, "repo.InsertUsersRepo")
		return &data, err
	}
	return &data, nil
}

func (r *Repository) InsertTransactionAuthRepo(ctx context.Context, tx *gorm.DB, transactionAuth model.TransactionAuth) (*model.TransactionAuth, error) {
	ctx, span := r.Trace(ctx, "svc.repo.InsertTransactionAuthRepo", oteltrace.WithAttributes())
	defer span.End()
	// * tx
	if tx == nil {
		tx = r.DB().Ctx()
	}
	if err := tx.WithContext(ctx).Create(&transactionAuth).Error; err != nil {
		util.RecordSpanError(span, err, "repo.InsertTransactionAuthRepo")
		return &transactionAuth, err
	}
	return &transactionAuth, nil
}

func (r *Repository) GetPreRegisterByPreRegisterUidRepo(ctx context.Context, tx *gorm.DB, preRegistUid, transactionAuthUid string) (*model.TransactionAuth, error) {
	_, span := r.Trace(ctx, "svc.repo.GetPreRegisterByPreRegisterUidRepo", oteltrace.WithAttributes())
	defer span.End()
	if tx == nil {
		tx = r.DB().Ctx()
	}
	var modelTransactionAuth *model.TransactionAuth
	if err := tx.Model(&model.TransactionAuth{}).
		Preload("PreRegister").
		Where("pre_register_uid = ? AND id = ?", preRegistUid, transactionAuthUid).
		Take(&modelTransactionAuth).Error; err != nil {
		return nil, err
	}
	return modelTransactionAuth, nil
}
