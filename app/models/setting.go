package models

import (
	"time"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

const (
	SettingStyleDefault        = ""
	SettingStyleHidden         = "H"
	SettingStyleChecks         = "C"
	SettingStyleVerticalChecks = "VC"
	SettingStyleOrderedChecks  = "OC"
	SettingStyleRadios         = "R"
	SettingStyleVerticalRadios = "VR"
	SettingStyleSelect         = "S"
	SettingStyleMultiSelect    = "MS"
	SettingStyleTextarea       = "T"
	SettingStyleNumeric        = "N"
	SettingStyleDecimal        = "D"
	SettingStyleMonth          = "TM"
	SettingStyleDate           = "TD"
	SettingStyleTime           = "TT"
)

type SettingItem struct {
	Name    string
	Value   string
	Display string
}

type SettingGroup struct {
	Name  string     `json:"name"`
	Items []*Setting `json:"items"`
}

type SettingCategory struct {
	Name   string          `json:"name"`
	Groups []*SettingGroup `json:"groups"`
}

type Setting struct {
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

func (s *Setting) String() string {
	return toString(s)
}

func (s *Setting) Readonly(role string) bool {
	return s.Editor < role
}

func (s *Setting) MaskedValue() string {
	n := str.RuneCount(s.Value)
	switch {
	case n > 8:
		return str.Left(s.Value, 4) + str.Repeat("*", n-4)
	case n > 6:
		return str.Left(s.Value, 3) + str.Repeat("*", n-3)
	case n > 4:
		return str.Left(s.Value, 2) + str.Repeat("*", n-2)
	default:
		return str.Repeat("*", n)
	}
}

func (s *Setting) DisplayValue() string {
	if s.Value != "" {
		if s.Secret {
			return s.MaskedValue()
		}

		switch s.Style {
		case SettingStyleNumeric:
			return num.Comma(num.Atol(s.Value))
		case SettingStyleDecimal:
			return num.Comma(num.Atof(s.Value))
		}
	}
	return s.Value
}

func (s *Setting) Values() []string {
	return str.FieldsByte(s.Value, '\t')
}

func (s *Setting) IsSameMeta(n *Setting) bool {
	return s.Name == n.Name &&
		s.Style == n.Style && s.Order == n.Order &&
		s.Required == n.Required && s.Secret == n.Secret &&
		s.Viewer == n.Viewer && s.Editor == n.Editor &&
		s.Validation == n.Validation
}
