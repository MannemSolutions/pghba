package hba

import (
	"github.com/mannemsolutions/pghba/pkg/arg_list_comp"
)

/*
Objects of this type can be used as a reader, getting one rule at a time.
*/

type Rules struct {
	connTypes *arg_list_comp.ALC
	databases arg_list_comp.ALC
	users arg_list_comp.ALC
	addresses arg_list_comp.ALC
	method Method
	options Options
}

func NewRules(connTypes string, databases string, users string, addresses string, mask string, method string, options string) (Rules, error) {
	opts, _, err := NewOptionsFromString(options)
	if err != nil {
		return Rules{}, err
	}
	rules := Rules{
		method: NewMethod(method),
		options: opts,
		connTypes: arg_list_comp.NewALC(connTypes).Unique(),
	}
	for _, sConnType := range .ToList() {
		rules.connTypes = append(rules.connTypes, NewConnType(sConnType))
	}

	for _, sDatabase := range arg_list_comp.NewALC(databases).ToSortedArray().ToList() {
		rules.databases = append(rules.databases, Database(sDatabase))
	}

	for _, sUser := range arg_list_comp.NewALC(users).ToSortedArray().ToList() {
		rules.users = append(rules.users, User(sUser))
	}
	for _, sAddresses := range arg_list_comp.NewALC(addresses).ToSortedArray().ToList() {
		address, err := NewAddress(sAddresses)
		if err != nil {
			return Rules{}, err
		}
		err = address.SetMask(mask)
		if err != nil {
			return Rules{}, err
		}
		rules.addresses = append(rules.addresses, address)
	}
	return rules, nil
}

