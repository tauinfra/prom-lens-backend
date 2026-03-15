package request

type CreateMenuRequest struct {
	ParentID   *uint  `json:"parentId"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Component  string `json:"component"`
	Title      string `json:"title" binding:"required"`
	Icon       string `json:"icon"`
	Rank       int    `json:"rank"`
	ShowParent *bool  `json:"showParent"`
	ShowLink   *bool  `json:"showLink"`
	Creator    string `json:"creator"`
}

type UpdateMenuRequest struct {
	ParentID   *uint   `json:"parentId"`
	Name       *string `json:"name"`
	Path       *string `json:"path"`
	Component  *string `json:"component"`
	Title      *string `json:"title"`
	Icon       *string `json:"icon"`
	Rank       *int    `json:"rank"`
	ShowParent *bool   `json:"showParent"`
	ShowLink   *bool   `json:"showLink"`
}
