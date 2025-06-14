# 🔗 Action Link Checker

A GitHub Action that automatically checks for broken links in your files.

## ✨ Features

- 🔍 **Smart URL Detection**: Automatically extracts URLs from markdown links `[text](url)` and bare URLs
- ⚡ **Concurrent Checking**: Configurable number of concurrent workers for fast link validation
- ⏱️ **Timeout Control**: Configurable timeout for each link check
- 🔄 **Fallback Strategy**: Uses HEAD request first, falls back to GET if needed
- 📊 **Clear Reporting**: Displays working and broken links with clear visual indicators
- ❌ **CI/CD Integration**: Fails the workflow when broken links are found

## 🚀 Usage

### Basic Usage

```yaml
name: Check Links
on: [push, pull_request]

jobs:
  link-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check links in README
        uses: yankeexe/action-link-checker@v1
        with:
          file_path: 'README.md'
```

### Advanced Configuration

```yaml
name: Check Links
on: [push, pull_request]

jobs:
  link-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check links with custom settings
        uses: yankeexe/action-link-checker@v1
        with:
          file_path: 'docs/README.md'
          concurrent_workers: '50'
          timeout_seconds: '10'
```

### Multiple Files

```yaml
name: Check Links
on: [push, pull_request]

jobs:
  link-check:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        file: ['README.md', 'docs/CONTRIBUTING.md', 'docs/API.md']
    steps:
      - uses: actions/checkout@v4
      - name: Check links in ${{ matrix.file }}
        uses: yankeexe/action-link-checker@v1
        with:
          file_path: ${{ matrix.file }}
```

## 📋 Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `file_path` | Path to the file to parse for links | ✅ Yes | - |
| `concurrent_workers` | Number of concurrent workers for checking links | ❌ No | `30` |
| `timeout_seconds` | Timeout in seconds for each link check | ❌ No | `5` |

## 📤 Outputs

This action doesn't produce outputs but will:
- ✅ **Pass**: When all links are working
- ❌ **Fail**: When broken links are found (exits with code 1)

### 📝 Example Output

```
✅ Working URLs:
- https://github.com/yankeexe/action-link-checker
- https://docs.github.com/en/actions
- https://golang.org/

❌ Invalid URLs:
- https://broken-link-example.com/404
- https://another-broken-link.invalid/
```

## 🔗 Supported Link Formats

The action detects and validates:

### Markdown Links
```markdown
[GitHub](https://github.com)
[Documentation](https://docs.example.com)
```

### Bare URLs
```markdown
Visit https://example.com for more info
Check out https://github.com/user/repo
```

## ⚡ Performance Tips

### Optimize for Large Files
```yaml
- name: Check links with higher concurrency
  uses: yankeexe/action-link-checker@v1
  with:
    file_path: 'large-document.md'
    concurrent_workers: '50'
    timeout_seconds: '15'
```

### Handle Slow Links
```yaml
- name: Check links with extended timeout
  uses: yankeexe/action-link-checker@v1
  with:
    file_path: 'README.md'
    timeout_seconds: '30'
```

## 💡 Common Use Cases

### Documentation Sites
Perfect for checking links in documentation files, README files, and markdown-based websites.

### Blog Posts
Validate external links in blog posts and articles before publishing.

### API Documentation
Ensure all API endpoint links and external references are working.

## 🔧 Troubleshooting

### Action Fails with "Error reading file"
- Ensure the `file_path` is correct and the file exists
- Check that your workflow has checked out the repository (`actions/checkout`)

### Timeout Issues
- Increase `timeout_seconds` for slow-responding servers
- Some servers may block automated requests; this is expected behavior

### Rate Limiting
- Reduce `concurrent_workers` if you encounter rate limiting
- Some APIs may have strict rate limits for automated tools

## 🛠️ How It Works

1. The action reads the specified file and extracts URLs using regex pattern matching
2. It identifies both markdown-style links `[text](url)` and bare URLs
3. For each unique URL, it first attempts a HEAD request (faster, less bandwidth)
4. If the HEAD request fails, it falls back to a GET request
5. Results are collected and reported, with the action failing if any broken links are found

## 📈 Limitations

- The action currently only checks a single file at a time (use matrix strategy for multiple files)
- Some websites block automated requests, which may result in false negatives


## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

⭐ If this action helped you, please consider giving it a star!
