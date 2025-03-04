package assessment

import (
	"slices"
)

type AssessmentType string
type Assessments []AssessmentType
func (assessments *Assessments) Any(a ...AssessmentType) bool {
	for _, v := range a {
		if assessments.Enabled(v) {
			return true
		}
	}
	return false
}
func (assessments *Assessments) Enabled(a AssessmentType) bool {
	return slices.Contains(*assessments, a)
}
func (assessments *Assessments) Enable(a AssessmentType) {
	if !assessments.Enabled(a) {
		*assessments = append(*assessments, a)
	}
}
func (assessments *Assessments) Disable(a AssessmentType) {
	if assessments.Enabled(a) {
		for i, v := range *assessments {
			if v == a {
				*assessments = append((*assessments)[:i], (*assessments)[i+1:]...)
				break
			}
		}
	}
}
const (
	ScaAssessmentType       AssessmentType = "sca"
	SastAssessmentType      AssessmentType = "sast"
	IastAssessmentType      AssessmentType = "iast"
	DastAssessmentType      AssessmentType = "dast"
	SbomAssessmentType      AssessmentType = "sbom"
	VulnAssessmentType      AssessmentType = "vuln"
	MisconfigAssessmentType AssessmentType = "misconfig"
	RbacAssessmentType      AssessmentType = "rbac"
	LicenceAssessmentType   AssessmentType = "licence"
	SecretsAssessmentType   AssessmentType = "secrets"
	CspAssessmentType       AssessmentType = "csp"
)

var (
	AllAssessmentsTypes = Assessments{
		ScaAssessmentType,
		SastAssessmentType,
		IastAssessmentType,
		DastAssessmentType,
		SbomAssessmentType,
		VulnAssessmentType,
		MisconfigAssessmentType,
		RbacAssessmentType,
		LicenceAssessmentType,
		SecretsAssessmentType,
		CspAssessmentType,
	}	

	AssessmentTypeName = map[AssessmentType]string{
		ScaAssessmentType:       "Software Composition Analysis",
		SastAssessmentType:      "Static Application Security Testing",
		IastAssessmentType:      "Interactive Application Security Testing",
		DastAssessmentType:      "Dynamic Application Security Testing",
		SbomAssessmentType:      "Software Bill of Materials",
		VulnAssessmentType:      "Vulnerability Assessment",
		MisconfigAssessmentType: "Misconfiguration Assessment",
		RbacAssessmentType:      "Role Based Access Control Assessment",
		LicenceAssessmentType:   "Licence Compliance Assessment",
		SecretsAssessmentType:   "Secrets Management Assessment",
		CspAssessmentType:       "Cloud Service Provider Assessment",
	}
	AssessmentAssetMap = map[AssessmentType][]interface{}{
		ScaAssessmentType:       {Repository{}, Artefacts{}, Packages{}},
		SastAssessmentType:      {Repository{}, Artefacts{}},
		IastAssessmentType:      {Repository{}, Artefacts{}},
		DastAssessmentType:      {Ips{}, Hostnames{}, Domains{}, Endpoints{}, Ports{}, Api{}, Brand{}},
		SbomAssessmentType:      {Repository{}, Artefacts{}, Packages{}},
		VulnAssessmentType:      {Repository{}, Artefacts{}, Ips{}, Hostnames{}, Domains{}, Endpoints{}, Ports{}, Api{}, Brand{}},
		MisconfigAssessmentType: {Repository{}, Artefacts{}, Ips{}, Hostnames{}, Domains{}, Endpoints{}, Ports{}, Api{}, Brand{}},
		RbacAssessmentType:      {Repository{}, Artefacts{}, Ips{}, Hostnames{}, Domains{}, Endpoints{}, Ports{}, Api{}, Brand{}},
		LicenceAssessmentType:   {Repository{}, Artefacts{}, Ips{}, Hostnames{}, Domains{}, Endpoints{}, Ports{}, Api{}, Brand{}},
		SecretsAssessmentType:   {Repository{}, Artefacts{}, Packages{}},
		CspAssessmentType:       {Accounts{}},
	}
)




