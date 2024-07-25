package providers

import (
	"slices"

	"example.com/custom_notify/pkg/providers/smtp"
	"example.com/custom_notify/pkg/types"
	"github.com/acarl005/stripansi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

// ProviderOptions is configuration for notify providers
type ProviderOptions struct {
	SMTP []*smtp.Options `yaml:"smtp,omitempty"`
}

// Provider is an interface implemented by providers
type Provider interface {
	Send(message, CliFormat string) error
}

type Client struct {
	providers       []Provider
	providerOptions *ProviderOptions
	options         *types.Options
}

func New(providerOptions *ProviderOptions, options *types.Options) (*Client, error) {

	client := &Client{providerOptions: providerOptions, options: options}

	if providerOptions.SMTP != nil && (len(options.Providers) == 0 || slices.Contains(options.Providers, "smtp")) {

		provider, err := smtp.New(providerOptions.SMTP, options.IDs)
		if err != nil {
			return nil, errors.Wrap(err, "could not create smtp provider client")
		}
		client.providers = append(client.providers, provider)
	}

	return client, nil
}

func (p *Client) Send(message string) error {

	// strip unsupported color control chars
	message = stripansi.Strip(message)

	for _, v := range p.providers {
		if err := v.Send(message, p.options.MessageFormat); err != nil {
			for _, v := range multierr.Errors(err) {
				// gologger.Error().Msgf("%s", v)
				log.Error(v)
			}
		}
	}

	return nil
}
