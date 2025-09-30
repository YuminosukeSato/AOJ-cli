# AOJ CLI

A command-line interface tool for [Aizu Online Judge (AOJ)](https://onlinejudge.u-aizu.ac.jp/), inspired by atcoder-cli. This tool streamlines your competitive programming workflow by managing authentication, downloading problem sets, running tests locally, and submitting solutions directly from your terminal.

## Features

- 🔐 **Secure Authentication**: Login to AOJ with secure session management
- 📥 **Problem Download**: Fetch problem descriptions and test cases
- 🧪 **Local Testing**: Run your solutions against sample test cases locally
- 📤 **Direct Submission**: Submit your solutions to AOJ from the command line
- 🗂️ **Project Organization**: Automatically organize problems in a structured directory layout
- 🔄 **Session Persistence**: Maintains login sessions across multiple uses
- 🛡️ **Secure Storage**: Credentials stored securely in your local configuration

## Installation

### Via Homebrew (macOS/Linux)

```bash
brew install aoj-cli
```

### From Source

#### Prerequisites
- Go 1.21 or higher
- Git

```bash
# Clone the repository
git clone https://github.com/YuminosukeSato/AOJ-cli.git
cd AOJ-cli

# Build the binary
go build -o aoj ./cmd/aojcli

# Move to PATH (optional)
sudo mv aoj /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/YuminosukeSato/AOJ-cli/cmd/aojcli@latest
```

## Quick Start

### 1. Login to AOJ

```bash
aoj login
# Enter your username and password when prompted
```

### 2. Initialize a Problem

```bash
# Initialize a specific problem
aoj init ITP1_1_A

# Initialize with a custom directory name
aoj init ITP1_1_A --dir hello-world
```

### 3. Test Your Solution

```bash
# Run tests for the current problem
aoj test main.cpp

# Run tests with custom input
aoj test main.py --case 1
```

### 4. Submit Your Solution

```bash
# Submit your solution
aoj submit main.cpp

# Submit with specific language
aoj submit solution.py --lang Python3
```

## Commands

### `aoj login`
Authenticate with AOJ and save your session locally.

```bash
aoj login
```

### `aoj logout`
Clear your local session.

```bash
aoj logout
```

### `aoj init <problem-id>`
Initialize a new problem directory with test cases.

```bash
aoj init ITP1_1_A
aoj init --contest ITP1  # Initialize all problems in a contest
```

Options:
- `--dir, -d`: Custom directory name
- `--contest, -c`: Initialize entire contest

### `aoj test <file>`
Run your solution against sample test cases.

```bash
aoj test main.cpp
aoj test solution.py --case 2
```

Options:
- `--case, -c`: Run specific test case
- `--timeout, -t`: Set execution timeout (default: 2s)

### `aoj submit <file>`
Submit your solution to AOJ.

```bash
aoj submit main.cpp
aoj submit solution.py --lang Python3
```

Options:
- `--lang, -l`: Specify programming language
- `--wait, -w`: Wait for judge result

### `aoj status`
Check submission status.

```bash
aoj status  # Latest submission
aoj status --all  # All recent submissions
```

### `aoj config`
Manage configuration settings.

```bash
aoj config set default-lang cpp
aoj config get default-lang
aoj config list
```

## Configuration

Configuration file is stored at `~/.config/aoj/config.toml`.

### Example Configuration

```toml
[user]
default_language = "cpp"
default_template = "~/.config/aoj/templates/main.cpp"

[test]
timeout = 2000  # milliseconds
diff_mode = "unified"  # unified, split, or simple

[submit]
wait_result = true
open_browser = false
```

## Directory Structure

When you initialize a problem, AOJ CLI creates the following structure:

```
ITP1_1_A/
├── problem.md      # Problem description
├── samples/        # Sample test cases
│   ├── 1.in
│   ├── 1.out
│   ├── 2.in
│   └── 2.out
└── main.cpp        # Your solution file (from template)
```

## Templates

Create custom templates for different languages:

```bash
# Create a template directory
mkdir -p ~/.config/aoj/templates

# Add your templates
cat > ~/.config/aoj/templates/main.cpp << 'EOF'
#include <iostream>
using namespace std;

int main() {
    // Your code here
    return 0;
}
EOF

# Set as default
aoj config set default-template ~/.config/aoj/templates/main.cpp
```

## Development

### Prerequisites

- Go 1.21+
- Task (optional, for automation)

### Setup

```bash
# Clone the repository
git clone https://github.com/YuminosukeSato/AOJ-cli.git
cd AOJ-cli

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o aoj ./cmd/aojcli
```

### Using Task (Recommended)

```bash
# Install Task
brew install go-task/tap/go-task

# Run common tasks
task build      # Build the binary
task test       # Run tests
task lint       # Run linters
task dev        # Run with hot reload
```

### Project Structure

```
.
├── cmd/aojcli/         # Entry point
├── internal/
│   ├── cli/            # Command implementations
│   ├── domain/         # Business logic
│   │   ├── entity/     # Domain entities
│   │   ├── model/      # Value objects
│   │   └── repository/ # Repository interfaces
│   ├── infrastructure/ # External services
│   │   └── repository/ # Repository implementations
│   └── usecase/        # Application logic
├── pkg/
│   ├── cerrors/        # Error handling
│   ├── config/         # Configuration management
│   └── logger/         # Logging utilities
└── test/               # Integration tests
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md).

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go best practices
- Run `task lint` before committing
- Update documentation as needed

## Troubleshooting

### Common Issues

#### Login fails with "Invalid credentials"
- Verify your username and password
- Check if AOJ website is accessible
- Clear session with `aoj logout` and try again

#### "Command not found" after installation
- Ensure the binary is in your PATH
- Try using the full path to the binary
- Restart your terminal session

#### Test cases not downloading
- Check your internet connection
- Verify the problem ID is correct
- Try logging in again

### Debug Mode

Enable debug logging for troubleshooting:

```bash
AOJ_LOG_LEVEL=debug aoj test main.cpp
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [atcoder-cli](https://github.com/Tatamo/atcoder-cli)
- Thanks to [Aizu Online Judge](https://onlinejudge.u-aizu.ac.jp/) for providing the platform
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework

## Support

- **Issues**: [GitHub Issues](https://github.com/YuminosukeSato/AOJ-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/YuminosukeSato/AOJ-cli/discussions)
- **Wiki**: [Project Wiki](https://github.com/YuminosukeSato/AOJ-cli/wiki)

## Roadmap

- [ ] Support for more programming languages
- [ ] Interactive problem selector
- [ ] Contest mode with timer
- [ ] Performance statistics tracking
- [ ] Integration with code editors (VSCode, Vim)
- [ ] Parallel test execution
- [ ] Custom judge programs
- [ ] Problem recommendation system

---

Made with ❤️ by [Yuminosuke Sato](https://github.com/YuminosukeSato)