package domain

type MedicoListaEspera struct {
	IdMedico       int     `json:"idMedico"`
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	Nombres         string `json:"nombres"`
	Dmedico         *string `json:"dmedico"`
}
