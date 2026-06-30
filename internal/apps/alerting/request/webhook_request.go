package request

type RouteMatcherInput struct {
	Label    string `json:"label" binding:"required"`
	Operator string `json:"operator"`
	Value    string `json:"value" binding:"required"`
}

type RouteInput struct {
	Enabled  *bool               `json:"enabled"`
	Priority *int                `json:"priority"`
	Continue *bool               `json:"continue"`
	Matchers []RouteMatcherInput `json:"matchers"`
}

type CreateWebhookRequest struct {
	Name        string      `json:"name" binding:"required"`
	URL         string      `json:"url" binding:"required"`
	Description string      `json:"description"`
	Enabled     *bool       `json:"enabled"`
	Route       *RouteInput `json:"route"`
}

type UpdateWebhookRequest struct {
	Name        *string     `json:"name"`
	URL         *string     `json:"url"`
	Description *string     `json:"description"`
	Enabled     *bool       `json:"enabled"`
	Route       *RouteInput `json:"route"`
}

type VerifyWebhookRequest struct {
	URL  string `json:"url" binding:"required"`
	Name string `json:"name"`
}
