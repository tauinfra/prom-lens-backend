package request

type CreateCredentialRequest struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	SecretType      string `json:"secretType" binding:"required"`
	BaseURL         string `json:"baseURL"`
	Username        string `json:"username"`
	EncryptedSecret string `json:"encryptedSecret" binding:"required"`
}

type UpdateCredentialRequest struct {
	Name            *string `json:"name"`
	SecretType      *string `json:"secretType"`
	BaseURL         *string `json:"baseURL"`
	Username        *string `json:"username"`
	EncryptedSecret *string `json:"encryptedSecret"`
}
