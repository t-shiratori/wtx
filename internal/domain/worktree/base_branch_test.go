package worktree

import "testing"

func TestResolveBaseBranch(t *testing.T) {
	testPatterns := []struct {
		name       string
		fromFlag   string
		fromConfig string
		want       string
	}{
		{
			name:       "from flag has priority",
			fromFlag:   "feature/from-flag",
			fromConfig: "main",
			want:       "feature/from-flag",
		},
		{
			name:       "use config when flag is empty",
			fromFlag:   "",
			fromConfig: "develop",
			want:       "develop",
		},
		{
			name:       "fallback to HEAD",
			fromFlag:   "",
			fromConfig: "",
			want:       "HEAD",
		},
	}

	for _, testPattern := range testPatterns {
		t.Run(testPattern.name, func(t *testing.T) {
			got := ResolveBaseBranch(testPattern.fromFlag, testPattern.fromConfig)
			if got != testPattern.want {
				t.Fatalf("expected %q, got %q", testPattern.want, got)
			}
		})
	}
}
