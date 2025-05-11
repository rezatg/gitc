# 🤝 Contributing to `gitc`

Thank you for your interest in contributing to **`gitc`** — an AI-powered CLI tool for generating commit messages from your git diffs.

We welcome all types of contributions, whether you're fixing bugs, suggesting features, improving documentation, or writing tests.


## 📌 Table of Contents
* [Code of Conduct](#-code-of-conduct)
* [Getting Started](#-getting-started)
* [Development Setup](#-development-setup)
* [Making Contributions](#-making-contributions)
* [Commit Message Guidelines](#-commit-message-guidelines)
* [Pull Request Process](#-pull-request-process)
* [Feature Suggestions & Bugs](#-feature-suggestions--bugs)

## 📜 Code of Conduct
We follow a [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). Please treat everyone respectfully and kindly.

## ⚙️ Getting Started
1. **Fork** the repository.
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/<your-username>/gitc.git
   cd gitc
   ```
3. Create a new branch:
   ```bash
   git checkout -b feature/my-feature
   ```

## 🛠 Development Setup
  Make sure you have:
    * Go ≥ **1.18**
    * Git
    * (Optional) OpenAI API Key for testing
  
  Install dependencies:
    ```bash
    go mod tidy
    ```

Build the project:
  ```bash
  go build -o gitc ./cmd/gitc
  ```

Run the tool:
  ```bash
  ./gitc --help
  ```

Run tests:
  ```bash
  go test ./...
  ```

## ✍️ Making Contributions
You can contribute in the following ways:
  * 🐛 Bug Fixes
  * 📄 Documentation Improvements
  * 🚀 New Features
  * ✅ Tests and Coverage
  * 💡 Suggesting Ideas and Discussions

If unsure, [open a discussion](https://github.com/rezatg/gitc/discussions) or [create an issue](https://github.com/rezatg/gitc/issues) before starting work.

## 🧾 Commit Message Guidelines
We follow [Conventional Commits](https://www.conventionalcommits.org) to ensure readable, semantic commit history.

Example:
  ```bash
  feat(config): add support for Gemini provider
  fix(cli): handle missing config gracefully
  docs(readme): update installation instructions
  ```

> You can even use `gitc` itself to generate commit messages:
  ```bash
  gitc -a --commit-type feat
  ```

## 🚀 Pull Request Process
sure your PR targets the `main` branch.
  2. Make sure all tests pass.
  3. Write a meaningful title and description.
  4. Link related issues (e.g., `Closes #42`).
  5. Add relevant labels if possible.
  6. Wait for code review and address feedback.

## 💡 Feature Suggestions & Bugs
* 💬 Found a bug? [Open an issue](https://github.com/rezatg/gitc/issues)
* 💡 Have a feature idea? [Start a discussion](https://github.com/rezatg/gitc/discussions)
* 🙌 Want to help but don’t know where to start? Look for [good first issues](https://github.com/rezatg/gitc/labels/good%20first%20issue)

## 🫶 Thank You!
Your time, effort, and ideas make `gitc` better. We're thrilled to have you here 🙌
