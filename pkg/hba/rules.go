package hba

import (
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
)

/*
Objects of this type can be used as a reader, getting one rule at a time.
*/

type Rules struct {
	connTypes gnrtr.Gnrtr
	databases gnrtr.Gnrtr
	users gnrtr.Gnrtr
	addresses gnrtr.Gnrtr
	mask      string
	current Rule
}

//func getUnique(s string) gnrtr {
//	alc := gnrtr.NewGnrtr(s)
//	if alc == nil {
//		return gnrtr.StrToGnrtr(s)
//	}
//
//}

func NewRules(rowNum int, connTypes string, databases string, users string, addresses string, mask string, method string, options string) (Rules, error) {
	opts, _, err := NewOptionsFromString(options)
	if err != nil {
		return Rules{}, err
	}
	rules := Rules{
		connTypes: gnrtr.NewGnrtr(connTypes).Unique(),
		databases: gnrtr.NewGnrtr(databases).Unique(),
		users: gnrtr.NewGnrtr(users).Unique(),
		addresses: gnrtr.NewGnrtr(addresses).Unique(),
		mask: mask,
		current: Rule{
			rowNum: rowNum,
			connType: ConnTypeUnknown,
			method: NewMethod(method),
			options: opts,
		},
	}
	rules.databases.Next()
	rules.users.Next()
	rules.addresses.Next()
	return rules, nil
}

func (rs Rules) Current() (next Rule, err error){
	if rs.current.connType == ConnTypeUnknown {
		rs.current.connType = NewConnType(rs.connTypes.Current())
	}
	rs.current.database = Database(rs.databases.Current())
	rs.current.user = User(rs.users.Current())
	if rs.current.address.Unset() {

		if rs.current.address, err = NewAddress(rs.addresses.Current()); err != nil {
			return Rule{}, err
		}
		if err = rs.current.address.SetMask(rs.mask); err != nil {
			return Rule{}, err
		}
	}

	return rs.current.Clone(), nil
}

func (rs Rules) Next() (next Rule, done bool, err error) {
	if ct, done := rs.connTypes.Next(); ! done {
		rs.current.connType = NewConnType(ct)
		if next, err = rs.Current(); err != nil {
			return Rule{}, false, err
		} else {
			return next, false, nil
		}
	}

	if d, done := rs.databases.Next(); ! done {
		rs.current.database = Database(d)
		if next, err = rs.Current(); err != nil {
			return Rule{}, false, err
		} else {
			return next, false, nil
		}
	}

	if u, done := rs.users.Next(); ! done {
		rs.current.user = User(u)
		if next, err = rs.Current(); err != nil {
			return Rule{}, false, err
		} else {
			return next, false, nil
		}
	}

	if a, done := rs.addresses.Next(); ! done {
		if rs.current.address, err = NewAddress(a); err != nil {
			return Rule{}, false, err
		}
		if err = rs.current.address.SetMask(rs.mask); err != nil {
			return Rule{}, false, err
		}
		if next, err = rs.Current(); err != nil {
			return Rule{}, false, err
		} else {
			return next, false, nil
		}
	}
	return Rule{}, true, nil
}

