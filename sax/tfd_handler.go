package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
	"github.com/sucksens/gocfdi-transform/models"
)

// TFD11Handler handles parsing of Timbre Fiscal Digital 1.1 complement.
type TFD11Handler struct {
	config HandlerConfig
}

// NewTFD11Handler creates a new TFD11Handler.
func NewTFD11Handler(config HandlerConfig) *TFD11Handler {
	return &TFD11Handler{config: config}
}

// TransformFromBytes parses a TFD 1.1 XML byte slice.
func (h *TFD11Handler) TransformFromBytes(xmlBytes []byte) (interface{}, error) {
	return h.TransformFromString(string(xmlBytes))
}

// TransformFromString parses a TFD 1.1 XML string.
func (h *TFD11Handler) TransformFromString(xmlStr string) (*models.TFD11, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

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
			if se.Name.Local == "TimbreFiscalDigital" {
				return h.transformTFD(se)
			}
		}
	}

	return nil, errors.New("TimbreFiscalDigital element not found")
}

func (h *TFD11Handler) transformTFD(se xml.StartElement) (*models.TFD11, error) {
	version := getAttrValue(se, "Version")
	if version != "1.1" {
		return nil, errors.New("incorrect type of TFD, this handler only supports TFD version 1.1")
	}

	return &models.TFD11{
		Version:          version,
		NoCertificadoSAT: getAttrValue(se, "NoCertificadoSAT"),
		UUID:             strings.ToUpper(getAttrValue(se, "UUID")),
		FechaTimbrado:    getAttrValue(se, "FechaTimbrado"),
		RfcProvCert:      getAttrValue(se, "RfcProvCertif"),
		SelloCFD:         helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "SelloCFD")),
		SelloSAT:         helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "SelloSAT")),
	}, nil
}
