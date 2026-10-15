package domain

import "testing"

func TestPrimeraAtencionPDF(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/ftp", nil)
	ruta, ubi, err := r.PrimeraAtencionPDF(3619840, "abc-uuid")
	if err != nil {
		t.Fatal(err)
	}
	want := "/ftp/primera-atencion/3619840/abc-uuid.pdf"
	if ruta != want || ubi != want {
		t.Fatalf("got ruta=%q ubi=%q want %q", ruta, ubi, want)
	}
}

func TestPatronPersonalizado(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/docs", map[string]string{
		PatronArchivoPrimeraAtencion: "pa/{id}/firmas/{uuid}.{ext}",
	})
	ruta, _, err := r.PrimeraAtencionPDF(10, "u1")
	if err != nil {
		t.Fatal(err)
	}
	want := "/docs/pa/10/firmas/u1.pdf"
	if ruta != want {
		t.Fatalf("got %q want %q", ruta, want)
	}
}

func TestRechazaDotDot(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/ftp", map[string]string{
		PatronArchivoPrimeraAtencion: "../etc/{uuid}.pdf",
	})
	if _, _, err := r.PrimeraAtencionPDF(1, "x"); err == nil {
		t.Fatal("expected error for .. in pattern")
	}
}

func TestPrimeraAtencionPDFNombrado(t *testing.T) {
	r := NuevaRutasArchivoRemoto("/ftp", nil)
	ruta, ubi, err := r.PrimeraAtencionPDFNombrado(3619840, "PA_3619840_EX_R12.pdf")
	if err != nil {
		t.Fatal(err)
	}
	want := "/ftp/primera-atencion/3619840/PA_3619840_EX_R12.pdf"
	if ruta != want || ubi != want {
		t.Fatalf("got ruta=%q ubi=%q want %q", ruta, ubi, want)
	}
}
