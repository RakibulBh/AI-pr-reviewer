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
	default:
		return fmt.Errorf("action not supported: %v", action)
	}

	return nil
}

// Private methods

func (g *GithubUsecase) reviewPullRequest(ctx context.Context, event *github.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetName()
	repo := event.GetRepo().GetName()
	pullNumber := event.GetPullRequest().GetNumber()
	installationID := event.Installation.GetID()
	commitID := event.GetPullRequest().GetMergeCommitSHA() // Head SHA

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

	for _, file := range files {
		if file.GetPatch() == "" {
			continue // Skip binary or unchanged files
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
