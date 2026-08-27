package domain

// AuthClaims son los datos verificados extraídos de un token Bearer JWT
// válido.
type AuthClaims struct {
	Subject    string
	IdEmpleado int
}

type Menu struct {
	IdListGrupo int    `json:"idListGrupo"`
	Texto       string `json:"texto"`
	KeyIconWeb  string `json:"keyIconWeb"`
	ClaveWeb    string `json:"claveWeb"`
	Indice      int    `json:"indice"`
	Estado      bool   `json:"estado"`
	NroSubMenu  int    `json:"nroSubMenu"`
}

type MenuPermiso struct {
	Opciones    string `json:"opciones"`
	Indice      int    `json:"indice"`
	Texto       string `json:"texto"`
	Menu        string `json:"menu"`
	IdListGrupo int    `json:"idListGrupo"`
	KeyIconWeb  string `json:"keyIconWeb"`
	Estado      bool   `json:"estado"`
	ClaveWeb    string `json:"claveWeb"`
	Agregar     bool   `json:"agregar"`
	Modificar   bool   `json:"modificar"`
	Eliminar    bool   `json:"eliminar"`
}

type AuthMenus struct {
	Menus    []Menu        `json:"menus"`
	Permisos []MenuPermiso `json:"permisos"`
}

type UserProfile struct {
	IdEmpleado      int    `json:"idEmpleado"`
	Username        string `json:"username"`
	Nombres         string `json:"nombres"`
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	NombreCompleto  string `json:"nombreCompleto"`
	Foto            string `json:"foto"`
	Rol             string `json:"rol"`
}
