# MAGI NFT

ERC-1155 compliant multi-token contract for the Magi Network.

## Overview

MAGI NFT is a multi-token implementation following the ERC-1155 standard. It supports both fungible and non-fungible tokens within a single contract, with features for pausability, ownership management, batch operations, and enforced supply limits.

## Token Configuration

Token properties are configured at initialization via the `init` payload:

| Property    | Type   | Description                                           |
|-------------|--------|-------------------------------------------------------|
| name        | string | Collection name (e.g., "Magi NFT")                    |
| symbol      | string | Collection symbol (e.g., "MNFT")                      |
| baseUri     | string | Base URI for token metadata                           |
| trackMinted | bool   | Optional: if true, burned tokens cannot be re-minted  |

Example init payload:
```json
{"name": "Magi NFT", "symbol": "MNFT", "baseUri": "https://api.magi.network/metadata/"}
```

With permanent burn tracking:
```json
{"name": "Magi NFT", "symbol": "MNFT", "baseUri": "https://api.magi.network/metadata/", "trackMinted": true}
```

## Features

- **ERC-1155 Compliant**: Multi-token standard interface supporting both fungible and non-fungible tokens
- **Supply Enforcement**: Required `maxSupply` on mint (1 = unique NFT, >1 = editioned)
- **Soulbound Tokens**: Non-transferable tokens for badges, credentials, memberships
- **Batch Operations**: Efficient batch minting, burning, and transfers
- **Operator Approval**: Approve operators to manage all your tokens
- **Custom URIs**: Set individual URIs per token or use baseUri fallback
- **Mintable**: Owner can mint single or batch tokens
- **Burnable**: Owner can burn single or batch tokens
- **Pausable**: Owner can pause/unpause all transfers
- **Ownership Transfer**: Owner can transfer contract ownership

## ERC-1155 Compliance

### Standard Functions (ERC-1155 Core)

| Function | Type | Standard |
|----------|------|----------|
| `balanceOf` | Query | ERC-1155 |
| `balanceOfBatch` | Query | ERC-1155 |
| `safeTransferFrom` | Action | ERC-1155 |
| `safeBatchTransferFrom` | Action | ERC-1155 |
| `setApprovalForAll` | Action | ERC-1155 |
| `isApprovedForAll` | Query | ERC-1155 |
| `uri` | Query | ERC-1155 |
| `TransferSingle` | Event | ERC-1155 |
| `TransferBatch` | Event | ERC-1155 |
| `ApprovalForAll` | Event | ERC-1155 |
| `URI` | Event | ERC-1155 |

### Extensions

| Feature | Source |
|---------|--------|
| `supportsInterface` | ERC-165 Standard |
| `totalSupply` / `exists` | ERC-1155 Supply Extension |
| `mint` / `mintBatch` | Common pattern (not in spec) |
| `burn` / `burnBatch` | Common pattern (not in spec) |
| `pause` / `unpause` | OpenZeppelin Pausable |
| `changeOwner` / `getOwner` | OpenZeppelin Ownable |
| `maxSupply` | Custom |
| `trackMinted` / `totalMinted` | Custom |
| `soulbound` / `isSoulbound` | Inspired by EIP-5192 |
| `setURI` / `setBaseURI` | Custom |
| `getInfo` | Custom |

## Supply Management

The first mint of a token requires `maxSupply` to define the token type:

### Unique NFTs (1/1)
```json
{"to": "hive:collector", "id": "art-001", "amount": 1, "maxSupply": 1, "data": ""}
```
Once minted, no more copies can ever be created.

### Editioned NFTs
```json
{"to": "hive:collector", "id": "art-002", "amount": 50, "maxSupply": 100, "data": ""}
```
Up to 100 total can be minted across multiple transactions.

### Subsequent Mints
For existing tokens, `maxSupply` is optional:
```json
{"to": "hive:collector", "id": "art-002", "amount": 25, "data": ""}
```
Uses the stored maxSupply. If provided, it must match.

### Supply Rules
- `maxSupply` is **required** on first mint (1 = unique, >1 = editioned)
- First mint locks the max supply for that token ID (cannot be changed)
- Subsequent mints can omit `maxSupply` (uses stored value)
- Minting fails if total would exceed max supply
- Burning decreases total supply (but max supply stays the same)

### Burn Behavior (trackMinted)

| Setting | Behavior | Use Case |
|---------|----------|----------|
| `trackMinted: false` (default) | Burned tokens can be re-minted up to maxSupply | Gaming items, editions |
| `trackMinted: true` | Burned tokens are gone forever | Collectibles, certificates |

