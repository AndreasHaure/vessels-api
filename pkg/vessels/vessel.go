package vessels

type Vessel struct {
	IMO       int64  `json:"imo"`
	Name      string `json:"name"`
	Flag      string `json:"flag"`
	YearBuilt int64  `json:"year_built"`
	Owner     string `json:"owner"`
}
