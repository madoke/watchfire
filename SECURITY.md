# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | Yes                |

## Reporting a Vulnerability

If you discover a security vulnerability in Watchfire, please report it
responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please send an email to **security@watchfire.io** with:

- A description of the vulnerability
- Steps to reproduce the issue
- Any potential impact you've identified
- Suggested fix (if you have one)

## Response Timeline

- **Acknowledgment**: Within 48 hours of your report
- **Initial assessment**: Within 1 week
- **Fix and disclosure**: We aim to release a fix within 30 days of a confirmed
  vulnerability, depending on complexity

## Disclosure Policy

- We will coordinate disclosure with you
- We will credit reporters in our release notes (unless you prefer anonymity)
- We ask that you give us reasonable time to address the issue before public
  disclosure

## Security Best Practices

Watchfire runs coding agents in sandboxed environments. Key security
considerations:

- Agents are sandboxed to restrict filesystem access
- Agents cannot access `~/.ssh`, `~/.aws`, `~/.gnupg`, or `.env` files
- `.git/hooks` are blocked to prevent hook injection (Seatbelt only)
- All agent sessions run in isolated git worktrees

### Sandbox Platform Matrix

| Platform | Backend | Dependencies | Capabilities |
|----------|---------|-------------|--------------|
| **macOS** | Seatbelt (`sandbox-exec`) | Built-in | Full policy: read-only FS, writable exceptions, denied paths, regex file patterns, network allowed |
| **Linux (kernel 5.13+)** | Landlock LSM | None (kernel feature) | Read-only FS, writable exceptions via rules. No regex file patterns. Network allowed. |
| **Linux (older)** | Bubblewrap (`bwrap`) | `apt install bubblewrap` | Read-only FS via mount namespaces, writable bind mounts, denied paths via tmpfs overlay. No regex file patterns. |
| **Linux (no sandbox)** | None | — | Warning logged, agent runs unsandboxed |
| **Windows** | None | — | Agent runs unsandboxed — no sandbox available |
| **Other OS** | None | — | Warning logged, agent runs unsandboxed |

### Sandbox Selection

Sandbox backend can be configured at three levels (highest priority wins):

1. **CLI flag**: `--sandbox <backend>` or `--no-sandbox`
2. **Project setting**: `project.yaml` → `sandbox: "auto"`
3. **Global setting**: `settings.yaml` → `defaults.default_sandbox: "auto"`

Valid backends: `auto`, `seatbelt`, `landlock`, `bwrap`, `none`

### Known Limitations

- **Linux**: `.env` file write protection and `.git/hooks` blocking are not enforced
  (Seatbelt supports regex file patterns; Landlock and bubblewrap do not)
- **Landlock**: Requires kernel 5.13+ with `CONFIG_SECURITY_LANDLOCK=y`
- **Bubblewrap**: Requires user namespaces enabled (default on modern distros)

If you find a way to bypass these protections, we consider that a critical
security issue and ask that you report it immediately.
