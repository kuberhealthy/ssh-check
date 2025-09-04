package main

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	kh "github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// Private key for checking ssh
	sshPrivateKeyEnv = os.Getenv("SSH_PRIVATE_KEY")
	sshPrivateKey    string

	// Username for ssh auth
	usernameEnv = os.Getenv("SSH_USERNAME")
	username    string

	// exclusion list for node names we explicitly don't want to check
	sshExcludeEnv = os.Getenv("SSH_EXCLUDE_LIST")
	sshExclude    string

	// K8s config file for the client.
	kubeConfigFile = filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// K8s client used for the check.
	client *kubernetes.Clientset
)

func init() {
	parseInputValues()
}

func main() {
	var failed bool
	errs := []string{}

	ctx := context.Background()

	client, err := createClient(kubeConfigFile)
	if err != nil {
		errs = append(errs, err.Error())
		err = kh.ReportFailure(errs)
		if err != nil {
			log.Fatalln("error reporting to kuberhealthy: ", err)
		}
	}

	log.Infoln("Kubernetes client created.")

	// get list of nodes from cluster
	allNodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})

	// parse nodes filtering out not ready nodes/excluded nodes, and ssh to them
	for _, node := range allNodes.Items {
		if check_excludes(node) {
			continue
		}

		// parse conditions and filtering out any nodes that aren't ready
		for _, condition := range node.Status.Conditions {
			if condition.Type != "Ready" {
				continue
			}
			if condition.Status == "True" {
				err = ssh_check(node)
				if err != nil {
					errs = append(errs, err.Error())
					err = kh.ReportFailure(errs)
					if err != nil {
						log.Fatalln("error reporting to kuberhealthy: ", err)
					}
					failed = true
				}
			}
		}

		// report success if no issues were found during ssh checks
		if !failed {
			err = kh.ReportSuccess()
			if err != nil {
				log.Infoln("error reporting success to kuberhealthy: ", err)
			}
		}
	}
}

// parseInputValues parses all incoming environment variables for the program into globals and fatals on errors.
func parseInputValues() {
	if len(sshPrivateKeyEnv) == 0 {
		log.Fatalln("SSH key required")
	}
	sshPrivateKey = sshPrivateKeyEnv
	log.Infoln("Parsed SSH_PRIVATE_KEY")

	if len(usernameEnv) == 0 {
		log.Fatalln("Username required")
	}
	username = usernameEnv
	log.Infoln("Parsed SSH_USERNAME: ", username)

	if len(sshExcludeEnv) != 0 {
		sshExclude = sshExcludeEnv
		log.Infoln("Parsed SSH_EXCLUDE_LIST: ", sshExclude)
	}

}

// find nodes we want to exclude, return true if it's found in exclude list
func check_excludes(node v1.Node) bool {
	for _, host := range strings.Split(sshExclude, " ") {
		if host == node.Name {
			return true
		}
	}
	return false
}

// ssh to node IP, and report results
func ssh_check(node v1.Node) error {
	signer, err := ssh.ParsePrivateKey([]byte(sshPrivateKey))
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	log.Infoln("Attempting Connection to: ", node.Name)
	// get ip of node, ssh to it
	for _, address := range node.Status.Addresses {
		if address.Type == "InternalIP" {
			conn, err := ssh.Dial("tcp", address.Address+":22", config)
			if err != nil {
				return err
			}

			defer conn.Close()
		}
	}

	return nil
}

func createClient(kubeConfigFile string) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
