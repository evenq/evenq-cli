package imports

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/evenq/evenq-cli/src/shared/csv"
	"github.com/evenq/evenq-cli/src/shared/events"
	"github.com/evenq/evenq-cli/src/shared/util"
	"github.com/manifoldco/promptui"
)

const PKKey = "__evenqpk"
const TSKey = "__evenqts"

func getEventName(ctx context.Context) (string, error) {
	p := promptui.Prompt{
		Label: "Please enter the name for your event",
	}

	name, err := p.Run()
	if err != nil {
		return "", err
	}

	data, ok := events.Get(ctx, name)
	if !ok {
		p2 := promptui.Select{
			Label:    "Event not found. Create it?",
			Items:    []string{"Yes", "No", "Try Again"},
			HideHelp: true,
		}

		_, create, err := p2.Run()
		if err != nil {
			return "", err
		}

		if create == "Try Again" {
			return getEventName(ctx)
		} else if create == "Yes" {
			_, ok := events.Create(ctx, name)
			if !ok {
				return "", errors.New("failed to create event")
			} else {
				fmt.Printf("Created %v successfuly", name)
			}
		} else {
			return "", errors.New("can not proceed without event")
		}
	} else {
		fmt.Printf("found %v from %v with %v records.\n", name, data.CreatedAt.Format("2006-01-02"), data.EventStats.TotalCount)
	}

	return name, nil
}

func validatePath(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func mapHeaders(ctx context.Context, file string) (map[string]string, error) {
	hasPK := false
	hasTS := false
	hmap := map[string]string{}

	hh, err := csv.ReadHeaders(file)
	if err != nil {
		return hmap, err
	}

	p2 := promptui.Select{
		Label:    fmt.Sprintf("Is this your header row: %v?", strings.Join(hh, ",")),
		Items:    []string{"Yes", "No"},
		HideHelp: true,
	}

	_, resp, _ := p2.Run()
	if resp == "No" {
		return hmap, errors.New("your CSV file needs to have a header row")
	}

	for _, header := range hh {
		res, err := getHeaderMapping(ctx, header, &hasTS, &hasPK)
		if err != nil {
			return hmap, err
		}
		hmap[header] = res
	}

	return hmap, nil
}

func getHeaderMapping(ctx context.Context, h string, hasTS *bool, hasPK *bool) (string, error) {
	const (
		yesOpt       = "Yes"
		yesRenameOpt = "Yes, Rename"
		tsOpt        = "Yes, As Timestamp"
		pkOpt        = "Yes, As Partition Key"
		noOpt        = "No"
	)

	items := []string{yesOpt, yesRenameOpt}

	if *hasTS == false {
		items = append(items, tsOpt)
	}
	if *hasPK == false {
		items = append(items, pkOpt)
	}

	items = append(items, noOpt)

	p2 := promptui.Select{
		Label:        fmt.Sprintf("Import the %v column", util.BlueText(h)),
		Items:        items,
		HideHelp:     true,
		HideSelected: true,
	}

	_, resp, err := p2.Run()
	if err != nil {
		return "", err
	}

	switch resp {
	case yesOpt:
		fmt.Printf("import %v\n", util.BlueText(h))
		return h, nil
	case yesRenameOpt:
		p := promptui.Prompt{
			Label:   fmt.Sprintf("import %v as", util.BlueText(h)),
			Default: h,
		}
		name, err := p.Run()
		if err != nil {
			return "", err
		}

		fmt.Printf("import %v as %v\n", util.BlueText(h), util.BlueText(name))

		return name, nil
	case tsOpt:
		fmt.Printf("import %v as Timestamp\n", util.BlueText(h))
		*hasTS = true
		return TSKey, nil
	case pkOpt:
		fmt.Printf("import %v as Partition Key\n", util.BlueText(h))
		*hasPK = true
		return PKKey, nil
	}

	return "", nil
}
