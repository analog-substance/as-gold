package gold

import (
	"context"
	"errors"
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/go-git/go-git/v5"
	"log"
	"path"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/analog-substance/as-gold/pkg/types"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v56/github"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const GitHubFolderPath = "github.com"

type SolidGold struct {
	*types.Group
	*types.GitHub `json:"github"`
	githubToken   string
	githubClient  *github.Client
}

func NewSolidGold() *SolidGold {
	return &SolidGold{
		Group:  types.NewGroup(),
		GitHub: types.NewGitHub(),
	}
}

func FromJSONFile(filepath string) *SolidGold {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		s := NewSolidGold()
		err := s.ToJSONFile(filepath)
		util.CheckErr(err)
	}

	jsonFile, err := os.Open(filepath)
	util.CheckErr(err)
	defer jsonFile.Close()
	byteAr, err := io.ReadAll(jsonFile)
	util.CheckErr(err)
	return FromJSON(byteAr)
}

func FromJSON(byteAr []byte) *SolidGold {
	var solidGold SolidGold
	err := json.Unmarshal(byteAr, &solidGold)
	util.CheckErr(err)
	if solidGold.GitHub == nil {
		solidGold.GitHub = types.NewGitHub()
	}
	return &solidGold
}

func (s *SolidGold) ToJSON() []byte {
	goldJSON, err := json.Marshal(s)
	util.CheckErr(err)
	return goldJSON
}

func (s *SolidGold) SetGitHubAccessToken(accessToken string) {
	s.githubToken = accessToken
}

func (s *SolidGold) GithubClient() *github.Client {
	if s.githubClient == nil {
		s.githubClient = github.NewClient(nil)

		if s.githubToken != "" {
			s.githubClient = s.githubClient.WithAuthToken(s.githubToken)
		}
	}

	return s.githubClient
}

func (s *SolidGold) ToJSONFile(goldFile string) error {
	return os.WriteFile(goldFile, s.ToJSON(), 0644)
}

func (s *SolidGold) ProcessPath(path string) *types.Group {

	gitRepos, err := findGitDirs(path)
	util.CheckErr(err)

	if _, err := os.Stat(fmt.Sprintf("%s/.git", path)); !os.IsNotExist(err) {
		gitRepos = append(gitRepos, path)
	}

	for _, repoPath := range gitRepos {
		s.ProcessRepo(err, repoPath)
	}

	s.Group.MergeDuplicate()
	return s.Group
}

func (s *SolidGold) ProcessRepo(err error, path string) {
	r, err := git.PlainOpen(path)
	util.CheckErr(err)

	cIter, err := r.Log(&git.LogOptions{All: true})
	util.CheckErr(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		s.Group.FindOrCreateHuman(c.Author.Name, c.Author.Email)
		s.Group.FindOrCreateHuman(c.Committer.Name, c.Committer.Email)
		return nil
	})
}

func (s *SolidGold) Merge(primaryID string, otherIDs ...string) {
	s.Group.MergeIDs(primaryID, otherIDs...)
}

func (s *SolidGold) ConsumeBreachFiles(breachFilePaths ...string) {
	for _, breachFilePath := range breachFilePaths {
		tsvFile, err := os.Open(breachFilePath)
		util.CheckErr(err)
		reader := csv.NewReader(tsvFile)
		reader.Comma = '\t'
		reader.FieldsPerRecord = -1

		tsvData, err := reader.ReadAll()
		util.CheckErr(err)

		tsvFile.Close()

		for _, cols := range tsvData {
			email := ""
			password := ""
			if len(cols) >= 1 {
				email = cols[0]
			}

			if len(cols) >= 2 {
				password = cols[1]
			}

			human := s.Group.FindOrCreateHumanByEmail(email)
			human.AddPassword(password)
		}
	}
}

func (s *SolidGold) ConsumeGophishFiles(gophishFilePaths ...string) {
	for _, gophishFile := range gophishFilePaths {
		csvFile, err := os.Open(gophishFile)
		util.CheckErr(err)
		reader := csv.NewReader(csvFile)
		//reader.Comma = ','
		reader.FieldsPerRecord = -1

		csvData, err := reader.ReadAll()
		util.CheckErr(err)

		csvFile.Close()

		for _, cols := range csvData {

			firstName := cols[0]
			lastName := cols[1]
			email := cols[2]
			position := cols[3]

			if strings.ToLower(firstName) == "First Name" {
				continue
			}

			human := s.Group.FindOrCreateHumanByEmail(email)
			human.AddName(fmt.Sprintf("%s %s", strings.TrimSpace(firstName), strings.TrimSpace(lastName)))
			human.AddRole(position)

		}
	}
}

