package types

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

// QueryPool liquidity query endpoint supported by the liquidity querier
const (
	QueryPool  = "pool"
	QueryPools = "pools"
)

// QueryPoolParams is the query parameters for 'custom/liquidity'
type QueryPoolParams struct {
	PoolId uint64 `json:"pool_id" yaml:"pool_id"`
}

// return params of Liquidity Pool Query
func NewQueryPoolParams(poolId uint64) QueryPoolParams {
	return QueryPoolParams{
		PoolId: poolId,
	}
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/liquidity/pools'
type QueryPoolsParams struct {
	Page, Limit int
}

// return params of Liquidity Pools Query
func NewQueryPoolsParams(page, limit int) QueryPoolsParams {
	return QueryPoolsParams{page, limit}
}
