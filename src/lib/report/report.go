package report

import (
	"encoding/json"
	sarif "github.com/owenrumney/go-sarif/v2/sarif"
)

type Report struct {
	Report *sarif.Report
	Error  error
}
func (r *Report) Json(ident bool) string {
	if ident {
		jsonBytes, err := json.MarshalIndent(r.Report, "", "  ")
		if err != nil {
			return err.Error()
		}
		return string(jsonBytes)
	}
	jsonBytes, err := json.Marshal(r.Report)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}
func (r *Report) NewRun(toolName, informationURI string) *sarif.Run {
	run := sarif.NewRunWithInformationURI(toolName, informationURI)
	run.PropertyBag = *sarif.NewPropertyBag()
	run.PropertyBag.AddString("toolName", toolName)
	r.Report.AddRun(run)
	return run
}
func NewReport() *Report {
	report, err := sarif.New(sarif.Version210)
	report.PropertyBag = *sarif.NewPropertyBag()
	report.PropertyBag.AddString("toolName", "ydk-go")
	report.PropertyBag.AddString("toolVersion", "0.0.1")
	report.PropertyBag.AddString("toolUrl", "")
	report.PropertyBag.AddString("createdAt", "")
	return &Report{
		Report: report,
		Error:  err,
	}
}