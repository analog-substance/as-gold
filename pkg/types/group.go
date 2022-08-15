package types

import (
	"sort"
	"strings"
)

type Group struct {
	Members []*Human
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) FindOrCreateHuman(name, email string) *Human {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	human := g.FindHumanByEmail(email)
	if human != nil {
		human.AddName(name)
		return human
	}

	if strings.Contains(name, " ") {
		human := g.FindHumanByName(name)
		if human != nil {
			human.AddEmail(email)
			return human
		}
	}

	newHuman := NewHuman()
	newHuman.AddEmail(email)
	newHuman.AddName(name)

	g.Members = append(g.Members, newHuman)
	return newHuman
}

func (g *Group) FindOrCreateHumanByEmail(email string) *Human {
	email = strings.TrimSpace(email)

	human := g.FindHumanByEmail(email)
	if human != nil {
		return human
	}

	newHuman := NewHuman()
	newHuman.AddEmail(email)

	g.Members = append(g.Members, newHuman)
	return newHuman
}

func (g *Group) FindHumanByEmail(email string) *Human {
	email = strings.TrimSpace(email)

	for _, human := range g.Members {
		for _, e := range human.Emails {
			if strings.EqualFold(email, e) {
				return human
			}
		}
	}

	return nil
}

func (g *Group) FindHumanByName(name string) *Human {
	name = strings.TrimSpace(name)

	for _, human := range g.Members {
		for _, n := range human.Names {
			if strings.EqualFold(name, n) {
				return human
			}
		}
	}

	return nil
}

func (g *Group) MergeDuplicate() {

	newGroup := Group{}

	for _, human1 := range g.Members {
		nh := newGroup.FindOrCreateHuman(human1.Names[0], human1.Emails[0])
		nh.Merge(human1)
	}

	g.Members = newGroup.Members
	g.Sort()
}

func (g *Group) Append(human *Human) {
	g.Members = append(g.Members, human)
}

func (g *Group) Sort() {

	sort.Slice(g.Members, func(i, j int) bool {
		if len(g.Members[i].Names) > 0 {
			if len(g.Members[j].Names) > 0 {
				return g.Members[i].Names[0] < g.Members[j].Names[0]
			}
			// if we have a name and they don't, we should be first.
			return true
		}

		if len(g.Members[i].Emails) > 0 {
			if len(g.Members[j].Emails) > 0 {
				return g.Members[i].Emails[0] < g.Members[j].Emails[0]
			}
			// if we have an email and they don't, we should be first.
			return true
		}
		// dont care....
		return false

	})
}

func (g *Group) IndexOf(id string) int {
	for index, human1 := range g.Members {
		if human1.UUID == id {
			return index
		}
	}

	return -1
}

func (g *Group) MemberByID(id string) *Human {
	for _, human1 := range g.Members {
		if human1.UUID == id {
			return human1
		}
	}

	return nil
}

func (g *Group) MergeIDs(primaryID string, otherIDs ...string) {
	primaryHuman := g.MemberByID(primaryID)
	if primaryHuman != nil {
		for _, otherID := range otherIDs {
			other := g.MemberByID(otherID)
			primaryHuman.Merge(other)
		}
		g.RemoveIDs(otherIDs...)
	}
}

func (g *Group) RemoveIDs(otherIDs ...string) {
	for _, otherID := range otherIDs {
		index := g.IndexOf(otherID)
		if index == -1 {
			continue
		}

		newMembers := g.Members[:index]
		newMembers = append(newMembers, g.Members[index+1:]...)
		g.Members = newMembers
	}
}

func (g *Group) FindWithEmailDomains(domains ...string) *Group {
	subGroup := &Group{}
	for _, human1 := range g.Members {
	Email:
		for _, email := range human1.Emails {
			for _, domain := range domains {
				if strings.HasSuffix(strings.ToLower(email), strings.ToLower(domain)) {
					subGroup.Append(human1)
					break Email
				}
			}
		}
	}

	return subGroup
}

func (g *Group) FindWithString(searchStrings ...string) *Group {
	subGroup := &Group{}
Email:
	for _, human1 := range g.Members {
		for _, email := range human1.Emails {
			for _, search := range searchStrings {
				if strings.Contains(email, strings.ToLower(search)) {
					subGroup.Append(human1)
					continue Email
				}
			}
		}
		for _, name := range human1.Names {
			for _, search := range searchStrings {
				if strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
					subGroup.Append(human1)
					continue Email
				}
			}
		}
	}

	subGroup.Sort()
	return subGroup
}

func (g *Group) FindWithPasswords() *Group {
	subGroup := &Group{}
	for _, human1 := range g.Members {
		if len(human1.Passwords) > 0 {
			subGroup.Append(human1)
		}
	}

	return subGroup
}
