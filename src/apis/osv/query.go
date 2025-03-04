package osv

// QueryBatchRequest represents the payload for the OSV API batch query.
type QueryBatchRequest struct {
	Queries []Query `json:"queries"`
}

// Query represents a single query in the batch.
type Query struct {
	Commit  string        `json:"commit,omitempty"`
	Version string        `json:"version,omitempty"`
	Package *PackageQuery `json:"package,omitempty"`
}

// PackageQuery represents the package query structure.
type PackageQuery struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
	PURL      string `json:"purl,omitempty"`
}

// QueryBatchResponse represents the response from the OSV API batch query.
type QueryBatchResponse struct {
	Results []QueryResult `json:"results"`
}

// QueryResult represents a single result in the batch query response.
type QueryResult struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
}

// Vulnerability represents a single vulnerability returned by the API.
type Vulnerability struct {
	ID       string `json:"id"`
	Modified string `json:"modified"`
	Details  string `json:"details"`
}
