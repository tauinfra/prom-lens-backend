package dto

type TektonTaskRunDTO struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Succeeded string `json:"succeeded"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"createdAt"`
}
