package abci

import (
	"github.com/tendermint/classic/crypto/merkle"
	"github.com/tendermint/go-amino-x"
)

var Package = amino.RegisterPackage(amino.NewPackage(
	"github.com/tendermint/classic/abci/types",
	"abci",
	amino.GetCallersDirname(),
).
	WithGoPkgName("abci").
	WithDependencies(
		merkle.Package,
	).
	WithTypes(

		/*
			pkg.Type{ // example of overriding type name.
				Type:             reflect.TypeOf(EmbeddedSt5{}),
				Name:             "EmbeddedSt5NameOverride",
				PointerPreferred: false,
			},
		*/

		// request types
		RequestBase{},
		RequestEcho{},
		RequestFlush{},
		RequestInfo{},
		RequestSetOption{},
		RequestInitChain{},
		RequestQuery{},
		RequestBeginBlock{},
		RequestCheckTx{},
		RequestDeliverTx{},
		RequestEndBlock{},
		RequestCommit{},

		// response types
		ResponseBase{},
		ResponseException{},
		ResponseEcho{},
		ResponseFlush{},
		ResponseInfo{},
		ResponseSetOption{},
		ResponseInitChain{},
		ResponseQuery{},
		ResponseBeginBlock{},
		ResponseCheckTx{},
		ResponseDeliverTx{},
		ResponseEndBlock{},
		ResponseCommit{},

		// error types
		StringError(""),

		// misc types
		ConsensusParams{},
		BlockParams{},
		EvidenceParams{},
		ValidatorParams{},
		ValidatorUpdate{},
		LastCommitInfo{},
		VoteInfo{},
		//Validator{},
		//Violation{},

		// events
		SimpleEvent(""),

		// mocks
		MockHeader{},

		// Params (abci/types/params.go)
	))