When `trackMinted` is enabled:
- `totalMinted` tracks all tokens ever minted (never decreases)
- Supply check uses `totalMinted` instead of `totalSupply`
- After burning, those slots cannot be re-minted

## Soulbound Tokens

Soulbound tokens are non-transferable tokens that become bound to the recipient. The contract owner can distribute soulbound tokens, but once received, recipients cannot transfer them further. They can still be burned.

### Use Cases
- **Achievement Badges**: Gaming accomplishments bound to player
- **Certifications**: Educational credentials bound to recipient
- **Memberships**: Non-transferable membership tokens
- **Identity Tokens**: KYC/verification tokens
- **Proof of Attendance**: Event participation tokens (POAPs)

### Distribution Workflow

The owner can pre-mint soulbound tokens and distribute them later:

```json
// Step 1: Owner mints 10 soulbound badges to themselves
{"to": "hive:owner", "id": "badge-001", "amount": 10, "maxSupply": 100, "soulbound": true, "data": ""}

// Step 2: Owner distributes badges to recipients
{"from": "hive:owner", "to": "hive:player", "id": "badge-001", "amount": 1, "data": ""}

// Step 3: Recipient CANNOT transfer further (bound to them)
```

Or mint directly to recipients:
```json
{"to": "hive:player", "id": "badge-001", "amount": 1, "maxSupply": 1, "soulbound": true, "data": ""}
```

### Batch Minting with Mixed Soulbound

```json
{
  "to": "hive:owner",
  "ids": ["badge", "item"],
  "amounts": [1, 5],
  "maxSupplies": [1, 100],
  "soulbound": [true, false],
  "data": ""
}
```

### Soulbound Rules
- `soulbound` can only be set on the **first mint** of a token (like maxSupply)
- Once set, the soulbound status cannot be changed
- **Owner can transfer** soulbound tokens (for distribution)
- **Recipients cannot transfer** soulbound tokens (bound to them)
- Soulbound tokens **can be burned** (for revocation scenarios)
- Use `isSoulbound` query to check a token's status

### Soulbound vs Regular Tokens

| Workflow | Regular Token | Soulbound Token |
|----------|---------------|-----------------|
| Pre-mint to owner | ✅ Supported | ✅ Supported |
| Owner distributes | ✅ Supported | ✅ Supported |
| Recipient transfers | ✅ Supported | ❌ Blocked |
| Burn | ✅ Supported | ✅ Supported |

## Functions

### Actions (State-Changing)

| Function               | Payload                                                        | Access |
|------------------------|----------------------------------------------------------------|--------|
| `init`                 | `{"name": string, "symbol": string, "baseUri": string, "trackMinted": bool}` | ContractOwner |
| `mint`                 | `{"to": string, "id": string, "amount": uint64, "maxSupply": uint64, "soulbound": bool, "data": string}` | Owner |
| `mintBatch`            | `{"to": string, "ids": []string, "amounts": []uint64, "maxSupplies": []uint64, "soulbound": []bool, "data": string}` | Owner |
| `burn`                 | `{"from": string, "id": string, "amount": uint64}`             | Owner |
| `burnBatch`            | `{"from": string, "ids": []string, "amounts": []uint64}`       | Owner |
| `safeTransferFrom`     | `{"from": string, "to": string, "id": string, "amount": uint64, "data": string}` | Owner/Operator |
| `safeBatchTransferFrom`| `{"from": string, "to": string, "ids": []string, "amounts": []uint64, "data": string}` | Owner/Operator |
| `setApprovalForAll`    | `{"operator": string, "approved": bool}`                       | Any |
| `setURI`               | `{"id": string, "uri": string}`                                | Owner |
| `setBaseURI`           | `{"baseUri": string}`                                          | Owner |
| `pause`                | -                                                              | Owner |
| `unpause`              | -                                                              | Owner |
| `changeOwner`          | `{"newOwner": string}`                                         | Owner |

### Queries (Read-Only)

