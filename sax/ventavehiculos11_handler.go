package sax

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

type VentaVehiculos11Handler struct {
	*BaseHandler
	builder *ModelBuilder
}

func NewVentaVehiculos11Handler(config HandlerConfig) *VentaVehiculos11Handler {
	return &VentaVehiculos11Handler{
		BaseHandler: NewBaseHandler(config),
		builder:     NewModelBuilder(config),
	}
}

// ProcessVentaVehiculosElement procesa el elemento VentaVehiculos
func (h *VentaVehiculos11Handler) ProcessVentaVehiculosElement(se xml.StartElement, decoder *xml.Decoder) (*models.VentaVehiculos11Data, error) {
	// Validar version
	if err := h.ValidateVersion(se, "1.1"); err != nil {
		return nil, err
	}
	// Crear estructura de datos
	data := &models.VentaVehiculos11Data{
		Version:             "1.1",
		ClaveVehicular:      h.builder.ExtractString(se, "ClaveVehicular"),
		Niv:                 h.builder.ExtractString(se, "Niv"),
		InformacionAduanera: []models.InformacionAduanera{},
		Partes:              []models.Parte{},
	}

	err := ProcessChildElements(decoder, "VentaVehiculos", func(childSE xml.StartElement, chilDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "InformacionAduanera":
			data.InformacionAduanera = append(data.InformacionAduanera, h.transformInformacionAduanera(childSE))
		case "Parte":
			data.Partes = append(data.Partes, h.transformParte(childSE, chilDecoder))
		}
		return nil
	})

	return data, err
}

// TransformFromString parses a Venta Vehiculos 1.1 XML string.
func (h *VentaVehiculos11Handler) TransformString(xmlString string) (*models.VentaVehiculos11Data, error) {
	//Crear el decoder XML para parsear la cadena
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	// Buscar el elemento VentaVehiculos
	se, err := FindElement(decoder, "VentaVehiculos")
	if err != nil {
		if err == io.EOF {
			return &models.VentaVehiculos11Data{
				InformacionAduanera: []models.InformacionAduanera{},
				Partes:              []models.Parte{},
			}, nil
		}
		return nil, err
	}
	return h.ProcessVentaVehiculosElement(*se, decoder)
}

// transformInformacionAduanera transforma el elemento InformacionAduanera
func (h *VentaVehiculos11Handler) transformInformacionAduanera(se xml.StartElement) models.InformacionAduanera {
	return models.InformacionAduanera{
		Numero: h.builder.ExtractString(se, "Numero"),
		Fecha:  h.builder.ExtractString(se, "Fecha"),
		Aduana: h.builder.ExtractString(se, "Aduana"),
	}
}

// transformParte transforma el elemento Parte
func (h *VentaVehiculos11Handler) transformParte(se xml.StartElement, decoder *xml.Decoder) models.Parte {
	parte := models.Parte{
		NoIdentificacion:    h.builder.ExtractString(se, "NoIdentificacion"),
		Cantidad:            h.builder.ExtractString(se, "Cantidad"),
		Unidad:              h.builder.ExtractString(se, "Unidad"),
		Descripcion:         h.builder.ExtractString(se, "Descripcion"),
		ValorUnitario:       h.builder.ExtractString(se, "ValorUnitario"),
		Importe:             h.builder.ExtractString(se, "Importe"),
		InformacionAduanera: []models.InformacionAduanera{},
	}

	ProcessChildElements(decoder, "Parte", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		if childSE.Name.Local == "InformacionAduanera" {
			parte.InformacionAduanera = append(parte.InformacionAduanera, h.transformInformacionAduanera(childSE))
		}
		return nil
	})

	return parte
}
