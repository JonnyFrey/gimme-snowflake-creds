package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/hashicorp/go-hclog"
)

var (
	validate *validator.Validate

	oauthParams = []string{
		"okta-org",
		"client-id",
		"issuer-url",
		"redirect-uri",
	}

	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
)

type Configuration struct {
	Profile      string
	Default      string `mapstructure:"default"`
	Account      string `mapstructure:"account" validate:"required"`
	Database     string `mapstructure:"database" validate:"required"`
	Warehouse    string `mapstructure:"warehouse" validate:"required"`
	Schema       string `mapstructure:"schema"`
	OAuth        bool   `mapstructure:"oauth"`
	Generic      bool   `mapstructure:"generic"`
	OktaOrg      string `mapstructure:"okta-org" validate:"required,url"`
	ODBCPath     string `mapstructure:"odbc-path" validate:"required"`
	ODBCDriver   string `mapstructure:"odbc-driver" validate:"required"`
	ClientID     string `mapstructure:"client-id" validate:"required"`
	Role         string `mapstructure:"role" validate:"required"`
	IssuerURL    string `mapstructure:"issuer-url" validate:"required,url"`
	RedirectURI  string `mapstructure:"redirect-uri" validate:"required,uri"`
	Username     string `mapstructure:"username" validate:"required,email"`
	Password     string
	HomeDir      string
	Logger       hclog.Logger
	ColorSuccess string
	ColorFailure string
}

type Credentials struct {
	ExpiresIn   int
	AccessToken string
}

func LoadDefaults(p *Configuration) error {
	p.ColorSuccess = colorGreen
	p.ColorFailure = colorRed

	return nil
}

func ValidateConfiguration(p *Configuration) error {
	validate = validator.New()

	// Register function to get tag name from mapstructure tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("mapstructure"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(p)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if !p.OAuth && contains(oauthParams, err.Field()) {
				return nil
			} else {
				fmt.Println(string(p.ColorFailure), "Parameter", err.Field(), "is required")
			}
		}

		return err
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
