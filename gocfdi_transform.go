// Package gocfdi_transform nos da handlers para transformar un CFDI en una estructura de datos facil de manejar en go.
//
// This library supports CFDI versions 3.2, 3.3, and 4.0, along with common complements
// like Timbre Fiscal Digital (TFD), Pagos, and VentaVehiculos.
//
// Example usage:
//
//	handler := sax.NewCFDI40Handler(sax.NewDefaultConfig())
//	data, err := handler.TransformFromFile("invoice.xml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("UUID: %s\n", data.TFD11[0].UUID)
package gocfdi_transform

import (
	"github.com/sucksens/gocfdi-transform/models"
	"github.com/sucksens/gocfdi-transform/sax"
)

// Version of the library
const Version = "0.1.0"

// Re-export main types for convenience
type (
	// HandlerConfig is the configuration for handlers
	HandlerConfig = sax.HandlerConfig

	// CFDI40Data is the parsed CFDI 4.0 data
	CFDI40Data = models.CFDI40Data

	// Pagos20Data is the parsed Pagos 2.0 data
	Pagos20Data = models.Pagos20Data
)

// NewDefaultConfig creates a new default handler configuration.
func NewDefaultConfig() HandlerConfig {
	return sax.NewDefaultConfig()
}

// NewCFDI40Handler creates a new CFDI 4.0 handler.
func NewCFDI40Handler(config HandlerConfig) *sax.CFDI40Handler {
	return sax.NewCFDI40Handler(config)
}

// NewPagos20Handler creates a new Pagos 2.0 handler.
func NewPagos20Handler(config HandlerConfig) *sax.Pagos20Handler {
	return sax.NewPagos20Handler(config)
}
