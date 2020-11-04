package tests

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
	"golang.org/x/crypto/ssh"
)

const (
	awsbiImageTag = "epiphanyplatform/awsbi:0.0.1"
	awsksImageTag = "epiphanyplatform/awsks:0.0.1"
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

func setupPlan(t *testing.T, suffix, sharedPath, awsAccessKey, awsSecretKey string) {
	if err := generateRsaKeyPair(sharedPath, "test_vms_rsa"); err != nil {
		t.Fatalf("wasnt able to create rsa file: %s", err)
	}

	initCommand := []string{
		"init",
		"M_VMS_COUNT=0",
		"M_PUBLIC_IPS=false",
		fmt.Sprintf("M_NAME=awsks-module-tests-%s", suffix),
		"M_VMS_RSA=test_vms_rsa",
	}

	initOpts := &docker.RunOptions{
		Command: initCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, initOpts)

	planCommand := []string{"plan",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	planOpts := &docker.RunOptions{
		Command: planCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, planOpts)

	applyCommand := []string{"apply",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	applyOpts := &docker.RunOptions{
		Command: applyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, applyOpts)
}

func cleanupPlan(t *testing.T, sharedPath, awsAccessKey, awsSecretKey string) {
	planDestroyCommand := []string{"plan-destroy",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	planDestroyOpts := &docker.RunOptions{
		Command: planDestroyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, planDestroyOpts)

	destroyCommand := []string{"destroy",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	destroyOpts := &docker.RunOptions{
		Command: destroyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, destroyOpts)
}

func TestPlan(t *testing.T) {
	sharedPath, err := setup("plan")
	if err != nil {
		t.Fatalf("setup() failed with: %v", err)
	}

	awsAccessKey, awsSecretKey := getAwsCreds(t)

	setupPlan(t, "plan", sharedPath, awsAccessKey, awsSecretKey)

	tests := []struct{
		name                   string
		initParams             []string
		wantPlanOutputLastLine string
		wantTfPlanLocation     string
	}{
		{
			name: "plan",
			initParams: []string{"M_NAME=awsks-module-tests-plan"},
			wantPlanOutputLastLine: `Plan: 43 to add, 0 to change, 0 to destroy.`,
			wantTfPlanLocation: "awsks/terraform-apply.tfplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommand := []string{"init"}
			initCommand = append(initCommand, tt.initParams...)

			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, initOpts)

			planCommand := []string{"plan",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			planOpts := &docker.RunOptions{
				Command: planCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotPlanOutput := docker.Run(t, awsksImageTag, planOpts)
			gotPlanOutputLastLine, err := getLastLineFromMultilineString(gotPlanOutput)
			if err != nil {
				t.Fatalf("reading last line from multiline failed with: %v", err)
			}

			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			tfPlanLocation := path.Join(sharedPath, tt.wantTfPlanLocation)
			if _, err := os.Stat(tfPlanLocation); os.IsNotExist(err) {
				t.Fatalf("missing tfplan file: %s", tfPlanLocation)
			}
		})
	}

	cleanupPlan(t, sharedPath, awsAccessKey, awsSecretKey)
	cleanup(sharedPath)
}

func setup(suffix string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := path.Join(wd, fmt.Sprintf("%s-%s", "shared", suffix))
	return p, os.MkdirAll(p, os.ModePerm)
}

func cleanup(sharedPath string) error {
	return os.RemoveAll(sharedPath)
}

func normStr(s string) string {
	return strings.TrimSpace(s)
}

func getLastLineFromMultilineString(s string) (string, error) {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF {
			return string(line), nil
		}
	}
}

func generateRsaKeyPair(directory, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
}

func getAwsCreds(t *testing.T) (awsAccessKey, awsSecretKey string) {
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
	if len(awsAccessKey) == 0 {
		t.Fatalf("expected non-empty AWS_ACCESS_KEY environment variable")
	}

	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	if len(awsSecretKey) == 0 {
		t.Fatalf("expected non-empty AWS_SECRET_KEY environment variable")
	}

	return
}


