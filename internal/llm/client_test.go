package llm

import (
	"testing"
	"time"
)

func TestExtractWaitTime(t *testing.T) {
	tests := []struct {
		name         string
		errorBody    string
		wantDuration time.Duration
		wantOk       bool
	}{
		{
			name: "Valid rate limit error with 29 seconds",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Rate limit of 10 per 60s exceeded for UserByModelByMinute. Please wait 29 seconds before retrying.",
					"details": "Rate limit of 10 per 60s exceeded for UserByModelByMinute. Please wait 29 seconds before retrying."
				}
			}`,
			wantDuration: 29 * time.Second,
			wantOk:       true,
		},
		{
			name: "Valid rate limit error with 60 seconds",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Please wait 60 seconds before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 60 * time.Second,
			wantOk:       true,
		},
		{
			name: "Valid rate limit error with 5 seconds",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Please wait 5 second before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 5 * time.Second,
			wantOk:       true,
		},
		{
			name: "Different error code",
			errorBody: `{
				"error": {
					"code": "SomeOtherError",
					"message": "Please wait 29 seconds before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 0,
			wantOk:       false,
		},
		{
			name: "Rate limit error without wait time",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Rate limit exceeded",
					"details": ""
				}
			}`,
			wantDuration: 0,
			wantOk:       false,
		},
		{
			name:         "Invalid JSON",
			errorBody:    `This is not valid JSON`,
			wantDuration: 0,
			wantOk:       false,
		},
		{
			name:         "Empty string",
			errorBody:    "",
			wantDuration: 0,
			wantOk:       false,
		},
		{
			name: "Malformed JSON with rate limit",
			errorBody: `{
				"error": {
					"code": "RateLimitReached"
				}
			}`,
			wantDuration: 0,
			wantOk:       false,
		},
		{
			name: "Wait time under 1 hour (300s)",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Please wait 300 seconds before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 300 * time.Second,
			wantOk:       true,
		},
		{
			name: "Wait time under 1 hour (3599s)",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Please wait 3599 seconds before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 3599 * time.Second,
			wantOk:       true,
		},
		{
			name: "Daily rate limit (over 1 hour - 3601s)",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Please wait 3601 seconds before retrying.",
					"details": ""
				}
			}`,
			wantDuration: 0, // Don't retry for daily limits
			wantOk:       false,
		},
		{
			name: "Daily rate limit (real example - 5631s)",
			errorBody: `{
				"error": {
					"code": "RateLimitReached",
					"message": "Rate limit of 50 per 86400s exceeded for UserByModelByDay. Please wait 5631 seconds before retrying.",
					"details": "Rate limit of 50 per 86400s exceeded for UserByModelByDay. Please wait 5631 seconds before retrying."
				}
			}`,
			wantDuration: 0, // Don't retry for daily limits
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDuration, gotOk := extractWaitTime(tt.errorBody)
			if gotOk != tt.wantOk {
				t.Errorf("extractWaitTime() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotDuration != tt.wantDuration {
				t.Errorf("extractWaitTime() gotDuration = %v, want %v", gotDuration, tt.wantDuration)
			}
		})
	}
}
