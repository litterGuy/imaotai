package reqfunc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetLocationByAddress(address string) ([]Geocodes, error) {
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=0a7f4baae0a5e37e6f90e4dc88e3a10d&output=json&address=%s", address)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rt := new(LocationResp)
	err = json.Unmarshal(body, &rt)
	if err != nil {
		return nil, err
	}

	if rt.Status != "1" {
		return nil, errors.New(rt.Info)
	}
	return rt.Geocodes, nil
}

type LocationResp struct {
	Status   string     `json:"status"`
	Info     string     `json:"info"`
	Infocode string     `json:"infocode"`
	Count    string     `json:"count"`
	Geocodes []Geocodes `json:"geocodes"`
}
type Neighborhood struct {
	Name []interface{} `json:"name"`
	Type []interface{} `json:"type"`
}
type Building struct {
	Name []interface{} `json:"name"`
	Type []interface{} `json:"type"`
}
type Geocodes struct {
	FormattedAddress string        `json:"formatted_address"`
	Country          string        `json:"country"`
	Province         string        `json:"province"`
	Citycode         string        `json:"citycode"`
	City             string        `json:"city"`
	District         string        `json:"district"`
	Township         []interface{} `json:"township"`
	Neighborhood     Neighborhood  `json:"neighborhood"`
	Building         Building      `json:"building"`
	Adcode           string        `json:"adcode"`
	Street           []interface{} `json:"street"`
	Number           []interface{} `json:"number"`
	Location         string        `json:"location"`
	Level            string        `json:"level"`
}
