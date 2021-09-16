package hba

type Method int

const (
	MethodCert Method = iota
	MethodGss
	MethodIdent
	MethodKrb5
	MethodLdap
	MethodMd5
	MethodPam
	MethodPassword
	MethodPeer
	MethodRadius
	MethodReject
	MethodScramSha256
	MethodSspi
	MethodTrust
	MethodUnknown
)

var (
	toMethod = map[string]Method{
		"cert": MethodCert,
		"gss": MethodGss,
		"ident": MethodIdent,
		"krb5": MethodKrb5,
		"ldap": MethodLdap,
		"md5": MethodMd5,
		"pam": MethodPam,
		"password": MethodPassword,
		"peer": MethodPeer,
		"radius": MethodRadius,
		"reject": MethodReject,
		"scram-sha-256": MethodScramSha256,
		"sspi": MethodSspi,
		"trust": MethodTrust,
	}
)

func methods() (s []string) {
	i := 0;
	keys := make ([]string, 0, len(toMethod))
	for key := range toMethod {
		keys[i] = key
	}
	return keys
}

func NewMethod(str string) (m Method) {
	if m, exists := toMethod[str]; exists {
		return m
	}
	return MethodUnknown
}
