package eventCommand

import (
	"errors"
	"fmt"

	"github.com/evenq/evenq-cli/src/shared/api"
	"github.com/evenq/evenq-cli/src/shared/events"
	"github.com/evenq/evenq-cli/src/shared/util"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func RunCreate(c *cli.Context) error {
	if ctx, ok := api.StartAuth(c.Context); ok {
		c.Context = ctx
	} else {
		fmt.Println("could not authenticate")
		return nil
	}

	p := promptui.Prompt{
		Label: "Please enter a name for your event",
	}
	name, _ := p.Run()

	data, ok := events.Get(c.Context, name)
	if ok {
		fmt.Printf("%v was created on %v. Has %v records.\n", name, data.CreatedAt.Format("2006-01-02"), data.EventStats.TotalCount)
		return nil
	}

	s := util.Spinner("Creating Event")

	res, ok := events.Create(c.Context, name)
	if !ok {
		s.Stop()
		return errors.New("failed to create event")
	}

	s.Stop()

	fmt.Printf("event %v created successfully!\n", util.BlueText(res.ID))

	return nil
}
