package ports

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/model"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/cache"
	"github.com/tomatosAt/reskill-go-react/pkg/database"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type Repository interface {
	Module() string
	AppCfg() *config.Config
	Log() *logrus.Entry
	DB() *database.Client
	Cache() *cache.Redis
	Trace(ctx context.Context, spanName string, attributes ...trace.SpanStartOption) (context.Context, trace.Span)
	// pre register
	GetPreRegisterByEmailTelNoRepo(ctx context.Context, tx *gorm.DB, email, telNo string) (*model.PreRegister, error)
	InsertPreRegisterRepo(ctx context.Context, tx *gorm.DB, data model.PreRegister) (*model.PreRegister, error)
	InsertTransactionAuthRepo(ctx context.Context, tx *gorm.DB, transactionAuth model.TransactionAuth) (*model.TransactionAuth, error)
	GetPreRegisterByPreRegisterUidRepo(ctx context.Context, tx *gorm.DB, preRegistUid, transactionAuthUid string) (*model.TransactionAuth, error)
	SetAuthSession(uid, uidPreRegister string, s dto.AuthSession) error
	// user
	GetAuthCtxRepo(ctx context.Context) (*dto.AuthSession, error)
	RepeatByUsernameRepo(ctx context.Context, tx *gorm.DB, username string) bool
	InsertUsersAuthRepo(ctx context.Context, tx *gorm.DB, userAuth model.User) (*model.User, error)
	GetUsersAuthRepo(ctx context.Context, tx *gorm.DB, username, password string) (*model.User, error)
	// session
	GetAuthSession(uid, accountId string) (*dto.AuthSession, error)
	ClearAuthSession(uid, accountId string) error
	// auth
	GetPreRegisterByUserPassEmailRepo(ctx context.Context, tx *gorm.DB, username, password, email string) (*model.PreRegister, error)
	UpdatePreRegisterByUserPassEmailRepo(ctx context.Context, tx *gorm.DB, preRegisterID string, dataUpdate map[string]interface{}) error
}

type Service interface {
	// pre register
	DOBCheckSVC(ctx context.Context, dob string) (string, error)
	TitleNameCheckSVC(ctx context.Context, titleNameTH, titleNameENG string) error
	CheckFormatPreRegisterSVC(ctx context.Context, data *dto.PreRegisterDataBasePayload) error
	PreRegisterSVC(ctx context.Context, data dto.PreRegisterDataBasePayload) (*dto.ResponsePreRegister, int, error)
	// session
	DecryptStringToStructService(ctx context.Context, code string) (dto.PreRegisterConvern, error)
	CreateSessionService(ctx context.Context, code string) (dto.Session, error)
	GetProfilesService(ctx context.Context) (dto.ResponseUserProfiles, int, error)
	// auth
	CheckFormatUsernameSVC(ctx context.Context, data *dto.UserPasswordPayload) error
	LoginUserPassService(ctx context.Context, userPass dto.UserPasswordPayload) (*dto.ResponsePreRegister, int, error)
	LogoutService(ctx context.Context) (int, error)
}
