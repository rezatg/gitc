# Changelog

## [0.2.0] - 2025-05-15
### Added
- Experimental support for Grok (xAI) and DeepSeek AI providers.
- New `--url` flag for custom API endpoints.
- Interactive mode for commit message preview and editing.

### Changed
- Updated README with accurate provider status and improved clarity.
- Revised config structure to remove `open_ai` field.

### Fixed
- API key persistence issues in configuration.
- Improved validation for configuration settings.

## [0.1.1] - 2025-04-01
- Initial release with OpenAI support.

## [0.2.2] - 2025-06-18
### Changed
- Refactored commit message prompt for clarity, brevity, and better LLM compatibility.
- Improved formatting rules for single-line and multi-line commit messages.
- Enhanced summary/body structure with strict formatting guidelines.

### Added
- Pretty-printed `git commit` command output using `-m` flags with line continuation (`\`).
- Detection for single-line vs multi-line commit messages with clean display logic.

### Fixed
- Prevented missing space between version number and summary in commit messages.
- Eliminated redundant spacing and fixed formatting edge cases in output.