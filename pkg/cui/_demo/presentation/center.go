package main

import "github.com/malivvan/vv/pkg/cui"

// Center returns a new primitive which shows the provided primitive in its
// center, given the provided primitive's size.
func Center(width, height int, p cui.Primitive) cui.Primitive {
	subFlex := cui.NewFlex()
	subFlex.SetDirection(cui.FlexRow)
	subFlex.AddItem(cui.NewBox(), 0, 1, false)
	subFlex.AddItem(p, height, 1, true)
	subFlex.AddItem(cui.NewBox(), 0, 1, false)

	flex := cui.NewFlex()
	flex.AddItem(cui.NewBox(), 0, 1, false)
	flex.AddItem(subFlex, width, 1, true)
	flex.AddItem(cui.NewBox(), 0, 1, false)

	return flex
}
