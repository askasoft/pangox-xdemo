package models

import (
	"time"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

const (
	ConfigStyleDefault        = ""
	ConfigStyleHidden         = "H"
	ConfigStyleChecks         = "C"
	ConfigStyleVerticalChecks = "VC"
	ConfigStyleOrderedChecks  = "OC"
	ConfigStyleRadios         = "R"
	ConfigStyleVerticalRadios = "VR"
	ConfigStyleSelect         = "S"
	ConfigStyleMultiSelect    = "MS"
	ConfigStyleTextarea       = "T"
	ConfigStyleNumeric        = "N"
	ConfigStyleDecimal        = "D"
	ConfigStyleMonth          = "TM"
	ConfigStyleDate           = "TD"
	ConfigStyleTime           = "TT"
)

type Config struct {
	Name       string    `gorm:"size:64;not null;primaryKey"`
	Value      string    `gorm:"not null"`
	Style      string    `gorm:"size:2;not null"`
	Order      int       `gorm:"not null"`
	Required   bool      `gorm:"not null"`
	Secret     bool      `gorm:"not null"`
	Viewer     string    `gorm:"size:1;not null"`
	Editor     string    `gorm:"size:1;not null"`
	Validation string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;<-:create" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime:false" json:"updated_at"`
}

func (c *Config) String() string {
	return toString(c)
}

func (c *Config) Readonly(role string) bool {
	return c.Editor < role
}

func (c *Config) MaskedValue() string {
	n := str.RuneCount(c.Value)
	switch {
	case n > 8:
		return str.Left(c.Value, 4) + str.Repeat("*", n-4)
	case n > 6:
		return str.Left(c.Value, 3) + str.Repeat("*", n-3)
	case n > 4:
		return str.Left(c.Value, 2) + str.Repeat("*", n-2)
	default:
		return str.Repeat("*", n)
	}
}

func (c *Config) DisplayValue() string {
	if c.Value != "" {
		if c.Secret {
			return c.MaskedValue()
		}

		switch c.Style {
		case ConfigStyleNumeric:
			return num.Comma(num.Atol(c.Value))
		case ConfigStyleDecimal:
			return num.Comma(num.Atof(c.Value))
		}
	}
	return c.Value
}

func (c *Config) Values() []string {
	return str.FieldsByte(c.Value, '\t')
}

func (c *Config) IsSameMeta(n *Config) bool {
	return c.Name == n.Name &&
		c.Style == n.Style && c.Order == n.Order &&
		c.Required == n.Required && c.Secret == n.Secret &&
		c.Viewer == n.Viewer && c.Editor == n.Editor &&
		c.Validation == n.Validation
}
