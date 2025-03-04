package defectdojo


// #region Languages
type DefectDojoLanguageModel struct {
	Id         int    `json:"id"`
	Files      int    `json:"files"`
	Blank      int    `json:"blank"`
	Comment    int    `json:"comment"`
	Code       int    `json:"code"`
	Created    string `json:"created"`
	LanguageId int    `json:"language"`
	ProductId  int    `json:"product"`
	UserId     int    `json:"user"`
}
type DefectDojoLanguageResponse struct {
	DefectDojoResponse
	Results []DefectDojoLanguageModel `json:"results"`
}
// #endregion


type DefectDojoErrorResponse struct {
	Code        int      `json:"code,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Description []string `json:"description,omitempty"`
	Message     string   `json:"message,omitempty"`
}

type DefectDojoResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []interface{}
}