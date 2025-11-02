<s>[{{T .Loc "sitename"}}] Notice about two-factor email authentication</s>

<div>
	<p>Thank you for using {{T .Loc "sitename"}}.</p>
	<p>A login request has been received at {{TIME .Now}}.</p>
	<p>Please continue the two-step authentication within {{.Expires}} minutes.</p>
	<p>Passcode: {{.Passcode}}</p>
	<br>
	{{HTML (T .Loc "email.signature")}}
</div>
