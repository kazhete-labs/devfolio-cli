package model

// User is a public GitHub profile snapshot.
type User struct {
	Login       string
	Name        string
	Bio         string
	Company     string
	Blog        string
	Location    string
	AvatarURL   string
	HTMLURL     string
	PublicRepos int
	Followers   int
	Following   int
}

// Repo is a public repository used for portfolio + scoring.
type Repo struct {
	Name            string
	FullName        string
	Description     string
	HTMLURL         string
	Language        string
	StargazersCount int
	ForksCount      int
	Fork            bool
	Archived        bool
	Topics          []string
	DefaultBranch   string
	README          string
	Score           READMEScore
}

// READMEScore is the quality scorecard for a single README.
type READMEScore struct {
	Total    int
	Max      int
	Checks   []CheckResult
	Grade    string
	Summary  string
}

// CheckResult is one scoring rule outcome.
type CheckResult struct {
	ID      string
	Label   string
	Passed  bool
	Weight  int
	Detail  string
}

// Portfolio is the full generate result.
type Portfolio struct {
	User           User
	Repos           []Repo
	Languages      map[string]int
	AverageScore   float64
	GeneratedAtUTC string
}
