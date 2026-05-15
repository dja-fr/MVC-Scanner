package analyzer

type FileData struct {
	Path         string   `json:"path"`
	SelfTokens   int      `json:"selfTokens"`
	TotalTokens  int      `json:"totalTokens"` // MVC
	Grade        string   `json:"grade"`
	Dependencies []string `json:"dependencies"`
	TopDep       string   `json:"topContributor"`
	TopDepTokens int      `json:"topContributorTokens"`
}
