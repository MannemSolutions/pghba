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
	// toMethod contains a map from all defined authentication methods as string to their individual type.
	// This is done using a variable because maps can't be constants in Go.
	toMethod = map[string]Method{
		"cert":          MethodCert,
		"gss":           MethodGss,
		"ident":         MethodIdent,
		"krb5":          MethodKrb5,
		"ldap":          MethodLdap,
		"md5":           MethodMd5,
		"pam":           MethodPam,
		"password":      MethodPassword,
		"peer":          MethodPeer,
		"radius":        MethodRadius,
		"reject":        MethodReject,
		"scram-sha-256": MethodScramSha256,
		"sspi":          MethodSspi,
		"trust":         MethodTrust,
	}
	// fromMethod contains a map from defined authentication methods to their string representation.
	fromMethod = map[Method]string{}
)

//func methods() (s []string) {
//	i := 0;
//	keys := make ([]string, 0, len(toMethod))
//	for key := range toMethod {
//		keys[i] = key
//	}
//	return keys
//}

// Based on whether str maps to something in the toMethod variable, NewMethod returns the value belonging to key str,
// otherwise return MethodUnknown
func NewMethod(str string) (m Method) {
	if m, exists := toMethod[str]; exists {
		return m
	}
	return MethodUnknown
}

// If fromMethod is uninitialized, derive it from toMethod.
// If Method m maps to a value, return it as string. Otherwise return "unknown_method"
// TODO Q: it seems strange to mix OO with normal functions?
func (m Method) String() string {
	if len(fromMethod) == 0 {
		fromMethod = make(map[Method]string)
		for s, m := range toMethod {
			fromMethod[m] = s
		}
	}
	if s, exists := fromMethod[m]; exists {
		return s
	}
	return "unknown_method"
}
