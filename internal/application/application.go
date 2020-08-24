package application

import (
	"context"
	"log"
	"net/http"
	"sync"

	config "github.com/golang-friends/members/internal/config"
	"github.com/google/go-github/v32/github"
)

type Application struct {
	Config *config.Config
	Client *http.Client
}

func (app *Application) GetConfigFromGitHub() config.Config {
	client := github.NewClient(app.Client)

	admins, _, err := client.Organizations.ListMembers(context.Background(), app.Config.Orgname, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       "admin",
	})

	if err != nil {
		panic(err)
	}

	members, _, err := client.Organizations.ListMembers(context.Background(), app.Config.Orgname, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       "member",
	})

	if err != nil {
		panic(err)
	}

	var config config.Config

	for _, admin := range admins {
		config.Admins = append(config.Admins, admin.GetLogin())
	}

	for _, member := range members {
		config.Members = append(config.Members, member.GetLogin())
	}

	return config
}

func (app *Application) Update(dryRun bool) error {
	cli := github.NewClient(app.Client)
	var wg sync.WaitGroup

	for _, admin := range app.Config.Admins {
		if dryRun {
			log.Printf("Adding %v as admin", admin)
			continue
		}

		wg.Add(1)
		// needs to capture
		admin := admin
		go func() {
			defer wg.Done()
			editMembership(cli, app.Config.Orgname, admin, "admin")
		}()
	}

	for _, member := range app.Config.Members {
		if dryRun {
			log.Printf("Adding %v as member", member)
			continue
		}

		wg.Add(1)
		// needs to capture
		member := member
		go func() {
			defer wg.Done()
			editMembership(cli, app.Config.Orgname, member, "member")
		}()
	}

	wg.Wait()
	return nil
}

func editMembership(cli *github.Client, orgname, username, role string) {
	log.Printf("Adding %v as %v", username, role)
	_, _, err := cli.Organizations.EditOrgMembership(
		context.Background(), username, orgname, &github.Membership{
			Role: github.String(role),
		})
	if err != nil {
		log.Printf("EditOrgMembership (%v) has failed for %v: %v", role, username, err)
	}
}

func NewApplication(cfg *config.Config, client *http.Client) *Application {
	return &Application{Config: cfg, Client: client}
}
