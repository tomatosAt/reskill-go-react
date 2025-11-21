package repositories

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	// "github.com/tomatosAt/reskill-go-react/pkg/cache"
)

const (
	AuthorizedContext = "auth"
	SessionTicketKey  = "session_ticket"
	SessionKey        = "session_auth"
)

var (
	sessionKey = func(uid string, accountId string) string {
		return fmt.Sprintf("session:%s:%s", uid, accountId)
	}
)

func (r *Repository) SetAuthSession(uid string, uidPreRegister string, s dto.AuthSession) error {
	s.ExpiresAt = time.Now().Add(config.CachingShortDuration)
	return r.cache.Set(sessionKey(uid, uidPreRegister), s, config.CachingShortDuration)
}

func (r *Repository) GetAuthSession(uid, accountId string) (*dto.AuthSession, error) {
	var session dto.AuthSession
	if err := r.cache.Get(fmt.Sprintf("session:%s:%s", uid, accountId), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) ClearAuthSession(uid, accountId string) error {
	return r.cache.Del(sessionKey(uid, accountId))
}

func (r *Repository) VerifyAuthSession(uid, accountId string, s *dto.AuthSession) error {
	if err := r.cache.Get(sessionKey(uid, accountId), s); err != nil {
		logrus.Errorln("get cache auth error ->", err)
		return err
	}
	if s.ExpiresAt.Sub(time.Now()) < time.Minute*15 {
		logrus.Infoln("expand time cache", s.ExpiresAt.Sub(time.Now()), (s.ExpiresAt.Sub(time.Now()) < time.Minute*15))
		r.cache.Client.Expire(sessionKey(uid, accountId), config.SessionTimeOut)
	}
	return nil
}
