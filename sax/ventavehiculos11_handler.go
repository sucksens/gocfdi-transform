package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

type VentaVehiculos11Handler struct {
	config HandlerConfig
}

func NewVentaVehiculos11Handler(config HandlerConfig) *VentaVehiculos11Handler {
	return &VentaVehiculos11Handler{config: config}
}

func (h *VentaVehiculos11Handler) ProcessVentaVehiculosElement(se xml.StartElement, decoder *xml.Decoder) (*models.VentaVehiculos11Data, error) {
	version := strings.TrimSpace(getAttrValue(se, "Version"))
	if version != "1.1" {
		return nil, errors.New("incorrect type of Venta Vehiculos, this handler only supports Venta Vehiculos version 1.1")
	}
	data := &models.VentaVehiculos11Data{
		Version:             version,
		ClaveVehicular:      getAttrValue(se, "ClaveVehicular"),
		Niv:                 getAttrValue(se, "Niv"),
		InformacionAduanera: []models.InformacionAduanera{},
		Partes:              []models.Parte{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "InformacionAduanera":
				data.InformacionAduanera = append(data.InformacionAduanera, h.transformInformacionAduanera(t, decoder))
			case "Parte":
				data.Partes = append(data.Partes, h.transformParte(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "VentaVehiculos" {
				return data, nil
			}
		}
	}
	return data, nil
}

func (h *VentaVehiculos11Handler) transformBytes(xmlBytes []byte) (*models.VentaVehiculos11Data, error) {
	return h.transformString(string(xmlBytes))
}

func (h *VentaVehiculos11Handler) transformString(xmlString string) (*models.VentaVehiculos11Data, error) {
	data := &models.VentaVehiculos11Data{
		InformacionAduanera: []models.InformacionAduanera{},
		Partes:              []models.Parte{},
	}

	decoder := xml.NewDecoder(strings.NewReader(xmlString))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "InformacionAduanera":
				data.InformacionAduanera = append(data.InformacionAduanera, h.transformInformacionAduanera(se, decoder))
			case "Parte":
				data.Partes = append(data.Partes, h.transformParte(se, decoder))
			}
		case xml.EndElement:
			if se.Name.Local == "VentaVehiculos" {
				return data, nil
			}
		}
	}
	return data, nil
}

func (h *VentaVehiculos11Handler) transformInformacionAduanera(se xml.StartElement, decoder *xml.Decoder) models.InformacionAduanera {
	return models.InformacionAduanera{
		Numero: getAttrValue(se, "Numero"),
		Fecha:  getAttrValue(se, "Fecha"),
		Aduana: getAttrValue(se, "Aduana"),
	}
}

func (h *VentaVehiculos11Handler) transformParte(se xml.StartElement, decoder *xml.Decoder) models.Parte {
	parte := models.Parte{
		NoIdentificacion:    getAttrValue(se, "NoIdentificacion"),
		Cantidad:            getAttrValue(se, "Cantidad"),
		Unidad:              getAttrValue(se, "Unidad"),
		Descripcion:         getAttrValue(se, "Descripcion"),
		ValorUnitario:       getAttrValue(se, "ValorUnitario"),
		Importe:             getAttrValue(se, "Importe"),
		InformacionAduanera: []models.InformacionAduanera{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return parte
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "InformacionAduanera" {
				parte.InformacionAduanera = append(parte.InformacionAduanera, h.transformInformacionAduanera(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "Parte" {
				return parte
			}
		}
	}
}
