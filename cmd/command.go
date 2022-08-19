package cmd

import (
	"context"
	"fmt"
	"net/url"
	"notifier/config"
	"notifier/library"
	"notifier/service"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func BuildCommand(ctx context.Context, conf *config.Notifier) *cobra.Command {
	return &cobra.Command{
		Use:   "notifier",
		Short: "Notifier sender",
		Long:  `Notifier sends notifications to the provided URL`,
		Run: func(command *cobra.Command, args []string) {
			params, err := buildFlags(command)
			if err != nil {
				log.Error("error building the flags")
				return
			}

			if err := executeNotifierPipeline(ctx, conf, *params); err != nil {
				log.Error(err)
				return
			}
		},
	}
}

func ConfigFlags(command *cobra.Command, conf *config.Notifier) error {
	command.PersistentFlags().
		StringP("url", "u", "", "send the notifications to this URL")

	command.PersistentFlags().
		DurationP("interval", "i", conf.Interval, "notification interval")

	err := command.MarkPersistentFlagRequired("url")
	if err != nil {
		return fmt.Errorf("setting flag as required %v", err)
	}
	return nil
}

func buildFlags(command *cobra.Command) (*Params, error) {
	cmdURL, err := command.Flags().GetString("url")
	if err != nil {
		return nil, err
	}
	interval, err := command.Flags().GetDuration("interval")
	if err != nil {
		return nil, err
	}

	return &Params{
		URL:      cmdURL,
		Interval: interval,
	}, nil

}

func executeNotifierPipeline(ctx context.Context, conf *config.Notifier, params Params) error {
	definedURL, err := url.ParseRequestURI(params.URL)
	if err != nil {
		return fmt.Errorf("error parsing the URL")
	}
	httpClient := library.NewClient(definedURL, conf.HTTPClient)

	processor := service.NewProcessor(os.Stdin, httpClient, params.Interval, conf.MaxAllowedErrors, conf.BufferSize)
	errList := processor.Process(ctx)
	err = waitForPipeline(errList, conf.MaxAllowedErrors)
	if err != nil {
		return fmt.Errorf("max allowed errors reached")
	}
	return nil
}

func waitForPipeline(errList <-chan error, maxErrors int) error {
	errCounter := 0
	for err := range errList {
		if err != nil {
			log.Error(err)
			errCounter++
			if errCounter > maxErrors {
				return err
			}
		}
	}
	return nil
}
