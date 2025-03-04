package osv

// OSV represents the root structure of the schema.
type OSV struct {
	SchemaVersion string   `json:"schema_version" yaml:"schema_version"`
	ID            string   `json:"id" yaml:"id"`
	Modified      string   `json:"modified" yaml:"modified"`
	Published     string   `json:"published" yaml:"published"`
	Withdrawn     string   `json:"withdrawn,omitempty" yaml:"withdrawn,omitempty"`
	Aliases       []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Related       []string `json:"related,omitempty" yaml:"related,omitempty"`
	Summary       string   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Details       string   `json:"details,omitempty" yaml:"details,omitempty"`
	Severity      []struct {
		Type  string `json:"type" yaml:"type"`
		Score string `json:"score" yaml:"score"`
	} `json:"severity,omitempty" yaml:"severity,omitempty"`
	Affected   []Affected  `json:"affected,omitempty" yaml:"affected,omitempty"`
	References []Reference `json:"references,omitempty" yaml:"references,omitempty"`
	Credits    []Credit    `json:"credits,omitempty" yaml:"credits,omitempty"`
}

// Affected represents the software versions affected by the vulnerability.
type Affected struct {
	Package           Package     `json:"package" yaml:"package"`
	Ranges            []Range     `json:"ranges,omitempty" yaml:"ranges,omitempty"`
	Versions          []string    `json:"versions,omitempty" yaml:"versions,omitempty"`
	DatabaseSpecific  interface{} `json:"database_specific,omitempty" yaml:"database_specific,omitempty"`
	EcosystemSpecific interface{} `json:"ecosystem_specific,omitempty" yaml:"ecosystem_specific,omitempty"`
}

// Package represents the package details in the affected structure.
type Package struct {
	Name      string `json:"name" yaml:"name"`
	Ecosystem string `json:"ecosystem" yaml:"ecosystem"`
	PURL      string `json:"purl,omitempty" yaml:"purl,omitempty"`
}

// Range represents a range of versions affected.
type Range struct {
	Type             string      `json:"type" yaml:"type"`
	Events           []Event     `json:"events" yaml:"events"`
	Repo             string      `json:"repo,omitempty" yaml:"repo,omitempty"`
	DatabaseSpecific interface{} `json:"database_specific,omitempty" yaml:"database_specific,omitempty"`
}

// Event represents individual events in a version range.
type Event struct {
	Introduced   string `json:"introduced,omitempty" yaml:"introduced,omitempty"`
	Fixed        string `json:"fixed,omitempty" yaml:"fixed,omitempty"`
	LastAffected string `json:"last_affected,omitempty" yaml:"last_affected,omitempty"`
	Limit        string `json:"limit,omitempty" yaml:"limit,omitempty"`
}

// Reference represents external links related to the vulnerability.
type Reference struct {
	Type string `json:"type" yaml:"type"`
	URL  string `json:"url" yaml:"url"`
}

// Credit represents individuals or organizations credited for reporting or fixing the vulnerability.
type Credit struct {
	Name    string   `json:"name" yaml:"name"`
	Contact []string `json:"contact,omitempty" yaml:"contact,omitempty"`
}
