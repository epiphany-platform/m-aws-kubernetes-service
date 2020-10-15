package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
)

func TestInit(t *testing.T) {
	tests := []struct{
		name               string
		initParams         []string
		stateLocation      string
		stateContent       string
		wantOutput         string
		wantConfigLocation string
		wantConfigContent  string
		wantStateContent   string
	}{
		{
			name: "init with defaults",
			initParams: nil,
			stateLocation: "state.yml",
			stateContent: ``,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: unset
  region: eu-central-1
  public_subnet_id: unset
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: unset
  region: eu-central-1
  public_subnet_id: unset
`,
			wantStateContent: `
kind: state
awsks:
  status: initialized
`,
		},
		{
			name: "init with variables",
			initParams: []string{"M_NAME=value1", "M_VPC_ID=value2", "M_REGION=value3", "M_PUBLIC_SUBNET_ID=value4"},
			stateLocation: "state.yml",
			stateContent: ``,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: value1
  vpc_id: value2
  region: value3
  public_subnet_id: value4
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: value1
  vpc_id: value2
  region: value3
  public_subnet_id: value4
`,
			wantStateContent: `
kind: state
awsks:
  status: initialized
`,
		},
		{
			name: "init with state",
			initParams: nil,
			stateLocation: "state.yml",
			stateContent: `
kind: state
awsbi:
  status: applied
  name: epiphany
  instance_count: 0
  region: eu-central-1
  use_public_ip: false
  force_nat_gateway: true
  rsa_pub_path: "/shared/vms_rsa.pub"
  output:
    private_ip.value: []
    public_ip.value: []
    public_subnet_id.value: subnet-0137cf1e7921c1551
    vpc_id.value: vpc-0baa2c4e9e48e608c
`,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: vpc-0baa2c4e9e48e608c
  region: eu-central-1
  public_subnet_id: subnet-0137cf1e7921c1551
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: vpc-0baa2c4e9e48e608c
  region: eu-central-1
  public_subnet_id: subnet-0137cf1e7921c1551
`,
			wantStateContent: `
kind: state
awsbi:
  status: applied
  name: epiphany
  instance_count: 0
  region: eu-central-1
  use_public_ip: false
  force_nat_gateway: true
  rsa_pub_path: "/shared/vms_rsa.pub"
  output:
    private_ip.value: []
    public_ip.value: []
    public_subnet_id.value: subnet-0137cf1e7921c1551
    vpc_id.value: vpc-0baa2c4e9e48e608c
awsks:
  status: initialized
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sharedPath, err := setup("init")
			if err != nil {
				t.Fatalf("setup() failed with: %v", err)
			}
			defer cleanup(sharedPath)

			stateLocation := path.Join(sharedPath, tt.stateLocation)
			if err := ioutil.WriteFile(stateLocation, []byte(normStr(tt.stateContent)), 0644); err != nil {
				t.Fatalf("wasnt able to save state file: %s", err)
			}

			command := []string{"init"}
			command = append(command, tt.initParams...)

			runOpts := &docker.RunOptions{
				Command: command,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			output := docker.Run(t, awsksImageTag, runOpts)
			if diff := deep.Equal(normStr(output), normStr(tt.wantOutput)); diff != nil {
				t.Error(diff)
			}

			configLocation := path.Join(sharedPath, tt.wantConfigLocation)
			if _, err := os.Stat(configLocation); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", configLocation)
			}

			gotConfigContent, err := ioutil.ReadFile(configLocation)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}

			if diff := deep.Equal(normStr(string(gotConfigContent)), normStr(tt.wantConfigContent)); diff != nil {
				t.Error(diff)
			}

			gotStateContent, err := ioutil.ReadFile(stateLocation)
			if err != nil {
				t.Errorf("wasnt able to read form state file: %v", err)
			}

			if diff := deep.Equal(normStr(string(gotStateContent)), normStr(tt.wantStateContent)); diff != nil {
				t.Error(diff)
			}
		})
	}
}
