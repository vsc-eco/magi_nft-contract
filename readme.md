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
- **Token Properties**: Attach arbitrary JSON metadata to tokens at mint time
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
| `setProperties` / `getProperties` | Custom |
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

## Token Properties

Each token can have arbitrary JSON properties attached to it. Properties are set during the first mint and can be updated later by the contract owner.

### Setting Properties at Mint

Include `properties` in the mint payload (any valid JSON value):

```json
{"to": "hive:collector", "id": "art-001", "amount": 1, "maxSupply": 1, "properties": {"rarity": "legendary", "power": 42}, "data": ""}
```

Properties can be any valid JSON — objects, arrays, strings, numbers, or booleans:

```json
{"to": "hive:collector", "id": "tag-001", "amount": 1, "maxSupply": 1, "properties": "hello, world", "data": ""}
```

### Batch Minting with Properties

```json
{
  "to": "hive:owner",
  "ids": ["sword", "shield"],
  "amounts": [1, 1],
  "maxSupplies": [1, 1],
  "properties": [{"damage": 50}, {"defense": 30}],
  "data": ""
}
```

### Updating Properties

The contract owner can update properties after mint:

```json
{"id": "art-001", "properties": {"rarity": "legendary", "power": 99, "enchanted": true}}
```

### Reading Properties

Anyone can query token properties:

```json
{"id": "art-001"}
```

Returns `{"properties": ...}` with the stored JSON, or `{"properties": null}` if no properties are set.

### Properties Rules
- Properties can only be set on the **first mint** of a token (like maxSupply and soulbound)
- Subsequent mints of the same token ignore the `properties` field
- The contract **owner** can update properties at any time via `setProperties`
- Anyone can read properties via `getProperties`
- No length limitations (gas cost is the natural limit)
- Any valid JSON value is accepted

## Functions

### Actions (State-Changing)

| Function               | Payload                                                        | Access |
|------------------------|----------------------------------------------------------------|--------|
| `init`                 | `{"name": string, "symbol": string, "baseUri": string, "trackMinted": bool}` | ContractOwner |
| `mint`                 | `{"to": string, "id": string, "amount": uint64, "maxSupply": uint64, "soulbound": bool, "properties": any, "data": string}` | Owner |
| `mintBatch`            | `{"to": string, "ids": []string, "amounts": []uint64, "maxSupplies": []uint64, "soulbound": []bool, "properties": []any, "data": string}` | Owner |
| `burn`                 | `{"from": string, "id": string, "amount": uint64}`             | Owner |
| `burnBatch`            | `{"from": string, "ids": []string, "amounts": []uint64}`       | Owner |
| `safeTransferFrom`     | `{"from": string, "to": string, "id": string, "amount": uint64, "data": string}` | Owner/Operator |
| `safeBatchTransferFrom`| `{"from": string, "to": string, "ids": []string, "amounts": []uint64, "data": string}` | Owner/Operator |
| `setApprovalForAll`    | `{"operator": string, "approved": bool}`                       | Any |
| `setProperties`        | `{"id": string, "properties": any}`                            | Owner |
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
| `getProperties`    | `{"id": string}`                             | `{"properties": any\|null}`  |
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
│   ├── helpers_test.go    # Test utilities
│   ├── init_test.go       # Initialization tests
│   ├── mint_test.go       # Mint/MintBatch tests
│   ├── burn_test.go       # Burn/BurnBatch tests
│   ├── transfer_test.go   # Transfer tests
│   ├── balance_test.go    # Balance query tests
│   ├── supply_test.go     # Supply management tests
│   ├── approval_test.go   # Operator approval tests
│   ├── uri_test.go        # URI management tests
│   ├── soulbound_test.go  # Soulbound token tests
│   ├── properties_test.go # Token properties tests
│   ├── trackminted_test.go# Track minted tests
│   ├── admin_test.go      # Admin/ownership tests
│   ├── lifecycle_test.go  # Lifecycle tests
│   ├── erc165_test.go     # ERC-165 interface tests
│   ├── benchmark_test.go  # RC consumption benchmark
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
| setProperties          | 300-400 |
| mint                   | 360-400 |
| mintBatch              | 450+    |
| safeBatchTransferFrom  | 430+    |
| setURI                 | 855     |
| init                   | 1313    |

## Benchmark: Real-World Scenario

RC consumption for a realistic NFT workflow (see `test/benchmark_test.go`):

| Step | RC (1000 = 1 HBD in the wallet) |
|------|---:|
| Init contract | 1,313 |
| Mint 100 unique NFTs with basic properties (2 batches of 50) — **total** | 107,776 |
| Mint 100 unique NFTs with basic properties (2 batches of 50) — **avg per batch** | 53,888 |
| Mint 10,000 editions of 1 NFT with basic properties | 1,783 |
| Transfer 10 unique NFTs (single) — **total** | 3,080 |
| Transfer 10 unique NFTs (single) — **avg per transfer** | 308 |
| Transfer 50 unique NFTs (batch) | 6,994 |
| Transfer 500 editions | 401 |
| Transfer 1,000 editions | 407 |
| Burn 5 unique NFTs (single) — **total** | 1,085 |
| Burn 5 unique NFTs (single) — **avg per burn** | 217 |
| Burn 20 unique NFTs (batch) | 1,767 |
| Burn 100 editions | 269 |
| Burn 1,000 editions | 268 |