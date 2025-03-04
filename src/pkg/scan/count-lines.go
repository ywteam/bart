package scan
import (
	"github.com/ywteam/ydk-go/lib/assessment"
)

type CountLinesAssessment struct {
	
}
func (a *CountLinesAssessment) Start() {
	
}
func (a *CountLinesAssessment) Types() []assessment.AssessmentType {
	return []assessment.AssessmentType{assessment.ScaAssessmentType}
}