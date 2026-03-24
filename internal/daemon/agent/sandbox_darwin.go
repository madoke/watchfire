//go:build darwin

package agent

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const profileTemplate = `(version 1)
(deny default)

; READ ACCESS - Allow all, deny sensitive
(allow file-read* (subpath "/"))

; DENY sensitive credential paths
(deny file-read* (subpath "%s/.ssh"))
(deny file-read* (subpath "%s/.aws"))
(deny file-read* (subpath "%s/.gnupg"))
(deny file-read* (literal "%s/.netrc"))
(deny file-read* (literal "%s/.npmrc"))

; DENY protected user folders (prevents TCC prompts)
(deny file-read* (subpath "%s/Desktop"))
(deny file-read* (subpath "%s/Documents"))
(deny file-read* (subpath "%s/Downloads"))
(deny file-read* (subpath "%s/Music"))
(deny file-read* (subpath "%s/Movies"))
(deny file-read* (subpath "%s/Pictures"))
; NOTE: Keychains allowed - required for Claude Code auth

; WRITE ACCESS - Project + Claude config + caches + temp
(allow file-write* (subpath "%s"))
(allow file-write* (subpath "%s/.claude"))
(allow file-write* (literal "%s/.claude.json"))
(allow file-write* (subpath "%s/Library/Caches/claude-cli-nodejs"))
(allow file-write* (subpath "/tmp"))
(allow file-write* (subpath "/private/tmp"))
(allow file-write* (subpath "/var/folders"))
(allow file-write* (subpath "/private/var/folders"))

; PACKAGE MANAGER CACHES
(allow file-write* (subpath "%s/.npm"))
(allow file-write* (subpath "%s/.yarn"))
(allow file-write* (subpath "%s/.pnpm-store"))
(allow file-write* (subpath "%s/.cache"))
(allow file-write* (subpath "%s/.local/share/pnpm"))
(allow file-write* (subpath "%s/Library/Caches/npm"))
(allow file-write* (subpath "%s/Library/Caches/yarn"))

; CLI TOOL CONFIG (Vercel, Firebase, gcloud, etc.)
(allow file-write* (subpath "%s/Library/Application Support"))

; DEV TOOL CACHES
(allow file-write* (subpath "%s/.cargo"))
(allow file-write* (subpath "%s/go"))
(allow file-write* (subpath "%s/.rustup"))

; PROTECTED - Block writes even in project
(deny file-write* (regex #"/\.env($|[^/]*)"))
(deny file-write* (subpath "%s/.git/hooks"))

; NETWORK, DEVICES, PROCESS, IPC
(allow network*)
(allow file-read* (subpath "/dev"))
(allow file-write* (subpath "/dev"))
(allow process-exec*)
(allow process-fork)
(allow process-info*)
(allow signal)
(allow mach*)
(allow sysctl*)
(allow ipc*)
(allow file-ioctl)
`

// GenerateProfile generates a macOS sandbox-exec profile for the given policy.
// If Trace is true, a (trace ...) directive is prepended to log denied operations.
func GenerateProfile(homeDir, projectDir string, trace bool) string {
	profile := fmt.Sprintf(profileTemplate,
		// DENY sensitive credential paths (5 args: homeDir)
		homeDir, homeDir, homeDir, homeDir, homeDir,
		// DENY protected user folders (6 args: homeDir)
		homeDir, homeDir, homeDir, homeDir, homeDir, homeDir,
		// WRITE ACCESS - Project + Claude config (4 args: projectDir, homeDir, homeDir, homeDir)
		projectDir, homeDir, homeDir, homeDir,
		// PACKAGE MANAGER CACHES (7 args: homeDir)
		homeDir, homeDir, homeDir, homeDir, homeDir, homeDir, homeDir,
		// CLI TOOL CONFIG (1 arg: homeDir)
		homeDir,
		// DEV TOOL CACHES (3 args: homeDir)
		homeDir, homeDir, homeDir,
		// PROTECTED - .git/hooks (1 arg: projectDir)
		projectDir,
	)
	if trace {
		profile = "(trace \"/tmp/watchfire-sandbox-trace.sb\")\n" + profile
	}
	return profile
}

// platformDefaults returns macOS-specific path additions.
func platformDefaults(homeDir string) PlatformDefaults {
	return PlatformDefaults{
		ExtraWritable: []string{
			"/private/tmp",
			"/var/folders",
			"/private/var/folders",
			filepath.Join(homeDir, "Library", "Caches", "claude-cli-nodejs"),
			filepath.Join(homeDir, "Library", "Caches", "npm"),
			filepath.Join(homeDir, "Library", "Caches", "yarn"),
			filepath.Join(homeDir, "Library", "Application Support"),
		},
		ExtraDenied: []string{
			filepath.Join(homeDir, "Desktop"),
			filepath.Join(homeDir, "Documents"),
			filepath.Join(homeDir, "Downloads"),
			filepath.Join(homeDir, "Music"),
			filepath.Join(homeDir, "Movies"),
			filepath.Join(homeDir, "Pictures"),
		},
	}
}

// spawnSandboxedPlatform creates a sandboxed exec.Cmd using macOS sandbox-exec.
func spawnSandboxedPlatform(policy SandboxPolicy, command string, args ...string) (*exec.Cmd, string, error) {
	profile := GenerateProfile(policy.HomeDir, policy.ProjectDir, policy.Trace)

	// Write profile to temp file
	tmpFile, err := os.CreateTemp("", "watchfire-sandbox-*.sb")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create sandbox profile: %w", err)
	}
	if _, err := tmpFile.WriteString(profile); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return nil, "", fmt.Errorf("failed to write sandbox profile: %w", err)
	}
	_ = tmpFile.Close()

	// Build: sandbox-exec -f <tmpfile> <command> <args...>
	sandboxArgs := []string{"-f", tmpFile.Name(), command}
	sandboxArgs = append(sandboxArgs, args...)
	cmd := exec.Command("sandbox-exec", sandboxArgs...)

	cmd.Dir = policy.ProjectDir
	env := buildBaseEnv(policy.ProjectDir)

	// Ensure Homebrew paths are in PATH
	path := os.Getenv("PATH")
	for _, p := range []string{"/opt/homebrew/bin", "/usr/local/bin"} {
		if !strings.Contains(path, p) {
			path = p + ":" + path
		}
	}
	env = setEnv(env, "PATH", path)
	cmd.Env = env

	return cmd, tmpFile.Name(), nil
}

// spawnSandboxedWithBackend routes to the requested backend on macOS.
func spawnSandboxedWithBackend(backend string, policy SandboxPolicy, command string, args ...string) (*exec.Cmd, string, error) {
	switch backend {
	case SandboxSeatbelt:
		return spawnSandboxedPlatform(policy, command, args...)
	default:
		log.Printf("[sandbox] Backend %q not available on macOS, falling back to seatbelt", backend)
		return spawnSandboxedPlatform(policy, command, args...)
	}
}
