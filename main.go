package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chzyer/readline"
)

const (
	defExtension   = ".yml"
	defKubeconfDir = ".kubeconfig"
	k9sExecutable  = "k9s"
)

type TabComplete struct {
	clusterList []string
}

func main() {
	kubeconfigPath, err := selectKubeconfig()
	if err != nil {
		panic(err)
	}

	err = runK9s(kubeconfigPath)
	if err != nil {
		panic(err)
	}
}

func selectKubeconfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %v", err)
	}

	kubeconfDir := path.Join(homeDir, defKubeconfDir)
	dirItems, err := os.ReadDir(kubeconfDir)
	if err != nil {
		return "", fmt.Errorf("error reading kubeconfig directory: %v", err)
	}

	var kubeconfNames []string
	for _, v := range dirItems {
		if v.IsDir() || !strings.HasSuffix(v.Name(), defExtension) {
			continue
		}
		kubeconfNames = append(kubeconfNames, v.Name())
	}

	if len(kubeconfNames) == 0 {
		return "", fmt.Errorf("no valid kubeconfig files found")
	}

	selectedCluster, err := getUserInput("Enter cluster name: ", kubeconfNames)
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	kubeconfigPath := path.Join(kubeconfDir, selectedCluster) + defExtension

	return kubeconfigPath, nil
}

func getUserInput(prompt string, autoCompleteList []string) (string, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:       prompt,
		AutoComplete: TabComplete{clusterList: autoCompleteList},
	})
	if err != nil {
		return "", err
	}

	userInput, err := l.Readline()
	if err != nil {
		return "", err
	}

	return userInput, nil
}

func (t TabComplete) Do(line []rune, pos int) ([][]rune, int) {
	var nbCharacter = len(line)
	if nbCharacter == 0 {
		// No completion
		return [][]rune{[]rune("")}, 0
	}

	var strArray [][]rune
	for _, s := range t.clusterList {
		if strings.HasPrefix(s, string(line)) {
			var completionCandidate = s[nbCharacter:]
			completionCandidate = strings.TrimSuffix(completionCandidate, defExtension)
			strArray = append(strArray, []rune(completionCandidate))
		}
	}

	return strArray, nbCharacter
}

func runK9s(kubeconfigPath string) error {
	fmt.Println("🐶 Handing off to k9s...")
	cmd := exec.Command(k9sExecutable, "--kubeconfig", kubeconfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
