package api

type Schedular struct {
	Name     string `json:"name,required" validate:"min=2,nonzero"`
	URL      string `json:"url,required" validate:"min=2,nonzero"`
	Interval string `json:"interval,required" validate:"min=2,nonzero"`
}
