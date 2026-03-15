package request

type CreateNamespaceRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateNamespaceRequest struct {
	Labels map[string]string `json:"labels,omitempty"`
}
