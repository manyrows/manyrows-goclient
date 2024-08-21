package manyrowsclient

const (
	defaultPageSize = 50
	firstPage       = 0
	maxPageSize     = 100
	minPageSize     = 1
)

type PageRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type PageResource struct {
	PageRequest
	Total int64 `json:"total"`
}

func (p *PageRequest) GetLimit() int {
	return p.GetSize()
}

func (p *PageRequest) GetPage() int {
	if p.Page < firstPage {
		return firstPage
	}
	return p.Page
}

func (p *PageRequest) GetSize() int {
	if p.Size < minPageSize || p.Size > maxPageSize {
		return defaultPageSize
	}
	return p.Size
}

func (p *PageRequest) GetOffset() int {
	return p.GetPage() * p.GetSize()
}
