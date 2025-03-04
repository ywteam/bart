package assessment

import "github.com/ywteam/ydk-go/lib/report"

type AssessmentResult struct {
	AssessmentType AssessmentType
	Report         *report.Report
	Metrics        map[string]any
}
type IAssessment interface {
	Start() []AssessmentResult
	Types() []AssessmentType
	// Any(a ...T) bool
	// Enabled(a T) bool
	// Enable(a T)
	// Disable(a T)
}

type Assessment struct {
	Request *AssessmentRequest
	Result  []AssessmentResult
}
func (a *Assessment) Start() {
	for _, assessmentType := range AllAssessmentsTypes {
		if a.Request.Assessments.Enabled(assessmentType) {
			switch assessmentType {
			case ScaAssessmentType:
				// a.Result = append(a.Result, ScaAssessment(a.Request))
			case SastAssessmentType:
				// a.Result = append(a.Result, SastAssessment(a.Request))
			case IastAssessmentType:
				// a.Result = append(a.Result, IastAssessment(a.Request))
			case DastAssessmentType:
				// a.Result = append(a.Result, DastAssessment(a.Request))
			case SbomAssessmentType:
				// a.Result = append(a.Result, SbomAssessment(a.Request))
			case VulnAssessmentType:
				// a.Result = append(a.Result, VulnAssessment(a.Request))
			case MisconfigAssessmentType:
				// a.Result = append(a.Result, MisconfigAssessment(a.Request))
			case RbacAssessmentType:
				// a.Result = append(a.Result, RbacAssessment(a.Request))
			case LicenceAssessmentType:
				// a.Result = append(a.Result, LicenceAssessment(a.Request))
			case SecretsAssessmentType:
				// a.Result = append(a.Result, SecretsAssessment(a.Request))
			}
		}
	}
}

func NewAssessment(request *AssessmentRequest) *Assessment {
	return &Assessment{
		Request: request,
		Result:  nil,		
	}
}

