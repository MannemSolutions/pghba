package hba

import (
	"fmt"
	"github.com/mannemsolutions/pghba/pkg/arg_list_comp"
	"regexp"
	"strings"
)

/*
Objects of this type can be used as a reader, getting one rule at a time.
*/

type Rules struct {
	rowNum uint
	comments Comments
	str      string
	connType ConnTypes
	database Databases
	user Users
	address Addresses
	method Method
	options Options
}

func NewRules(connTypes string, databases string, users string, addresses string, mask string, method string, options string) (Rules, error) {
	var rules Rules
	ct := arg_list_comp.NewALC(connTypes)
	for sConnType :=

	ct := NewConnType(connType)
	mtd := NewMethod(method)
	db := Database(database)
	usr := User(user)
	addr, err := NewAddress(address)
	if err != nil {
		return Rule{}, err
	}
	if mask != "" {
		err = addr.SetMask(mask)
		if err != nil {
			return Rule{}, err
		}
	}
	opts, _, err := NewOptionsFromString(options)
	if err != nil {
		return Rule{}, err
	}

	if ct == ConnTypeUnknown  || mtd == MethodUnknown {
		return Rule{}, fmt.Errorf("new Rule has an invalid connection type (%s) or method (%s)", connType, method)
	}
	return Rule{
		connType: ct,
		method: mtd,
		database: db,
		user: usr,
		address: addr,
		options: opts,
	}, nil
}

