package middlerware

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/middleware"
	"github.com/tomatosAt/reskill-go-react/model"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/cache"
	"github.com/tomatosAt/reskill-go-react/pkg/database"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

type FrontWebMiddleware struct {
	Skipper    *middleware.SkipperPath
	Store      *cache.Redis
	PrivateKey *rsa.PrivateKey
	DB         *database.Client
	// OneId      *one_id.Client
	Cache *cache.Redis
}

func NewFWAuthMiddleware(skipper *middleware.SkipperPath, db *database.Client, store *cache.Redis, privateKey *rsa.PrivateKey, cache *cache.Redis) *FrontWebMiddleware {
	return &FrontWebMiddleware{
		Skipper:    skipper,
		Store:      store,
		PrivateKey: privateKey,
		DB:         db,
		// OneId:      oneId,
		Cache: cache,
	}
}

func GetAuthCtxMiddleware(ctx context.Context) (*dto.SessionClaims, error) {
	authSession, ok := ctx.Value(config.AuthorizeContext).(dto.SessionClaims)
	if !ok {
		return nil, errors.New("unauthorized claims")
	}
	return &authSession, nil
}

func (fw FrontWebMiddleware) Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if fw.Skipper != nil && fw.Skipper.Test(ctx) {
			return ctx.Next()
		}

		token, err := util.GetCookie(ctx, "access_token")
		if err != nil || token == nil {
			_, err := util.GetCookie(ctx, "refresh_token")
			if err != nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "token is required.")
			}
			newClaims, err := fw.GenerateAccessTokenRepo(ctx)
			if err != nil || newClaims == nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "token is invalid.")
			}
			authCtx := context.WithValue(ctx.Context(), config.AuthorizeContext, *newClaims)
			ctx.SetUserContext(authCtx)
			return ctx.Next()
		}
		claims, err := fw.VerifyAccessTokenJWT(*token)

		if err != nil {
			_, err := util.GetCookie(ctx, "refresh_token")
			if err != nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "token is expired or invalid.")
			}
			newClaims, err := fw.GenerateAccessTokenRepo(ctx)
			if err != nil || newClaims == nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "token is invalid.")
			}
			authCtx := context.WithValue(ctx.Context(), config.AuthorizeContext, *newClaims)
			ctx.SetUserContext(authCtx)
			return ctx.Next()
		}

		authCtx := context.WithValue(ctx.Context(), config.AuthorizeContext, *claims)
		ctx.SetUserContext(authCtx)

		return ctx.Next()
	}
}

func (fw FrontWebMiddleware) VerifyAccessTokenJWT(tokenString string) (*dto.SessionClaims, error) {
	var sessionClaims dto.SessionClaims
	token, err := fw.DecodeJWTAccessToken(tokenString, &sessionClaims)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.SessionClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (fw FrontWebMiddleware) DecodeJWTAccessToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := util.ParseWithClaims(tokenString, claims, fw.PrivateKey)
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token, nil
	}
	return nil, errors.New("decode token error")
}

func (fw FrontWebMiddleware) GenerateAccessTokenRepo(ctx *fiber.Ctx) (*dto.SessionClaims, error) {

	refreshToken, _ := util.GetCookie(ctx, "refresh_token")
	if refreshToken == nil {
		return nil, errors.New("token is invalid")
	}
	// tokenOne, err := fw.RefreshTokens(ctx, refreshToken)
	// if err != nil {
	// 	return nil, errors.New("token is invalid")
	// }
	// if tokenOne.AccountId == "" {
	// 	return nil, errors.New("invalid account ID in token")
	// }
	// userId, err := fw.GetUserByIdRepo(ctx.UserContext(), tokenOne.AccountId)
	// if err != nil {
	// 	return nil, errors.New("token is invalid")
	// }

	// token, subject, err := fw.GenerateNewAccessTokenRepo(ctx.UserContext(), userId.Id.String())
	// if err != nil {

	// 	return nil, errors.New("token is invalid")
	// }

	// newClaims, err := fw.VerifyAccessTokenJWT(token)
	// if err != nil || newClaims == nil {
	// 	return nil, errors.New("token is invalid")
	// }
	// fw.Cache.Set(fmt.Sprintf(config.CachingAccessTokenOne, subject), tokenOne, config.CachingTokenExpire)

	// timeExpireToken := time.Now().Add(time.Hour * 12)
	// timeExpireRefreshToken := time.Now().AddDate(0, 3, 0)

	// if token != "" {
	// 	util.SetCookie(ctx, "access_token", token, &timeExpireToken)
	// }
	// if tokenOne.RefreshTokenPKG != "" {
	// 	util.SetCookie(ctx, "refresh_token", tokenOne.RefreshTokenPKG, &timeExpireRefreshToken)
	// }
	// return newClaims, nil
	return nil, nil
}

func (fw FrontWebMiddleware) GenerateNewAccessTokenRepo(ctx context.Context, userID string) (string, string, error) {
	sid := util.GenUniqueIdV7()
	issAt := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location())
	expiresAt := issAt.Add(24 * time.Hour)
	claims := util.JWTStandardClaims{
		StandardClaims: &jwt.StandardClaims{
			Audience:  userID,
			IssuedAt:  issAt.Unix(),
			Subject:   sid.String(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token, err := util.EncodeJWTAccessToken(claims, fw.PrivateKey)
	if err != nil {
		return "", "", err
	}
	return token, sid.String(), nil
}

func (fw FrontWebMiddleware) GetUserByIdRepo(ctx context.Context, accountId string) (*model.User, error) {
	var users model.User
	if err := fw.DB.Ctx().WithContext(ctx).Where("account_id = ?", accountId).Take(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// func (fw FrontWebMiddleware) RefreshTokens(ctx *fiber.Ctx, refreshToken *string) (*one_id.ResponseOauthToken, error) {
// 	token, err := fw.OneId.RefreshTokenPKG(ctx.UserContext(), refreshToken)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }
