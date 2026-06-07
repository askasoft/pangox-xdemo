package mymodels

import (
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pangox-xdemo/app/models"
)

type AuditLog struct {
	models.AuditLog

	Params sqx.JSONAnyArray `gorm:"type:json" json:"params,omitempty"`
}
