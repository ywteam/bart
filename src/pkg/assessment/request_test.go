package assessment_test

import (
	"testing"
	"github.com/ywteam/ydk-go/pkg/assessment"
)

var req = &assessment.AssessmentRequest{}

func TestAddRepository(t *testing.T) {
	repo := assessment.Repository{
		Url:    "https://github.com/owner/repo",
		Owner:  "owner",
		Name:   "repo",
		Branch: "main",
		Commit: "abc123",
	}
	req.Add(repo)

	if len(req.Repositories) != 1 {
		t.Errorf("expected 1 repository, got %d", len(req.Repositories))
	}
	if req.Repositories[0] != repo {
		t.Errorf("expected %v, got %v", repo, req.Repositories[0])
	}
	req.AddRepository(
		"https://github.com",
		"owner",
		"repo",
		"main",
		"abc123",
	)
	if len(req.Repositories) != 2 {
		t.Errorf("expected 2 repositories, got %d", len(req.Repositories))
	}
	req.Add(assessment.Repository{
		Url:    "https://github.com/owner/repo",
		Owner:  "owner",
		Name:   "repo",
		Branch: "main",
		Commit: "abc123",
	})
	if len(req.Repositories) != 3 {
		t.Errorf("expected 3 repositories, got %d", len(req.Repositories))
	}
}

func TestAddArtefacts(t *testing.T) {
	// req := &AssessmentRequest{}
	artefacts := assessment.Artefacts{
		RemoteUrl: []string{"https://example.com/artifact1", "https://example.com/artifact2"},
	}
	req.Add(artefacts)

	if len(req.Artefacts.RemoteUrl) != 2 {
		t.Errorf("expected 2 artefacts, got %d", len(req.Artefacts.RemoteUrl))
	}
}

func TestAddIps(t *testing.T) {
	// req := &AssessmentRequest{}
	ips := assessment.Ips{
		Addresses: []string{"192.168.1.1", "192.168.1.2"},
	}
	req.Add(ips)

	if len(req.Ips.Addresses) != 2 {
		t.Errorf("expected 2 IP addresses, got %d", len(req.Ips.Addresses))
	}
}

func TestAddHostnames(t *testing.T) {
	// req := &AssessmentRequest{}
	hostnames := assessment.Hostnames{
		Hostname: []string{"example.com", "test.com"},
	}
	req.Add(hostnames)

	if len(req.Hostnames.Hostname) != 2 {
		t.Errorf("expected 2 hostnames, got %d", len(req.Hostnames.Hostname))
	}
}

func TestAddDomains(t *testing.T) {
	// req := &AssessmentRequest{}
	domains := assessment.Domains{
		Names: []string{"example", "test"},
		Tdl:   []string{"com", "org"},
	}
	req.Add(domains)

	if len(req.Domains.Names) != 2 {
		t.Errorf("expected 2 domain names, got %d", len(req.Domains.Names))
	}
	if len(req.Domains.Tdl) != 2 {
		t.Errorf("expected 2 TDLs, got %d", len(req.Domains.Tdl))
	}
}

func TestAddEndpoints(t *testing.T) {
	// req := &AssessmentRequest{}
	endpoints := assessment.Endpoints{
		Url: []string{"https://example.com/api1", "https://example.com/api2"},
	}
	req.Add(endpoints)

	if len(req.Endpoints.Url) != 2 {
		t.Errorf("expected 2 endpoints, got %d", len(req.Endpoints.Url))
	}
}

func TestAddAccounts(t *testing.T) {
	// req := &AssessmentRequest{}
	account := assessment.Accounts{
		Username: "user",
		Email:    "user@example.com",
		Password: "password",
	}
	req.Add(account)

	if len(req.Accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(req.Accounts))
	}
	if req.Accounts[0] != account {
		t.Errorf("expected %v, got %v", account, req.Accounts[0])
	}
}

func TestAddSecrets(t *testing.T) {
	// req := &AssessmentRequest{}
	secret := assessment.Secrets{
		Secret: "mysecret",
		Type:   "password",
	}
	req.Add(secret)

	if len(req.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(req.Secrets))
	}
	if req.Secrets[0] != secret {
		t.Errorf("expected %v, got %v", secret, req.Secrets[0])
	}
}

func TestAddPackages(t *testing.T) {
	// req := &AssessmentRequest{}
	pkg := assessment.Packages{
		Ecosystem: "npm",
		Name:      "package",
		Version:   "1.0.0",
		Commit:    "abc123",
		Url:       "https://example.com/package",
	}
	req.Add(pkg)

	if len(req.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(req.Packages))
	}
	if req.Packages[0] != pkg {
		t.Errorf("expected %v, got %v", pkg, req.Packages[0])
	}
}

func TestAddPorts(t *testing.T) {
	// req := &AssessmentRequest{}
	ports := assessment.Ports{
		Hosts: []string{"host1", "host2"},
		Start: 80,
		End:   8080,
	}
	req.Add(ports)

	if len(req.Ports) != 1 {
		t.Errorf("expected 1 ports entry, got %d", len(req.Ports))
	}
	if !comparePorts(req.Ports[0], ports) {
		t.Errorf("expected %v, got %v", ports, req.Ports[0])
	}
}

func comparePorts(a, b assessment.Ports) bool {
	if a.Start != b.Start || a.End != b.End {
		return false
	}
	if len(a.Hosts) != len(b.Hosts) {
		return false
	}
	for i := range a.Hosts {
		if a.Hosts[i] != b.Hosts[i] {
			return false
		}
	}
	return true
}

func TestAddApis(t *testing.T) {
	// req := &AssessmentRequest{}
	apis := assessment.Api{
		Endpoint: []string{"https://example.com/api1", "https://example.com/api2"},
	}
	req.Add(apis)

	if len(req.Apis.Endpoint) != 2 {
		t.Errorf("expected 2 APIs, got %d", len(req.Apis.Endpoint))
	}
}

func TestAddBrand(t *testing.T) {
	// req := &AssessmentRequest{}
	brand := assessment.Brand{
		Name:  "brand",
		Links: []string{"https://example.com", "https://example.com/about"},
	}
	req.Add(brand)

	if req.Brands[0].Name != brand.Name || !compareStringSlices(req.Brands[0].Links, brand.Links) {
		t.Errorf("expected %v, got %v", brand, req.Brands[0])
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestJson(t *testing.T) {
	t.Log(req.Json(true))
}
