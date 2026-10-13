package usecase

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	estadoBandejaEnCurso  = "En curso"
	estadoBandejaAtendido = "Atendido"
)

// Solo IdAtencion real (no IdCuentaAtencion: pueden diferir).
var clavesIdAtencionBandeja = []string{
	"IdAtencion",
	"idAtencion",
}

var clavesEstadoBandeja = []string{
	"Estado",
	"estado",
	"EstadoAtencion",
	"estadoAtencion",
	"DescEstado",
}

func idsAtencionDesdeBandeja(items []map[string]any) []int {
	seen := make(map[int]struct{}, len(items))
	ids := make([]int, 0, len(items))
	for _, row := range items {
		id := idAtencionDesdeFila(row)
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func enriquecerEstadoBandeja(items []map[string]any, conBorrador, firmados map[int]bool) {
	for _, row := range items {
		id := idAtencionDesdeFila(row)
		if id <= 0 {
			continue
		}
		if firmados[id] || estadoBandejaEsAtendido(row) {
			row["Estado"] = estadoBandejaAtendido
			continue
		}
		if !conBorrador[id] {
			continue
		}
		row["Estado"] = estadoBandejaEnCurso
	}
}

func idAtencionDesdeFila(row map[string]any) int {
	for _, key := range clavesIdAtencionBandeja {
		if id := enteroDesdeValor(row[key]); id > 0 {
			return id
		}
	}
	for key, value := range row {
		for _, want := range clavesIdAtencionBandeja {
			if strings.EqualFold(key, want) {
				if id := enteroDesdeValor(value); id > 0 {
					return id
				}
			}
		}
	}
	return 0
}

func estadoBandejaEsAtendido(row map[string]any) bool {
	raw := textoDesdeFila(row, clavesEstadoBandeja...)
	if raw == "" {
		return false
	}
	lower := strings.ToLower(raw)
	return strings.Contains(lower, "atend")
}

func textoDesdeFila(row map[string]any, keys ...string) string {
	for _, key := range keys {
		if text := textoDesdeValor(row[key]); text != "" {
			return text
		}
	}
	for key, value := range row {
		for _, want := range keys {
			if strings.EqualFold(key, want) {
				if text := textoDesdeValor(value); text != "" {
					return text
				}
			}
		}
	}
	return ""
}

func textoDesdeValor(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case []byte:
		return strings.TrimSpace(string(v))
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
}

func enteroDesdeValor(value any) int {
	switch v := value.(type) {
	case nil:
		return 0
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return n
	case []byte:
		n, err := strconv.Atoi(strings.TrimSpace(string(v)))
		if err != nil {
			return 0
		}
		return n
	default:
		n, err := strconv.Atoi(strings.TrimSpace(fmt.Sprintf("%v", v)))
		if err != nil {
			return 0
		}
		return n
	}
}
