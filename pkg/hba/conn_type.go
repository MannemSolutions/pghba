package hba

type ConnType int

const (
	ConnTypeLocal ConnType = iota
	ConnTypeHost
	ConnTypeHostSsl
	ConnTypeHostNoSsl
	ConnTypeHostGssEnc
	ConnTypeHostNoGssEnc
	ConnTypeUnknown
)

func NewConnType(str string) (ct ConnType) {
	toConnType := map[string]ConnType{
		"local": ConnTypeLocal,
		"host": ConnTypeHost,
		"hostssl": ConnTypeHostSsl,
		"hostnossl": ConnTypeHostNoSsl,
		"hostgssenc": ConnTypeHostGssEnc,
		"hostnogssenc": ConnTypeHostNoGssEnc,
	}
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