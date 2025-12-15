package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/FredrikMWold/radix-tui/internal/radix"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/equinor/radix-cli/pkg/cache"
	radixconfig "github.com/equinor/radix-cli/pkg/config"
)

type SelectedApplication string

func SelectApplication(application string) tea.Cmd {
	return func() tea.Msg {
		return SelectedApplication(application)
	}
}

// Applications contains loaded apps and the context they were loaded for
type Applications struct {
	Apps    []string
	Context string
}

// ApplicationsError is sent when loading applications fails
type ApplicationsError struct {
	Err     error
	Context string
}

// AuthWaiting indicates we expect to perform interactive login.
type AuthWaiting struct{}

// AuthOK indicates cached provider/token is available; proceed to load apps.
type AuthOK struct{}

// CheckAuth decides if we should show auth-waiting or start loading immediately.
func CheckAuth() tea.Cmd {
	return func() tea.Msg {
		authCacheFilename := fmt.Sprintf("%s/auth.json", radixconfig.RadixConfigDir)
		global := cache.New(authCacheFilename, "global")
		if prov, ok := global.GetItem("auth_provider_type"); !ok || prov != "msal_interactive" {
			return AuthWaiting{}
		}
		// Check msal cache content under the interactive namespace
		msal := cache.New(authCacheFilename, "msal_interactive")
		if content, ok := msal.GetItem("msal"); !ok || len(content) == 0 {
			return AuthWaiting{}
		}
		_ = time.Second // reserved for future TTL checks
		return AuthOK{}
	}
}

// LoginInteractive triggers an interactive login, then signals AuthOK.
type LoggedIn struct{}

func LoginInteractive() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		client, err := radix.New(false)
		if err == nil {
			_ = client.LoginInteractive(ctx)
		}
		return LoggedIn{}
	}
}

// LoadCachedApplications emits Applications from disk cache (if present) for the current context.
// This lets the UI start fast without waiting for network or auth.
func LoadCachedApplications() tea.Cmd {
	return func() tea.Msg {
		// Determine current context (same as getContext)
		cfg, err := radixconfig.GetRadixConfig()
		var ctxName string
		if err == nil && cfg != nil && cfg.CustomConfig != nil {
			ctxName = cfg.CustomConfig.Context
		}
		cacheFile := fmt.Sprintf("%s/radix-tui-cache.json", radixconfig.RadixConfigDir)
		c := cache.New(cacheFile, "applications")
		key := fmt.Sprintf("apps_%s", ctxName)
		if raw, ok := c.GetItem(key); ok && len(raw) > 0 {
			var apps []string
			if err := json.Unmarshal([]byte(raw), &apps); err == nil && len(apps) > 0 {
				sort.Strings(apps)
				return Applications{Apps: apps, Context: ctxName}
			}
		}
		return nil
	}
}

func GetApplications() tea.Msg {
	// Get current context for error reporting
	cfg, _ := radixconfig.GetRadixConfig()
	var ctxName string
	if cfg != nil && cfg.CustomConfig != nil {
		ctxName = cfg.CustomConfig.Context
	}

	ctx := context.Background()
	client, err := radix.New(false)
	if err != nil {
		return ApplicationsError{Err: err, Context: ctxName}
	}
	apps, err := client.ListApplications(ctx)
	if err != nil {
		return ApplicationsError{Err: err, Context: ctxName}
	}
	sort.Strings(apps)
	// Persist to cache for fast next startup
	func() {
		cacheFile := fmt.Sprintf("%s/radix-tui-cache.json", radixconfig.RadixConfigDir)
		c := cache.New(cacheFile, "applications")
		key := fmt.Sprintf("apps_%s", ctxName)
		if b, err := json.Marshal(apps); err == nil {
			// Keep for 7 days; we'll always refresh on startup anyway
			c.SetItem(key, string(b), 7*24*time.Hour)
		}
	}()
	return Applications{Apps: apps, Context: ctxName}
}

type Application struct {
	Jobs         []Job         `json:"jobs"`
	Environments []Environment `json:"environments"`
	Name         string        `json:"name"`
}

type Job struct {
	Name         string   `json:"name"`
	TriggeredBy  string   `json:"triggeredBy"`
	Environments []string `json:"environments"`
	Pipeline     string   `json:"pipeline"`
	Status       string   `json:"status"`
	Created      string   `json:"created"`
}

type Environment struct {
	BranchMapping string `json:"branchMapping"`
	Name          string `json:"name"`
	Status        string `json:"status"`
}

func GetApplicationData(application string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		client, err := radix.New(false)
		if err != nil {
			fmt.Println(err)
			return Application{}
		}
		app, err := client.GetApplication(ctx, application)
		if err != nil {
			fmt.Println(err)
			return Application{}
		}
		// Map API model to local struct expected by views
		var jobs []Job
		if app.Jobs != nil {
			for _, j := range app.Jobs {
				if j == nil {
					continue
				}
				jobs = append(jobs, Job{
					Name:         stringDeref(j.Name),
					TriggeredBy:  j.TriggeredBy,
					Environments: j.Environments,
					Pipeline:     j.Pipeline,
					Status:       j.Status,
					Created: func() string {
						if j.Created != nil {
							return j.Created.String()
						}
						return ""
					}(),
				})
			}
		}
		var envs []Environment
		if app.Environments != nil {
			for _, e := range app.Environments {
				if e == nil {
					continue
				}
				envs = append(envs, Environment{
					BranchMapping: e.BranchMapping,
					Name:          stringDeref(e.Name),
					Status:        e.Status,
				})
			}
		}
		return Application{
			Jobs:         jobs,
			Environments: envs,
			Name:         stringDeref(app.Name),
		}
	}
}

func stringDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ContextChanged is sent when the context is successfully changed
type ContextChanged string

// SetContext changes the radix context and saves it to config
func SetContext(ctx string) tea.Cmd {
	return func() tea.Msg {
		if !radixconfig.IsValidContext(ctx) {
			return nil
		}
		radixConfig, err := radixconfig.GetRadixConfig()
		if err != nil {
			return nil
		}
		radixConfig.CustomConfig.Context = ctx
		if err := radixconfig.Save(radixConfig); err != nil {
			return nil
		}
		return ContextChanged(ctx)
	}
}
