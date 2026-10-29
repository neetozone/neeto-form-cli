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
| `--verbose` | Expand every field of every record instead of printing a table. Pretty output only. |

Precedence if multiple are set: `--toon` > `--quiet` > `--json` > pretty.
`--verbose` shapes pretty output alone, so it changes nothing under the other three.

## Output modes & response envelope

**Pretty (default on a TTY)** — tables for arrays, key-value for objects,
breadcrumbs appended. Not intended for machine consumption. A field holding a
list of `label`/`value` pairs, which is how a submission carries its answers,
becomes one column per label, in the order the API returned them; the record's
own fields claim their columns first and the answers fill the remaining ones, up
to the shared renderer's column limit.

A table therefore drops whatever does not fit, so re-run the command with
`--verbose` whenever the answer being looked for is not in the table. That is
the only pretty output carrying a later question on a long form, and the only
one carrying a nested record field such as a submission's `user_agent`. Under
`--json` or `--toon` every field is already there, so no re-run is needed.

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
| `setup cursor` / `windsurf` / `copilot` / `gemini` / `codex` | Write NeetoForm rule files into the current project directory; re-run after an upgrade to refresh them. |

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
- `Multiple subdomains authenticated (acme, beta); specify --subdomain.` — pick one.
- `Not authenticated for "foo". Authenticated subdomains: acme, beta.` — bad `--subdomain`.
- `required flag(s) "xxx" not set` (from cobra) — missing required flag.
- API errors come through with the server's message body; inspect the
  JSON envelope (or the `--quiet` payload) for `error` / `errors` / `notice`
  keys and any suggestions the API returns.

## Product-specific commands

`neetoform commands` remains the authoritative, always current catalog, with every
flag. The sections below cover what ships today and what to know before running it.

Form-scoped resources live under `forms`: `neetoform forms <resource> <verb> <form-id> [<id>]`.
Workspace resources are top level. Enum flags are checked client-side, so a bad value fails
before any request with `Invalid --<flag> "x". Valid values: ...`. Commands that take a nested
payload (question attributes, conditions, rule actions, theme properties) read it from
`--json-file <path>`; scalar flags passed alongside override the file's keys. Delete commands
run at once with no confirmation prompt, so confirm with the user before running one.

### `forms`

| Command | Purpose |
|---|---|
| `forms list [--status <s>] [--search <words>] [--order asc\|desc] [--page N] [--page-size N]` | List forms, oldest first unless `--order desc`. `--search` keeps forms whose title contains every word, archived ones included unless `--status` narrows. |
| `forms show <form-id>` | One form with its `questions`, in form order. |
| `forms submissions list <form-id> [--submitted-after <date>] [--submitted-before <date>] [--order asc\|desc] [--page N] [--page-size N]` | Completed submissions, oldest first unless `--order desc`. `--submitted-after` is inclusive and `--submitted-before` exclusive; both take `YYYY-MM-DD` (midnight UTC) or a full ISO 8601 date-time. |

`--status` takes exactly `active`, `archived` or `favorite`, and `--order` exactly `asc` or
`desc`; anything else is rejected client-side before a request is made.

Leaving out `--page` returns every matching record in one response. Pass `--page` when you want
one bounded page, such as `--order desc --page 1 --page-size 100` for the newest 100.

A form record carries `id`, `title`, `state`, `is_published`, `is_archived`, `is_disabled`,
`is_suspended`, `submissions_count`, `attempt_url`, `created_by`, `created_at` and `updated_at`.
`forms show` adds `questions`; each has `id`, `label`, `kind`, `slug`, `position`, `is_required`
and `is_hidden`.

A submission carries `id`, `created_at`, `field_values`, `user_agent` and `responses`. Each
entry in `responses` has `label`, `position`, `kind`, `slug` and `value`, in the form's current
question order. A question the person skipped has no entry, so to build columns, key answers by
`slug` and take the column list from `forms show`.

#### Submissions

