package tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
)

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
