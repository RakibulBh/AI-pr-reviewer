package usecase

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
)

type GithubUsecase struct {
	repository *repository.GithubRepository
	gemini     *repository.GeminiRepository
	appID      int64
	privateKey *rsa.PrivateKey
}

const OPENED_ACTION = "opened"
const REOPENED_ACTION = "reopened"

func NewGithubUsecase(repository *repository.GithubRepository, gemini *repository.GeminiRepository, appID int64, privateKey *rsa.PrivateKey) *GithubUsecase {
	return &GithubUsecase{
		repository: repository,
		gemini:     gemini,
		appID:      appID,
		privateKey: privateKey,
	}
}

func (g *GithubUsecase) PullRequestReviewer(ctx context.Context, event *github.PullRequestEvent) error {
	action := event.GetAction()

	switch action {
	case OPENED_ACTION:
		err := g.reviewPullRequest(ctx, event)
		if err != nil {
			return err
		}
		slog.Info("pull request review completed successfully")
	case REOPENED_ACTION:
		err := g.reviewPullRequest(ctx, event)
		if err != nil {
			return err
		}
		slog.Info("pull request review completed successfully")
	default:
		slog.Info("recieved an action which is not supported yet", "action", action)
	}

	return nil
}

// Private methods

func (g *GithubUsecase) reviewPullRequest(ctx context.Context, event *github.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	pullNumber := event.GetPullRequest().GetNumber()
	installationID := event.Installation.GetID()
	commitID := event.GetPullRequest().GetHead().GetSHA()

	// Debug logging
	slog.Info("Processing PR review",
		"owner", owner,
		"repo", repo,
		"pullNumber", pullNumber,
		"installationID", installationID,
		"commitID", commitID)

	// Create new bot client for this request
	jwt, err := g.generateJWT()
	if err != nil {
		return fmt.Errorf("error generating jwt token for github client: %v", err)
	}
	installationToken, err := g.getInstallationToken(ctx, installationID, jwt)
	if err != nil {
		return fmt.Errorf("error retrieving installation token from github: %v", err)

	}
	client := github.NewClient(oauth2.NewClient(context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: installationToken})))

	// Loop until there are no more pages of files to review, pages start from 1
	pageCount := 1
	for {
		files, err := g.repository.ListPullRequestFiles(ctx, client, owner, repo, pullNumber, pageCount)
		if err != nil {
			slog.Error("error fetching diffs", "error", err, "owner", owner, "repo", repo, "pullNumber", pullNumber)
			return err
		}

		// If there are no files break the loop
		if len(files) <= 0 {
			break
		}

		// Parse the files to send to the LLM
		formattedDiffs := g.formatFilesForLLM(files)
		reviews, err := g.gemini.GetCodeReviews(formattedDiffs)
		if err != nil {
			slog.Error("error getting code reviews from LLM", "error", err)
			return err
		}
		slog.Info("reviews have been created by the LLM", "number_of_reviews", len(reviews))

		// Create each review
		for _, review := range reviews {

			comment := &github.PullRequestComment{
				Body:     &review.Body,
				CommitID: &commitID,
				Path:     &review.Path,
				Line:     &review.Line,
			}

			time.Sleep(time.Second * 5)

			err = g.repository.CreateReviewComments(ctx, client, owner, repo, pullNumber, comment)
			if err != nil {
				slog.Error("error creating review comment", "error", err, "comment", comment)
			}
		}

		pageCount++
		time.Sleep(time.Second * 15)
	}

	return nil
}

func (g *GithubUsecase) formatFilesForLLM(files []*github.CommitFile) string {
	var formattedFiles []string
	const maxDiffLines = 500 // Maximum lines in a diff to send to LLM

	for _, file := range files {
		if file.GetPatch() == "" {
			continue // Skip binary or unchanged files
		}

		// Skip binary files
		if g.isBinaryFile(file) {
			slog.Info("Skipping binary file", "filename", file.GetFilename())
			continue
		}

		// Skip config files, env files, and other non-reviewable files
		if g.isConfigFile(file) {
			slog.Info("Skipping config/env file", "filename", file.GetFilename())
			continue
		}

		// Skip very large diffs that would overwhelm the LLM
		patchLines := strings.Split(file.GetPatch(), "\n")
		if len(patchLines) > maxDiffLines {
			slog.Info("Skipping file with large diff",
				"filename", file.GetFilename(),
				"lines", len(patchLines),
				"maxLines", maxDiffLines)
			continue
		}

		formatted := fmt.Sprintf(`
			FILE: %s
			STATUS: %s (+%d, -%d)
			DIFF:
			%s
			---END FILE---`,
			file.GetFilename(),
			file.GetStatus(),
			file.GetAdditions(),
			file.GetDeletions(),
			g.formatPatchWithLineNumbers(file.GetPatch()),
		)

		formattedFiles = append(formattedFiles, formatted)
	}

	return strings.Join(formattedFiles, "\n\n")
}