| Command | Purpose |
|---|---|
| `forms submissions list <form-id> [--status completed\|incomplete\|all] [--search T] ...` | `--status` defaults to `completed` on the server; `--search` matches answer text. |
| `forms submissions show <form-id> <submission-id>` | One submission. |
| `forms submissions delete <form-id> <submission-id>` / `bulk-delete <form-id> --ids a,b` | Move to the trash (restorable in the app). |
| `forms submissions send-responses <form-id> --ids a,b --email-question-id <q> --subject S --body <html> [--cc-emails] [--bcc-emails]` | Emails respondents a copy of their answers. Confirm content and recipients first. |
| `forms submissions export <form-id> --type csv\|xlsx\|pdf [--ids a,b] [--recipient-email E] [--include-hidden-fields] [--include-meta-data] [--include-question-descriptions]` | Emails a download link to a team member; never returns the file. With an API key `--recipient-email` is required. |


#### Form lifecycle

| Command | Purpose |
|---|---|
| `forms create --title <t> [--style single_page\|one_question_at_a_time] [--created-by-email <e>]` | New blank draft form. `--created-by-email` is needed only with an API key. |
| `forms update <form-id> --title <t>` | Rename. |
| `forms clone <form-id>` | Copy a form as a new draft. |
| `forms publish <form-id>` / `forms reset <form-id>` | Publish the draft, or discard draft changes back to the published version. |
| `forms archive\|unarchive\|enable\|disable\|favorite\|unfavorite <form-id>...` | Bulk state change, one or more ids. `enable`/`disable` replace any deadline or limit. |
| `forms delete <form-id>` / `forms bulk-delete <form-id>...` | Move forms to the trash. |
| `forms counts` | Form counts by status. |
| `forms transfer <form-id>... --to <subdomain>` | Move forms to another workspace. Needs a CLI login (not an API key) and admin rights in the destination. |
| `forms insights <form-id> [--start-date D --end-date D] [--time-zone Z]` | Form insights (visits, starts, completions) for a date range. Pass both dates or neither. |

Bulk commands are all or nothing: an unknown or forbidden id changes nothing. Question, logic and
settings edits change the draft; run `forms publish` for respondents to see them.

#### Questions

| Command | Purpose |
|---|---|
| `question-kinds list` | Every question kind and the attributes it accepts. Read this before building a question. |
| `forms questions list\|show <form-id> [<question-id>]` | Read questions. |
| `forms questions create <form-id> --kind <k> [--position N] [--label] [--description] [--placeholder] [--required] [--json-file <attrs.json>]` | Add one question. |
| `forms questions bulk-create <form-id> --json-file <file>` | File holds `{"questions":[{"kind":"email","attributes":{"label":"Work email"}}]}`. |
| `forms questions update <form-id> <question-id> [same attribute flags]` | Change attributes. Options, rows and products are replaced as a list: keep an item by passing its `id`. |
| `forms questions clone\|delete <form-id> <question-id>` | Copy or remove. |
| `forms questions reorder <form-id> --ids a,b,c` | Needs every element id of the form, including logo, welcome screen and submit button. |

#### Conditional logic and automation rules

| Command | Purpose |
|---|---|
| `forms conditions list\|show\|delete\|clone <form-id> [<id>]` | Show/hide, calculation and thank-you-page conditions. |
| `forms conditions create\|update <form-id> [<id>] --json-file <file> [--enabled]` | File holds `{"conditions":{"join_type":"and","items":[...]},"actions":{"items":[...]}}`, both objects, not arrays; each is replaced as a whole on update. A choice question's `value` is the option's `alt_id` from `forms questions show`. |
| `forms conditions reorder <form-id> --ids a,b` | Every condition id in the new order. |
| `forms automation-rules list <form-id> [--search T]`, `show\|delete\|clone <form-id> <id>`, `reorder <form-id> --ids a,b` | Manage rules. |
| `forms automation-rules create\|update <form-id> [<id>] --name N [--description] [--status enabled\|disabled] [--performer any\|<role>] [--events submission_created,custom_field_updated] [--delay-value N --delay-unit minutes\|hours\|days] [--json-file <file>]` | Rule `conditions` and `actions` come from the file. Actions email, text and call endpoints on real submissions. `show` never returns API action header values, so resend them on update. |

