package gormx

const (
	TreePathDelimiter = "."
	DefaultPageSize   = 100
)

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

// PaginationParam 分页查询条件
type PaginationParam struct {
	Pagination bool `query:"-"`        // 是否使用分页查询
	OnlyCount  bool `query:"-"`        // 是否仅查询count
	Current    int  `query:"current"`  // 当前页
	PageSize   int  `query:"pageSize"` // 页大小
}

// GetCurrent 获取当前页
func (a PaginationParam) GetCurrent() int {
	if a.Current == 0 {
		return 1
	}
	return a.Current
}

// GetPageSize 获取页大小
func (a PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = DefaultPageSize
	}
	return pageSize
}

// QueryOptions 查询可选参数项
type QueryOptions struct {
	SelectFields []string
	OmitFields   []string
	OrderFields  OrderByParams
}

// Direction 排序方向
type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

// OrderByParam 排序字段
type OrderByParam struct {
	Field     string
	Direction Direction
}

type OrderByParams []OrderByParam

func (a OrderByParams) ToSQL() string {
	if len(a) == 0 {
		return ""
	}
	var sql string
	for _, v := range a {
		sql += v.Field + " " + string(v.Direction) + ","
	}
	return sql[:len(sql)-1]
}
