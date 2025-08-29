# AI PR Reviewer 🤖

An intelligent GitHub App that provides automated code reviews using Google's Gemini AI. Get comprehensive, actionable feedback on your pull requests with the expertise of a senior software engineer.

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Issues](https://img.shields.io/github/issues/RakibulBh/AI-pr-reviewer)](https://github.com/RakibulBh/AI-pr-reviewer/issues)

## ✨ Features

- **Intelligent Code Reviews**: Leverages Google Gemini AI to provide comprehensive code analysis
- **GitHub Integration**: Seamlessly integrates as a GitHub App with webhook support
- **Automated Comments**: Automatically comments on pull requests with detailed feedback
- **Multi-file Analysis**: Reviews all changed files in a pull request with proper context
- **Severity Indicators**: Categorizes feedback as BLOCKING, IMPORTANT, NIT, or QUESTION
- **Educational Feedback**: Explains the reasoning behind each suggestion
- **Performance Optimized**: Handles large PRs with pagination and rate limiting

## 🚀 Quick Start

### Prerequisites

- Go 1.24.5 or higher
- GitHub App credentials
- Google Gemini API key

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/RakibulBh/AI-pr-reviewer.git
   cd AI-pr-reviewer
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

## ⚙️ Configuration

Create a `.env` file with the following variables:

```env
# Environment
ENV=development
PORT=8080

# GitHub App Configuration
APP_ID=your_github_app_id
GITHUB_BOT_PRIVATE_KEY=your_private_key_content
GITHUB_REPO_WEBHOOK_SECRET=your_webhook_secret

# AI Configuration
GEMINI_API_KEY=your_gemini_api_key
```

### GitHub App Setup

1. Create a new GitHub App in your organization/account settings
2. Set the webhook URL to `https://your-domain.com/github/webhook`
3. Enable the following permissions:
   - **Pull requests**: Read & Write
   - **Contents**: Read
   - **Metadata**: Read
4. Subscribe to **Pull request** events
5. Generate and download the private key

## 📖 Usage

### Basic Usage

Once installed and configured, the AI PR Reviewer will automatically:

1. **Monitor Pull Requests**: Listens for `opened` and `reopened` PR events
2. **Analyze Changes**: Reviews all modified files in the pull request
3. **Generate Comments**: Creates inline comments with detailed feedback
4. **Provide Context**: Explains reasoning behind each suggestion

### Example Review Comments

The AI reviewer provides structured feedback like:

````
IMPORTANT: Consider using a more specific error type here instead of generic error.

The current implementation returns a generic error which makes debugging difficult for consumers of this API. Consider creating a custom error type like `ValidationError` that includes the field name and validation rule that failed.

Example:
```go
type ValidationError struct {
    Field string
    Rule  string
    Value interface{}
}
````

This improves error handling and makes the API more developer-friendly.

```

### Supported Languages

The AI reviewer works with any programming language but is optimized for:
- Go
- JavaScript/TypeScript
- Python
- Java
- C#
- And more...

## 🏗️ Architecture

```

├── cmd/
│ └── main.go # Application entry point
├── internal/
│ ├── config/ # Configuration management
│ ├── delivery/http/ # HTTP handlers and routes
│ ├── model/ # Data models
│ ├── repository/ # External service integrations
│ ├── usecase/ # Business logic
│ └── utils/ # Utility functions
└── docs/ # Documentation and credentials

````

### Key Components

- **GitHub Controller**: Handles webhook events from GitHub
- **GitHub Repository**: Manages GitHub API interactions
- **Gemini Repository**: Interfaces with Google's Gemini AI
- **GitHub Usecase**: Orchestrates the review process

## 🔧 Development

### Running Locally

1. **Start the development server**
   ```bash
   go run cmd/main.go
````

2. **Use ngrok for webhook testing**

   ```bash
   ngrok http 8080
   ```

3. **Update your GitHub App webhook URL** to the ngrok URL

### Building for Production

```bash
# Build binary
go build -o bin/ai-pr-reviewer cmd/main.go

# Run binary
./bin/ai-pr-reviewer
```

### Docker Deployment

```bash
# Build image
docker build -t ai-pr-reviewer .

# Run container
docker run -p 8080:8080 --env-file .env ai-pr-reviewer
```

## 📊 Review Quality

The AI reviewer is engineered to provide:

- **Comprehensive Analysis**: Reviews code quality, performance, security, and maintainability
- **Educational Feedback**: Explains the "why" behind each suggestion
- **Actionable Suggestions**: Provides specific code examples and improvements
- **Context Awareness**: Considers the broader system implications
- **Best Practices**: Enforces industry standards and conventions

## 🤝 Contributing

We welcome contributions!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Google Gemini AI](https://ai.google.dev/) for powering the intelligent reviews
- [GitHub API](https://docs.github.com/en/rest) for seamless integration
- [Go Chi](https://github.com/go-chi/chi) for the lightweight HTTP router

## 📞 Support

- 📧 Email: [rakibulbhuiyan.dev@gmail.com](mailto:rakibulbhuiyan.dev@gmail.com)
- 🐛 Issues: [GitHub Issues](https://github.com/RakibulBh/AI-pr-reviewer/issues)

---

**Made with ❤️ by developers, for developers**
