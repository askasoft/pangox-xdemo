package settings

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/doc/csvx"
	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func loadSettingList(c *xin.Context, actor string) []*models.Setting {
	tt := tenant.FromCtx(c)
	au := tenant.AuthUser(c)

	settings, err := tt.ListSettingsByRole(app.SDB(), actor, au.Role)
	if err != nil {
		panic(err)
	}

	return settings
}

func disableSettingSuperSecret(c *xin.Context, settings []*models.Setting) {
	au := tenant.AuthUser(c)

	if au.IsSuper() {
		for _, stg := range settings {
			stg.Secret = false
		}
	}
}

func buildSettingCategories(locale string, settings []*models.Setting) []*models.SettingCategory {
	scs := []*models.SettingCategory{}

	sks := tbs.GetBundle(locale).GetSection("setting.category").Keys()
	for _, sk := range sks {
		sgs := []*models.SettingGroup{}
		gks := str.Fields(tbs.GetText(locale, "setting.category."+sk))
		for _, gk := range gks {
			items := []*models.Setting{}
			for _, stg := range settings {
				if str.StartsWith(stg.Name, gk) {
					items = append(items, stg)
				}
			}
			if len(items) > 0 {
				sgs = append(sgs, &models.SettingGroup{Name: gk, Items: items})
			}
		}
		if len(sgs) > 0 {
			scs = append(scs, &models.SettingCategory{Name: sk, Groups: sgs})
		}
	}

	return scs
}

func getSettingItemList(locale string, name string) *linkedhashmap.LinkedHashMap[string, string] {
	value := tbs.GetText(locale, "setting.list."+name)
	if value == "" {
		return nil
	}

	m := &linkedhashmap.LinkedHashMap[string, string]{}
	if err := m.UnmarshalJSON(str.UnsafeBytes(value)); err != nil {
		panic(fmt.Errorf("invalid setting list [%s] %s: %w", locale, name, err))
	}
	return m
}

func bindSettingLists(c *xin.Context, h xin.H, settings []*models.Setting) {
	lists := map[string]any{}

	for _, stg := range settings {
		list := getSettingItemList(c.Locale, stg.Name)
		if list != nil {
			lists[stg.Name] = list
		}
	}

	h["Lists"] = lists
}

func SettingIndex(c *xin.Context) {
	settings := loadSettingList(c, "viewer")

	disableSettingSuperSecret(c, settings)

	scs := buildSettingCategories(c.Locale, settings)

	h := middles.H(c)
	h["Settings"] = scs
	bindSettingLists(c, h, settings)

	c.HTML(http.StatusOK, "admin/settings/settings", h)
}

