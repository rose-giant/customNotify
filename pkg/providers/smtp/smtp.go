package smtp

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Provider struct {
	SMTP    []*Options `yaml:"smtp,omitempty"`
	counter int
}

type Options struct {
	ID              string   `yaml:"id,omitempty"`
	Server          string   `yaml:"smtp_server,omitempty"`
	Username        string   `yaml:"smtp_username,omitempty"`
	Password        string   `yaml:"smtp_password,omitempty"`
	FromAddress     string   `yaml:"from_address,omitempty"`
	SMTPCC          []string `yaml:"smtp_cc,omitempty"`
	SMTPFormat      string   `yaml:"smtp_format,omitempty"`
	Subject         string   `yaml:"subject,omitempty"`
	HTML            bool     `yaml:"smtp_html,omitempty"`
	DisableStartTLS bool     `yaml:"smtp_disable_starttls,omitempty"`
}

func New(options []*Options, ids []string) (*Provider, error) {
	provider := &Provider{}

	for _, o := range options {
		if len(ids) == 0 || slices.Contains(ids, o.ID) {
			provider.SMTP = append(provider.SMTP, o)
		}
	}

	provider.counter = 0

	return provider, nil
}

func buildUrl(password, username, server, fromAddr, subject string, toAddr []string, html, disableTLS bool) string {
	return fmt.Sprintf("smtp://%s:%s@%s:587/?fromAddress=%s&toAddresses=%s&subject=%s&UseHTML=%s&UseStartTLS=%s",
		username, url.QueryEscape(password), server, fromAddr, strings.Join(toAddr, ","), subject,
		strconv.FormatBool(html), strconv.FormatBool(!disableTLS))
}

func (p *Provider) Send(message, CliFormat string) error {
	var SmtpErr error
	p.counter++
	for _, pr := range p.SMTP {
		// msg := utils.FormatMessage(message, utils.SelectFormat(CliFormat, pr.SMTPFormat), p.counter)
		url := buildUrl(pr.Password, pr.Username, pr.Server, pr.FromAddress, pr.Subject,
			pr.SMTPCC, pr.HTML, pr.DisableStartTLS)
		err := shoutrrr.Send(url, message)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("failed to send smtp notification for id: %s ", pr.ID))
			SmtpErr = multierr.Append(SmtpErr, err)
			continue
		}
		fmt.Printf("smtp notification sent for id: %s \n", pr.ID)
	}
	return SmtpErr
}
