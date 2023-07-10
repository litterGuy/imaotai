package reqfunc

type ImaotaiResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type LoginResp struct {
	ImaotaiResp
	Data *LoginUser `json:"data"`
}
type LoginUser struct {
	UserID       int64  `json:"userId"`
	UserName     string `json:"userName"`
	Mobile       string `json:"mobile"`
	VerifyStatus int    `json:"verifyStatus"`
	IDCode       string `json:"idCode"`
	IDType       int    `json:"idType"`
	Token        string `json:"token"`
	UserTag      int    `json:"userTag"`
	Cookie       string `json:"cookie"`
	Did          string `json:"did"`
}
