package dto

import "valyria-backend/internal/apps/dragon/model"

func ToCredentialDTO(item model.Credential) CredentialDTO {
	return CredentialDTO{
		ID:         item.ID,
		Name:       item.Name,
		Type:       item.Type,
		SecretType: item.SecretType,
		BaseURL:    item.BaseURL,
		Username:   item.Username,
		Creator:    item.Creator,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}
