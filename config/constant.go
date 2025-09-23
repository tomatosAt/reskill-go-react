package config

import "time"

// Default configuration

const (
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvTest        = "test"

	ServerListen       = "0.0.0.0"
	ServerPort         = "8000"
	ServerTimeoutRead  = "15s"
	ServerTimeoutWrite = "15s"
	ServerTimeoutIdle  = "60s"

	MariadbHost = "127.0.0.1"
	MariadbPort = "3306"

	RedisHost      = "127.0.0.1"
	RedisPort      = "6379"
	RedisDb        = "0"
	ReadBufferSize = 8192
)

// Defined constant store

const (
	CachingExtremeShortDuration = time.Minute * 5
	CachingShortDuration        = time.Hour
	CachingMediumDuration       = time.Hour * 4
	CachingLongDuration         = time.Hour * 8
	CachingTokenExpire          = time.Hour * 24

	StoreBlockingState       = "block:ip:%s" // ip
	StoreBlockingPattern     = "block:ip:*"
	StoreSMSProviderIndex    = "config:sms:provider"
	CachingAccessTokenOne    = "cache:one_id:access_token:%s"
	CachingAccountDetailsOne = "cache:one_id:account_details:%s"
	// CachingRaOrgChartDetails            = "cache:ra:orgchart_details:%s"
	// CachingRaAccountDetails             = "cache:ra:account_details:%s"
	// CachingRaAccountDetailsByEmployeeID = "cache:ra:account_details_by_employee_id:%s"
	// CachingAccessTokenRbac              = "cache:rbac:access_token:%s"
	// CachingRbacDetails                  = "cache:rbac:details:%s"
	// CachingProfile    = "cache:profile:%s"
	// CachingPermission = "cache:permissson:%s"
)

// Tracer Attribute

const (
	EventCacheHit  = "Cache HIT"
	EventCacheMiss = "Cache MISS"
)

// Authorized context key

const (
	B2BAuthorizeContext = "b2b_auth"
	AuthorizeContext    = "auth"
)
