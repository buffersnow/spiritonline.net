package web

import "encoding/xml"

type Context map[string]any

func (m Context) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Context"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range m {
		if err := e.EncodeElement(v, xml.StartElement{
			Name: xml.Name{Local: k},
		}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

func (m *Context) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	result := make(Context)

	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			var value string
			if err := d.DecodeElement(&value, &elem); err != nil {
				return err
			}
			result[elem.Name.Local] = value

		case xml.EndElement:
			if elem.Name == start.Name {
				*m = result
				return nil
			}
		}
	}
}

func (d Details) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var sysErr string
	if d.Err != nil {
		sysErr = d.Err.Error()
	}

	out := baseTransmittableError{
		Version: d.Version,
		Code:    d.Code,
		Message: d.Message,
		Error:   sysErr,
		Context: d.Context,
	}

	return e.EncodeElement(out, start)
}

func (e *Details) ContextXML() string {
	if e.Context == nil {
		return ""
	}

	out, err := xml.MarshalIndent(Context(e.Context), "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(out)
}