func SettingSave(c *xin.Context) {
	tt := tenant.FromCtx(c)

	settings := loadSettingList(c, "editor")

	usettings := checkPostSettings(c, settings)
	if len(c.Errors) > 0 {
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	if len(usettings) > 0 {
		detail := buildSettingDetails(c, settings, usettings)
		if !saveSettings(c, usettings, models.AL_SETTINGS_UPDATE, detail) {
			return
		}
		tt.PurgeSettings()
	}

	c.JSON(http.StatusOK, xin.H{"success": tbs.GetText(c.Locale, "success.saved")})
}

func buildSettingDetails(c *xin.Context, settings []*models.Setting, usettings []*models.Setting) string {
	scs := buildSettingCategories(c.Locale, settings)

	ads := linkedhashmap.NewLinkedHashMap[string, any]()
	for _, sc := range scs {
		scn := tbs.GetText(c.Locale, "setting.category.label."+sc.Name)
		for _, sg := range sc.Groups {
			sgn := tbs.GetText(c.Locale, "setting.group.label."+sg.Name)
			for _, ci := range sg.Items {
				i := asg.IndexFunc(usettings, func(stg *models.Setting) bool {
					return stg.Name == ci.Name
				})
				if i < 0 {
					continue
				}

				stg := usettings[i]

				sin := scn + " / " + sgn + " / " + tbs.GetText(c.Locale, "setting."+stg.Name)

				var siv any = stg.DisplayValue()

				lm := getSettingItemList(c.Locale, stg.Name)
				if lm != nil && !lm.IsEmpty() {
					if stg.IsMultiple() {
						vs := stg.Values()
						lbs := make([]string, 0, len(vs))
						for _, v := range vs {
							lbs = append(lbs, lm.SafeGet(v, v))
						}
						siv = lbs
					} else {
						siv = lm.SafeGet(stg.Value, stg.Value)
					}
				}

				ads.Set(sin, siv)
			}
		}
	}

	return jsonx.Stringify(ads)
}

func validateSetting(c *xin.Context, stg *models.Setting) bool {
	v := &stg.Value

	switch stg.Style {
	case models.SettingStyleNumeric:
		*v = str.RemoveByte(str.Strip(*v), ',')
		if *v != "" && !str.IsNumeric(*v) {
			c.AddError(&args.ParamError{
				Param:   stg.Name,
				Label:   tbs.GetText(c.Locale, "setting."+stg.Name, stg.Name),
				Message: tbs.GetText(c.Locale, "error.param.numeric"),
			})
			return false
		}
	case models.SettingStyleDecimal:
		*v = str.RemoveByte(str.Strip(*v), ',')
		if *v != "" && !str.IsDecimal(*v) {
			c.AddError(&args.ParamError{
				Param:   stg.Name,
				Label:   tbs.GetText(c.Locale, "setting."+stg.Name, stg.Name),
				Message: tbs.GetText(c.Locale, "error.param.decimal"),
			})
			return false
		}
	case models.SettingStyleTextarea:
		// we keep space for textarea value
	default:
		*v = str.Strip(*v)
	}

	validation := ""
	if stg.Required {
		validation = "required"
	}
	if stg.Validation != "" {
		validation += str.If(stg.Required, ",", "omitempty,") + stg.Validation
	}
	if validation != "" {
		var vv any
		switch stg.Style {
		case models.SettingStyleNumeric:
			vv = num.Atol(*v)
		case models.SettingStyleDecimal:
			vv = num.Atof(*v)
		default:
			vv = v
		}

		if err := app.VAD.Field(stg.Name, vv, validation); err != nil {
			args.AddBindErrors(c, err, "setting.")
			return false
		}
	}

	if *v == "" {
		return true
	}

	lm := getSettingItemList(c.Locale, stg.Name)
	if lm != nil && !lm.IsEmpty() {
		var ok bool

		if stg.IsMultiple() {
			vs := str.FieldsByte(*v, '\t')
			ok = lm.ContainsAll(vs...)
		} else {
			ok = lm.Contains(*v)
		}

		if !ok {
			c.AddError(&args.ParamError{
				Param:   stg.Name,
				Label:   tbs.GetText(c.Locale, "setting."+stg.Name, stg.Name),
				Message: tbs.GetText(c.Locale, "error.param.invalid"),
			})
			return false
		}
	}

	return true
}

func checkPostSettings(c *xin.Context, settings []*models.Setting) (usettings []*models.Setting) {
	var vs []string
	var v string
	var ok bool

	for _, stg := range settings {
		if stg.IsMultiple() {
			vs, ok = c.GetPostFormArray(stg.Name)
			if ok {
				vs = str.RemoveEmpties(asg.Clone(vs))
				v = str.Join(vs, "\t")
			}
		} else {
			v, ok = c.GetPostForm(stg.Name)
		}

		if !ok || v == stg.Value || v == stg.DisplayValue() {
			// skip unknown or unmodified value
			continue
		}

		stg.Value = v
		usettings = append(usettings, stg)
	}

	for _, ustg := range usettings {
		validateSetting(c, ustg)
	}

	return
}

func saveSettings(c *xin.Context, settings []*models.Setting, action string, detail string) bool {
	tt := tenant.FromCtx(c)
	au := tenant.AuthUser(c)

	err := app.SDB().Transaction(func(tx *sqlx.Tx) error {
		if err := tt.SaveSettingsByRole(tx, au, settings, c.Locale); err != nil {
			return err
		}
		return tt.AddAuditLog(tx, c, action, detail)
	})
	if err == nil {
		return true
	}

	c.AddError(err)

	sc := gog.If(schema.IsUnsavedSettingItemsError(err), http.StatusBadRequest, http.StatusInternalServerError)
	c.JSON(sc, middles.E(c))
	return false
}

func SettingExport(c *xin.Context) {
	settings := loadSettingList(c, "editor")

	disableSettingSuperSecret(c, settings)

	scs := buildSettingCategories(c.Locale, settings)

	c.SetAttachmentHeader("settings.csv")
	_, _ = c.Writer.WriteString(string(iox.BOM))

	if err := exportSettings(c.Writer, c.Locale, scs); err != nil {
		c.Logger.Error(err)
	}
}

func exportSettings(w io.Writer, locale string, scs []*models.SettingCategory) error {
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	defer cw.Flush()

	if err := cw.Write([]string{"Name", "Value", "Display"}); err != nil {
		return err
	}

	for _, sc := range scs {
		scn := tbs.GetText(locale, "setting.category.label."+sc.Name)
		for _, sg := range sc.Groups {
			sgn := tbs.GetText(locale, "setting.group.label."+sg.Name)
			for _, ci := range sg.Items {
				disp := fmt.Sprintf("%s / %s / %s", scn, sgn, tbs.GetText(locale, "setting."+ci.Name, ci.Name))
				if err := cw.Write([]string{ci.Name, ci.DisplayValue(), disp}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func SettingImport(c *xin.Context) {
	mfh, err := c.FormFile("file")
	if err != nil {
		err = tbs.Error(c.Locale, "csv.error.required")
		c.AddError(err)
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	uf, err := mfh.Open()
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}
	defer uf.Close()

	var csvstgs []*models.SettingItem
	if err := csvx.ScanReader(uf, &csvstgs); err != nil {
		err = tbs.Error(c.Locale, "csv.error.data")
		c.AddError(err)
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	settings := loadSettingList(c, "editor")

	usettings := checkCsvSettings(c, settings, csvstgs)
	if len(c.Errors) > 0 {
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	if len(usettings) > 0 {
		detail := buildSettingDetails(c, settings, usettings)
		if !saveSettings(c, usettings, models.AL_SETTINGS_IMPORT, detail) {
			return
		}

		tenant.FromCtx(c).PurgeSettings()
	}

	c.JSON(http.StatusOK, xin.H{"success": tbs.GetText(c.Locale, "success.imported")})
}

func checkCsvSettings(c *xin.Context, settings []*models.Setting, csvstgs []*models.SettingItem) (usettings []*models.Setting) {
	stgmaps := map[string]*models.Setting{}
	for _, stg := range settings {
		stgmaps[stg.Name] = stg
	}

	for _, ci := range csvstgs {
		stg, ok := stgmaps[ci.Name]
		if !ok {
			msg := tbs.Format(c.Locale, "setting.import.invalid", ci.Name)
			c.AddError(errors.New(msg))
			continue
		}

		// drop '\r' (because csv reader drop '\r')
		ci.Value = str.RemoveByte(ci.Value, '\r')
		stg.Value = str.RemoveByte(stg.Value, '\r')

		if ci.Value == stg.Value || ci.Value == stg.DisplayValue() {
			// skip unmodified value
			continue
		}

		stg.Value = ci.Value
		usettings = append(usettings, stg)
	}

	for _, ustg := range usettings {
		validateSetting(c, ustg)
	}

	return
}
