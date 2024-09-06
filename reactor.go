package hs2

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proxy"
)

type ConsensusReactor struct {
	proxyApp     proxy.AppConns
	height       int64
	proposerAddr []byte
}

func NewConsensusReactor(proxyApp proxy.AppConns, height int64) *ConsensusReactor {
	proposerAddr, _ := hex.DecodeString("59CD8C9232F9C6576FCB4BF7BD8551B24D2172E7")
	return &ConsensusReactor{proxyApp: proxyApp, height: height, proposerAddr: proposerAddr}
}

func (c *ConsensusReactor) Run(ctx context.Context) {
	timer := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer.C:
			proposalRes := c.CallPrepareProposal(ctx, c.height+1)
			processRes := c.CallProcessProposal(ctx, c.height+1, proposalRes.GetTxs())
			fmt.Println(processRes.GetStatus().String())
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(c.height+1))
			finalizeRes := c.CallFinalizeBlock(ctx, c.height+1, proposalRes.GetTxs(), b)
			fmt.Println("app hash:", finalizeRes.AppHash)
		case <-ctx.Done():
			return
		}
	}
}

func (c *ConsensusReactor) CallPrepareProposal(ctx context.Context, height int64) *abcitypes.ResponsePrepareProposal {
	req := &abcitypes.RequestPrepareProposal{
		Txs:             make([][]byte, 0),
		MaxTxBytes:      94371840,
		Height:          height,
		ProposerAddress: c.proposerAddr,
	}
	res, err := c.proxyApp.Consensus().PrepareProposal(ctx, req)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return res
}

func (c *ConsensusReactor) CallProcessProposal(ctx context.Context, height int64, txs [][]byte) *abcitypes.ResponseProcessProposal {
	req := &abcitypes.RequestProcessProposal{
		Txs:             txs,
		Height:          height,
		ProposerAddress: c.proposerAddr,
	}
	res, err := c.proxyApp.Consensus().ProcessProposal(ctx, req)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return res
}

func (c *ConsensusReactor) CallFinalizeBlock(ctx context.Context, height int64, txs [][]byte, hash []byte) *abcitypes.ResponseFinalizeBlock {
	req := &abcitypes.RequestFinalizeBlock{
		Txs:             txs,
		Hash:            hash,
		Height:          height,
		ProposerAddress: c.proposerAddr,
	}
	res, err := c.proxyApp.Consensus().FinalizeBlock(ctx, req)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return res
}