#### Settings (all under `forms`, all take `<form-id>`)

| Command | Purpose |
|---|---|
| `thank-you-pages list\|show\|create\|update\|delete` | `--kind custom_message\|redirect_to_url`, `--message`, `--redirect-url`, `--primary` and more. The primary page cannot be deleted. |
| `email-notifications show\|update <form-id> <admin\|respondent>` | `--enabled`, `--subject`, `--message`, `--to/--cc/--bcc a,b`, `--reply-to-question-id` (admin), `--reply-to` and `--send-to-question-id` (respondent). |
| `unique-submission show\|update --method no_track\|cookie_track\|ip_track\|record_track [--question-id]` | Duplicate prevention. |
| `availability show\|update [--method enable\|disable\|based_on_submission_count\|based_on_time\|based_on_limit_and_time] [--deadline] [--limit] [--closed-message]` | Open, close, deadline, limit. |
| `access-control show\|update --auth no\|basic\|email\|session\|invited [--password] [--allowed-keywords a,b]` | Gate the form. The password is never returned. |
| `access-control-invitations list\|create\|delete --emails a,b` | Invitee list for `--auth invited`. No email is sent. |
| `discount-codes list\|update [--ids a,b]` | Which workspace codes the form accepts. `update` replaces the list and requires `--ids`; `--ids=` clears it. |
| `tracking-codes show\|update [--google-analytics] [--google-tag-manager] [--meta-pixel]` | Pass `""` to remove one. |
| `neeto-chat show\|update --enabled=<bool>` | NeetoChat widget. |
| `message-templates list [--type email\|sms] [--status active\|inactive]`, `show\|delete <id>`, `create --name --type --body [--subject] [--status]`, `update <id>` | Email and SMS templates. |
| `theme update --theme-id <id>` | Apply a theme. |
| `webhooks list\|create --endpoint URL [--events submit] [--secret S]` | Form webhooks. |

#### Sharing

| Command | Purpose |
|---|---|
| `forms share show <form-id> [--embed inline\|floating-popup\|element-popup\|iframe]` | Public link, custom domain links and embed codes; `--embed` prints one snippet. Links work only on a published form. |
| `forms public-link update <form-id> [--public-link-id <slug>]` | New public link; the old URL stops working at once, so read `share show` first. |
| `forms prefilled-links create <form-id> --answer <question-id>=<value>... [--answers-json '<array>'] [--discount-code C]` | Build a prefilled URL. Nothing is saved. Option labels must match exactly. |
| `forms email-invitations create <form-id> --emails a,b [--subject] [--body] [--button-text] [--reply-to] [--from] [--send-at <iso>]` | Emails real people. Confirm recipients and never retry a success. |

### Workspace resources

