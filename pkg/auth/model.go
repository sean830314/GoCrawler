package auth

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	Authorized   bool
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
