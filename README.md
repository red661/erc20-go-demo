# ERC20 Go Demo

使用 Go 通过 go-ethereum 与链交互，部署并操作一个基于 OpenZeppelin 的 ERC20 合约（支持 owner 增发）。仓库已包含合约与对应的 Go 绑定代码，无需单独编译合约即可直接运行。

## 功能

- **部署合约**：指定名称、符号、初始发行量（最小单位）
- **增发代币**：仅合约 owner 可调用
- **转账代币**：从私钥账户向目标地址转账

## 环境要求

- Go（建议与 `go.mod` 一致或更高版本）
- 已有可用的以太坊节点 RPC（本地或测试网/主网），且账户有足够 gas 费
- 一个 EOA 私钥（十六进制，前缀可有可无 `0x`）
- 可选：Node.js（用于安装 OpenZeppelin 依赖，若不改合约可忽略）

## 目录结构（关键）

```text
erc20-go-demo/
  cmd/
    common.go      # 加载配置、创建签名器/客户端等公共逻辑
    deploy.go      # 部署合约入口（main）
    mint.go        # 增发入口（main，需填写合约地址）
    transfer.go    # 转账入口（main，需填写合约地址）
  contracts/
    MyToken.sol    # 基于 OZ 的 ERC20 + Ownable（owner 可增发）
  internal/token/
    token.go       # 已生成的 Go 合约绑定
  go.mod / go.sum
  package.json     # 仅用于安装 @openzeppelin/contracts 以支持 Solidity import
```

## 安装与准备

1. 拉取依赖

```bash
cd erc20-go-demo
go mod download

# 如需修改合约再编译，可安装 OZ 依赖（可选）
npm install
```

1. 新建 `.env`（根目录）

```dotenv
# RPC 与账户
RPC_URL=https://sepolia.infura.io/v3/<YOUR_KEY>
PRIVATE_KEY=0xabc...your_private_key
CHAIN_ID=11155111         # 例如 Sepolia=11155111, OP Sepolia=11155420 等

# 代币元数据与初始发行量（最小单位，通常等于 10^decimals）
TOKEN_NAME=MyToken
TOKEN_SYMBOL=MTK
INITIAL_SUPPLY_WEI=1000000000000000000000  # 例：1000 * 10^18

# 可选：转账所需
TO_ADDRESS=0xYourReceiverAddress
TRANSFER_AMOUNT_WEI=1000000000000000000   # 1 * 10^18
```

提示：本合约默认 `decimals=18`（OpenZeppelin ERC20）。`*_WEI` 字段请按最小单位填写。

## 使用方法

重要：`cmd` 目录下包含多个 `main()`，请按子命令分别运行，使用 `go run` 显式传入对应文件与 `common.go`，避免一次性构建整个 `cmd` 目录导致多个入口冲突。

### 1. 部署合约

```bash
go run ./cmd/deploy.go ./cmd/common.go
```

输出示例：

```text
🚀 Deploying MyToken... tx=0x...
📜 Contract address: 0xYourContractAddress
Token: MyToken (MTK) decimals=18
TotalSupply: 1000000000000000000000
Deployer balance: 1000000000000000000000
✅ Deployed.
```

请记录合约地址，后续增发与转账需要填写到对应文件中。

### 2. 增发代币（owner）

编辑 `cmd/mint.go`，将占位符合约地址替换为部署得到的地址：

```go
contractAddr := common.HexToAddress("<PUT-YOUR-CONTRACT-ADDRESS>")
```

`mint.go` 中默认示例为增发 `10 * 10^18`，可根据需要修改。

执行：

```bash
go run ./cmd/mint.go ./cmd/common.go
```

### 3. 转账代币

编辑 `cmd/transfer.go`，替换合约地址占位符：

```go
contractAddr := common.HexToAddress("<PUT-YOUR-CONTRACT-ADDRESS>")
```

并在 `.env` 中设置 `TO_ADDRESS` 与 `TRANSFER_AMOUNT_WEI`。

执行：

```bash
go run ./cmd/transfer.go ./cmd/common.go
```

## 合约说明

- 源码：`contracts/MyToken.sol`
- 依赖：`@openzeppelin/contracts`（已在 `package.json` 声明）
- 特性：
  - 标准 ERC20 接口（`name/symbol/decimals/totalSupply/balanceOf/transfer/...`）
  - `owner` 可调用 `mint(address,uint256)` 增发

## 重新生成 Go 绑定（可选）

仓库已提供 `internal/token/token.go`。若你修改了 `contracts/MyToken.sol`，可使用 `abigen` 重新生成绑定（需本地安装 `abigen` 与 `solc`，并确保 `node_modules` 可被编译器访问）：

```bash
# 确保已安装 OZ 依赖
npm install

# 简单方式（abigen 会调用 solc）：
abigen --sol contracts/MyToken.sol \
  --pkg token \
  --out internal/token/token.go
```

如遇 import 路径解析问题，请根据本地环境为 `solc`/`abigen` 配置 include/allow-paths 或改用 Foundry/Hardhat 编译后再用 `abigen --bin --abi` 生成。

## 常见问题

- 运行报多个 `main` 冲突：请按示例使用 `go run ./cmd/<subcmd>.go ./cmd/common.go`，不要对整个 `cmd` 目录直接 `go build`。
- 交易卡住或失败：检查 `RPC_URL` 可用性、`CHAIN_ID` 与目标网络一致、账户余额是否足够支付 gas。
- 小数与单位：本示例按 18 位小数，数值以最小单位（`wei`）传入。

## 免责声明

本仓库仅用于学习与演示。上主网前请充分审计与测试，谨慎保管私钥与资金。
