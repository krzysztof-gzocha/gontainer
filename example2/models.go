package main

import (
	"github.com/gomponents/gontainer/pkg"
)

type Wallet struct {
	Value uint
}

func NewWallet(value uint, c pkg.Container) *Wallet {
	return &Wallet{Value: value}
}

type Person struct {
	Name   string
	Wallet *Wallet
}

func NewPerson(name string, wallet *Wallet) *Person {
	return &Person{Name: name, Wallet: wallet}
}
