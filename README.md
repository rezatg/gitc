# gitc - ✨ AI-Powered Git Commit Messages
[![Go Reference](https://pkg.go.dev/badge/github.com/rezatg/gotube#section-readme.svg)](https://pkg.go.dev/github.com/rezatg/gitc)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rezatg/gitc?logo=go)](go.mod)
[![Sourcegraph](https://sourcegraph.com/github.com/rezatg/gitc/-/badge.svg)](https://sourcegraph.com/github.com/rezatg/gitc?badge)
[![Discussions](https://img.shields.io/github/discussions/rezatg/gitc?color=58a6ff&label=Discussions&logo=github)](https://github.com/rezatg/gitc/discussions)


**AI-Commit is a command-line tool that leverages AI to generate professional Git commit messages based on staged changes. It supports Conventional Commits, Gitmoji, and customizable commit message conventions, making it ideal for developers who want to streamline their commit workflow. Powered by OpenAI, it analyzes your git diff and produces clear, concise, and context-aware commit messages.**

<br>
<p align="center">
  <img src="./logo.png" alt="logo project" style="height: auto border-radius: 5px;box-shadow: 0 4px 8px rgba(0,0,0,0.2);">
</p>


## 📦 Installation
### Prerequisites :
- Go: Version **1.18** or higher (required for building from source).
- Git: Required for retrieving staged changes.
- OpenAI API Key: Required for AI-powered commit message generation. Set it via the `AI_API_KEY` environment variable or in the config file.

#### Using Go:
```bash
go install github.com/rezatg/gitc@latest
```

### Manual Install
1. Download binary from [releases](https://github.com/rezatg/gitc/releases)
2. `chmod +x gitc`
3. Move to `/usr/local/bin`

## Verify Installation
After installation, verify the tool is installed correctly and check its version:

```bash
gitc --version
```

## 🚀 Features
AI-Commit streamlines your Git workflow by automating professional commit message creation with AI. Its robust feature set ensures flexibility and precision for developers and teams.

### AI and Commit Generation
- **AI-Powered Commit Messages**: Generates high-quality commit messages using OpenAI's API, analyzing staged git changes for context-aware results.
- **Multilingual Support**: Creates commit messages in multiple languages (e.g., English, Persian, Russian) to suit global teams.
- **Extensible AI Providers**: Supports OpenAI with plans for Anthropic and other providers, ensuring future-proofing.

### Commit Standards and Customization
- **Conventional Commits**: Adheres to standard commit types (`feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`, `revert`, `init`, `security`) for semantic versioning.
- **Gitmoji Integration**: Optionally adds Gitmoji emojis (e.g., ✨ for `feat`, 🚑 for `fix`) for visually appealing commits.
- **Custom Commit Conventions**: Supports JSON-based custom prefixes (e.g., JIRA ticket IDs) for tailored commit formats.

### Git Integration
- **Optimized Git Diff Processing**: Automatically retrieves and filters staged git diff, excluding irrelevant files (e.g., `node_modules/*`, `*.lock`).
- **Configurable Exclusions**: Customize file exclusion patterns via config file to focus on relevant changes.

### Configuration and Networking
- **Flexible Configuration**: Customize via CLI flags, environment variables, or a JSON config file (`~/.gitc/config.json`).
- **Proxy Support**: Configurable proxy settings for API requests in restricted environments.
- **Timeout and Redirect Control**: Adjustable timeouts and HTTP redirect limits for reliable API interactions.
- **Environment Variable Support**: Simplifies setup for sensitive data (e.g., API keys) and common settings.

### Performance and Reliability
- **Fast Processing**: Leverages `sonic` for rapid JSON parsing and `fasthttp` for efficient HTTP requests.
- **Error Handling**: Robust validation and error messages ensure reliable operation.

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

## Environment Variables
```bash
export OPENAI_API_KEY="sk-your-key-here"
export GITC_LANGUAGE="fa"
export GITC_MODEL="gpt-4"
```

# 💻 Basic Usage
```bash
# Stage your changes first
git add .

# Generate commit message
gitc

# With Gitmoji
gitc --emoji

# Specific language
gitc --lang fa

# Custom commit type
gitc --commit-type fix
```

## 📚 Full Options
The following CLI flags are available for the `ai-commit` command and its `config` subcommand. All flags can also be set via environment variables or the `~/.gitc/config.json` file.

| Flag | Alias | Description | Default | Environment Variable | Example |
|------|-------|-------------|---------|----------------------|---------|
| `--provider` | `-a` | AI provider to use (e.g., `openai`, `anthropic`) | `openai` | `AI_PROVIDER` | `--provider openai` |
| `--model` | - | OpenAI model for commit message generation | `gpt-4o-mini` | - | `--model gpt-4o` |
| `--lang` | - | Language for commit messages (e.g., `en`, `fa`, `ru`) | `en` | `GITC_LANGUAGE` | `--lang fa` |
| `--timeout` | - | Request timeout in seconds | `10` | - | `--timeout 15` |
| `--maxLength` | - | Maximum length of the commit message | `200` | - | `--maxLength 150` |
| `--api-key` | `-k` | API key for the AI provider | - | `AI_API_KEY` | `--api-key sk-xxx` |
| `--proxy` | `-p` | Proxy URL for API requests | - | `GITC_PROXY` | `--proxy http://proxy.example.com:8080` |
| `--commit-type` | `-t` | Commit type for Conventional Commits (e.g., `feat`, `fix`) | - | `GITC_COMMIT_TYPE` | `--commit-type feat` |
| `--custom-convention` | `-C` | Custom commit message convention (JSON format) | - | `GITC_CUSTOM_CONVENTION` | `--custom-convention '{"prefix": "JIRA-123"}'` |
| `--emoji` | `-g` | Add Gitmoji to the commit message | `false` | `GITC_GITMOJI` | `--emoji` |
| `--max-redirects` | `-r` | Maximum number of HTTP redirects | `5` | `GITC_MAX_REDIRECTS` | `--max-redirects 10` |
| `--config` | `-c` | Path to the configuration file | `~/.gitc/config.json` | `GITC_CONFIG_PATH` | `--config ./my-config.json` |

### Notes:
- Flags for the `config` subcommand are similar but exclude defaults, as they override the config file.
- The `--custom-convention` flag expects a JSON string with a `prefix` field (e.g., `{"prefix": "JIRA-123"}`).
- Environment variables take precedence over config file settings but are overridden by CLI flags.


## 🤖 AI Providers

Currently, `ai-commit` supports the following AI providers. Additional providers (e.g., Anthropic) are planned for future releases.

| Provider | Supported Models | Required Configuration | Notes |
|----------|------------------|------------------------|-------|
| OpenAI   | `gpt-4o`, `gpt-4o-mini`, `gpt-3.5-turbo`, etc. | `api_key`, `model`, `url` (optional) | Default provider. Requires a valid OpenAI API key. |
| Anthropic | - | - | Coming soon. |