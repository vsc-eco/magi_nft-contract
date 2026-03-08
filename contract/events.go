package main

import (
	"magi_nft/sdk"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ==========================================
// MAGI NFT - ERC-1155 Event Emission (tinyjson)
// ==========================================

// ==============
// Init Event
// ==============

func emitInit(owner, name, symbol, baseURI string) {
	event := InitEvent{
		Type: "init_magi_nft",
		Attributes: InitAttributes{
			Owner:   owner,
			Name:    name,
			Symbol:  symbol,
			BaseURI: baseURI,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==================
// TransferSingle Event (ERC-1155)
// ==================

func emitTransferSingle(operator, from, to, id string, value uint64) {
	event := TransferSingleEvent{
		Type: "TransferSingle",
		Attributes: TransferSingleAttributes{
			Operator: operator,
			From:     from,
			To:       to,
			Id:       id,
			Value:    value,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==================
// TransferBatch Event (ERC-1155)
// ==================

func emitTransferBatch(operator, from, to string, ids []string, values []uint64) {
	event := TransferBatchEvent{
		Type: "TransferBatch",
		Attributes: TransferBatchAttributes{
			Operator: operator,
			From:     from,
			To:       to,
			Ids:      ids,
			Values:   values,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// ApprovalForAll Event (ERC-1155)
// ======================

func emitApprovalForAll(account, operator string, approved bool) {
	event := ApprovalForAllEvent{
		Type: "ApprovalForAll",
		Attributes: ApprovalForAllAttributes{
			Account:  account,
			Operator: operator,
			Approved: approved,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// URI Event (ERC-1155)
// ======================

func emitURI(value, id string) {
	event := URIEvent{
		Type: "URI",
		Attributes: URIAttributes{
			Value: value,
			Id:    id,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Owner Change Event
// ======================

func emitOwnerChange(previousOwner, newOwner string) {
	event := OwnerChangeEvent{
		Type:       "ownerChange",
		Attributes: OwnerChangeAttributes{PreviousOwner: previousOwner, NewOwner: newOwner},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Pause Events
// ======================

func emitPaused(by string) {
	event := PausedEvent{
		Type:       "paused",
		Attributes: PausedAttributes{By: by},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

func emitUnpaused(by string) {
	event := UnpausedEvent{
		Type:       "unpaused",
		Attributes: UnpausedAttributes{By: by},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// BaseURI Change Event
// ======================

func emitTemplateMint(templateId string, copyIds []string) {
	event := TemplateMintEvent{
		Type: "templateMint",
		Attributes: TemplateMintAttributes{
			TemplateId: templateId,
			CopyIds:    copyIds,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

func emitTokenCreated(tokenId string, maxSupply uint64, soulbound bool) {
	event := TokenCreatedEvent{
		Type: "tokenCreated",
		Attributes: TokenCreatedAttributes{
			TokenId:   tokenId,
			MaxSupply: maxSupply,
			Soulbound: soulbound,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

func emitBaseURIChange(previousURI, newURI string) {
	event := BaseURIChangeEvent{
		Type:       "baseUriChange",
		Attributes: BaseURIChangeAttributes{PreviousURI: previousURI, NewURI: newURI},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}
