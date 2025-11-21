package middlerware

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/patcharp/golib/v2/helper"
	"github.com/patcharp/golib/v2/util/httputil"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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

func (fw FrontWebMiddleware) NewAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if fw.Skipper != nil && fw.Skipper.Test(ctx) {
			return ctx.Next()
		}
		// 1) Get token
		token := httputil.GetTokenFromHeader(ctx, "Bearer")
		if token == "" {
			token = httputil.GetTokenFromCookie(ctx, "token")
			if token == "" {
				return helper.HttpInvalidRequest(ctx, 401, nil, "authorized token not found")
			}
		}

		// 2) Try verify access token
		claims, err := fw.VerifyAccessToken(token)
		// 3) If access token expired -> refresh flow
		if err != nil {
			// get refresh token from session or cookie
			refreshToken := httputil.GetTokenFromCookie(ctx, "refresh_token")
			if refreshToken == "" {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "access token expired")
			}
			// decode refresh token
			refreshClaims, err := fw.VerifyRefreshToken(refreshToken)
			if err != nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "invalid refresh token")
			}
			// verify session by refreshClaims
			var authSession dto.AuthSession
			if err := fw.VerifyAuthSession(uuid.FromStringOrNil(refreshClaims.Subject), refreshClaims.Audience, &authSession); err != nil {
				return helper.HttpInvalidRequest(ctx, 401, nil, "invalid authorization session")
			}
			// generate new access token
			newToken, _, err := util.GenerateNewAccessTokenRepo(
				authSession.Uid.String(),
				authSession.PreRegisterUid,
				fw.PrivateKey,
			)
			if err != nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "refresh failed")
			}
			// verify new access token (claims)
			claims, err = fw.VerifyAccessToken(newToken)
			if err != nil {
				return util.HttpError(ctx, fiber.StatusUnauthorized, "token invalid after refresh")
			}
		}
		// 4) verify session (only once)
		var authSession dto.AuthSession
		if err := fw.VerifyAuthSession(uuid.FromStringOrNil(claims.Subject), claims.Audience, &authSession); err != nil {
			return helper.HttpInvalidRequest(ctx, 401, nil, "invalid authorization session")
		}
		// 5) attach session
		authCtx := context.WithValue(ctx.UserContext(), config.AuthorizeContext, authSession)
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

func (fw FrontWebMiddleware) GetUserByIdRepo(ctx context.Context, accountId string) (*model.User, error) {
	var users model.User
	if err := fw.DB.Ctx().WithContext(ctx).Where("account_id = ?", accountId).Take(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

//	func (fw FrontWebMiddleware) RefreshTokens(ctx *fiber.Ctx, refreshToken *string) (*one_id.ResponseOauthToken, error) {
//		token, err := fw.OneId.RefreshTokenPKG(ctx.UserContext(), refreshToken)
//		if err != nil {
//			return nil, err
//		}
//		return token, nil
//	}

func (fw FrontWebMiddleware) VerifyRefreshToken(tokenString string) (*SessionClaims, error) {
	var sessionClaims SessionClaims
	token, err := util.DecodeJWTAccessToken(tokenString, &sessionClaims, fw.PrivateKey)
	if err != nil {
		return nil, err
	}
	// แปลง claims เหมือน VerifyAccessToken
	claims, ok := token.(*SessionClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (fw FrontWebMiddleware) VerifyAccessToken(tokenString string) (*SessionClaims, error) {
	var sessionClaims SessionClaims
	token, err := util.DecodeJWTAccessToken(tokenString, &sessionClaims, fw.PrivateKey)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.(*SessionClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func (fw FrontWebMiddleware) VerifyAuthSession(uid uuid.UUID, accountId string, s *dto.AuthSession) error {
	sessionKey := fmt.Sprintf("session:%s:%s", uid.String(), accountId)
	if err := fw.Cache.Get(sessionKey, s); err != nil {
		logrus.Errorln("get cache auth error ->", err)
		return err
	}
	if s.ExpiresAt.Sub(time.Now()) < time.Minute*15 {
		logrus.Infoln("expand time cache", s.ExpiresAt.Sub(time.Now()), (s.ExpiresAt.Sub(time.Now()) < time.Minute*15))
		fw.Cache.Client.Expire(sessionKey, config.SessionTimeOut)
	}
	return nil
}
