package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"
	"github.com/spf13/cobra"
)

var (
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy the project",
		Run:   runDeploy,
	}

	dockerEmail    string
	dockerUsername string
	dockerPassword string

	kubernetesApiServer string

	kubernetesCaPath  string
	kubernetesCrtPath string
	kubernetesKeyPath string

	kubectlVersion string

	kubernetesResourcesDirectoryPath string
	templatedResourcesDirectoryPath  string
	removeResourceFilesAfterUse      bool
)

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&dockerEmail, "docker-email", "", "email to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", "", "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", "", "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&kubernetesApiServer, "kubernetes-api-server", "https://api.g8s.fra-1.giantswarm.io", "kubernetes api to deploy to")

	deployCmd.Flags().StringVar(&kubernetesCaPath, "kubernetes-ca-path", "", "path to kubernetes ca file")
	deployCmd.Flags().StringVar(&kubernetesCrtPath, "kubernetes-crt-path", "", "path to kubernetes certificate file")
	deployCmd.Flags().StringVar(&kubernetesKeyPath, "kubernetes-key-path", "", "path to kubernetes key file")

	deployCmd.Flags().StringVar(&kubectlVersion, "kubectl-version", "1.5.2", "kubectl version")

	deployCmd.Flags().StringVar(&kubernetesResourcesDirectoryPath, "kubernetes-resources-directory-path", "./kubernetes", "directory holding kubernetes resources")
	deployCmd.Flags().StringVar(&templatedResourcesDirectoryPath, "templated-resources-directory-path", "./kubernetes-templated", "directory holding templated kubernetes resources")
	deployCmd.Flags().BoolVar(&removeResourceFilesAfterUse, "remove-resource-files-after-use", true, "whether to remove templated kubernetes resource files after use")
}

func runDeploy(cmd *cobra.Command, args []string) {
	if err := utils.TemplateKubernetesResources(kubernetesResourcesDirectoryPath, templatedResourcesDirectoryPath, sha); err != nil {
		log.Fatalf("could not template kubernetes resources: %v\n", err)
	}

	dockerLogin := commands.Command{
		Name: "docker-login",
		Args: []string{
			"docker",
			"login",
			fmt.Sprintf("--email=%v", dockerEmail),
			fmt.Sprintf("--username=%v", dockerUsername),
			fmt.Sprintf("--password=%v", dockerPassword),
			registry,
		},
	}

	dockerPush := commands.Command{
		Name: "docker-push",
		Args: []string{
			"docker",
			"push",
			fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
		},
	}

	kubectlClusterInfo := commands.Command{
		Name: "kubectl-cluster-info",
		Args: []string{
			"docker",
			"run",
			"--rm",
			"-v", fmt.Sprintf("%v:/ca.pem", kubernetesCaPath),
			"-v", fmt.Sprintf("%v:/crt.pem", kubernetesCrtPath),
			"-v", fmt.Sprintf("%v:/key.pem", kubernetesKeyPath),
			fmt.Sprintf("giantswarm/kubectl:%v", kubectlVersion),
			fmt.Sprintf("--server=%v", kubernetesApiServer),
			"--certificate-authority=/ca.pem",
			"--client-certificate=/crt.pem",
			"--client-key=/key.pem",
			"cluster-info",
		},
	}

	templatedResourcesDirectoryAbsolutePath, err := filepath.Abs(templatedResourcesDirectoryPath)
	if err != nil {
		log.Fatalf("could not get absolute path for templated resources directory: %v\n", err)
	}
	kubectlApply := commands.Command{
		Name: "kubectl-apply",
		Args: []string{
			"docker",
			"run",
			"--rm",
			"-v", fmt.Sprintf("%v:/ca.pem", kubernetesCaPath),
			"-v", fmt.Sprintf("%v:/crt.pem", kubernetesCrtPath),
			"-v", fmt.Sprintf("%v:/key.pem", kubernetesKeyPath),
			"-v", fmt.Sprintf("%v:/kubernetes", templatedResourcesDirectoryAbsolutePath),
			fmt.Sprintf("giantswarm/kubectl:%v", kubectlVersion),
			fmt.Sprintf("--server=%v", kubernetesApiServer),
			"--certificate-authority=/ca.pem",
			"--client-certificate=/crt.pem",
			"--client-key=/key.pem",
			"apply", "-f", "/kubernetes",
		},
	}

	commands.RunCommands([]commands.Command{
		dockerLogin,
		dockerPush,
		kubectlClusterInfo,
		kubectlApply,
	})

	if removeResourceFilesAfterUse {
		if err := os.RemoveAll(templatedResourcesDirectoryPath); err != nil {
			log.Fatalf("could not remove templated resources directory: %v\n", err)
		}
	}
}