package util

type PaginationParam struct {
	CurrentPage uint32 `json:"current_page"`
	PageSize    uint32 `json:"page_size"`
}

type PaginationResponse struct {
	CurrentPage uint32 `json:"current_page"`
	PageSize    uint32 `json:"page_size"`
	TotalItems  uint32 `json:"total_items"`
	TotalPages  uint32 `json:"total_pages"`
}

func GeneratePaginationResponse(currentPage, pageSize, totalItems uint32) PaginationResponse {
	if pageSize == 0 {
		pageSize = 1
	}

	return PaginationResponse{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  (totalItems + pageSize - 1) / pageSize,
	}
}
