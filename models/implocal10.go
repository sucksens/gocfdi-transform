package models

type ImpuestosLocales10 struct {
	Version            string             `json:"version"`
	TotaldeRetenciones string             `json:"total_de_retenciones"`
	TotaldeTraslados   string             `json:"total_de_traslados"`
	RetencionesLocales []RetencionLocal10 `json:"retenciones_locales,omitempty"`
	TrasladosLocales   []TrasladoLocal10  `json:"traslados_locales,omitempty"`
}

type RetencionLocal10 struct {
	ImpLocRetenido  string `json:"imp_loc_retenido"`
	TasadeRetencion string `json:"tasa_de_retencion"`
	Importe         string `json:"importe"`
}

type TrasladoLocal10 struct {
	ImpLocTrasladado string `json:"imp_loc_trasladado"`
	TasadeTraslado   string `json:"tasa_de_traslado"`
	Importe          string `json:"importe"`
}
