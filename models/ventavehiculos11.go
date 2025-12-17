package models

// VentaVehiculos11Data representa los datos de un complemento Venta Vehículos 1.1.
type VentaVehiculos11Data struct {
	Version             string                `json:"version"`
	ClaveVehicular      string                `json:"clave_vehicular"`
	Niv                 string                `json:"niv"`
	InformacionAduanera []InformacionAduanera `json:"informacion_aduanera"`
	Partes              []Parte               `json:"partes"`
}

// InformacionAduanera representa la información aduanera de un complemento Venta Vehículos 1.1.
type InformacionAduanera struct {
	Numero string `json:"numero"`
	Fecha  string `json:"fecha"`
	Aduana string `json:"aduana"`
}

// Parte representa un parte de un complemento Venta Vehículos 1.1.
type Parte struct {
	NoIdentificacion    string                `json:"no_identificacion"`
	Cantidad            string                `json:"cantidad"`
	Unidad              string                `json:"unidad"`
	Descripcion         string                `json:"descripcion"`
	ValorUnitario       string                `json:"valor_unitario"`
	Importe             string                `json:"importe"`
	InformacionAduanera []InformacionAduanera `json:"informacion_aduanera"`
}
