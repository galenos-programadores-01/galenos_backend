package domain

import "testing"

func TestPrimeraAtencionPDF(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/galenos-pro", nil)
	ruta, ubi, err := r.PrimeraAtencionPDF(42)
	if err != nil {
		t.Fatal(err)
	}
	want := "/galenos-pro/emergencia/primera-atencion/42.pdf"
	if ruta != want || ubi != want {
		t.Fatalf("got ruta=%q ubi=%q want %q", ruta, ubi, want)
	}
}

func TestPrimeraAtencionRecetaPDF(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/galenos-pro", nil)
	ruta, ubi, err := r.PrimeraAtencionRecetaPDF(42, 98765)
	if err != nil {
		t.Fatal(err)
	}
	want := "/galenos-pro/emergencia/primera-atencion/42_98765.pdf"
	if ruta != want || ubi != want {
		t.Fatalf("got ruta=%q ubi=%q want %q", ruta, ubi, want)
	}
}

func TestPatronPersonalizado(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/docs", map[string]string{
		PatronArchivoPrimeraAtencion: "pa/{idPrimeraAtencion}/firma.pdf",
	})
	ruta, _, err := r.PrimeraAtencionPDF(10)
	if err != nil {
		t.Fatal(err)
	}
	want := "/docs/pa/10/firma.pdf"
	if ruta != want {
		t.Fatalf("got %q want %q", ruta, want)
	}
}

func TestRechazaVarsIncompletas(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/ftp", map[string]string{
		PatronArchivoPrimeraAtencion: "x/{idPrimeraAtencion}/{uuid}.pdf",
	})
	if _, _, err := r.PrimeraAtencionPDF(1); err == nil {
		t.Fatal("expected incomplete vars error")
	}
}

func TestRechazaDotDot(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/ftp", map[string]string{
		PatronArchivoPrimeraAtencion: "../etc/{idPrimeraAtencion}.pdf",
	})
	if _, _, err := r.PrimeraAtencionPDF(1); err == nil {
		t.Fatal("expected error for .. in pattern")
	}
}

func TestParseNombreArchivoLotePA(t *testing.T) {
	pa, err := ParseNombreArchivoLotePA("42.pdf")
	if err != nil || !pa.EsPrincipal || pa.IdPrimeraAtencion != 42 {
		t.Fatalf("pa: %+v err=%v", pa, err)
	}
	rec, err := ParseNombreArchivoLotePA("42_98765.pdf")
	if err != nil || rec.EsPrincipal || rec.IdPrimeraAtencion != 42 || rec.IdReceta != 98765 {
		t.Fatalf("rec: %+v err=%v", rec, err)
	}
}
