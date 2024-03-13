# Query GitHub Issues

This is a GPTScript tool that uses the GitHub API to search for issues and pull requests across GitHub.

Set your GitHub access token to the `GPTSCRIPT_GITHUB_TOKEN` environment variable in order to use this tool.
If the variable is not set, the tool will attempt to make unauthenticated API requests.

## Example

```yaml
tools: github.com/gptscript-ai/query-github-issues

What was the most recently assigned issue to g-linville?
```

## License

This tool is available under the Apache License 2.0. See [LICENSE](LICENSE) for more information.
