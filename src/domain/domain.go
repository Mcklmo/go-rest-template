package domain

type (
	Output struct {
		Body struct {
			Output []string `json:"output"`
		}
	}
	Input struct {
		Body struct {
			Date string `json:"date" format:"date"`
		}
	}
	Store interface{}
)
