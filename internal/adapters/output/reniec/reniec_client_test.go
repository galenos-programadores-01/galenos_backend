package reniec

import (
	"reflect"
	"testing"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

func valorCampo(d domain.ReniecDatos, campo string) string {
	return reflect.ValueOf(d).FieldByName(campo).String()
}

// Fixture: respuesta real de obtenerDatosCompletos (49 campos).
var fixtureCompleto = []string{
	"0000", "", "76114703", "5", "JARA", "BENITES", "SIN DATOS", "JULIA DAYANA",
	"92", "33", "14", "01", "06", "000", "AMERICA", "PERU", "LIMA", "LIMA",
	"COMAS", "SIN DATOS", "1", "22", "2", "14", "01", "06", "LIMA", "LIMA",
	"COMAS", "07/03/2006", "RENE", "NATALIA LUZMILA", "09/03/2011", "07/08/2023",
	"", "SIN DATOS", "CA D1 MZ. F1 LT.4", "SIN DATOS", "SIN DATOS", "SIN DATOS",
	"SIN DATOS", "CASUARINAS DE COLLIQUE", "SIN DATOS", "SIN DATOS", "SIN DATOS",
	"SIN DATOS", "07", "foto", "firma",
}

// Fixture: respuesta real de obtenerDatosBasicos (23 campos).
var fixtureBasico = []string{
	"0000", "", "JARA", "BENITES", "SIN DATOS", "JULIA DAYANA", "92", "33",
	"14", "01", "06", "000", "AMERICA", "PERU", "LIMA", "LIMA", "COMAS",
	"SIN DATOS", "CA D1 MZ. F1 LT.4", "2", "07/03/2006", "07/08/2023", "76114703",
}

func TestInterpretarCompleto(t *testing.T) {
	d := interpretarDatos(fixtureCompleto, "completo")
	esperado := map[string]string{
		"Sexo":                   "FEMENINO",
		"EstadoCivil":            "SOLTERO",
		"ApellidoPaterno":        "JARA",
		"ApellidoMaterno":        "BENITES",
		"Nombres":                "JULIA DAYANA",
		"FechaNacimiento":        "2006-03-07",
		"Departamento":           "LIMA",
		"Provincia":              "LIMA",
		"Distrito":               "COMAS",
		"Direccion":              "CA D1 MZ. F1 LT.4",
		"Ubigeo":                 "140106",
		"DepartamentoNacimiento": "LIMA",
		"ProvinciaNacimiento":    "LIMA",
		"DistritoNacimiento":     "COMAS",
		"NombrePadre":            "RENE",
		"NombreMadre":            "NATALIA LUZMILA",
	}
	for campo, want := range esperado {
		if got := valorCampo(d, campo); got != want {
			t.Errorf("%s = %q, se esperaba %q", campo, got, want)
		}
	}
}

func TestInterpretarBasico(t *testing.T) {
	d := interpretarDatos(fixtureBasico, "basico")
	esperado := map[string]string{
		"Sexo":            "FEMENINO",
		"ApellidoPaterno": "JARA",
		"Nombres":         "JULIA DAYANA",
		"FechaNacimiento": "2006-03-07",
		"Departamento":    "LIMA",
		"Distrito":        "COMAS",
		"Direccion":       "CA D1 MZ. F1 LT.4",
		"Ubigeo":          "140106",
	}
	for campo, want := range esperado {
		if got := valorCampo(d, campo); got != want {
			t.Errorf("%s = %q, se esperaba %q", campo, got, want)
		}
	}
}

func TestEsError(t *testing.T) {
	if esError(fixtureCompleto) {
		t.Error("una respuesta con codigo 0000 no debe considerarse error")
	}
	if !esError([]string{"5114", "EL DOCUMENTO INGRESADO NO EXISTE EN RENIEC"}) {
		t.Error("el codigo 5114 debe considerarse error")
	}
	if !esError([]string{"9000", "x"}) {
		t.Error("el codigo 9000 debe considerarse error")
	}
	if esError([]string{"", ""}) {
		t.Error("sin codigo de error no debe considerarse error")
	}
}

func TestSexoDesdeCodigo(t *testing.T) {
	casos := map[string]string{
		"1": "MASCULINO", "2": "FEMENINO", "M": "MASCULINO",
		"F": "FEMENINO", "": "", "9": "",
	}
	for in, want := range casos {
		if got := sexoDesdeCodigo(in); got != want {
			t.Errorf("sexoDesdeCodigo(%q) = %q, se esperaba %q", in, got, want)
		}
	}
}