func findGitDirs(dir string) ([]string, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".git") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (s *SolidGold) ConsumeGithubOrgs(includeMembers bool, orgs ...string) {
	s.GitHub.Organizations = util.UniqueSlice(append(s.GitHub.Organizations, orgs...))
	// maybe save

	orgs = util.UniqueSlice(orgs)
	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{Type: "sources"}
		if org == "" {
			continue
		}

		for {
			repos, resp, err := s.GithubClient().Repositories.ListByOrg(context.Background(), org, opt)

			if err != nil {
				log.Println("Error encountered while getting an organization's repo list", err)
				continue
			}

			for _, repo := range repos {
				gitCloneURL(path.Join(GitHubFolderPath, repo.GetFullName()), repo.GetSSHURL())
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		usernames := []string{}
		if includeMembers {
			memberListOpt := &github.ListMembersOptions{}
			for {

				members, resp, _ := s.GithubClient().Organizations.ListMembers(context.Background(), org, memberListOpt)
				for _, member := range members {
					if member.GetEmail() != "" {
						h := s.FindOrCreateHuman(member.GetLogin(), member.GetEmail())

						if member.GetName() != "" {
							h.AddName(member.GetName())
						}
					}

					usernames = append(usernames, member.GetLogin())
				}

				if resp.NextPage == 0 {
					break
				}

				opt.Page = resp.NextPage
			}
			s.ConsumeGithubUsers(false, usernames...)
		}

	}
}

func (s *SolidGold) ConsumeGithubUsers(includeOrgs bool, users ...string) {
	s.GitHub.Users = util.UniqueSlice(append(s.GitHub.Users, users...))

	users = util.UniqueSlice(users)
	for _, user := range users {
		if user == "" {
			continue
		}
		opt := &github.RepositoryListOptions{Type: "owner"}
		for {

			repos, resp, err := s.GithubClient().Repositories.List(context.Background(), user, opt)
			if err != nil {
				log.Println("Error encountered while getting a users repo list", err)
				continue
			}

			for _, repo := range repos {
				if !*repo.Fork {
					gitCloneURL(path.Join(GitHubFolderPath, repo.GetFullName()), repo.GetCloneURL())
				}
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		if includeOrgs {
			orgs, _, err := s.GithubClient().Organizations.List(context.Background(), user, nil)
			if err != nil {
				log.Println("Error encountered while getting a user's org list", err)
				continue
			}

			orgNames := []string{}
			for _, org := range orgs {
				orgNames = append(orgNames, org.GetName())
			}
			s.ConsumeGithubOrgs(false, orgNames...)
		}
	}
}

func (s *SolidGold) UpdateGithub(includeMembers, includeOrgs bool) {
	orgOrUsers, _ := os.ReadDir(GitHubFolderPath)
	for _, orgOrUser := range orgOrUsers {
		if orgOrUser.IsDir() {
			orgOrUserPath := path.Join(GitHubFolderPath, orgOrUser.Name())
			repos, _ := os.ReadDir(orgOrUserPath)
			for _, repo := range repos {
				if repo.IsDir() {
					repoPath := path.Join(orgOrUserPath, repo.Name())
					repoInst, err := git.PlainOpen(repoPath)
					if err != nil {
						log.Printf("Error encountered while attempting to open repo %s: %s", repoPath, err)
						continue
					}
					w, err := repoInst.Worktree()
					if err != nil {
						log.Printf("Error encountered while attempting to get work tree %s: %s", repoPath, err)
						continue
					}
					err = w.Pull(&git.PullOptions{RemoteName: "origin"})
					if err != nil {
						if !errors.Is(err, git.NoErrAlreadyUpToDate) {
							log.Printf("Error encountered while attempting to pull repo %s: %s", repoPath, err)
						}
						continue
					}
					log.Printf("Successfully updated repo %s", repoPath)
				}
			}
		}
	}

	s.ConsumeGithubUsers(includeOrgs, s.GitHub.Users...)
	s.ConsumeGithubOrgs(includeMembers, s.GitHub.Organizations...)
}

func gitCloneURL(path, repoURL string) {
	repoExists, err := exists(path)
	if err != nil {
		return
	}
	if !repoExists {
		log.Printf("cloning %s to %s", repoURL, path)
		_, _ = git.PlainClone(path, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
