package api

import (
	"net/http"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func MyIP(c *xin.Context) {
	c.String(http.StatusOK, c.ClientIP())
}

func MyHeader(c *xin.Context) {
	c.IndentedJSON(http.StatusOK, c.Request.Header)
}

func MyGet(c *xin.Context) {
	c.IndentedJSON(http.StatusOK, c.Querys())
}

func MyPost(c *xin.Context) {
	c.IndentedJSON(http.StatusOK, c.PostForms())
}

func MyPets(c *xin.Context) {
	q := str.Strip(c.Query("q"))

	if q == "" {
		c.JSON(http.StatusOK, []any{})
		return
	}

	pqa := &args.PetQueryArg{Name: q}
	pqa.Limit = 20

	tt := tenant.Get(c)
	pets, err := tt.FindPets(app.SDB(), pqa)
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}

	c.JSON(http.StatusOK, pets)
}
