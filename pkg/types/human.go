package types

import (
	"github.com/google/uuid"
	"slices"
	"sort"
	"strings"
)

type Human struct {
	UUID      string
	Names     []string
	Emails    []string
	Passwords []string
	Usernames []string
	Roles     []string
	URLs      []string
	Phones    []string
}

func NewHuman() *Human {
	return &Human{UUID: uuid.New().String()}
}

func (h *Human) AddName(name string) {
	h.Names = addToSlice(h.Names, name)
}

func (h *Human) AddEmail(email string) {
	email = strings.ToLower(strings.Trim(email, "“” "))
	h.Emails = addToSlice(h.Emails, email)
}

func (h *Human) AddPassword(password string) {
	h.Passwords = addToSlice(h.Passwords, password)
}

func (h *Human) AddUsername(username string) {
	username = strings.ToLower(strings.Trim(username, "“” "))
	h.Usernames = addToSlice(h.Usernames, username)
}

func (h *Human) AddRole(roleName string) {
	roleName = strings.ToLower(strings.Trim(roleName, "“” "))
	h.Roles = addToSlice(h.Roles, roleName)
}

func (h *Human) AddURL(urlToAdd string) {
	urlToAdd = strings.ToLower(strings.Trim(urlToAdd, "“” "))
	h.URLs = addToSlice(h.URLs, urlToAdd)
}

func (h *Human) AddPhone(phone string) {
	h.Phones = addToSlice(h.Phones, phone)
}

func (h *Human) Merge(otherHuman *Human) {
	for _, n := range otherHuman.Names {
		h.AddName(n)
	}

	for _, e := range otherHuman.Emails {
		h.AddEmail(e)
	}

	for _, n := range otherHuman.Passwords {
		h.AddPassword(n)
	}

	for _, n := range otherHuman.Usernames {
		h.AddUsername(n)
	}

	for _, n := range otherHuman.Roles {
		h.AddRole(n)
	}

	for _, n := range otherHuman.URLs {
		h.AddURL(n)
	}

	for _, n := range otherHuman.Phones {
		h.AddPhone(n)
	}
}

func addToSlice(slice []string, item string) []string {
	item = strings.TrimSpace(item)
	if item == "" || slices.Contains(slice, item) {
		return slice
	}
	slice = append(slice, item)

	sort.Strings(slice)
	return slice
}
