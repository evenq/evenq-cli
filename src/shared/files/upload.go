package files

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cheggaaa/pb/v3"
)

var accClient *s3.Client

func GetAccelerateClient() *s3.Client {
	if accClient != nil {
		return accClient
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		return nil
	}

	accClient = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UseAccelerate = true
	})

	return accClient
}

func Upload(r io.Reader, size int64, url string) error {
	client := http.Client{
		Timeout: time.Hour,
	}

	fmt.Println("Uploading File...")

	// start new bar
	bar := pb.Full.Start64(size)
	bar.SetWidth(100)
	// create proxy reader
	barReader := bar.NewProxyReader(r)

	req, err := http.NewRequest("PUT", url, barReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.ContentLength = size

	resp, err := client.Do(req)
	if err != nil {
		bar.Finish()
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with statuscode: %v", resp.StatusCode)
	}

	bar.Finish()

	return nil
}
