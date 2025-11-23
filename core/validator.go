package core

import (
	"fmt"

	"github.com/Ether-Security/leviathan/utils"
	"github.com/go-playground/validator/v10"
)

func (r *Runner) ValidateInputType() bool {
	inputType, err := retrieveType(r.Target)
	if err != nil || inputType != r.Workflow.Validator {
		utils.Logger.Error().
			Str("want", r.Workflow.Validator).
			Str("got", inputType).
			Msg("Wrong input type")
		return false
	}
	return true
}

func retrieveType(raw string) (string, error) {
	var err error
	var inputType string
	v := validator.New()

	err = v.Var(raw, "required,url")
	if err == nil {
		inputType = "url"
	}

	err = v.Var(raw, "required,ipv4")
	if err == nil {
		inputType = "ip"
	}

	err = v.Var(raw, "required,fqdn")
	if err == nil {
		inputType = "domain"
	}

	err = v.Var(raw, "required,hostname")
	if err == nil {
		inputType = "domain"
	}

	err = v.Var(raw, "required,cidr")
	if err == nil {
		inputType = "cidr"
	}

	err = v.Var(raw, "required,uri")
	if err == nil {
		inputType = "url"
	}

	if inputType == "" {
		return "unknown", fmt.Errorf("unrecognized input")
	}

	return inputType, nil
}
