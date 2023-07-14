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

type SessionResp struct {
	ImaotaiResp
	Data SessionData `json:"data"`
}
type SessionItemList struct {
	Content   string `json:"content"`
	ItemCode  string `json:"itemCode"`
	JumpURL   string `json:"jumpUrl"`
	Picture   string `json:"picture"`
	PictureV2 string `json:"pictureV2"`
	Title     string `json:"title"`
}
type SessionData struct {
	ItemList  []SessionItemList `json:"itemList"`
	SessionID int               `json:"sessionId"`
}

type ShopResp struct {
	ImaotaiResp
	Data ShopRespData `json:"data"`
}

type ShopRespItems struct {
	Count               int    `json:"count"`
	MaxReserveCount     int    `json:"maxReserveCount"`
	DefaultReserveCount int    `json:"defaultReserveCount"`
	ItemID              string `json:"itemId"`
	Inventory           int    `json:"inventory"`
	OwnerName           string `json:"ownerName"`
}
type ShopRespShops struct {
	ShopID string          `json:"shopId"`
	Items  []ShopRespItems `json:"items"`
}

type ShopRespData struct {
	Shops     []ShopRespShops `json:"shops"`
	ValidTime int64           `json:"validTime"`
}

type ShopItemBean struct {
	ShopRespItems
	ShopID string `json:"shopId"`
}

type ShopBean struct {
	Address       string   `json:"address"`
	City          int      `json:"city"`
	CityName      string   `json:"cityName"`
	District      int      `json:"district"`
	DistrictName  string   `json:"districtName"`
	FullAddress   string   `json:"fullAddress"`
	Lat           float64  `json:"lat"`
	Layaway       bool     `json:"layaway"`
	Lng           float64  `json:"lng"`
	Name          string   `json:"name"`
	OpenEndTime   string   `json:"openEndTime"`
	OpenStartTime string   `json:"openStartTime"`
	Province      int      `json:"province"`
	ProvinceName  string   `json:"provinceName"`
	ShopID        string   `json:"shopId"`
	Tags          []string `json:"tags"`
	TenantName    string   `json:"tenantName"`
	Distance      float64  `json:"distance"`
}

type ShopResourceResp struct {
	ImaotaiResp
	Data ShopResourceData `json:"data"`
}

type ShopResourceData struct {
	MtshopsPc MtshopsPc `json:"mtshops_pc"`
}

type MtshopsPc struct {
	Url string `json:"url"`
}
