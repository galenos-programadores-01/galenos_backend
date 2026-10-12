package domain

import "encoding/json"

// ListadoUpsResponse es la respuesta del servicio REST del MINSA que lista
// las Unidades Productoras de Servicios (UPS) de un establecimiento
// (listadoUps/{codigo_renipress}).
type ListadoUpsResponse struct {
	Codigo  string          `json:"codigo"`
	Mensaje string          `json:"mensaje"`
	Datos   ListadoUpsItems `json:"datos"`
}

// ListadoUpsItem es una UPS del establecimiento devuelta por MINSA.
type ListadoUpsItem struct {
	CodUps      string `json:"codUps"`
	Descripcion string `json:"descripcion"`
}

// ListadoUpsItems decodifica el bloque "datos" de la respuesta. MINSA lo
// devuelve como arreglo cuando hay resultados y como cadena (p. ej. "") cuando
// no los hay, por eso implementa UnmarshalJSON.
type ListadoUpsItems []ListadoUpsItem

func (l *ListadoUpsItems) UnmarshalJSON(b []byte) error {
	raw := string(b)
	if raw == "null" || (len(raw) > 0 && raw[0] == '"') {
		return nil
	}
	var items []ListadoUpsItem
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}
	*l = ListadoUpsItems(items)
	return nil
}
