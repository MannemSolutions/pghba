package hba

type ConnType int
type ConnTypes []ConnType

const (
	ConnTypeLocal ConnType = iota
	ConnTypeHost
	ConnTypeHostSsl
	ConnTypeHostNoSsl
	ConnTypeHostGssEnc
	ConnTypeHostNoGssEnc
	ConnTypeUnknown
)

var (
	toConnType = map[string]ConnType{
		"local":        ConnTypeLocal,
		"host":         ConnTypeHost,
		"hostssl":      ConnTypeHostSsl,
		"hostnossl":    ConnTypeHostNoSsl,
		"hostgssenc":   ConnTypeHostGssEnc,
		"hostnogssenc": ConnTypeHostNoGssEnc,
	}
	fromConnType = map[ConnType]string{}
	allConnTypes []string
)

func AllConnTypes() (all []string){
	if len(allConnTypes) == 0 {
		for ct := range toConnType {
			allConnTypes = append(allConnTypes, ct)
		}
	}
	return allConnTypes
}

func NewConnType(str string) (ct ConnType) {
	if ct, exists := toConnType[str]; exists {
		return ct
	}
	return ConnTypeUnknown
}

func (ct ConnType) Compare(other ConnType) int {
	if ct < other {
		return -1
	} else if ct > other {
		return 1
	} else {
		return 0
	}
}

func (ct ConnType) String() string {
	if len(fromConnType) == 0 {
		fromConnType = make(map[ConnType]string)
		for s, ct := range toConnType {
			fromConnType[ct] = s
		}
	}
	if s, exists := fromConnType[ct]; exists {
		return s
	}
	return "unknown_conn_type"
}
