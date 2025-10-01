package types

import "testing"

func TestNewAuthor(t *testing.T) {
	tests := []struct {
		name       string
		login      string
		fullName   string
		profileURL string
		isBot      bool
		want       *Author
	}{
		{
			name:       "Regular user",
			login:      "octocat",
			fullName:   "The Octocat",
			profileURL: "https://github.com/octocat",
			isBot:      false,
			want: &Author{
				Login:      "octocat",
				Name:       "The Octocat",
				ProfileURL: "https://github.com/octocat",
				IsBot:      false,
			},
		},
		{
			name:       "Bot user",
			login:      "github-actions[bot]",
			fullName:   "GitHub Actions",
			profileURL: "https://github.com/apps/github-actions",
			isBot:      true,
			want: &Author{
				Login:      "github-actions[bot]",
				Name:       "GitHub Actions",
				ProfileURL: "https://github.com/apps/github-actions",
				IsBot:      true,
			},
		},
		{
			name:       "User with empty name",
			login:      "testuser",
			fullName:   "",
			profileURL: "https://github.com/testuser",
			isBot:      false,
			want: &Author{
				Login:      "testuser",
				Name:       "",
				ProfileURL: "https://github.com/testuser",
				IsBot:      false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthor(tt.login, tt.fullName, tt.profileURL, tt.isBot)

			if got.Login != tt.want.Login {
				t.Errorf("NewAuthor() Login = %v, want %v", got.Login, tt.want.Login)
			}
			if got.Name != tt.want.Name {
				t.Errorf("NewAuthor() Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.ProfileURL != tt.want.ProfileURL {
				t.Errorf("NewAuthor() ProfileURL = %v, want %v", got.ProfileURL, tt.want.ProfileURL)
			}
			if got.IsBot != tt.want.IsBot {
				t.Errorf("NewAuthor() IsBot = %v, want %v", got.IsBot, tt.want.IsBot)
			}
		})
	}
}
