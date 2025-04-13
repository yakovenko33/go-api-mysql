package request

type RefreshTokensRequestData struct {
	UserId       string
	RefreshToken string
	UserAgent    string
}

type RefreshTokens struct {
	requestData *RefreshTokensRequestData
}

func CreateRefreshTokens(requestData *RefreshTokensRequestData) *RefreshTokens {
	return &RefreshTokens{
		requestData: requestData,
	}
}

func (m *RefreshTokens) GetValidationErrors() map[string]string {
	return nil
}

func (m *RefreshTokens) GetUserAgent() string {
	return m.requestData.UserAgent
}

func (m *RefreshTokens) GetRefreshToken() string {
	return m.requestData.RefreshToken
}

func (m *RefreshTokens) GetUserId() string {
	return m.requestData.UserId
}
