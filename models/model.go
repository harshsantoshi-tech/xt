package models

// OfferWidgetTag The models for the OfferWidgetTag and its nested components.
type OfferWidgetTag struct {
	Title       UikitText  `json:"title"`
	Style       UikitStyle `json:"style"`
	ButtonStyle UikitStyle `json:"buttonStyle"`
}

type UikitText struct {
	Value string `json:"value"`
}

type UikitStyle struct {
	Color    string  `json:"color,omitempty"`
	FontSize float64 `json:"font_size,omitempty"`
	Value    string  `json:"value,omitempty"`
}

type UikitBackground struct {
	Color    string  `json:"color"`
	FontSize float64 `json:"font_size"`
}

// Response The model for the complete API response.
type Response struct {
	Response ResponseModel  `json:"response"`
	Compete  OfferWidgetTag `json:"compete"`
}

// ResponseModel A generic response model.
type ResponseModel struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Greeting A model for the JSON post request.
type Greeting struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Time    string `json:"time"`
}
