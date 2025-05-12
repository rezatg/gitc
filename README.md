<div align="center">
  <img src="https://github.com/rezatg/gitc/blob/master/assets/logo.jpg" alt="gitc AI-Powered Commits" style="clip-path: inset(35px 0 35px 0);margin: 0; padding: 0px, border-radius: 5px;box-shadow: 0 4px 8px rgba(0,0,0,0.2);">
</div>

# ✨ gitc - AI-Powered Git Commit Messages

[![Go Reference](https://pkg.go.dev/badge/github.com/rezatg/gitc)](https://pkg.go.dev/github.com/rezatg/gitc)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rezatg/gitc?logo=go)](go.mod)
[![Sourcegraph](https://sourcegraph.com/github.com/rezatg/gitc/-/badge.svg)](https://sourcegraph.com/github.com/rezatg/gitc?badge)
[![Discussions](https://img.shields.io/github/discussions/rezatg/gitc?color=58a6ff&label=Discussions&logo=github)](https://github.com/rezatg/gitc/discussions)
[![Downloads](https://img.shields.io/github/downloads/rezatg/gitc/total?color=blue)](https://github.com/rezatg/gitc/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/rezatg/gitc)

<div align="center">
  <a href="#-installation">Installation</a> •
  <a href="#-features">Features</a> •
  <a href="#-configuration">Configuration</a> •
  <a href="#-basic-usage">Usage</a> •
  <a href="#-full-options">Full Options</a> •
  <a href="#-ai-providers">AI Providers</a>
</div>
<br>

> `gitc` is a fast, lightweight CLI tool that uses AI to generate clear, consistent, and standards-compliant commit messages — directly from your Git diffs. With built-in support for [Conventional Commits](https://www.conventionalcommits.org), [Gitmoji](https://gitmoji.dev), and fully customizable rules, `gitc` helps you and your team write better commits, faster

# 🚀 Features
gitc streamlines your Git workflow by automating professional commit message creation with AI. Its robust feature set ensures flexibility and precision for developers and teams.

- ### 🧠 AI and Commit Generation
  - **AI-Powered Commit Messages**: Generates high-quality commit messages using OpenAI's API, analyzing staged git changes for context-aware results.
  - **Multilingual Support**: Creates commit messages in multiple languages (e.g., English, Persian, Russian) to suit global teams.
  - **Extensible AI Providers**: Supports OpenAI with plans for Anthropic and other providers, ensuring future-proofing.

- ### 📝 Commit Standards and Customization
  - **Conventional Commits**: Adheres to standard commit types (`feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`, `revert`, `init`, `security`) for semantic versioning.
  - **Gitmoji Integration**: Optionally adds Gitmoji emojis (e.g., ✨ for `feat`, 🚑 for `fix`) for visually appealing commits.
  - **Custom Commit Conventions**: Supports JSON-based custom prefixes (e.g., JIRA ticket IDs) for tailored commit formats.

- ### 🔧 Git Integration
  - **Optimized Git Diff Processing**: Automatically retrieves and filters staged git diff, excluding irrelevant files (e.g., `node_modules/*`, `*.lock`).
  - **Configurable Exclusions**: Customize file exclusion patterns via config file to focus on relevant changes.

- ### ⚙️ Environment & Configuration
  - **Flexible Configuration**: Customize via CLI flags, environment variables, or a JSON config file (`~/.gitc/config.json`).
  - **Proxy Support**: Configurable proxy settings for API requests in restricted environments.
  - **Timeout and Redirect Control**: Adjustable timeouts and HTTP redirect limits for reliable API interactions.
  - **Environment Variable Support**: Simplifies setup for sensitive data (e.g., API keys) and common settings.

- ### ⚡️ Performance & Reliability
  - **Fast Processing**: Leverages [sonic](https://github.com/bytedance/sonic) for rapid JSON parsing and [fasthttp](https://github.com/valyala/fasthttp) for efficient HTTP requests.
  - **Error Handling**: Robust validation and error messages ensure reliable operation.

## 📦 Installation
### Prerequisites:
  - Go: Version **1.18** or higher (required for building from source).
  - Git: Required for retrieving staged changes.
  - OpenAI API Key: Required for AI-powered commit message generation. Set it via the `AI_API_KEY` environment variable or in the config file.

#### Quick Install:
  ```bash
  go install github.com/rezatg/gitc@latest
  ```

### Manual Install
  1. Download binary from [releases](https://github.com/rezatg/gitc/releases)
  2. `chmod +x gitc`
  3. Move to `/usr/local/bin`

### Verify Installation
  After installation, verify the tool is installed correctly and check its version:
  ```bash
  gitc --version
  ```

# 💻 Basic Usage
```bash
# 1. Stage your changes
git add . # or gitc -a

# 2. Generate perfect commit message
gitc

# Pro Tip: Add emojis and specify language
gitc --emoji --lang fa

# Custom commit type
gitc --commit-type fix
```

## Environment Variables
```bash
export OPENAI_API_KEY="sk-your-key-here"
export GITC_LANGUAGE="fa"
export GITC_MODEL="gpt-4"
```

# ⚙️ Configuration
Config File (`~/.gitc/config.json`) :
```json
{
  "provider": "openai",
  "max_length": 200,
  "proxy": "",
  "language": "en",
  "timeout": 10,
  "commit_type": "",
  "custom-convention": "",
  "use_gitmoji": false,
  "max_redirects": 5,
  "open_ai": {
    "api_key": "sk-your-key-here",
    "model": "gpt-4o-mini",
    "url": "https://api.openai.com/v1/chat/completions"
  }
}
```

# 📚 Full Options
The following CLI flags are available for the `ai-commit` command and its `config` subcommand. All flags can also be set via environment variables or the `~/.gitc/config.json` file.

| Flag | Alias | Description | Default | Environment Variable | Example |
|------|-------|-------------|---------|----------------------|---------|
| `--all` | `-a` | Stage all changes before generating commit message (equivalent to `git add .`) | `false` | `GITC_STAGE_ALL` | `-all` or `-a`
| `--provider` | - | AI provider to use (e.g., `openai`, `anthropic`) | `openai` | `AI_PROVIDER` | `--provider openai` |
| `--model` | - | OpenAI model for commit message generation | `gpt-4o-mini` | - | `--model gpt-4o` |
| `--lang` | - | Language for commit messages (e.g., `en`, `fa`, `ru`) | `en` | `GITC_LANGUAGE` | `--lang fa` |
| `--timeout` | - | Request timeout in seconds | `10` | - | `--timeout 15` |
| `--maxLength` | - | Maximum length of the commit message | `200` | - | `--maxLength 150` |
| `--api-key` | `-k` | API key for the AI provider | - | `AI_API_KEY` | `--api-key sk-xxx` |
| `--proxy` | `-p` | Proxy URL for API requests | - | `GITC_PROXY` | `--proxy http://proxy.example.com:8080` |
| `--commit-type` | `-t` | Commit type for Conventional Commits (e.g., `feat`, `fix`) | - | `GITC_COMMIT_TYPE` | `--commit-type feat` |
| `--custom-convention` | `-C` | Custom commit message convention (JSON format) | - | `GITC_CUSTOM_CONVENTION` | `--custom-convention '{"prefix": "JIRA-123"}'` |
| `--emoji` | `-g` | Add Gitmoji to the commit message | `false` | `GITC_GITMOJI` | `--emoji` |
| `--no-emoji` | - | Disables Gitmoji in commit messages (overrides `--emoji` and config file) | `false` | - | `--no-emoji`
| `--max-redirects` | `-r` | Maximum number of HTTP redirects | `5` | `GITC_MAX_REDIRECTS` | `--max-redirects 10` |
| `--config` | `-c` | Path to the configuration file | `~/.gitc/config.json` | `GITC_CONFIG_PATH` | `--config ./my-config.json` |

> [!NOTE]
> - Flags for the `config` subcommand are similar but exclude defaults, as they override the config file.
> - **Flags** > **Environment Variables** > **Config File** — This is the order of precedence when multiple settings are provided.
> - The `--custom-convention` flag expects a JSON string with a `prefix` field (e.g., `{"prefix": "JIRA-123"}`).
> - The `--version` flag displays the current tool version (e.g., `0.1.0`) and can be used to verify installation.
> - The `--all` flag (alias `-a`) stages all changes in the working directory before generating the commit message, streamlining the workflow. For example, `gitc -a --emoji` stages all changes and generates a commit message with Gitmoji.
> - Environment variables take precedence over config file settings but are overridden by CLI flags.
> - You can reset all configuration values to their defaults by using gitc config `gitc reset-config`.


## 🤖 AI Providers
`gitc` is designed to be AI-provider agnostic. While it currently supports OpenAI out of the box, support for additional providers is on the roadmap to ensure flexibility and future-proofing.

| Provider | Supported Models | Required Configuration | Status |
| --- | --- | --- | --- |
| **OpenAI** | `gpt-4o`, `gpt-4o-mini`, `gpt-3.5-turbo` | `api_key`, `model`, `url` (optional) | ✅ Supported (default) |
| **DeepSeek** | Coming Soon | - | 🔜 Planned |
| **Grok (xAI)** | Coming Soon | - | 🔜 Planned |
| **Gemini (Google)** | Coming Soon | - | 🔜 Planned |
| **Others** | - | - | 🧪 Under consideration |
> ℹ️ We're actively working on supporting multiple AI backends to give you more control, flexibility, and performance. Have a provider you'd like to see? [Open a discussion](https://github.com/rezatg/gitc/discussions)!

## 🤝 Contributing

We welcome contributions! Please check out the [contributing guide](CONTRIBUTING.md) before making a PR.

## ⭐️ Star History
[![Star History Chart](https://api.star-history.com/svg?repos=rezatg/gitc&type=Date)](https://www.star-history.com/#rezatg/gitc&Date)
