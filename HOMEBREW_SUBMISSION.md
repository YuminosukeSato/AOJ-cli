# Homebrew Core Submission Guide

## Steps to Submit AOJ CLI to Homebrew Core

### 1. Prerequisites

Before submitting to Homebrew Core, ensure:

- [ ] Project has stable releases with semantic versioning
- [ ] Formula is tested and working
- [ ] Project has good documentation
- [ ] License is specified (MIT ✓)
- [ ] CI/CD is set up
- [ ] Project is notable/useful to the community

### 2. Prepare the Formula

The formula file `aoj-cli.rb` is ready for submission. Before submitting:

1. **Create a stable release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Calculate SHA256:**
   ```bash
   curl -sL https://github.com/YuminosukeSato/AOJ-cli/archive/refs/tags/v1.0.0.tar.gz | sha256sum
   ```

3. **Update the formula with the actual SHA256**

### 3. Test the Formula Locally

```bash
# Install the formula locally
brew install --build-from-source ./aoj-cli.rb

# Test it works
aoj --help

# Uninstall for clean testing
brew uninstall aoj-cli
```

### 4. Submit to Homebrew Core

1. **Fork homebrew/homebrew-core**
2. **Create a new branch:**
   ```bash
   git checkout -b aoj-cli
   ```

3. **Add the formula:**
   ```bash
   cp aoj-cli.rb Formula/aoj-cli.rb
   ```

4. **Test with Homebrew's audit:**
   ```bash
   brew audit --new-formula Formula/aoj-cli.rb
   brew test Formula/aoj-cli.rb
   ```

5. **Commit and create PR:**
   ```bash
   git add Formula/aoj-cli.rb
   git commit -m "aoj-cli: new formula

   Command-line interface for Aizu Online Judge (AOJ)"
   git push origin aoj-cli
   ```

6. **Create Pull Request** to homebrew/homebrew-core

### 5. Alternative: Start with Homebrew Cask

If the main formula is rejected, consider submitting as a cask for binary distribution.

### 6. Submission Checklist

- [ ] Formula follows Homebrew style guidelines
- [ ] Formula has proper `desc` and `homepage`
- [ ] License is specified
- [ ] Dependencies are minimal and necessary
- [ ] Test block verifies the installation works
- [ ] No `brew tap` required in documentation
- [ ] Formula builds successfully from source
- [ ] Passes `brew audit` and `brew test`

## After Submission

Once submitted:
1. Monitor the PR for feedback from Homebrew maintainers
2. Address any requested changes promptly
3. Update documentation once merged

## Benefits of Homebrew Core

- Users can install with simple `brew install aoj-cli`
- Better discoverability
- Automatic updates via Homebrew
- No need to maintain separate tap repository