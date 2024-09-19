package types

type GitHub struct {
	Organizations []string `json:"organizations"`
	Users         []string `json:"users"`
}

func NewGitHub() *GitHub {
	return &GitHub{
		Organizations: []string{},
		Users:         []string{},
	}
}
