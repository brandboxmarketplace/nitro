package nitro

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/command"
)

func TestInit(t *testing.T) {
	type args struct {
		name    string
		cpus    string
		memory  string
		disk    string
		php     string
		db      string
		version string
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "installs the latest version",
			args: args{
				name:    "thisname",
				cpus:    "4",
				memory:  "4G",
				disk:    "20G",
				php:     "7.4",
				db:      "mysql",
				version: "5.7",
			},
			want: []Command{
				{
					Machine:   "thisname",
					Type:      "launch",
					Chainable: true,
					Input:     command.CloudInit,
					Args:      []string{"--name", "thisname", "--cpus", "4", "--mem", "4G", "--disk", "20G", "--cloud-init", "-"},
				},
				{
					Machine:   "thisname",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"thisname", "--", "sudo", "apt", "install", "-y", "php7.4", "php7.4-mbstring", "php7.4-cli", "php7.4-curl", "php7.4-fpm", "php7.4-gd", "php7.4-intl", "php7.4-json", "php7.4-mysql", "php7.4-opcache", "php7.4-pgsql", "php7.4-zip", "php7.4-xml", "php-xdebug", "php-imagick"},
				},
				{
					Machine:   "thisname",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"thisname", "--", "docker", "run", "-d", "--restart=always", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro", "mysql:5.7"},
				},
			},
		},
		{
			name: "installs a specific version of php and postgres",
			args: args{
				name:    "thisname",
				cpus:    "4",
				memory:  "4G",
				disk:    "20G",
				php:     "7.3",
				db:      "postgres",
				version: "11.5",
			},
			want: []Command{
				{
					Machine:   "thisname",
					Type:      "launch",
					Chainable: true,
					Input:     command.CloudInit,
					Args:      []string{"--name", "thisname", "--cpus", "4", "--mem", "4G", "--disk", "20G", "--cloud-init", "-"},
				},
				{
					Machine:   "thisname",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"thisname", "--", "sudo", "apt", "install", "-y", "php7.3", "php7.3-mbstring", "php7.3-cli", "php7.3-curl", "php7.3-fpm", "php7.3-gd", "php7.3-intl", "php7.3-json", "php7.3-mysql", "php7.3-opcache", "php7.3-pgsql", "php7.3-zip", "php7.3-xml", "php-xdebug", "php-imagick"},
				},
				{
					Machine:   "thisname",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"thisname", "--", "docker", "run", "-d", "--restart=always", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro", "postgres:11.5"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.name, tt.args.cpus, tt.args.memory, tt.args.disk, tt.args.php, tt.args.db, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