| Command | Purpose |
|---|---|
| `discount-codes list\|show\|delete`, `create --code C --value V [--kind percentage\|fixed] [--currency] [--expires-at] [--redeems-limit N] [--form-ids a,b]`, `update <id> [same flags]` | Workspace discount codes. On update, `--form-ids` replaces the list and leaving it out keeps it. |
| `discount-codes forms list` | Forms a code can apply to, with their currency; a fixed code works only on forms in its currency. |
| `discount-codes redemptions list <id>` / `discount-codes report <id> [--start-date] [--end-date]` | Individual redemptions, or totals. |
| `templates list [--search] [--category]`, `templates show <id-or-slug>` | Form templates. |
| `templates use <template-id> [--created-by-email E]` | New draft form from a template. `--created-by-email` only with an API key. |
| `roles list` | Organization roles, for `team-members --role`. |
| `audit-logs list [--created-after <iso>] [--created-before <iso>]` | Workspace activity; after is inclusive, before exclusive. |
| `themes list\|show\|delete`, `themes create --name N [--base-theme-id ID] [--property key=value]... [--background-image SIGNED_ID] [--json-file]`, `themes update <id> [--background-image SIGNED_ID \| --remove-background-image]` | System themes cannot change; base a new one on them. Changing a theme changes every form using it. Colours are hex. Background properties: `is_content_background_enabled` and `is_dominant_background_color_enabled` (`true`/`false`), `content_background_color` and `dominant_color` (hex), `background_image_position` (`center`, `top left`, `top center`, `top right`, `center left`, `center right`, `bottom left`, `bottom right`, `bottom center`), `background` (image overlay, integer -100 to 100, e.g. `--property background=-40`). |
| `images upload <path> [--content-type T]` | Upload a LOCAL image (png, jpeg, gif or svg, at most 1 MB) and print its `signed_id`. It reserves the upload, then sends the bytes to the storage address itself; that address is a credential and is never printed. Checks the file before any request. |
| `images create --url <https-url> [--filename N]` | Store an image from a public https URL and print its `signed_id`. Use `images upload` for a local file. |
| `webhooks list\|show\|delete`, `webhooks create --endpoint URL [--events submit] [--secret S]` | Workspace webhooks. The secret is never shown again. |
| `ip-restrictions show` | IP restriction status and allowed ranges. |
| `custom-domains list\|show\|create\|delete <hostname>` | After `create`, add a CNAME to `dns_target` (or an A record to `dns_target_ip` for a root domain) and poll `show` until `active`. |
| `api-keys list [--search]`, `api-keys delete <id>` | `delete` revokes a key and cannot be undone; CLI login only. Tokens are never printed. |

### `team-members`

| Command | Purpose |
|---|---|
| `team-members list [--email <email>] [--sort <field>] [--order asc\|desc] [--page N] [--page-size N]` | List active members. `--email` matches one exact address. `--sort` takes `created_at`, `updated_at`, `email`, `first_name` or `last_name`; `--order` alone sorts by `created_at`. |
| `team-members show <id>` | One member by UUID. |
| `team-members create --emails <a,b> --role <name> [--invited-by <admin-email>] [--send-invitation-email=false]` | Invite people. Both `--emails` and `--role` are required. |
| `team-members update <id> [--email] [--first-name] [--last-name] [--time-zone] [--role]` | Change only the flags passed. |
| `team-members delete <id>` | Deactivate a member. |
| `deactivated-team-members list [--page N] [--page-size N]` | Deactivated members. |
| `deactivated-team-members reactivate <id>` | Restore a deactivated member. |

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

### See what a form looks like

`neetoform forms outline <form-id> [--draft]` prints a compact text outline of the form as respondents see it: the theme (font, colours, background image, card, buttons), every element in order with its kind, required mark, placeholder and options, the logic rules in plain words, and the thank-you page. Read it after building or changing a form to check the result; pass `--draft` until the form is published.

### Theme background image flow

1. Get a `signed_id`: `neetoform images upload ./bg.png` (local file) or `neetoform images create --url https://example.com/bg.png`.
2. Attach it: `neetoform themes update <theme-id> --background-image <signed_id>` (also on `themes create`).
3. Tune it: `--property background_image_position="top left" --property background=-40`. `background_image_position` is one of center, top left, top center, top right, center left, center right, bottom left, bottom right or bottom center; `background` is the overlay, an integer from -100 (darker) to 100 (lighter); `is_dominant_background_color_enabled` (true or false) needs `dominant_color` (hex); `is_content_background_enabled` (true or false) with `content_background_color` (hex, alpha allowed).
4. Remove it: `neetoform themes update <theme-id> --remove-background-image`.
5. Scripting: `images upload` and `images create` with `--quiet` print only the signed_id.

The `signed_id` is the only value to carry between steps. Never print or log an upload address.
