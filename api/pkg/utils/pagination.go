package utils

// RequestPagination defines pagination query parameters for the API.
// @Param page query int false "Page number, starting from 1" minimum(1)
// @Param page_size query int false "Number of items per page" minimum(1) maximum(100)
// @Param order query string false "Field to order by"
// @Param sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Description Pagination parameters
type RequestPagination struct {
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100"`
	Order    string `json:"order" query:"order" validate:"omitempty"`
	Sort     string `json:"sort" query:"sort" validate:"omitempty,oneof=DESC ASC"`
}

type ResponsePagination struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalRecords int `json:"total_records"`
}

func NewPaginationRequest() RequestPagination {
	return RequestPagination{
		Page:     1,
		PageSize: 50,
		Order:    "id",
		Sort:     "DESC",
	}
}

func NewPaginationResponse() *ResponsePagination {
	return &ResponsePagination{
		Page:         1,
		PageSize:     50,
		TotalRecords: 0,
	}
}
