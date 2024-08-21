package manyrowsclient

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Entity struct {
	ID               uuid.UUID              `json:"id"`
	CollectionItemID *uuid.UUID             `json:"collectionItemId,omitempty"`
	CreatedAt        *time.Time             `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
	Status           int                    `json:"status"`
	Attributes       map[string]interface{} `json:"attributes"`
}

type clientOptions struct {
	baseURL string
	apiKey  string
}

func (q *clientOptions) SetBaseURL(baseURL string) {
	q.baseURL = baseURL
}

func (q *clientOptions) SetAPIKey(apiKey string) {
	q.apiKey = apiKey
}

type QueryRequest struct {
	clientOptions
	PageRequest
	Search                   string     `json:"search"`
	Sort                     string     `json:"sort"`
	SortDirection            string     `json:"sortDirection"`
	Status                   int        `json:"status"`
	Filters                  []Filter   `json:"filters"`
	ExpandSubEntities        bool       `json:"expandSubEntities"`
	RelFilter                *RelFilter `json:"relFilter,omitempty"`
	CollectionParentEntityId uuid.UUID  `json:"collectionParentEntityId"`
}

type GetOneRequest struct {
	clientOptions
	ID uuid.UUID
}

type DeleteOneRequest struct {
	clientOptions
	ID uuid.UUID
}

type CreateEntityRequest struct {
	clientOptions
	Attributes map[string]interface{} `json:"attributes"`
	Status     int                    `json:"status"`
}

type CreateEntityResponse struct {
	ID       uuid.UUID
	Location string
}

type CreateCollectionItemRequest struct {
	clientOptions
	Entity1ID uuid.UUID `json:"entity1Id"`
	Entity2ID uuid.UUID `json:"entity2Id"`
}

type UpdateRequest struct {
	clientOptions
	Attributes map[string]interface{} `json:"attributes"`
	Status     int                    `json:"status"`
}

type DeleteRequest struct {
	clientOptions
	IDs []uuid.UUID `json:"ids"`
}

type DeleteCollectionItemsRequest struct {
	clientOptions
	CollectionItemIds []uuid.UUID `json:"ids"`
}

type MoveCollectionItemRequest struct {
	clientOptions
	CollectionItemId uuid.UUID `json:"collectionItemId"`
	Index            int64     `json:"index"`
}

type RelFilter struct {
	RelDefID uuid.UUID `json:"relDefId"`
	EntityID uuid.UUID `json:"entityId"`
}

type QueryResponse struct {
	PageResource
	Items []Entity `json:"items"`
}

type ErrorInfo struct {
	Field    string `json:"field,omitempty"`
	Extra    any    `json:"extra,omitempty"`
	Message  string `json:"message,omitempty"`
	HttpCode int    `json:"httpCode,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type Filter struct {
	AttributeKey string `json:"attributeKey"`
	Value        any    `json:"value"`
}
