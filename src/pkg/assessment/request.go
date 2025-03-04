package assessment

import (
	"encoding/json"
	"errors"

	"github.com/ywteam/ydk-go/lib/logger"
)

type Repository struct {
	Url    string `json:"url"`
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}
type Artefacts struct {
	RemoteUrl []string `json:"remoteUrl"`
}
type Ips struct {
	Addresses []string `json:"addresses"`
}
type Hostnames struct {
	Hostname []string `json:"hostname"`
}
type Domains struct {
	Names []string `json:"names"`
	Tdl   []string `json:"tdl"`
}
type Endpoints struct {
	Url []string `json:"url"`
}
type Accounts struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Secrets struct {
	Secret string `json:"secret"`
	Type   string `json:"type"`
}
type Packages struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	/**
	* @example
	*	pkg:pypi/mlflow@0.4.0
	*	pkg:github.com/owner/repo@v1.0.0
	 */
	Url string `json:"url"`
}
type Ports struct {
	Hosts []string `json:"hosts"`
	Start int      `json:"start"`
	End   int      `json:"end"`
}
type Api struct {
	Endpoint []string `json:"url"`
	Openapi  string   `json:"openapi"`
}
type Brand struct {
	Name  string   `json:"name"`
	Links []string `json:"links"`
}

type CloudServiceProviderName string
const (
	AwsCloudServiceProviderName   CloudServiceProviderName = "aws"
	AzureCloudServiceProviderName CloudServiceProviderName = "azure"
	GcpCloudServiceProviderName   CloudServiceProviderName = "gcp"
)

type CloudServiceProvider struct {
	Name        CloudServiceProviderName `json:"name"`
	Credentials string                   `json:"credentials"`
}
type AssessmentRequest struct {
	Repositories         []Repository           `json:"repositories"`
	Artefacts            Artefacts              `json:"artefacts"`
	Ips                  Ips                    `json:"ips"`
	Hostnames            Hostnames              `json:"hostnames"`
	Domains              Domains                `json:"domains"`
	Endpoints            Endpoints              `json:"endpoints"`
	Accounts             []Accounts             `json:"accounts"`
	Secrets              []Secrets              `json:"secrets"`
	Packages             []Packages             `json:"packages"`
	Ports                []Ports                `json:"ports"`
	Apis                 Api                    `json:"apis"`
	Brands               []Brand                `json:"brands"`
	CloudServiceProvider []CloudServiceProvider `json:"csp"`
	Assessments          `json:"assessments"`
}

