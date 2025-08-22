package radix

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	radixapi "github.com/equinor/radix-cli/generated/radixapi/client"
	"github.com/equinor/radix-cli/generated/radixapi/client/application"
	"github.com/equinor/radix-cli/generated/radixapi/client/platform"
	rm "github.com/equinor/radix-cli/generated/radixapi/models"
	"github.com/equinor/radix-cli/pkg/cache"
	"github.com/equinor/radix-cli/pkg/client/auth"
	radixconfig "github.com/equinor/radix-cli/pkg/config"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// Client provides access to Radix API using the same auth flow as radix-cli.
type Client struct {
	API        *radixapi.Radixapi
	provider   auth.Provider
	triedLogin bool
}

// New creates an authenticated Radix API client using stored auth (or interactive on first use).
func New(verbose bool) (*Client, error) {
	// Ensure the Radix config directory exists (same as radix-cli does on startup)
	if _, err := os.Stat(radixconfig.RadixConfigDir); err != nil {
		if mkErr := os.MkdirAll(radixconfig.RadixConfigDir, os.ModePerm); mkErr != nil {
			return nil, mkErr
		}
	}

	// Prefer interactive auth provider by default (avoid device-code unless explicitly requested)
	authCacheFilename := fmt.Sprintf("%s/auth.json", radixconfig.RadixConfigDir)
	globalCache := cache.New(authCacheFilename, "global")
	globalCache.SetItem("auth_provider_type", "msal_interactive", 365*24*time.Hour)

	rc, err := radixconfig.GetRadixConfig()
	if err != nil {
		return nil, err
	}
	contextName := rc.CustomConfig.Context
	endpoint := getEndpoint("server-radix-api", "", contextName, "")

	provider, err := auth.New()
	if err != nil {
		return nil, err
	}

	transport := httptransport.New(endpoint, "/api/v1", []string{"https"})
	transport.DefaultAuthentication = provider
	transport.Debug = verbose

	return &Client{API: radixapi.New(transport, strfmt.Default), provider: provider}, nil
}

// ensureLogin attempts an API call and performs interactive login on 401/provider-not-set, then retries via fn.
func (c *Client) ensureLogin(ctx context.Context, fn func() error) error {
	if err := fn(); err != nil {
		// If not yet tried, always attempt an interactive login once and retry.
		if !c.triedLogin {
			c.triedLogin = true
			if lerr := c.loginViaHelper(ctx); lerr != nil {
				return err
			}
			return fn()
		}
		// Best-effort detection for repeated auth errors
		if strings.Contains(strings.ToLower(err.Error()), "auth provider not set") ||
			strings.Contains(strings.ToLower(err.Error()), "unauthorized") ||
			strings.Contains(strings.ToLower(err.Error()), "401") {
			if lerr := c.loginViaHelper(ctx); lerr == nil {
				return fn()
			}
		}
		return err
	}
	return nil
}

// ListApplications returns application names.
func (c *Client) ListApplications(ctx context.Context) ([]string, error) {
	var out []string
	var resp *platform.ShowApplicationsOK
	call := func() error {
		var err error
		resp, err = c.API.Platform.ShowApplications(platform.NewShowApplicationsParams(), nil)
		return err
	}
	if err := c.ensureLogin(ctx, call); err != nil {
		return nil, err
	}
	for _, a := range resp.Payload {
		if a.Name != nil {
			out = append(out, *a.Name)
		}
	}
	return out, nil
}

// GetApplication returns a full application document.
func (c *Client) GetApplication(ctx context.Context, app string) (*rm.Application, error) {
	var resp *application.GetApplicationOK
	call := func() error {
		params := application.NewGetApplicationParams()
		params.SetAppName(app)
		var err error
		resp, err = c.API.Application.GetApplication(params, nil)
		return err
	}
	if err := c.ensureLogin(ctx, call); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// TriggerBuildDeploy triggers a build-deploy pipeline for the given app/branch.
func (c *Client) TriggerBuildDeploy(ctx context.Context, app, branch, toEnv string) (*rm.JobSummary, error) {
	var resp *application.TriggerPipelineBuildDeployOK
	call := func() error {
		params := application.NewTriggerPipelineBuildDeployParams()
		params.SetAppName(app)
		params.SetPipelineParametersBuild(&rm.PipelineParametersBuild{
			GitRef:        branch,
			GitRefType:    "branch",
			ToEnvironment: toEnv,
		})
		var err error
		resp, err = c.API.Application.TriggerPipelineBuildDeploy(params, nil)
		return err
	}
	if err := c.ensureLogin(ctx, call); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// TriggerApplyConfig triggers an apply-config pipeline for the given app.
func (c *Client) TriggerApplyConfig(ctx context.Context, app string) (*rm.JobSummary, error) {
	var resp *application.TriggerPipelineApplyConfigOK
	call := func() error {
		params := application.NewTriggerPipelineApplyConfigParams()
		params.SetAppName(app)
		params.SetPipelineParametersApplyConfig(&rm.PipelineParametersApplyConfig{})
		var err error
		resp, err = c.API.Application.TriggerPipelineApplyConfig(params, nil)
		return err
	}
	if err := c.ensureLogin(ctx, call); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func getEndpoint(service, env, context, cluster string) string {
	zoneDomain, defaultEnv := getPatternForContext(context)
	if strings.TrimSpace(env) == "" {
		env = defaultEnv
	}

	if cluster != "" {
		return fmt.Sprintf("%s-%s.%s.%sradix.equinor.com", service, env, cluster, zoneDomain)
	}

	return fmt.Sprintf("%s-%s.%sradix.equinor.com", service, env, zoneDomain)
}

func getPatternForContext(context string) (string, string) {
	switch context {
	case radixconfig.ContextDevelopment:
		return "dev.", "qa"
	case radixconfig.ContextPlayground:
		return "playground.", "prod"
	case radixconfig.ContextPlatform2:
		return "c2.", "prod"
	case radixconfig.ContextProduction, radixconfig.ContextPlatform:
		return "", "prod"
	default:
		return "", "prod"
	}
}

// LoginInteractive forces an interactive login flow.
func (c *Client) LoginInteractive(ctx context.Context) error {
	c.triedLogin = true
	return c.loginViaHelper(ctx)
}

// LoginDeviceCode forces a device-code login flow.
func (c *Client) LoginDeviceCode(ctx context.Context) error {
	c.triedLogin = true
	return c.loginViaHelper(ctx)
}

// loginViaHelper starts a separate process of this binary to perform interactive login with no terminal output.
func (c *Client) loginViaHelper(ctx context.Context) error {
	// If we're already in the helper process, perform the interactive login directly
	// to avoid recursively spawning more helper processes.
	if os.Getenv("RADIX_TUI_LOGIN_HELPER") == "1" {
		// Force interactive login; other modes disabled by default.
		return c.provider.Login(ctx, true, false, false, "", "", "")
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, exe)
	cmd.Env = append(os.Environ(), "RADIX_TUI_LOGIN_HELPER=1")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
