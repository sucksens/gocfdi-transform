package cfdi40_test

import (
	"testing"

	"github.com/sucksens/gocfdi-transform/sax"
)

func TestCFDI40Handler_UsePagos20(t *testing.T) {
	xmlStr := `
<cfdi:Comprobante Version="4.0" Serie="P" Folio="123" Fecha="2023-10-27T12:00:00" NoCertificado="00001000000500000000" SubTotal="0" Moneda="XXX" Total="0" TipoDeComprobante="P" Exportacion="01" LugarExpedicion="01000" xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:pago20="http://www.sat.gob.mx/Pagos20">
    <cfdi:Emisor Rfc="TEST010203001" Nombre="EMISOR DE PRUEBA" RegimenFiscal="601"/>
    <cfdi:Receptor Rfc="TEST010203002" Nombre="RECEPTOR DE PRUEBA" DomicilioFiscalReceptor="01000" RegimenFiscalReceptor="601" UsoCFDI="CP01"/>
    <cfdi:Conceptos>
        <cfdi:Concepto ClaveProdServ="84111506" Cantidad="1" ClaveUnidad="ACT" Descripcion="Pago" ValorUnitario="0" Importe="0" ObjetoImp="01"/>
    </cfdi:Conceptos>
    <cfdi:Complemento>
        <pago20:Pagos Version="2.0">
            <pago20:Totales MontoTotalPagos="100.00"/>
            <pago20:Pago FechaPago="2023-10-27T12:00:00" FormaDePagoP="03" MonedaP="MXN" Monto="100.00">
            </pago20:Pago>
        </pago20:Pagos>
    </cfdi:Complemento>
</cfdi:Comprobante>
`

	t.Run("Parse Pagos20 when enabled", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UsePagos20()
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(data.Pagos20) == 0 {
			t.Fatal("Expected Pagos20 data, got none")
		}

		if data.Pagos20[0].Version != "2.0" {
			t.Errorf("Expected Pagos version 2.0, got %s", data.Pagos20[0].Version)
		}

		if data.Pagos20[0].Totales.MontoTotalPagos != "100.00" {
			t.Errorf("Expected MontoTotalPagos 100.00, got %s", data.Pagos20[0].Totales.MontoTotalPagos)
		}
	})

	t.Run("Do NOT parse Pagos20 when disabled", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig())
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(data.Pagos20) > 0 {
			t.Fatal("Expected NO Pagos20 data, but got some")
		}
	})
}

func TestCFDI40Handler_ParsePagos20FromFile(t *testing.T) {
	filePath := "../recursos/cfdi40_pagos.xml"

	t.Run("Parse Pagos20 from file", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UsePagos20()
		data, err := handler.TransformFromFile(filePath)
		if err != nil {
			t.Fatalf("Unexpected error reading file %s: %v", filePath, err)
		}

		if len(data.Pagos20) == 0 {
			t.Fatal("Expected Pagos20 data, got none")
		}

		pagoData := data.Pagos20[0]

		if pagoData.Version != "2.0" {
			t.Errorf("Expected Pagos version 2.0, got %s", pagoData.Version)
		}

		// Verify Totales
		if pagoData.Totales.MontoTotalPagos != "1160.00" {
			t.Errorf("Expected MontoTotalPagos 1160.00, got %s", pagoData.Totales.MontoTotalPagos)
		}
		if pagoData.Totales.TotalTrasladosBaseIVA16 != "1000.00" {
			t.Errorf("Expected TotalTrasladosBaseIVA16 1000.00, got %s", pagoData.Totales.TotalTrasladosBaseIVA16)
		}

		// Verify Pago
		if len(pagoData.Pagos) != 1 {
			t.Fatalf("Expected 1 Pago, got %d", len(pagoData.Pagos))
		}
		pago := pagoData.Pagos[0]
		if pago.Monto != "1160.00" {
			t.Errorf("Expected Pago Monto 1160.00, got %s", pago.Monto)
		}
		if pago.MonedaP != "MXN" {
			t.Errorf("Expected Pago MonedaP MXN, got %s", pago.MonedaP)
		}

		// Verify DoctoRelacionado
		if len(pago.DoctoRelacionado) != 1 {
			t.Fatalf("Expected 1 DoctoRelacionado, got %d", len(pago.DoctoRelacionado))
		}
		docto := pago.DoctoRelacionado[0]
		if docto.IdDocumento != "00000000-0000-0000-0000-000000000001" {
			t.Errorf("Expected IdDocumento 00000000-0000-0000-0000-000000000001, got %s", docto.IdDocumento)
		}
		if docto.ImpPagado != "1160.00" {
			t.Errorf("Expected ImpPagado 1160.00, got %s", docto.ImpPagado)
		}

		// Verify ImpuestosDR
		if len(docto.ImpuestosDR) == 0 {
			t.Fatal("Expected ImpuestosDR, got none")
		}
		if len(docto.ImpuestosDR[0].TrasladosDR) == 0 {
			t.Fatal("Expected TrasladosDR, got none")
		}
		trasladoDR := docto.ImpuestosDR[0].TrasladosDR[0]
		if trasladoDR.BaseDR != "1000.00" {
			t.Errorf("Expected BaseDR 1000.00, got %s", trasladoDR.BaseDR)
		}
		if trasladoDR.ImporteDR != "160.00" {
			t.Errorf("Expected ImporteDR 160.00, got %s", trasladoDR.ImporteDR)
		}

		// Verify ImpuestosP
		if len(pago.ImpuestosP) == 0 {
			t.Fatal("Expected ImpuestosP, got none")
		}
		if len(pago.ImpuestosP[0].TrasladosP) == 0 {
			t.Fatal("Expected TrasladosP, got none")
		}
		trasladoP := pago.ImpuestosP[0].TrasladosP[0]
		if trasladoP.BaseP != "1000.00" {
			t.Errorf("Expected BaseP 1000.00, got %s", trasladoP.BaseP)
		}
		if trasladoP.ImporteP != "160.00" {
			t.Errorf("Expected ImporteP 160.00, got %s", trasladoP.ImporteP)
		}
	})
}
