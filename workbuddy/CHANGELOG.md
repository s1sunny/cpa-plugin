# Changelog

## 0.5.0

### Added
- International (Global) WorkBuddy account support (`www.workbuddy.ai` domain)
- Domain-aware billing API routing: CN accounts → `codebuddy.cn`, Global → `workbuddy.ai`
- Expert trial pack claim API: `POST /plugins/workbuddy/trial` (Global only, one-time 250 credits / 14 days)
- Panel region badges: light green `CN` (daily checkin) + light orange `Global` (expert trial)
- "全部领取" batch claim button for Global accounts
- Auto-scheduler region branch: CN → daily checkin, Global → claim expert trial if unclaimed
- `wbAccount.region` and `wbAccount.trial_claimed` fields in accounts API response
- `hasTrialPack()` helper detects trial pack from `get-user-resource` packages

### Changed
- `billingBase` selection is now domain-driven via `billingBaseFor(sa)`
- `backendHeaders` Origin/Referer dynamically set per account domain via `originRefererFor(sa)`
- Panel card buttons: CN → 签到, Global → 领取专家加油包 / 已领取
- "全部签到" button only triggers CN accounts (Global accounts are skipped with a message)
- `runAutoCheckin` branches by region: CN daily checkin, Global trial claim

## 0.4.3

### Changed
- Panel import modal: white surface + dark text for readable contrast (was dark-on-dark)

## 0.4.2

### Changed
- Panel: credential import is a toolbar button (left of 刷新数据) opening a modal, instead of an always-visible card

## 0.4.1

### Added
- Panel **耗尽** badge + `exhausted` field on accounts API (shared with scheduler)
- Credential **import** API `POST /plugins/workbuddy/import` + panel paste UI
- Per-account check-in lock (multi-tab safe)
- `executor.count_tokens` stub (`input_tokens:0` — upstream has no API)
- LICENSE (MIT), VERSION file, GitHub Actions multi-arch release workflow

### Changed
- SSE cleanChunk strips empty `extra_fields` / `refusal` / `reasoning_content`
- Scheduler credits mode prefers non-exhausted accounts first

## 0.4.0

### Added
- CPA **Scheduler** capability with `scheduler_mode`: `off` (default) | `credits`
- Credits-aware multi-account pick using panel credit cache

## 0.3.18

### Fixed
- ConfigFields use SDK `ConfigFieldType*` constants

## 0.3.17

### Fixed
- `FrontendAuthProvider` set false; remove dead frontend-auth handlers

## 0.3.16

### Fixed
- Panel refresh toast + busy feedback

## 0.3.15

### Fixed
- Normalize OpenAI object `tool_choice` for CodeBuddy upstream