| Function           | Payload                                      | Response                     |
|--------------------|----------------------------------------------|------------------------------|
| `balanceOf`        | `{"account": string, "id": string}`          | `{"balance": uint64}`        |
| `balanceOfBatch`   | `{"accounts": []string, "ids": []string}`    | `{"balances": []uint64}`     |
| `totalSupply`      | `{"id": string}`                             | `{"totalSupply": uint64}`    |
| `maxSupply`        | `{"id": string}`                             | `{"maxSupply": uint64}`      |
| `totalMinted`      | `{"id": string}`                             | `{"totalMinted": uint64}`    |
| `exists`           | `{"id": string}`                             | `{"exists": bool}`           |
| `isSoulbound`      | `{"id": string}`                             | `{"soulbound": bool}`        |
| `isApprovedForAll` | `{"account": string, "operator": string}`    | `{"approved": bool}`         |
| `uri`              | `{"id": string}`                             | `{"uri": string}`            |
| `getOwner`         | -                                            | `{"owner": string}`          |
| `getInfo`          | -                                            | `{"name", "symbol", "baseUri", "trackMinted"}` |
| `isPaused`         | -                                            | `{"paused": bool}`           |
| `supportsInterface`| `{"interfaceId": string}`                    | `{"supported": bool}`        |

## Events

All events include `type`, `attributes`, and `tx` (transaction ID).

| Event Type       | Attributes                                    |
|------------------|-----------------------------------------------|
| `init_magi_nft`  | `owner`, `name`, `symbol`, `baseUri`          |
| `TransferSingle` | `operator`, `from`, `to`, `id`, `value`       |
| `TransferBatch`  | `operator`, `from`, `to`, `ids`, `values`     |
| `ApprovalForAll` | `account`, `operator`, `approved`             |
| `URI`            | `value`, `id`                                 |
| `baseUriChange`  | `previousUri`, `newUri`                       |
| `ownerChange`    | `previousOwner`, `newOwner`                   |
| `paused`         | `by`                                          |
| `unpaused`       | `by`                                          |

### ERC-1155 Event Compliance

- **Mint**: Emits `TransferSingle` or `TransferBatch` with `from: ""`
- **Burn**: Emits `TransferSingle` or `TransferBatch` with `to: ""`
- **SetURI**: Emits `URI` event when token URI is updated
- **SetBaseURI**: Emits `baseUriChange` event when base URI is updated

## Operator Approval Pattern

Unlike ERC-20's per-token allowances, ERC-1155 uses operator approval. An approved operator can transfer ALL tokens of ALL types on behalf of the approver.

### How It Works

```
1. User approves operator for all tokens
2. Operator can transfer any token type/amount from user
3. Approval remains until explicitly revoked
```

### Marketplace Integration

**Step 1: User Approves Marketplace**

```json
{
  "action": "setApprovalForAll",
  "payload": {"operator": "hive:marketplace", "approved": true}
}
```

**Step 2: Marketplace Executes Transfer**

```json
{
  "action": "safeTransferFrom",
  "payload": {
    "from": "hive:seller",
    "to": "hive:buyer",
    "id": "42",
    "amount": 1,
    "data": ""
  }
}
```

**Step 3: Revoke Approval (Optional)**

```json
{
  "action": "setApprovalForAll",
  "payload": {"operator": "hive:marketplace", "approved": false}
}
```

## URI Management

### Default Behavior

URIs are constructed as `baseUri + tokenId`:
- baseUri: `https://api.magi.network/metadata/`
- tokenId: `123`
- Result: `https://api.magi.network/metadata/123`

### Custom URIs

Set a specific URI for individual tokens:

```json
{
  "action": "setURI",
  "payload": {"id": "1", "uri": "https://custom.example.com/token1.json"}
}
```

Custom URIs take precedence over the baseUri pattern.

## Build

```bash
tinygo build -gc=custom -scheduler=none -panic=trap -no-debug -target=wasm-unknown -o test/artifacts/main.wasm ./contract
```

## Test

```bash
go test ./test/...
```

## Project Structure

```
magi_nft/
├── contract/
│   ├── main.go            # Entry point and state helpers
│   ├── token.go           # Exported WASM functions
│   ├── internal.go        # Internal helper functions
│   ├── types.go           # Type definitions
│   ├── types_tinyjson.go  # JSON serialization (tinyjson)
│   └── events.go          # Event emission
├── sdk/                   # VSC SDK bindings
├── test/
│   ├── basic_test.go      # ERC-1155 tests
│   ├── helpers_test.go    # Test utilities
│   └── artifacts/         # Compiled WASM
└── readme.md
```

## RC Consumption

| Function               | Avg RC  |
|------------------------|---------|
| Queries                | 100     |
| unpause                | 110     |
| pause                  | 122     |
| burn                   | 185     |
| setApprovalForAll      | 190     |
| changeOwner            | 205     |
| setBaseURI             | 230-360 |
| burnBatch              | 245     |
| safeTransferFrom       | 290-320 |
| mint                   | 360-400 |
| mintBatch              | 450+    |
| safeBatchTransferFrom  | 430+    |
| setURI                 | 855     |
| init                   | 1313    |