func (g *GithubUsecase) formatPatchWithLineNumbers(patch string) string {
	lines := strings.Split(patch, "\n")
	var formatted []string
	position := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			formatted = append(formatted, line) // Keep hunk headers as-is
			position = 0                        // Reset position counter
			continue
		}

		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, " ") {
			position++
			formatted = append(formatted, fmt.Sprintf("[POS:%d] %s", position, line))
		} else if strings.HasPrefix(line, "-") {
			formatted = append(formatted, line) // Deleted lines don't get positions
		} else {
			formatted = append(formatted, line)
		}
	}

	return strings.Join(formatted, "\n")
}

func (g *GithubUsecase) isBinaryFile(file *github.CommitFile) bool {
	filename := file.GetFilename()

	// Common binary file extensions
	binaryExtensions := []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".lib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".svg",
		".mp3", ".mp4", ".avi", ".mov", ".wav", ".flac",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".bin", ".dat", ".db", ".sqlite", ".sqlite3",
		".ttf", ".otf", ".woff", ".woff2",
		".jar", ".war", ".ear", ".class",
		".pyc", ".pyo", ".o", ".obj",
	}

	// Check file extension
	for _, ext := range binaryExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}

	// Check if patch contains binary file indicator
	patch := file.GetPatch()
	if strings.Contains(patch, "Binary files") ||
		strings.Contains(patch, "GIT binary patch") ||
		strings.Contains(patch, "cannot display: file marked as a binary type") {
		return true
	}

	return false
}

func (g *GithubUsecase) isConfigFile(file *github.CommitFile) bool {
	filename := strings.ToLower(file.GetFilename())

	// Config file patterns and extensions
	configPatterns := []string{
		// Environment files
		".env", ".env.local", ".env.development", ".env.production", ".env.test",
		".env.example", ".env.sample", ".env.template",

		// Package manager files
		"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "composer.lock",
		"pipfile.lock", "poetry.lock", "go.sum", "cargo.lock",

		// Config files
		".gitignore", ".gitattributes", ".dockerignore", ".eslintignore",
		".prettierignore", ".editorconfig", ".nvmrc", ".node-version",

		// CI/CD files
		".github/", ".gitlab-ci.yml", ".travis.yml", ".circleci/",
		"jenkinsfile", "dockerfile", "docker-compose.yml", "docker-compose.yaml",

		// IDE/Editor files
		".vscode/", ".idea/", "*.iml",

		// Build/Deploy configs
		"webpack.config.js", "rollup.config.js", "vite.config.js",
		"tsconfig.json", "jsconfig.json", "babel.config.js", ".babelrc",
		"tailwind.config.js", "postcss.config.js",
		"makefile", "cmake", "build.gradle", "pom.xml",

		// Linting/Formatting configs
		".eslintrc", ".prettierrc", ".stylelintrc", "tslint.json",
		".flake8", ".pylintrc", "pyproject.toml", "setup.cfg",

		// Documentation that doesn't need review
		"readme.md", "changelog.md", "license", "authors", "contributors",
		"code_of_conduct.md", "contributing.md", "security.md",
	}

	// Check exact filename matches
	for _, pattern := range configPatterns {
		if strings.Contains(filename, pattern) {
			return true
		}
	}

	// Check file extensions that are typically config
	configExtensions := []string{
		".toml", ".ini", ".cfg", ".conf", ".config", ".properties",
		".yaml", ".yml", ".json", ".xml",
	}

	for _, ext := range configExtensions {
		if strings.HasSuffix(filename, ext) {
			// Allow certain JSON/YAML files that might contain business logic
			if !strings.Contains(filename, "schema") &&
				!strings.Contains(filename, "spec") &&
				!strings.Contains(filename, "test") &&
				!strings.Contains(filename, "mock") {
				return true
			}
		}
	}

	return false
}

func (g *GithubUsecase) generateJWT() (string, error) {
	if g.privateKey == nil {
		return "", fmt.Errorf("private key is nil")
	}
	if g.appID == 0 {
		return "", fmt.Errorf("app ID is not set")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": strconv.FormatInt(g.appID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.privateKey)
}

func (g *GithubUsecase) getInstallationToken(ctx context.Context, installationID int64, jwtToken string) (string, error) {
	// Create GitHub client with JWT
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: jwtToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get installation access token
	installation, _, err := client.Apps.CreateInstallationToken(
		ctx, installationID, &github.InstallationTokenOptions{},
	)
	if err != nil {
		return "", err
	}

	return installation.GetToken(), nil
}
