---
name: neetoform
description: >
  Manage NeetoForm from the command line.
  Use when the user asks about operations exposed by the NeetoForm CLI.
---

## Prerequisites

Run `neetoform doctor` to check authentication and connectivity.
If not authenticated, run `neetoform login`.

## Authentication & multi-subdomain

Credentials for every logged-in subdomain are stored together in
`~/.config/neetoform/auth.json`. A command that talks to the API picks which
subdomain to use by these rules:

- 0 subdomains authenticated → every credential-using command errors with
  "Not authenticated. Run 'neetoform login' to authenticate.".
- 1 subdomain authenticated → that one is the implicit default; `--subdomain`
  may be omitted.
- 2+ subdomains authenticated → **`--subdomain <name>` is required** on every
  credential-using command, including `doctor`. The error lists every
  authenticated subdomain so the agent can offer a choice.

`login` / `logout` / `whoami` have dedicated behavior:

| Command | Behavior |
|---|---|
| `neetoform login --subdomain <name>` | Adds or refreshes the entry for `<name>`. No flag → prompts for the subdomain. |
| `neetoform logout --subdomain <name>` | Removes that one entry. |
| `neetoform logout --all` | Removes every entry. |
| `neetoform logout` (no flag) | Removes the only entry if exactly one is logged in; errors if multiple. |
| `neetoform whoami` | Lists every logged-in account. Marks the entry `(default)` when exactly one. |
| `neetoform whoami --subdomain <name>` | Shows just that one. |

## Global flags (persistent on every command)

| Flag | Purpose |
|---|---|
| `--subdomain <name>` | Select which logged-in subdomain the command targets. Required when multiple are logged in. |
| `--json` | Force JSON envelope output even on a TTY. |
| `--quiet` | Emit only the raw payload — no envelope, no breadcrumbs. For action commands (create/update), emits just the resource identifier; `delete` emits `success`. Designed for scripting. |
| `--toon` | Emit TOON (Token Optimized Output Notation). Preferred for feeding list/show output back to an LLM; ~30–60% fewer tokens than JSON. |

Precedence if multiple are set: `--toon` > `--quiet` > `--json` > pretty.

## Output modes & response envelope

**Pretty (default on a TTY)** — tables for arrays, key-value for objects,
breadcrumbs appended. Not intended for machine consumption.

**JSON envelope** (non-TTY, or `--json`):
```json
{
  "data": <resource body>,
  "breadcrumbs": [{ "label": "List", "command": "neetoform <resource> list" }],
  "pagination": {
    "current_page_number": 1,
    "total_pages": 10,
    "total_records": 250
  }
}
```
`breadcrumbs` is omitted when empty. `pagination` is present only for list
commands.

**Quiet** (`--quiet`) — `data` contents only, no envelope. For action
commands `PrintQuiet` unwraps a single-key wrapper and prints the first of
`sid` / `id` / `name`. For `delete` it prints `success`.

**TOON** (`--toon`) — same data as JSON, re-encoded into TOON. Shape is
equivalent but whitespace/keys are compressed. Parse by re-reading keys as
you would JSON.

### Pagination

List commands accept `--page` (1-indexed) and `--page-size` (max 100).
The envelope's `pagination` field always exposes:
`current_page_number`, `total_pages`, `total_records`. Agents should loop
by incrementing `--page` until `current_page_number == total_pages`.

## Discovery

The full, always-accurate command tree (including any flags added after
this skill was built) is available as JSON:

```bash
neetoform commands
```

Each catalog entry has `command`, `description`, optional `flags` (with
`name`, `type`, `default`, `description`, `required`), and `subcommands`.
Use this whenever a user asks about a flag or command not covered below.

## Diagnostics & IDE setup

| Command | Purpose |
|---|---|
| `doctor` | Auth check + API reachability + version. Uses `--subdomain` when multiple are logged in. |
| `version` | Print CLI version / commit / build date. |
| `update` | Update the CLI to the latest version (auto-detects brew / shell / PowerShell install). |
| `commands` | Emit the full command/flag catalog as JSON. |
| `completion zsh\|bash\|fish\|powershell` | Install shell completion. `--print` emits the script instead. |
| `setup claude` | Install NeetoForm plugin into Claude Code (`plugin.json`, hooks, this SKILL.md). |
| `setup cursor` / `windsurf` / `copilot` / `gemini` / `codex` | Write NeetoForm rule files into the current project directory; safe to re-run. |

## Environment variable override

Set `NEETOFORM_BASE_URL` to point the CLI at a staging or local server:

```bash
export NEETOFORM_BASE_URL=http://acme.lvh.me:8980
neetoform login --subdomain acme
```

## Error surface

Every command exits non-zero on failure and writes a single-line message to
stderr. Common errors the agent should expect:

- `Not authenticated. Run 'neetoform login' to authenticate.` — empty credential store.
- `Multiple subdomains authenticated (acme, beta); specify --subdomain or --all.` — pick one.
- `Not authenticated for "foo". Authenticated subdomains: acme, beta.` — bad `--subdomain`.
- `required flag(s) "xxx" not set` (from cobra) — missing required flag.
- API errors come through with the server's message body; inspect the
  JSON envelope (or the `--quiet` payload) for `error` / `errors` / `notice`
  keys and any suggestions the API returns.

## Product-specific commands

Two resource groups. `neetoform commands` remains the authoritative, always
current catalog; the tables below cover what ships today.

### `forms`

| Command | Purpose |
|---|---|
| `forms list [--status <s>] [--page N] [--page-size N]` | List forms, newest first. |
| `forms submissions list <form-id> [--page N] [--page-size N]` | List completed submissions for one form. |

`--status` takes exactly `active`, `archived` or `favorite`; anything else is
rejected client-side with `Invalid --status "x"; valid values: ...` before a
request is made.

A form record carries `id`, `title`, `state`, `is_published`, `is_archived`,
`is_disabled`, `is_suspended`, `submissions_count`, `attempt_url`, `created_by`,
`created_at` and `updated_at`. The `id` is the `<form-id>` every other command
takes.

A submission carries `id`, `created_at`, `field_values`, `user_agent` and
`responses`. Each entry in `responses` has `label` (the question), `position`,
`kind` (the field type, which decides the shape of `value`), `slug` and `value`.

### `team-members`

| Command | Purpose |
|---|---|
| `team-members list [--email <email>] [--page N] [--page-size N]` | List active members. `--email` matches one exact address. |
| `team-members show <id>` | One member by UUID. |
| `team-members create --emails <a,b> --role <name> [--send-invitation-email=false]` | Invite people. Both `--emails` and `--role` are required. |
| `team-members update <id> [--email] [--first-name] [--last-name] [--time-zone] [--role]` | Change only the flags passed. |
| `team-members delete <id>` | Deactivate a member. |

A member record carries `id`, `email`, `first_name`, `last_name`, `time_zone`,
`profile_image_url`, `active` and `organization_role`.

`--role` must match a role in the workspace exactly, including case; `Admin` and
`Standard` always exist. `create` is all-or-nothing, so a rejected address or
role leaves the workspace unchanged and the call is safe to retry once fixed. An
address that already belongs to the workspace is reactivated and its role
overwritten. `delete` is refused when it would remove the workspace's last admin.

Every command runs with the signed-in user's permissions, so a 403 means the
user's organization role does not allow the operation rather than a bad
credential.

Full reference: https://apidocs.neetoform.com/cli-reference/overview
