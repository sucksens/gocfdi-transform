package models

// Pagos20Data es la estructura de datos para el Pagos complement version 2.0.
type Pagos20Data struct {
	Version string    `json:"version"`
	Totales Totales20 `json:"totales"`
	Pagos   []Pago20  `json:"pago"`
}

// Totales20 es la estructura de datos para la sección de totales en Pagos 2.0.
type Totales20 struct {
	TotalRetencionesIVA         string `json:"total_retenciones_iva"`
	TotalRetencionesISR         string `json:"total_retenciones_isr"`
	TotalRetencionesIEPS        string `json:"total_retenciones_ieps"`
	TotalTrasladosBaseIVA16     string `json:"total_traslados_base_iva_16"`
	TotalTrasladosImpuestoIVA16 string `json:"total_traslados_impuesto_iva_16"`
	TotalTrasladosBaseIVA8      string `json:"total_traslados_base_iva_8"`
	TotalTrasladosImpuestoIVA8  string `json:"total_traslados_impuesto_iva_8"`
	TotalTrasladosBaseIVA0      string `json:"total_traslados_base_iva_0"`
	TotalTrasladosImpuestoIVA0  string `json:"total_traslados_impuesto_iva_0"`
	TotalTrasladosBaseIVAExento string `json:"total_traslados_base_iva_exento"`
	MontoTotalPagos             string `json:"monto_total_pagos"`
}

// Pago20 es la estructura de datos para un pago individual en Pagos 2.0.
type Pago20 struct {
	FechaPago        string               `json:"fecha_pago"`
	FormaDePagoP     string               `json:"forma_de_pago_p"`
	MonedaP          string               `json:"moneda_p"`
	TipoCambioP      string               `json:"tipo_cambio_p"`
	Monto            string               `json:"monto"`
	NumOperacion     string               `json:"num_operacion"`
	RfcEmisorCtaOrd  string               `json:"rfc_emisor_cta_ord"`
	NomBancoOrdExt   string               `json:"nom_banco_ord_ext"`
	CtaOrdenante     string               `json:"cta_ordenante"`
	RfcEmisorCtaBen  string               `json:"rfc_emisor_cta_ben"`
	CtaBeneficiario  string               `json:"cta_beneficiario"`
	TipoCadPago      string               `json:"tipo_cad_pago"`
	CertPago         string               `json:"cert_pago"`
	CadPago          string               `json:"cad_pago"`
	SelloPago        string               `json:"sello_pago"`
	DoctoRelacionado []DoctoRelacionado20 `json:"docto_relacionado"`
	ImpuestosP       []ImpuestosP         `json:"impuestos_p"`
}

// DoctoRelacionado20 es la estructura de datos para un documento relacionado en Pagos 2.0.
type DoctoRelacionado20 struct {
	IdDocumento      string        `json:"id_documento"`
	Serie            string        `json:"serie"`
	Folio            string        `json:"folio"`
	MonedaDR         string        `json:"moneda_dr"`
	EquivalenciaDR   string        `json:"equivalencia_dr"`
	NumParcialidad   string        `json:"num_parcialidad"`
	ImpSaldoAnt      string        `json:"imp_saldo_ant"`
	ImpPagado        string        `json:"imp_pagado"`
	ImpSaldoInsoluto string        `json:"imp_saldo_insoluto"`
	ObjetoImpDR      string        `json:"objecto_imp_dr"`
	ImpuestosDR      []ImpuestosDR `json:"impuestos_dr"`
}

// ImpuestosDR es la estructura de datos para los impuestos de un documento relacionado en Pagos 2.0.
type ImpuestosDR struct {
	RetencionesDR []ImpuestoDRItem `json:"retenciones_dr"`
	TrasladosDR   []ImpuestoDRItem `json:"traslados_dr"`
}

// ImpuestoDRItem es la estructura de datos para un ítem de impuesto en los impuestos de un documento relacionado en Pagos 2.0.
type ImpuestoDRItem struct {
	BaseDR       string `json:"base_dr"`
	ImpuestoDR   string `json:"impuesto_dr"`
	TipoFactorDR string `json:"tipo_factor_dr"`
	TasaOCuotaDR string `json:"tasa_o_cuota_dr"`
	ImporteDR    string `json:"importe_dr"`
}

// ImpuestosP es la estructura de datos para los impuestos en un pago individual en Pagos 2.0.
type ImpuestosP struct {
	RetencionesP []RetencionP `json:"retenciones_p"`
	TrasladosP   []TrasladoP  `json:"traslados_p"`
}

// RetencionP es la estructura de datos para una retención de impuesto en un pago individual en Pagos 2.0.
type RetencionP struct {
	ImpuestoP string `json:"impuesto_p"`
	ImporteP  string `json:"importe_p"`
}

// TrasladoP es la estructura de datos para un traslado de impuesto en un pago individual en Pagos 2.0.
type TrasladoP struct {
	BaseP       string `json:"base_p"`
	ImpuestoP   string `json:"impuesto_p"`
	TipoFactorP string `json:"tipo_factor_p"`
	TasaOCuotaP string `json:"tasa_o_cuota_p"`
	ImporteP    string `json:"importe_p"`
}
