package cmd

import (
	"fmt"
	"os"
	"regexp"
	"unicode"

	"github.com/keep-network/keep-core/cmd/config"
	"github.com/urfave/cli"
)

// ValidationItem stores validation results
type ValidationItem struct {
	IsValid  bool     `json:"is_valid"`
	Error    string   `json:"error,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// ValidationType contains the results of validating the config file
type ValidationType struct {
	Config *ValidationItem `json:"config,omitempty"`
}

// IsValid returns true if config file validation was successful
func (v ValidationType) IsValid() bool {
	if v.Config == nil {
		return false
	}
	if v.Config != nil && !v.Config.IsValid {
		return false
	}
	return true
}

type warningsType []string

// Example:  warnings.delete("")
func (warnings *warningsType) delete(selector string) {
	var prunedWarnings warningsType
	for _, str := range *warnings {
		if str != selector {
			prunedWarnings = append(prunedWarnings, str)
		}
	}
	*warnings = prunedWarnings
}

const (
	cryptoRegexPattern  = "([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$"
	wsURLRegexPattern   = "ws:\\/\\/.+:.+"
	invalidPasswordRule = "invalid password, failed (%s) rule"
)

// ValidateConfig validates the contents of the config file
func ValidateConfig(c *cli.Context) (err error) {

	configPath := c.String("config")
	configPath, err = GetConfigFilePath(configPath)
	if err != nil {
		fmt.Printf("error with config file path:%s\n", err)
	}

	config.KeepOpts, err = config.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("error reading config file:%s\n", err)
	}

	warnings := warningsType{}
	url := config.KeepOpts.Ethereum.URL

	accountAddress := config.KeepOpts.Ethereum.Account.Address
	accountKeyFile := config.KeepOpts.Ethereum.Account.KeyFile
	accountKeyFilePassword := config.KeepOpts.Ethereum.Account.KeyFilePassword
	contractKeepRandomBeaconAddress := config.KeepOpts.Ethereum.ContractAddresses["KeepRandomBeacon"]
	contractGroupContractAddress := config.KeepOpts.Ethereum.ContractAddresses["GroupContract"]

	if valid, urlWarning := isURL(url); !valid {
		warnings = append(warnings, urlWarning)
	}
	if valid, accountAddressWarning := isAddress(accountAddress, "accountAddress"); !valid {
		warnings = append(warnings, accountAddressWarning)
	}
	if valid, accountKeyFileWarning := isFile(accountKeyFile); !valid {
		warnings = append(warnings, accountKeyFileWarning)
	}
	if valid, pwdWarnings := isPassword(accountKeyFilePassword); !valid {
		warnings = append(warnings, pwdWarnings...)
	}
	if valid, randomBeaconAddressWarning := isAddress(contractKeepRandomBeaconAddress, "contractKeepRandomBeaconAddress"); !valid {
		warnings = append(warnings, randomBeaconAddressWarning)
	}
	if valid, groupContractAddressWarning := isAddress(contractGroupContractAddress, "contractGroupContractAddress"); !valid {
		warnings = append(warnings, groupContractAddressWarning)
	}

	//warnings.delete("")
	configValidation := ValidationItem{
		IsValid:  true,
		Warnings: warnings,
	}
	if len(warnings) > 0 {
		configValidation.IsValid = false
		configValidation.Error = "invalid password"
	}
	validation := ValidationType{}
	validation.Config = &configValidation

	if !validation.IsValid() {
		fmt.Println("Config validation failed:")
		for _, warning := range warnings {
			fmt.Printf("* %s\n", warning)
		}
		os.Exit(1)
	}
	fmt.Println("validate-config success!")
	return nil
}

func isURL(url string) (matched bool, warning string) {
	// returns true if this is a valid URL
	matched, err := regexp.MatchString(wsURLRegexPattern, url)
	if err != nil || !matched {
		warning = fmt.Sprintf("error matching URL (%s) - %v", url, err)
		return
	}
	return matched, warning
}

func isAddress(address, name string) (matched bool, warning string) {
	// returns true if this is a valid address
	matched, err := regexp.MatchString(cryptoRegexPattern, address)
	if err != nil || !matched {
		warning = fmt.Sprintf("error matching (%s) address (%s) - error: %v", name, address, err)
		return
	}
	return matched, warning
}

func isFile(path string) (found bool, warning string) {
	// returns true if this file exists
	if exist, err := FileExists(path); err != nil {
		warning = fmt.Sprintf("unable to read file - error: %v", err)
	} else if !exist {
		warning = fmt.Sprintf("file does not exist: %s", path)
	}
	return (len(warning) == 0), warning
}

func isPassword(password string) (valid bool, warnings warningsType) {
	// returns true if this is a valid password
	atLeast8, number, upper, special := checkPassword(password)
	if !atLeast8 {
		warnings = append(warnings, fmt.Sprintf(invalidPasswordRule, "atLeast8"))
	}
	if !number {
		warnings = append(warnings, fmt.Sprintf(invalidPasswordRule, "number"))
	}
	if !upper {
		warnings = append(warnings, fmt.Sprintf(invalidPasswordRule, "upper"))
	}
	if !special {
		warnings = append(warnings, fmt.Sprintf(invalidPasswordRule, "special"))
	}
	return (len(warnings) == 0), warnings
}

func checkPassword(s string) (atLeast8, number, upper, special bool) {
	characters := 0
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
			characters++
		case unicode.IsUpper(s):
			upper = true
			characters++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
			characters++
		case unicode.IsLetter(s) || s == ' ':
			characters++
		default:
			//return false, false, false, false
		}
	}
	atLeast8 = characters >= 8
	return
}
