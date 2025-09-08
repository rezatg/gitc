# Changelog

## [0.3.0] - 2025-06-25
### Added
- Pretty-printed `git commit` command output with `-m` flags and line continuation (`\`).
- Automatic detection of single-line vs multi-line commit messages with clean display formatting.

### Changed
- Refactored commit message prompt for clarity, brevity, and better LLM compatibility.
- Improved formatting rules for both single-line and multi-line commit messages.
- Enforced strict summary/body structure with consistent guidelines.

---

## [0.2.0] - 2025-05-15
### Added
- Experimental support for Grok (xAI) and DeepSeek AI providers.
- New `--url` flag for custom API endpoints.
- Interactive mode for commit message preview and editing.

### Changed
- Updated README with improved provider status and clarity.
- Revised config structure to remove `open_ai` field.

### Fixed
- API key persistence issues in configuration.
- Improved validation for configuration settings.
- Fixed missing space between version number and summary in commit messages.
- Eliminated redundant spacing and formatting edge cases in output.

---

## [0.1.1] - 2025-04-01
- Initial release with OpenAI support.
