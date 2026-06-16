package helpers

import maps0 "maps"

func mergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		maps0.Copy(result, m)
	}

	return result
}

var ValidationErrMsg = mergeMaps(
	FilterErrMSg,
	RegisterErrMsg,
	ServiceTypeErrMsg,
	ProviderErrMsg,
	UserErrMsg,
	UserRatingErrMsg,
	JobErrMsg,
)

var FilterErrMSg = map[string]string{
	"BaseFilters.Search": "must only be alphanumeric and spaces",
	"BaseFilters.Sort":   "must only be alpha and spaces",
}

var RegisterErrMsg = map[string]string{
	"RegisterRequest.Email":    "must be correct email format",
	"RegisterRequest.Password": "must have minimum 8 characters",
}

var ServiceTypeErrMsg = map[string]string{
	"ServiceTypeRequest.Name":  "must only be characters and spaces",
	"ServiceTypeRequest.Price": "must have base price",
}

var ProviderErrMsg = map[string]string{
	"ProviderRequest.Name": "must only be characters and spaces",
}

var UserErrMsg = map[string]string{
	"UserRequest.FirstName": "must have first name",
	"UserRequest.LastName":  "must have last name",
	"UserRequest.Email":     "must be correct email format",
	"UserRequest.Username":  "must have username",
	"UserRequest.Password":  "must have minimum 8 characters",
}

var UserRatingErrMsg = map[string]string{
	"UserRatingRequest.Comment": "must only be alphanumeric and spaces",
}

var JobErrMsg = map[string]string{
	"JobRequest.Name":        "must only be alphanumeric and spaces",
	"JobRequest.Description": "must only be alphanumeric and spaces",
}
