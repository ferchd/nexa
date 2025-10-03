package types

type SummaryStats struct {
    TotalChecks     int `json:"total_checks"`
    Successful      int `json:"successful"`
    Failed          int `json:"failed"`
    ExternalChecks  int `json:"external_checks"`
    CorporateChecks int `json:"corporate_checks"`
}