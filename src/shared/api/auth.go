package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/evenq/evenq-cli/src/shared/config"
	"github.com/evenq/evenq-cli/src/shared/util"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gopkg.in/guregu/null.v3"
)

type user struct {
	ID             string      `json:"id"`
	Name           null.String `json:"name"`
	Email          string      `json:"email"`
	PasswordHash   null.String `json:"-"`
	CreatedAt      time.Time   `json:"createdAt"`
	IsVerified     bool        `json:"isVerified"`
	IsSubChangelog bool        `json:"subChangelog"`
	OrgIDs         []string    `json:"orgIds"`
	IsOnboarded    bool        `json:"isOnboarded"`
}

const SignupUrl = "https://app.evenq.io/signup"

func RunLogin(c *cli.Context) error {
	_, ok := StartAuth(c.Context)
	if !ok {
		return errors.New("could not authenticate")
	}

	return nil
}

func StartAuth(ctx context.Context) (context.Context, bool) {
	ctx, ok := getAuthenticatedContext(ctx)
	if !ok {
		return nil, false
	}

	ctx, ok = runOrgPicker(ctx)
	if !ok {
		return nil, false
	}

	fmt.Printf("Logged in as %v on behalf of %v\n",
		util.BlueText(ctx.Value("email")),
		util.BlueText(ctx.Value("orgId")),
	)

	return ctx, true
}

func CheckAuth(ctx context.Context) (context.Context, bool) {
	if existing, ok := getExistingToken(); ok {
		s := util.Spinner("Authenticating...")
		if newCtx, ok := validateToken(ctx, existing); ok {
			s.Stop()
			return runOrgPicker(newCtx)
		}
		s.Stop()
	}

	return ctx, false
}

func getAuthenticatedContext(ctx context.Context) (context.Context, bool) {
	if existing, ok := getExistingToken(); ok {
		s := util.Spinner("Authenticating...")
		if newCtx, ok := validateToken(ctx, existing); ok {
			s.Stop()
			return newCtx, true
		}
		s.Stop()
	}

	askSignupQuestion()

	if token, ok := promptLogin(); ok {
		s := util.Spinner("Authenticating...")
		if newCtx, ok := validateToken(ctx, token); ok {
			s.Stop()
			return newCtx, true
		}
		s.Stop()
	}

	return nil, false
}

func askSignupQuestion() {
	prompt := promptui.Select{
		Label: "Do you already have an " + util.BlueText("evenq.io") + " account?",
		Items: []string{"Yes", "No"},
	}
	_, hasAccount, err := prompt.Run()
	if err != nil {
		return
	}

	if hasAccount == "No" {
		util.OpenUrl(SignupUrl)
	}
}

func runOrgPicker(ctx context.Context) (context.Context, bool) {
	orgs, ok := ctx.Value("orgs").([]string)
	if !ok {
		fmt.Println("orgs not an array of strings")
		return nil, false
	}

	if len(orgs) == 1 {
		ctx = context.WithValue(ctx, "orgId", orgs[0])
	} else {
		prompt := promptui.Select{
			Label: "Select your organization",
			Items: orgs,
		}

		_, orgId, err := prompt.Run()
		if err != nil {
			return nil, false
		}

		ctx = context.WithValue(ctx, "orgId", orgId)
	}

	return ctx, true
}

func getExistingToken() (string, bool) {
	return config.GetValue("authToken")
}

func promptLogin() (string, bool) {
	prompt := promptui.Prompt{
		Label: "Please enter your email address",
	}

	email, err := prompt.Run()
	if err != nil {
		return "", false
	}

	prompt = promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}

	password, err := prompt.Run()
	if err != nil {
		return "", false
	}

	s := util.Spinner("Authenticating...")
	token, err := performLogin(email, password)
	s.Stop()
	if err != nil {
		fmt.Println(err.Error())
		return promptLogin()
	}

	return token, true
}

func performLogin(email string, password string) (string, error) {
	data := map[string]string{
		"email":    email,
		"password": password,
	}

	out := map[string]interface{}{}

	resp, err := Post(context.Background(), "/auth/login", data, &out)
	if err != nil {
		return "", err
	}

	if errStr, ok := out["error"].(string); ok {
		return "", errors.New(errStr)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "evenq_user_auth" {
			return c.Value, nil
		}
	}

	log.Println("no cookie in response")

	return "", nil
}

func validateToken(ctx context.Context, token string) (context.Context, bool) {
	time.Sleep(2 * time.Second)

	out := user{}

	contextWithToken := setToken(context.Background(), token)

	err := Get(contextWithToken, "/me", &out)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	if out.ID != "" {
		if !out.IsVerified {
			log.Println("The email address on your account is not verified, please check your inbox.")
			return nil, false
		}

		// save to context
		ctx = setToken(ctx, token)
		ctx = context.WithValue(ctx, "email", out.Email)
		ctx = context.WithValue(ctx, "orgs", out.OrgIDs)

		// save to config file
		config.SetValue("authToken", token)

		return ctx, true
	}

	return nil, false
}

func getToken(ctx context.Context) (string, bool) {
	if t, ok := ctx.Value("authToken").(string); ok {
		return t, true
	}

	return "", false
}

func getOrg(ctx context.Context) string {
	org, _ := ctx.Value("orgId").(string)
	return org
}

func setToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, "authToken", token)
}
