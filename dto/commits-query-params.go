package dto

type CommitQueryParams struct {
	SHA   string
	Since string
	Until string
}

func (p CommitQueryParams) String() string {
	queryString := ""
	if p.SHA != "" {
		queryString += "sha=" + p.SHA
	}
	if p.Since != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += "since=" + p.Since
	}
	if p.Until != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += "until=" + p.Until
	}
	if queryString != "" {
		return "?" + queryString
	}
	return queryString
}
