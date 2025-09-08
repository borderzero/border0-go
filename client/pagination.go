package client

import (
	"context"
	"iter"
)

const defaultPageSize = 100

// paginatedResponse represents the standard schema of the http response body
// of a paginated endpoint in the Border0 API. This is only used internally.
type paginatedResponse[T any] struct {
	Pagination pagination `json:"pagination"`
	List       []T        `json:"list"`
}

// pagination represents an object that contains pagination metadata returned
// by all paginated endpoints of the Border0 API.
type pagination struct {
	CurrentPage    int `json:"current_page"`
	NextPage       int `json:"next_page"`
	TotalRecords   int `json:"total_records"`
	TotalPages     int `json:"total_pages"`
	RecordsPerPage int `json:"records_per_page"`
	ActualPageSize int `json:"actual_page_size"`
}

// fetchPageFunc defines a function that fetches a page of items.
// It should return the items, the next page number (0 if no more), and an error.
type fetchPageFunc[T any] func(ctx context.Context, api *APIClient, page, size int) ([]T, int, error)

// Paginator provides sequential access to paginated API resources.
type Paginator[T any] struct {
	api      *APIClient
	pageSize int
	nextPage int
	done     bool
	fetch    fetchPageFunc[T]
}

// PageResult represents the result of fetching a single
// page of a given resource from the Border0 API.
//
// This is currently only used by the pagination iterator.
type PageResult[T any] struct {
	Items []T
	Err   error
}

// newPaginator creates a new Paginator.
func newPaginator[T any](api *APIClient, fetch fetchPageFunc[T], pageSize int) *Paginator[T] {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	return &Paginator[T]{
		api:      api,
		pageSize: pageSize,
		nextPage: 1,
		done:     false,
		fetch:    fetch,
	}
}

// HasNext reports whether more pages are available.
func (p *Paginator[T]) HasNext() bool { return !p.done }

// Next fetches the next page of items. When no further pages exist, it marks the paginator as done
// and returns an empty slice and nil error.
func (p *Paginator[T]) Next(ctx context.Context) ([]T, error) {
	if p.done {
		var empty []T
		return empty, nil
	}
	items, next, err := p.fetch(ctx, p.api, p.nextPage, p.pageSize)
	if err != nil {
		return nil, err
	}
	if next <= 0 {
		p.done = true
	} else {
		p.nextPage = next
	}
	return items, nil
}

// Iter returns an iterator over pages of T. Each page result contains
// the result of retrieving the next page, including the items and an
// error. If an error is encountered, the iterator is finished.
func (p *Paginator[T]) Iter(ctx context.Context) iter.Seq[PageResult[T]] {
	return func(yield func(PageResult[T]) bool) {
		for p.HasNext() {
			items, err := p.Next(ctx)

			// if there's an error, yield the error and stop
			if err != nil {
				yield(PageResult[T]{Items: nil, Err: err})
				return
			}

			// no more items returned, stop without yielding
			if len(items) == 0 {
				return
			}

			// stop early if caller requests it (e.g. via break)
			if !yield(PageResult[T]{Items: items}) {
				return
			}
		}
	}
}
