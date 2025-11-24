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

func (r *Repository) RepeatByUsernameRepo(ctx context.Context, tx *gorm.DB, username string) bool {
	ctx, span := r.Trace(ctx, "svc.repo.RepeatByUsernameRepo", oteltrace.WithAttributes())
	defer span.End()
	// * tx
	if tx == nil {
		tx = r.DB().Ctx()
	}
	var modelPreRegister *model.PreRegister
	db := r.DB().Ctx().WithContext(ctx)
	err := db.Where("username = ?", username).First(&modelPreRegister).Error
	if err != nil {
		return false
	}
	return true
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

func (r *Repository) GetPreRegisterByUserPassEmailRepo(ctx context.Context, tx *gorm.DB, username, password, email string) (*model.PreRegister, error) {
	ctx, span := r.Trace(ctx, "svc.repo.GetPreRegisterByUserPassRepo")
	defer span.End()

	if tx == nil {
		tx = r.DB().Ctx()
	}
	query := tx.Model(&model.PreRegister{})
	if email != "" {
		query = query.Where("email = ?", email)
	}
	if username != "" {
		query = query.Where("username = ?", username)
	}
	//  password เป็น hashed
	query = query.Where("password = ?", password)
	var modelPreRegister model.PreRegister
	if err := query.First(&modelPreRegister).Error; err != nil {
		return nil, err
	}
	return &modelPreRegister, nil
}

func (r *Repository) UpdatePreRegisterByUserPassEmailRepo(ctx context.Context, tx *gorm.DB, preRegisterID string, dataUpdate map[string]interface{}) error {
	ctx, span := r.Trace(ctx, "svc.repo.UpdatePreRegisterByUserPassEmailRepo", oteltrace.WithAttributes())
	defer span.End()
	if tx == nil {
		tx = r.DB().Ctx()
	}
	if err := tx.Model(&model.PreRegister{}).
		Where("id = ?", preRegisterID).Updates(dataUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertUsersAuthRepo(ctx context.Context, tx *gorm.DB, userAuth model.User) (*model.User, error) {
	ctx, span := r.Trace(ctx, "svc.repo.InsertUsersAuthRepo", oteltrace.WithAttributes())
	defer span.End()
	// * tx
	if tx == nil {
		tx = r.DB().Ctx()
	}
	if err := tx.WithContext(ctx).Create(&userAuth).Error; err != nil {
		util.RecordSpanError(span, err, "repo.InsertUsersAuthRepo")
		return &userAuth, err
	}
	return &userAuth, nil
}

func (r *Repository) GetUsersAuthRepo(ctx context.Context, tx *gorm.DB, username, password string) (*model.User, error) {
	_, span := r.Trace(ctx, "svc.repo.GetPreRegisterByPreRegisterUidRepo", oteltrace.WithAttributes())
	defer span.End()
	if tx == nil {
		tx = r.DB().Ctx()
	}
	var user *model.User
	if err := tx.Model(&model.User{}).
		Where("username = ? AND password = ?", username, password).
		Take(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
