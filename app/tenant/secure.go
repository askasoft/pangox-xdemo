package tenant

import (
	"net"

	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/net/netx"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pangox-xdemo/app/utils/tbsutil"
	"github.com/askasoft/pangox/xwa/xpwds"
)

const (
	AuthMethodPassword = "P"
	AuthMethodLDAP     = "L"
	AuthMethodSAML     = "S"
)

func (tt *Tenant) IsLDAPLogin() bool {
	return tt.ConfigValue("secure_login_method") == AuthMethodLDAP
}

func (tt *Tenant) IsSAMLLogin() bool {
	return tt.ConfigValue("secure_login_method") == AuthMethodSAML && tt.ConfigValue("secure_saml_idpmeta") != ""
}

func (tt *Tenant) SecureClientCIDRs() []*net.IPNet {
	ipnets, _ := netx.ParseCIDRs(str.Fields(tt.ConfigValue("secure_client_cidr")))
	return ipnets
}

type PasswordPolicy struct {
	xpwds.PasswordPolicy
	Locale    string
	Strengthm *linkedhashmap.LinkedHashMap[string, string]
}

func (pp *PasswordPolicy) ValidatePassword(pwd string) []string {
	vs := pp.PasswordPolicy.ValidatePassword(pwd)
	if len(vs) > 0 {
		for i, v := range vs {
			vs[i] = pp.Strengthm.SafeGet(v, v)
		}
	}
	return vs
}

func (tt *Tenant) GetPasswordPolicy(loc string) *PasswordPolicy {
	pp := &PasswordPolicy{Locale: loc}
	pp.MinLength, pp.MaxLength = num.Atoi(tt.ConfigValue("password_policy_minlen"), 8), 64
	pp.Strengths = tt.ConfigValues("password_policy_strength")
	pp.Strengthm = tbsutil.GetLinkedHashMap(loc, "config.list.password_policy_strength")
	pp.Strengthm.Set(xpwds.PASSWORD_INVALID_LENGTH, tbs.Format(loc, "error.param.pwdlen", pp.MinLength, pp.MaxLength))
	return pp
}

func (tt *Tenant) ValidatePassword(loc, pwd string) []string {
	return tt.GetPasswordPolicy(loc).ValidatePassword(pwd)
}
