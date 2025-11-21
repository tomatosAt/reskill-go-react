package middlerware

import (
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

const (
	AuthorizedContext = "auth"
	SessionTicketKey  = "session_ticket"
	SessionKey        = "session_auth"
)

type AuthSession struct {
	Uid               uuid.UUID `json:"uid"`
	PreRegisterUid    uuid.UUID `json:"pre_register_uid"`
	AccountId         string    `json:"account_id"`
	ClientIp          string    `json:"client_ip"`
	UserAgent         string    `json:"user_agent"`
	TokenType         string    `json:"token_type"`
	Name              string    `json:"name"`
	RegisterStatus    string    `json:"register_status"`
	ExpiresAt         time.Time `json:"expires_at"`
	OneClientID       string    `json:"one_client_id"`
	OneClientSecret   string    `json:"one_client_secret"`
	OneRefCode        string    `json:"one_ref_code"`
	RedirectUri       string    `json:"redirect_uri"`
	RedirectTokenType string    `json:"redirect_token_type"`
	HashIdCardNum     string    `json:"hash_id_card_num"`
	LoginType         string    `json:"login_type"`
	// TempAccountUid uuid.UUID `json:"temp_account_uid"`
	// HaveTempAccount bool `json:"have_temp_account"`
	// JoinBizComplete bool `json:"join_biz_complete"`
	// JoinBizRequest bool `json:"join_biz_request"`
	//UpdateTempAccountRequest bool `json:"update_temp_account_request"`
	//UpdateTempAccountComplete bool `json:"update_temp_account_complete"`
	//UpdateMetaDataRequest bool `json:"update_meta_data_request"`
	//UpdateMetaDataComplete bool `json:"update_meta_data_complete"`
	//UpdateAddressRequest bool `json:"update_address_request"`
	//UpdateAddressComplete bool `json:"update_address_complete"`
	// InsertProfile bool `json:"insert_profile"`
	TransactionUid uuid.UUID `json:"transaction_uid"`
	// AddressUid uuid.UUID `json:"address_uid"`
	// MetaDataUid uuid.UUID `json:"meta_data_uid"`
	// BusinessUid uuid.UUID `json:"business_uid"`
	// MetaDataType string `json:"meta_data_type"`
	// FormUid string `json:"form_uid"`
	// ExchangeUid uuid.UUID `json:"exchange_uid"`
}

type SessionClaims struct {
	*jwt.StandardClaims
}

type SignInTicket struct {
	Id       uuid.UUID
	Salt     string
	ClientIp string
	Agent    string
}

type RefreshTokenClaims struct {
	SessionUid uuid.UUID  `json:"session_uid"`
	AccountId  string     `json:"account_id"`
	IssuedAt   time.Time  `json:"issued_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
}
type RefreshClaims struct {
	Subject  string `json:"sub"` // user id หรือ session id
	Audience string `json:"aud"` // pre-register uid หรือ project id
	*jwt.StandardClaims
}
