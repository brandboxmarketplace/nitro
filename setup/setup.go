package setup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
)

// FirstTime is used when there is no configuration file found in a users
// home/.nitro directory. We do not prompt for input such as memory, cpu,
// disk space in version 2 as that is defined and managed at the docker
// level. If anything fails, we return an error.
func FirstTime(home, env string, output terminal.Outputer) error {
	// TODO(jasonmccallister) consider prompts for which type(s) of database?
	c := config.Config{
		File: filepath.Join(home, ".nitro", env+".yaml"),
	}

	output.Info("Setting up Nitro...")

	mysql, err := confirm("Would you like to use MySQL", true)
	if err != nil {
		return err
	}

	if mysql {
		output.Pending("adding mysql 8.0 on port 3306")

		// add a default mysql database
		c.Databases = append(c.Databases, config.Database{
			Engine:  "mysql",
			Version: "8.0",
			Port:    "3306",
		})

		output.Done()
	}

	postgres, err := confirm("Would you like to use PostgreSQL", true)
	if err != nil {
		return err
	}

	if postgres {
		output.Pending("adding postgres 12 on port 5432")

		// add a default postgres database
		c.Databases = append(c.Databases, config.Database{
			Engine:  "postgres",
			Version: "12",
			Port:    "5432",
		})

		output.Done()
	}

	redis, err := confirm("Would you like to setup Redis", true)
	if err != nil {
		return err
	}

	if redis {
		output.Pending("adding redis 6379")

		c.Services.Redis = true

		output.Done()
	}

	// save the file
	if err := c.Save(); err != nil {
		return err
	}

	return nil
}

func confirm(msg string, fallback bool) (bool, error) {
	if fallback {
		fmt.Print(fmt.Sprintf("%s [Y/n]? ", msg))
	} else {
		fmt.Print(fmt.Sprintf("%s [y/N]? ", msg))
	}

	// prompt for confirmation
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil && err.Error() != "unexpected newline" {
		return false, fmt.Errorf("unable to provide a prompt, %w", err)
	}

	// check for empty newline and use the default
	if err != nil && err.Error() == "unexpected newline" {
		return fallback, nil
	}

	// read the confirmation from the input
	resp := strings.TrimSpace(input)
	for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
		if resp == answer {
			fallback = true
		}
	}

	return fallback, nil
}

func setupPostgres() (bool, error) {
	fmt.Print("Would you like to use PostgreSQL [Y/n]? ")

	// prompt the user for confirmation
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, fmt.Errorf("unable to provide a prompt, %w", err)
	}

	var confirm bool
	resp := strings.TrimSpace(response)
	for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
		if resp == answer {
			confirm = true
		}
	}

	return confirm, nil
}

func setupRedis() (bool, error) {
	fmt.Print("Would you like to use Redis [Y/n]? ")

	// prompt the user for confirmation
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, fmt.Errorf("unable to provide a prompt, %w", err)
	}

	var confirm bool
	resp := strings.TrimSpace(response)
	for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
		if resp == answer {
			confirm = true
		}
	}

	return confirm, nil
}
