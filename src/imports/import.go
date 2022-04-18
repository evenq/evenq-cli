package imports

import (
	"errors"
	"fmt"

	"github.com/evenq/evenq-cli/src/shared/api"
	"github.com/evenq/evenq-cli/src/shared/events"
	"github.com/evenq/evenq-cli/src/shared/files"
	"github.com/evenq/evenq-cli/src/shared/util"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/guregu/null.v3"
)

func Run(c *cli.Context) error {
	var err error

	filePath := c.Args().Get(0)
	fileStat, err := validatePath(filePath)
	if err != nil {
		return err
	}

	if ctx, ok := api.CheckAuth(c.Context); ok {
		c.Context = ctx
	} else {
		fmt.Println("please run", util.BlueText("evenq login"), "to authenticate")
		return nil
	}

	eventName, err := getEventName(c.Context)
	if err != nil {
		return err
	}

	hmap, err := mapHeaders(c.Context, filePath)
	if err != nil {
		return err
	}

	reader, size, err := files.Prepare(filePath)
	if err != nil {
		return err
	}

	s := util.Spinner("Creating Import...")
	imp, ok := events.CreateImport(c.Context, events.Import{
		EventID:    eventName,
		FileSize:   null.IntFrom(size),
		FileName:   null.StringFrom(fileStat.Name()),
		FileFormat: null.StringFrom("csv-gzip"), // we always gzip so we can hardcode this
	})
	s.Stop()
	if !ok {
		return errors.New("could not create import")
	}

	if !imp.UploadURL.Valid {
		return errors.New("import missing upload url")
	}

	err = files.Upload(reader, size, imp.UploadURL.String)
	if err != nil {
		return err
	}

	ok = events.StartImport(c.Context, eventName, imp.ID, hmap)
	if !ok {
		return errors.New("failed to queue import")
	}

	color.Green("Import Queued! We will email you when it's done.")

	return nil
}
