package bank

import (
	"github.com/stateset/stateset-blockchain/x/bank/exported"
	"github.com/stateset/stateset-blockchain/x/bank/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace

	TransactionBacking                         = exported.TransactionBacking
	TransactionBackingReturned                 = exported.TransactionBackingReturned
	TransactionChallenge                       = exported.TransactionChallenge
	TransactionChallengeReturned               = exported.TransactionChallengeReturned
	TransactionUpvote                          = exported.TransactionUpvote
	TransactionUpvoteReturned                  = exported.TransactionUpvoteReturned
	TransactionInterestArgumentCreation        = exported.TransactionInterestArgumentCreation
	TransactionInterestUpvoteReceived          = exported.TransactionInterestUpvoteReceived
	TransactionInterestUpvoteGiven             = exported.TransactionInterestUpvoteGiven
	TransactionRewardPayout                    = exported.TransactionRewardPayout
	TransactionInterestArgumentCreationSlashed = exported.TransactionInterestArgumentCreationSlashed
	TransactionInterestUpvoteReceivedSlashed   = exported.TransactionInterestUpvoteReceivedSlashed
	TransactionInterestUpvoteGivenSlashed      = exported.TransactionInterestUpvoteGivenSlashed
	TransactionStakeCreatorSlashed             = exported.TransactionStakeCreatorSlashed
	TransactionStakeCuratorSlashed             = exported.TransactionStakeCuratorSlashed

	TransactionCuratorReward = exported.TransactionCuratorReward

	SortAsc                    = exported.SortAsc
	SortDesc                   = exported.SortDesc
	QueryTransactionsByAddress = exported.QueryTransactionsByAddress
	QueryParams                = exported.QueryParams
	RouterKey                  = exported.RouterKey
)

var (
	GetFilters              = exported.GetFilters
	FilterByTransactionType = exported.FilterByTransactionType
	SortOrder               = exported.SortOrder
	Limit                   = exported.Limit
	Offset                  = exported.Offset
	FromModuleAccount       = exported.FromModuleAccount
	ToModuleAccount         = exported.ToModuleAccount
	ModuleCodec             = types.ModuleCodec
)

type (
	TransactionSetter                = exported.TransactionSetter
	Filter                           = exported.Filter
	SortOrderType                    = exported.SortOrderType
	Transaction                      = exported.Transaction
	QueryTransactionsByAddressParams = exported.QueryTransactionsByAddressParams
)