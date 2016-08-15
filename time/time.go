package time

import (
	"encoding/xml"
	"time"
)

// RFC3339Time allows RFC 3339 compliant XML un/marshaling.
type RFC3339Time struct {
	time.Time
}

	return nil
}

// MarshalXML satisfies the xml.Marshaler interface.
func (c *RFC3339Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(c.Time.Format(time.RFC3339), start)
}

// MarshalXMLAttr satsfies the xml.MarshalerAttr interface.
func (c *RFC3339Time) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{
		Name:  name,
		Value: c.Time.Format(time.RFC3339),
	}
	return attr, nil
}

// UnmarshalXML satisfies the xml.Unmarshaler interface.
func (c *RFC3339Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start) {
		return err
	}
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}

// UnmarshalXMLAttr satisfies the xml.UnmarshalerAttr interface.
func (c *RFC3339Time) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}
