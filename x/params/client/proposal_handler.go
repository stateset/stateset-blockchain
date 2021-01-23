package client

import (
	govclient "github.com/stateset/stateset-blockchain/x/gov/client"
	"github.com/stateset/stateset-blockchain/x/params/client/cli"
	"github.com/stateset/stateset-blockchain/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler in cmsdk
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
