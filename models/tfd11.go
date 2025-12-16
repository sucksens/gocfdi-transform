package models

// TFD11 es la estructura de datos para el Timbre Fiscal Digital version 1.1
type TFD11 struct {
	Version          string `json:"version"`
	NoCertificadoSAT string `json:"no_certificado_sat"`
	UUID             string `json:"uuid"`
	FechaTimbrado    string `json:"fecha_timbrado"`
	RfcProvCert      string `json:"rfc_prov_cert"`
	SelloCFD         string `json:"sello_cfd"`
	SelloSAT         string `json:"sello_sat"`
}
