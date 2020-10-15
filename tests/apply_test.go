package tests

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestApply(t *testing.T) {
	sharedPath, err := setup("apply")
	if err != nil {
		t.Fatalf("setup() failed with: %v", err)
	}

	awsAccessKey, awsSecretKey := getAwsCreds(t)

	setupPlan(t, "apply", sharedPath, awsAccessKey, awsSecretKey)

	tests := []struct{
		name       string
		initParams []string
	}{
		{
			name: "apply",
			initParams: []string{"M_NAME=awsks-module-tests-apply"},
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

			docker.Run(t, awsksImageTag, planOpts)

			applyCommand := []string{"apply",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			applyOpts := &docker.RunOptions{
				Command: applyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, applyOpts)

			kubeconfigCommand := []string{"kubeconfig"}

			kubeconfigOpts := &docker.RunOptions{
				Command: kubeconfigCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, kubeconfigOpts)

			kubectlOpts := &k8s.KubectlOptions{
				ConfigPath: fmt.Sprintf("%s/kubeconfig", sharedPath),
			}

			k8s.RunKubectl(t, kubectlOpts, "get", "all", "-A")

			planDestroyCommand := []string{"plan-destroy",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			planDestroyOpts := &docker.RunOptions{
				Command: planDestroyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, planDestroyOpts)

			destroyCommand := []string{"destroy",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			destroyOpts := &docker.RunOptions{
				Command: destroyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, destroyOpts)
		})
	}

	cleanupPlan(t, sharedPath, awsAccessKey, awsSecretKey)
	cleanup(sharedPath)
}
