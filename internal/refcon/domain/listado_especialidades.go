package domain

import "encoding/json"

// ListadoEspecialidadesResponse es la respuesta del servicio REST de
// interoperabilidad del MINSA que lista las especialidades vigentes
// (listadoEspecialidades).
type ListadoEspecialidadesResponse struct {
	Codigo string                     `json:"codigo"`
	Data   ListadoEspecialidadesItems `json:"data"`
}

// EspecialidadMinsa es una especialidad vigente devuelta por MINSA.
type EspecialidadMinsa struct {
	CodigoEspecialidad string `json:"codigo_especialidad"`
	Especialidad       string `json:"especialidad"`
}

// ListadoEspecialidadesItems decodifica el bloque "data" de la respuesta.
// MINSA lo devuelve como arreglo cuando hay resultados y como cadena (p. ej.
// "") cuando no los hay, por eso implementa UnmarshalJSON.
type ListadoEspecialidadesItems []EspecialidadMinsa

func (l *ListadoEspecialidadesItems) UnmarshalJSON(b []byte) error {
	raw := string(b)
	if raw == "null" || (len(raw) > 0 && raw[0] == '"') {
		return nil
	}
	var items []EspecialidadMinsa
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}
	*l = ListadoEspecialidadesItems(items)
	return nil
}
