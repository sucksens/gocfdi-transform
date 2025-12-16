// Package models contiene las estructuras de datos para los CFDI y los complementos.
package models

// CFDI40Data es la estructura de datos para el CFDI 4.0
// Incluye el CFDI40 y los TFD11 si los hay.
type CFDI40Data struct {
	CFDI40  CFDI40        `json:"cfdi40"`
	TFD11   []TFD11       `json:"tfd11,omitempty"`
	Pagos20 []Pagos20Data `json:"pagos20,omitempty"`
}

// CFDI40 es la estructura de datos para el CFDI 4.0
type CFDI40 struct {
	Version           string            `json:"version"`
	Serie             string            `json:"serie"`
	Folio             string            `json:"folio"`
	Fecha             string            `json:"fecha"`
	NoCertificado     string            `json:"no_certificado"`
	SubTotal          string            `json:"subtotal"`
	Descuento         string            `json:"descuento"`
	Total             string            `json:"total"`
	Moneda            string            `json:"moneda"`
	TipoCambio        string            `json:"tipo_cambio"`
	TipoComprobante   string            `json:"tipo_comprobante"`
	MetodoPago        string            `json:"metodo_pago"`
	FormaPago         string            `json:"forma_pago"`
	CondicionesPago   string            `json:"condiciones_pago"`
	LugarExpedicion   string            `json:"lugar_expedicion"`
	Exportacion       string            `json:"exportacion"`
	Sello             string            `json:"sello"`
	Certificado       string            `json:"certificado"`
	Confirmacion      string            `json:"confirmacion"`
	Emisor            Emisor40          `json:"emisor"`
	Receptor          Receptor40        `json:"receptor"`
	Conceptos         []Concepto40      `json:"conceptos"`
	Impuestos         Impuestos         `json:"impuestos"`
	Complementos      string            `json:"complementos"`
	Addendas          string            `json:"addendas"`
	CFDIsRelacionados []CFDIRelacionado `json:"cfdis_relacionados,omitempty"`
}

// Emisor40 es la estructura de datos para el emisor del CFDI 4.0
type Emisor40 struct {
	RFC              string `json:"rfc"`
	Nombre           string `json:"nombre"`
	RegimenFiscal    string `json:"regimen_fiscal"`
	FacAtrAdquirente string `json:"fac_atr_adquirente"`
}

// Receptor40 es la estructura de datos para el receptor del CFDI 4.0
type Receptor40 struct {
	RFC                     string `json:"rfc"`
	Nombre                  string `json:"nombre"`
	DomicilioFiscalReceptor string `json:"domicilio_fiscal_receptor"`
	ResidenciaFiscal        string `json:"residencia_fiscal"`
	NumRegIdTrib            string `json:"num_reg_id_trib"`
	RegimenFiscalReceptor   string `json:"regimen_fiscal_receptor"`
	UsoCFDI                 string `json:"uso_cfdi"`
}

// Concepto40 es la estructura de datos para un concepto del CFDI 4.0
type Concepto40 struct {
	ClaveProdServ    string              `json:"clave_prod_serv"`
	NoIdentificacion string              `json:"no_identificacion"`
	Cantidad         string              `json:"cantidad"`
	ClaveUnidad      string              `json:"clave_unidad"`
	Unidad           string              `json:"unidad"`
	Descripcion      string              `json:"descripcion"`
	ValorUnitario    string              `json:"valor_unitario"`
	Importe          string              `json:"importe"`
	Descuento        string              `json:"descuento"`
	ObjetoImp        string              `json:"objeto_imp"`
	Terceros         Terceros            `json:"terceros,omitempty"`
	Traslados        []TrasladoConcepto  `json:"traslados,omitempty"`
	Retenciones      []RetencionConcepto `json:"retenciones,omitempty"`
}

// Terceros es la estructura de datos para los terceros del CFDI 4.0
type Terceros struct {
	Nombre          string `json:"nombre,omitempty"`
	RFC             string `json:"rfc,omitempty"`
	DomicilioFiscal string `json:"domicilioFiscal,omitempty"`
	RegimenFiscal   string `json:"regimenFiscal,omitempty"`
}

// Impuestos es la estructura de datos para los impuestos del CFDI 4.0
type Impuestos struct {
	TotalImpuestosTrasladados string      `json:"total_impuestos_trasladados"`
	TotalImpuestosRetenidos   string      `json:"total_impuestos_retenidos"`
	Traslados                 []Traslado  `json:"traslados"`
	Retenciones               []Retencion `json:"retenciones"`
}

// Traslado es la estructura de datos para un traslado de impuesto del CFDI 4.0
type Traslado struct {
	Base       string `json:"base"`
	Impuesto   string `json:"impuesto"`
	TipoFactor string `json:"tipo_factor"`
	TasaOCuota string `json:"tasa_o_cuota"`
	Importe    string `json:"importe"`
}

// Retencion es la estructura de datos para una retención de impuesto del CFDI 4.0
type Retencion struct {
	Impuesto string `json:"impuesto"`
	Importe  string `json:"importe"`
}

// TrasladoConcepto es la estructura de datos para un traslado de impuesto en un concepto del CFDI 4.0
type TrasladoConcepto struct {
	Base       string `json:"base"`
	Impuesto   string `json:"impuesto"`
	TipoFactor string `json:"tipo_factor"`
	TasaOCuota string `json:"tasa_o_cuota"`
	Importe    string `json:"importe"`
}

// RetencionConcepto es la estructura de datos para una retención de impuesto en un concepto del CFDI 4.0
type RetencionConcepto struct {
	Impuesto string `json:"impuesto"`
	Importe  string `json:"importe"`
}

// CFDIRelacionado es la estructura de datos para un CFDI relacionado del CFDI 4.0
type CFDIRelacionado struct {
	UUID         string `json:"uuid"`
	TipoRelacion string `json:"tipo_relacion"`
}
