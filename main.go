package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "checkout-main-and-update":
		if err := checkoutMainAndUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "setup-gitignore-local":
		if err := setupGitignoreLocal(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: michael-cmd <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  checkout-main-and-update    checkout main/master and pull latest")
	fmt.Fprintln(os.Stderr, "  setup-gitignore-local       configure a repo-local .gitignore.local that never pushes to the remote")
}

func checkoutMainAndUpdate() error {
	if err := assertGitRepo(); err != nil {
		return err
	}

	if err := assertNoUnstagedChanges(); err != nil {
		return err
	}

	branch, err := detectDefaultBranch()
	if err != nil {
		return err
	}

	fmt.Printf("switching to %s...\n", branch)
	if err := git("checkout", branch); err != nil {
		return fmt.Errorf("checkout failed: %w", err)
	}

	fmt.Printf("pulling latest %s...\n", branch)
	if err := git("pull"); err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}

	fmt.Printf("up to date on %s\n", branch)
	return nil
}

func setupGitignoreLocal() error {
	if err := assertGitRepo(); err != nil {
		return err
	}

	root, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("could not determine repository root: %w", err)
	}

	localIgnore := filepath.Join(root, ".gitignore.local")

	// point git's per-repo excludes file at .gitignore.local so its patterns
	// are honoured without living in the tracked .gitignore
	fmt.Printf("configuring core.excludesFile -> %s\n", localIgnore)
	if err := git("config", "--local", "core.excludesFile", localIgnore); err != nil {
		return fmt.Errorf("failed to set core.excludesFile: %w", err)
	}

	if _, err := os.Stat(localIgnore); err == nil {
		fmt.Printf("%s already exists, leaving it untouched\n", localIgnore)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("could not check %s: %w", localIgnore, err)
	}

	// seed the file so it ignores itself and never reaches the remote
	fmt.Printf("creating %s\n", localIgnore)
	if err := os.WriteFile(localIgnore, []byte(".gitignore.local\n"), 0o644); err != nil {
		return fmt.Errorf("failed to create %s: %w", localIgnore, err)
	}

	fmt.Println("done — add local-only ignore patterns to .gitignore.local")
	return nil
}

func assertGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not inside a git repository")
	}
	return nil
}

func assertNoUnstagedChanges() error {
	// staged but not committed
	staged, err := gitOutput("diff", "--cached", "--name-only")
	if err != nil {
		return fmt.Errorf("could not check staged changes: %w", err)
	}

	// modified but not staged
	unstaged, err := gitOutput("diff", "--name-only")
	if err != nil {
		return fmt.Errorf("could not check unstaged changes: %w", err)
	}

	// untracked files
	untracked, err := gitOutput("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return fmt.Errorf("could not check untracked files: %w", err)
	}

	var problems []string
	if staged != "" {
		problems = append(problems, fmt.Sprintf("staged changes (not yet committed):\n%s", indent(staged)))
	}
	if unstaged != "" {
		problems = append(problems, fmt.Sprintf("unstaged changes:\n%s", indent(unstaged)))
	}
	if untracked != "" {
		problems = append(problems, fmt.Sprintf("untracked files:\n%s", indent(untracked)))
	}

	if len(problems) > 0 {
		msg := "refusing to switch branches — working tree is not clean:\n\n"
		for _, p := range problems {
			msg += p + "\n\n"
		}
		msg += "commit, stash, or discard your changes before running this command."
		return fmt.Errorf("%s", msg)
	}

	return nil
}

func detectDefaultBranch() (string, error) {
	// prefer the remote's HEAD symref (works even before first local checkout)
	out, err := gitOutput("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		// refs/remotes/origin/HEAD -> refs/remotes/origin/main
		parts := strings.Split(strings.TrimSpace(out), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// fall back to checking local branch existence
	for _, candidate := range []string{"main", "master"} {
		check := exec.Command("git", "rev-parse", "--verify", candidate)
		check.Stderr = nil
		if check.Run() == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not detect default branch (expected 'main' or 'master')")
}

func git(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func indent(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}