func (a *AssessmentRequest) IsValid() bool {
	return true
}
func (a *AssessmentRequest) AddRepository(url, owner, name, branch, commit string) *AssessmentRequest {
	return a.Add(Repository{
		Url:    url,
		Owner:  owner,
		Name:   name,
		Branch: branch,
		Commit: commit,
	})
}
func (a *AssessmentRequest) AddArtefact(remoteUrl ...string) *AssessmentRequest {
	return a.Add(&Artefacts{RemoteUrl: remoteUrl})
}
func (a *AssessmentRequest) AddIp(addresses ...string) *AssessmentRequest {
	return a.Add(&Ips{Addresses: addresses})
}
func (a *AssessmentRequest) AddHostname(hostname ...string) *AssessmentRequest {
	return a.Add(&Hostnames{Hostname: hostname})
}
func (a *AssessmentRequest) AddDomain(names, tdl []string) *AssessmentRequest {
	return a.Add(&Domains{Names: names, Tdl: tdl})
}
func (a *AssessmentRequest) AddEndpoint(url ...string) *AssessmentRequest {
	return a.Add(&Endpoints{Url: url})
}
func (a *AssessmentRequest) AddAccount(username, email, password string) *AssessmentRequest {
	return a.Add(&Accounts{Username: username, Email: email, Password: password})
}
func (a *AssessmentRequest) AddSecret(secret string, ttype string) *AssessmentRequest {
	return a.Add(&Secrets{Secret: secret, Type: ttype})
}
func (a *AssessmentRequest) AddPackage(ecosystem, name, version, commit, url string) *AssessmentRequest {
	return a.Add(&Packages{Ecosystem: ecosystem, Name: name, Version: version, Commit: commit, Url: url})
}
func (a *AssessmentRequest) AddPort(start, end int, hosts ...string) *AssessmentRequest {
	return a.Add(&Ports{Hosts: hosts, Start: start, End: end})
}
func (a *AssessmentRequest) AddApi(url ...string) *AssessmentRequest {
	return a.Add(&Api{Endpoint: url})
}
func (a *AssessmentRequest) AddOpenapi(openapi string) *AssessmentRequest {
	return a.Add(&Api{Openapi: openapi})
}
func (a *AssessmentRequest) AddCloudServiceProvider(name CloudServiceProviderName, credentials string) *AssessmentRequest {
	return a.Add(&CloudServiceProvider{Name: name, Credentials: credentials})
}
func (a *AssessmentRequest) AddBrand(name string, links ...string) *AssessmentRequest {
	return a.Add(&Brand{Name: name, Links: links})
}
func (a *AssessmentRequest) Add(args ...interface{}) *AssessmentRequest {
	for _, arg := range args {
		switch v := arg.(type) {
		case Repository:
			if a.Repositories == nil {
				a.Repositories = []Repository{}
			}
			a.Repositories = append(a.Repositories, v)
		case Artefacts:
			if a.Artefacts.RemoteUrl == nil {
				a.Artefacts.RemoteUrl = []string{}
			}
			a.Artefacts.RemoteUrl = append(a.Artefacts.RemoteUrl, v.RemoteUrl...)
		case Ips:
			if a.Ips.Addresses == nil {
				a.Ips.Addresses = []string{}
			}
			a.Ips.Addresses = append(a.Ips.Addresses, v.Addresses...)
		case Hostnames:
			if a.Hostnames.Hostname == nil {
				a.Hostnames.Hostname = []string{}
			}
			a.Hostnames.Hostname = append(a.Hostnames.Hostname, v.Hostname...)
			// a.Hostnames = append(a.Hostnames, v)
		case Domains:
			if a.Domains.Names == nil {
				a.Domains.Names = []string{}
			}
			if a.Domains.Tdl == nil {
				a.Domains.Tdl = []string{}
			}
			a.Domains.Names = append(a.Domains.Names, v.Names...)
			a.Domains.Tdl = append(a.Domains.Tdl, v.Tdl...)
			// a.Domains = append(a.Domains, v)
		case Endpoints:
			if a.Endpoints.Url == nil {
				a.Endpoints.Url = []string{}
			}
			a.Endpoints.Url = append(a.Endpoints.Url, v.Url...)
			// a.Endpoints = append(a.Endpoints, v)
		case Accounts:
			a.Accounts = append(a.Accounts, v)
		case Secrets:
			if a.Secrets == nil {
				a.Secrets = []Secrets{}
			}
			a.Secrets = append(a.Secrets, v)
			// a.Secrets = append(a.Secrets, v)
		case Packages:
			a.Packages = append(a.Packages, v)
		case Ports:
			a.Ports = append(a.Ports, v)
		case Api:
			if a.Apis.Endpoint == nil {
				a.Apis.Endpoint = []string{}
			}
			a.Apis.Endpoint = append(a.Apis.Endpoint, v.Endpoint...)
		case Brand:
			if a.Brands == nil {
				a.Brands = []Brand{}
			}
			a.Brands = append(a.Brands, v)
		case CloudServiceProvider:
			if a.CloudServiceProvider == nil {
				a.CloudServiceProvider = []CloudServiceProvider{}
			}
			a.CloudServiceProvider = append(a.CloudServiceProvider, v)
		default:
			return nil
		}
	}
	return a
}
func FactoryAssessmentRequest(args ...interface{}) (*AssessmentRequest, error) {
	var request = &AssessmentRequest{
		Repositories:         []Repository{},
		Artefacts:            Artefacts{},
		Ips:                  Ips{},
		Hostnames:            Hostnames{},
		Domains:              Domains{},
		Endpoints:            Endpoints{},
		Accounts:             []Accounts{},
		Secrets:              []Secrets{},
		Packages:             []Packages{},
		Ports:                []Ports{},
		Apis:                 Api{},
		Brands:               []Brand{},
		CloudServiceProvider: []CloudServiceProvider{},
	}
	for _, arg := range args {
		request.Add(arg)
	}
	if !request.IsValid() {
		return nil, errors.New("invalid assessment request")
	}
	return request, nil
}
func (a *AssessmentRequest) Json(indent bool) string {
	if indent {
		raw, _ := json.MarshalIndent(a, "", "  ")
		return string(raw)
	}
	raw, _ := json.Marshal(a)
	return string(raw)
}
func NewAssessmentRequest(assets ...interface{}) *AssessmentRequest {
	request, err := FactoryAssessmentRequest(assets...)
	if err != nil {
		logger.Error("Error creating assessment request %s", err)
		return nil
	}
	return request
}
func NewAssessmentRequestFromJson(raw string) (*AssessmentRequest, error) {
	var request AssessmentRequest
	if err := json.Unmarshal([]byte(raw), &request); err != nil {
		return nil, err
	}
	if !request.IsValid() {
		return nil, errors.New("invalid assessment request")
	}
	return &request, nil
}
