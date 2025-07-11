// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package OmniVRF

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ITaskMailboxTypesTaskParams is an auto generated low-level Go binding around an user-defined struct.
type ITaskMailboxTypesTaskParams struct {
	RefundCollector     common.Address
	AvsFee              *big.Int
	ExecutorOperatorSet OperatorSet
	Payload             []byte
}

// OmniVRFVRFTaskData is an auto generated low-level Go binding around an user-defined struct.
type OmniVRFVRFTaskData struct {
	TaskHash [32]byte
	Seed     *big.Int
}

// OperatorSet is an auto generated low-level Go binding around an user-defined struct.
type OperatorSet struct {
	Avs common.Address
	Id  uint32
}

// OmniVRFMetaData contains all meta data concerning the OmniVRF contract.
var OmniVRFMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_taskMailbox\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CALLBACK_GAS_PRICE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_CALLBACK_GAS_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_CALLBACK_GAS_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateTaskFee\",\"inputs\":[{\"name\":\"operatorSet\",\"type\":\"tuple\",\"internalType\":\"structOperatorSet\",\"components\":[{\"name\":\"avs\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"callbackGasDeposits\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decodeTaskData\",\"inputs\":[{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"taskData\",\"type\":\"tuple\",\"internalType\":\"structOmniVRF.VRFTaskData\",\"components\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"seed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getRandomness\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"randomness\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"handlePostTaskCreation\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"handlePostTaskResultSubmission\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestRandomness\",\"inputs\":[{\"name\":\"callbackContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requests\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskMailbox\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITaskMailbox\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatePreTaskCreation\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taskParams\",\"type\":\"tuple\",\"internalType\":\"structITaskMailboxTypes.TaskParams\",\"components\":[{\"name\":\"refundCollector\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"avsFee\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"executorOperatorSet\",\"type\":\"tuple\",\"internalType\":\"structOperatorSet\",\"components\":[{\"name\":\"avs\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"id\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatePreTaskResultSubmission\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taskHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"cert\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawCallbackGas\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CallbackFailed\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"callbackContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RandomnessFulfilled\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"randomness\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RandomnessRequested\",\"inputs\":[{\"name\":\"taskHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"callbackContract\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"seed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientCallbackGas\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCallbackGasLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTaskData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequestAlreadyFulfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequestNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[]}]",
	Bin: "0x60a0604052348015600e575f5ffd5b5060405161122d38038061122d833981016040819052602b91603b565b6001600160a01b03166080526066565b5f60208284031215604a575f5ffd5b81516001600160a01b0381168114605f575f5ffd5b9392505050565b6080516111a86100855f395f8181610379015261078a01526111a85ff3fe6080604052600436106100e4575f3560e01c8063a2de781411610087578063db6ecf6711610057578063db6ecf671461030b578063e34c3d3a1461032a578063e58606551461033d578063f42a9e1314610368575f5ffd5b8063a2de78141461029e578063a7fbe930146102b6578063b3d03582146102d5578063ba33565d146102eb575f5ffd5b806327b66026116100c257806327b660261461017957806339c33215146101b35780635ce2eec4146101d75780639d866985146101f6575f5ffd5b806305b19402146100e857806309a9d71c1461012357806309c5c4501461015a575b5f5ffd5b3480156100f3575f5ffd5b50610107610102366004610c0e565b6103b3565b6040805192151583526020830191909152015b60405180910390f35b34801561012e575f5ffd5b5061014261013d366004610d84565b610404565b6040516001600160601b03909116815260200161011a565b348015610165575f5ffd5b50610177610174366004610c0e565b50565b005b348015610184575f5ffd5b50610198610193366004610dd0565b61040c565b6040805182518152602092830151928101929092520161011a565b3480156101be575f5ffd5b506101c96207a12081565b60405190815260200161011a565b3480156101e2575f5ffd5b506101776101f1366004610c0e565b610433565b348015610201575f5ffd5b5061025c610210366004610c0e565b5f6020819052908152604090208054600182015460028301546003840154600485015460058601546006909601546001600160a01b0395861696949095169492939192909160ff169087565b604080516001600160a01b0398891681529790961660208801529486019390935260608501919091526080840152151560a083015260c082015260e00161011a565b3480156102a9575f5ffd5b506101c96402540be40081565b3480156102c1575f5ffd5b506101776102d0366004610e0a565b61053c565b3480156102e0575f5ffd5b506101c9620186a081565b3480156102f6575f5ffd5b50610177610305366004610ece565b50505050565b348015610316575f5ffd5b50610177610325366004610c0e565b6105e6565b6101c9610338366004610f4c565b610627565b348015610348575f5ffd5b506101c9610357366004610f74565b60016020525f908152604090205481565b348015610373575f5ffd5b5061039b7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161011a565b5f81815260208190526040812080548291906001600160a01b03166103eb57604051632589d98f60e11b815260040160405180910390fd5b600581015460069091015460ff90911694909350915050565b5f5b92915050565b604080518082019091525f8082526020820152818060200190518101906104069190610f94565b335f9081526001602052604090205481111561048e5760405162461bcd60e51b8152602060048201526015602482015274496e73756666696369656e74206465706f7369747360581b60448201526064015b60405180910390fd5b335f90815260016020526040812080548392906104ac908490610fdd565b90915550506040515f90339083908381818185875af1925050503d805f81146104f0576040519150601f19603f3d011682016040523d82523d5f602084013e6104f5565b606091505b50509050806105385760405162461bcd60e51b815260206004820152600f60248201526e151c985b9cd9995c8819985a5b1959608a1b6044820152606401610485565b5050565b6001600160a01b038216301461056557604051635c427cd960e01b815260040160405180910390fd5b60608101516040516313db301360e11b815230916327b660269161058c919060040161101e565b6040805180830381865afa9250505080156105c4575060408051601f3d908101601f191682019092526105c191810190610f94565b60015b6105e157604051633600154d60e21b815260040160405180910390fd5b505050565b5f8142604051602001610603929190918252602082015260400190565b604051602081830303815290604052805190602001205f1c905061053882826108f1565b5f6001600160a01b038316156106c057620186a082108061064a57506207a12082115b156106685760405163fe7036d760e01b815260040160405180910390fd5b5f6106786402540be40084611030565b90508034101561069b5760405163a2c23f0d60e01b815260040160405180910390fd5b335f90815260016020526040812080543492906106b9908490611047565b9091555050505b6040516bffffffffffffffffffffffff193360601b1660208201524460348201524260548201524360748201525f9060940160408051808303601f1901815282825280516020918201208383018352308452600182850152825180840184525f808252818401838152855160808101875233815280860183905280870188905286518451968101969096529051958501959095529195509290916060808301910160408051601f1981840301815291815291525162221dbd60e51b81529091506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690630443b7a0906107bf90849060040161105a565b6020604051808303815f875af11580156107db573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107ff91906110be565b6040805160e081018252338082526001600160a01b038b811660208085018281528587018e8152606087018d815243608089019081525f60a08a0181815260c08b018281528d8352828852918c90209a518b54908a166001600160a01b0319918216178c55955160018c01805491909a16961695909517909755915160028901555160038801555160048701555160058601805491151560ff1990921691909117905591516006909401939093558351928352820188905292975087917fee5eec2b3210f6ea36e12fd442b35d8918a76da7ff35c15dc4b8d69480752aa0910160405180910390a35050505092915050565b5f82815260208190526040902080546001600160a01b031661092657604051632589d98f60e11b815260040160405180910390fd5b600581015460ff161561094c5760405163533d99dd60e01b815260040160405180910390fd5b6006810182905560058101805460ff1916600117905560405183907f9b0aa3f92f46e24caa76b000bdf0dd495b9b390c320cf6585ae10a12b7d09edb906109969085815260200190565b60405180910390a260018101546001600160a01b0316156105e1576105e1835f81815260208181526040808320815160e08101835281546001600160a01b0390811682526001830154169381019390935260028101549183018290526003810154606084015260048101546080840152600581015460ff16151560a08401526006015460c0830152909190610a31906402540be40090611030565b82516001600160a01b03165f90815260016020526040812080549293508392909190610a5e908490610fdd565b9091555050602082015160408084015160c08501519151631f1f897f60e01b81526004810187905260248101929092526001600160a01b0390921691631f1f897f916044015f604051808303815f88803b158015610aba575f5ffd5b5087f193505050508015610acc575060015b6105e157610ad86110d5565b806308c379a003610b725750610aec6110ee565b80610af75750610b74565b82516001600160a01b03165f9081526001602052604081208054849290610b1f908490611047565b9250508190555082602001516001600160a01b0316847f58a5347ef2b6c6fcc342cf2326c32d04d662829bb952a924df81526f75ec629b83604051610b64919061101e565b60405180910390a350505050565b505b81516001600160a01b03165f9081526001602052604081208054839290610b9c908490611047565b9250508190555081602001516001600160a01b0316837f58a5347ef2b6c6fcc342cf2326c32d04d662829bb952a924df81526f75ec629b604051610c01906020808252600d908201526c2ab735b737bbb71032b93937b960991b604082015260600190565b60405180910390a3505050565b5f60208284031215610c1e575f5ffd5b5035919050565b634e487b7160e01b5f52604160045260245ffd5b6040810181811067ffffffffffffffff82111715610c5957610c59610c25565b60405250565b6080810181811067ffffffffffffffff82111715610c5957610c59610c25565b601f8201601f1916810167ffffffffffffffff81118282101715610ca557610ca5610c25565b6040525050565b80356001600160a01b0381168114610cc2575f5ffd5b919050565b5f60408284031215610cd7575f5ffd5b604051610ce381610c39565b809150610cef83610cac565b8152602083013563ffffffff81168114610d07575f5ffd5b6020919091015292915050565b5f82601f830112610d23575f5ffd5b813567ffffffffffffffff811115610d3d57610d3d610c25565b604051610d54601f8301601f191660200182610c7f565b818152846020838601011115610d68575f5ffd5b816020850160208301375f918101602001919091529392505050565b5f5f60608385031215610d95575f5ffd5b610d9f8484610cc7565b9150604083013567ffffffffffffffff811115610dba575f5ffd5b610dc685828601610d14565b9150509250929050565b5f60208284031215610de0575f5ffd5b813567ffffffffffffffff811115610df6575f5ffd5b610e0284828501610d14565b949350505050565b5f5f60408385031215610e1b575f5ffd5b610e2483610cac565b9150602083013567ffffffffffffffff811115610e3f575f5ffd5b830160a08186031215610e50575f5ffd5b604051610e5c81610c5f565b610e6582610cac565b815260208201356001600160601b0381168114610e80575f5ffd5b6020820152610e928660408401610cc7565b6040820152608082013567ffffffffffffffff811115610eb0575f5ffd5b610ebc87828501610d14565b60608301525080925050509250929050565b5f5f5f5f60808587031215610ee1575f5ffd5b610eea85610cac565b935060208501359250604085013567ffffffffffffffff811115610f0c575f5ffd5b610f1887828801610d14565b925050606085013567ffffffffffffffff811115610f34575f5ffd5b610f4087828801610d14565b91505092959194509250565b5f5f60408385031215610f5d575f5ffd5b610f6683610cac565b946020939093013593505050565b5f60208284031215610f84575f5ffd5b610f8d82610cac565b9392505050565b5f6040828403128015610fa5575f5ffd5b50604051610fb281610c39565b825181526020928301519281019290925250919050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561040657610406610fc9565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610f8d6020830184610ff0565b808202811582820484141761040657610406610fc9565b8082018082111561040657610406610fc9565b6020815260018060a01b0382511660208201526001600160601b0360208301511660408201525f604083015160018060a01b03815116606084015263ffffffff602082015116608084015250606083015160a080840152610e0260c0840182610ff0565b5f602082840312156110ce575f5ffd5b5051919050565b5f60033d11156110eb5760045f5f3e505f5160e01c5b90565b5f60443d10156110fb5790565b6040513d600319016004823e80513d602482011167ffffffffffffffff8211171561112557505090565b808201805167ffffffffffffffff811115611141575050505090565b3d840160031901828201602001111561115b575050505090565b61116a60208285010185610c7f565b50939250505056fea264697066735822122060d9f7a09506aa23ffa65edd506b5fa2b033ceeee19591f04667217390287da564736f6c634300081b0033",
}

// OmniVRFABI is the input ABI used to generate the binding from.
// Deprecated: Use OmniVRFMetaData.ABI instead.
var OmniVRFABI = OmniVRFMetaData.ABI

// OmniVRFBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OmniVRFMetaData.Bin instead.
var OmniVRFBin = OmniVRFMetaData.Bin

// DeployOmniVRF deploys a new Ethereum contract, binding an instance of OmniVRF to it.
func DeployOmniVRF(auth *bind.TransactOpts, backend bind.ContractBackend, _taskMailbox common.Address) (common.Address, *types.Transaction, *OmniVRF, error) {
	parsed, err := OmniVRFMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OmniVRFBin), backend, _taskMailbox)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OmniVRF{OmniVRFCaller: OmniVRFCaller{contract: contract}, OmniVRFTransactor: OmniVRFTransactor{contract: contract}, OmniVRFFilterer: OmniVRFFilterer{contract: contract}}, nil
}

// OmniVRF is an auto generated Go binding around an Ethereum contract.
type OmniVRF struct {
	OmniVRFCaller     // Read-only binding to the contract
	OmniVRFTransactor // Write-only binding to the contract
	OmniVRFFilterer   // Log filterer for contract events
}

// OmniVRFCaller is an auto generated read-only Go binding around an Ethereum contract.
type OmniVRFCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OmniVRFTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OmniVRFTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OmniVRFFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OmniVRFFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OmniVRFSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OmniVRFSession struct {
	Contract     *OmniVRF          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OmniVRFCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OmniVRFCallerSession struct {
	Contract *OmniVRFCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OmniVRFTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OmniVRFTransactorSession struct {
	Contract     *OmniVRFTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OmniVRFRaw is an auto generated low-level Go binding around an Ethereum contract.
type OmniVRFRaw struct {
	Contract *OmniVRF // Generic contract binding to access the raw methods on
}

// OmniVRFCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OmniVRFCallerRaw struct {
	Contract *OmniVRFCaller // Generic read-only contract binding to access the raw methods on
}

// OmniVRFTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OmniVRFTransactorRaw struct {
	Contract *OmniVRFTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOmniVRF creates a new instance of OmniVRF, bound to a specific deployed contract.
func NewOmniVRF(address common.Address, backend bind.ContractBackend) (*OmniVRF, error) {
	contract, err := bindOmniVRF(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OmniVRF{OmniVRFCaller: OmniVRFCaller{contract: contract}, OmniVRFTransactor: OmniVRFTransactor{contract: contract}, OmniVRFFilterer: OmniVRFFilterer{contract: contract}}, nil
}

// NewOmniVRFCaller creates a new read-only instance of OmniVRF, bound to a specific deployed contract.
func NewOmniVRFCaller(address common.Address, caller bind.ContractCaller) (*OmniVRFCaller, error) {
	contract, err := bindOmniVRF(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OmniVRFCaller{contract: contract}, nil
}

// NewOmniVRFTransactor creates a new write-only instance of OmniVRF, bound to a specific deployed contract.
func NewOmniVRFTransactor(address common.Address, transactor bind.ContractTransactor) (*OmniVRFTransactor, error) {
	contract, err := bindOmniVRF(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OmniVRFTransactor{contract: contract}, nil
}

// NewOmniVRFFilterer creates a new log filterer instance of OmniVRF, bound to a specific deployed contract.
func NewOmniVRFFilterer(address common.Address, filterer bind.ContractFilterer) (*OmniVRFFilterer, error) {
	contract, err := bindOmniVRF(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OmniVRFFilterer{contract: contract}, nil
}

// bindOmniVRF binds a generic wrapper to an already deployed contract.
func bindOmniVRF(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OmniVRFMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OmniVRF *OmniVRFRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OmniVRF.Contract.OmniVRFCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OmniVRF *OmniVRFRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OmniVRF.Contract.OmniVRFTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OmniVRF *OmniVRFRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OmniVRF.Contract.OmniVRFTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OmniVRF *OmniVRFCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OmniVRF.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OmniVRF *OmniVRFTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OmniVRF.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OmniVRF *OmniVRFTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OmniVRF.Contract.contract.Transact(opts, method, params...)
}

// CALLBACKGASPRICE is a free data retrieval call binding the contract method 0xa2de7814.
//
// Solidity: function CALLBACK_GAS_PRICE() view returns(uint256)
func (_OmniVRF *OmniVRFCaller) CALLBACKGASPRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "CALLBACK_GAS_PRICE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CALLBACKGASPRICE is a free data retrieval call binding the contract method 0xa2de7814.
//
// Solidity: function CALLBACK_GAS_PRICE() view returns(uint256)
func (_OmniVRF *OmniVRFSession) CALLBACKGASPRICE() (*big.Int, error) {
	return _OmniVRF.Contract.CALLBACKGASPRICE(&_OmniVRF.CallOpts)
}

// CALLBACKGASPRICE is a free data retrieval call binding the contract method 0xa2de7814.
//
// Solidity: function CALLBACK_GAS_PRICE() view returns(uint256)
func (_OmniVRF *OmniVRFCallerSession) CALLBACKGASPRICE() (*big.Int, error) {
	return _OmniVRF.Contract.CALLBACKGASPRICE(&_OmniVRF.CallOpts)
}

// MAXCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0x39c33215.
//
// Solidity: function MAX_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFCaller) MAXCALLBACKGASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "MAX_CALLBACK_GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0x39c33215.
//
// Solidity: function MAX_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFSession) MAXCALLBACKGASLIMIT() (*big.Int, error) {
	return _OmniVRF.Contract.MAXCALLBACKGASLIMIT(&_OmniVRF.CallOpts)
}

// MAXCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0x39c33215.
//
// Solidity: function MAX_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFCallerSession) MAXCALLBACKGASLIMIT() (*big.Int, error) {
	return _OmniVRF.Contract.MAXCALLBACKGASLIMIT(&_OmniVRF.CallOpts)
}

// MINCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0xb3d03582.
//
// Solidity: function MIN_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFCaller) MINCALLBACKGASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "MIN_CALLBACK_GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0xb3d03582.
//
// Solidity: function MIN_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFSession) MINCALLBACKGASLIMIT() (*big.Int, error) {
	return _OmniVRF.Contract.MINCALLBACKGASLIMIT(&_OmniVRF.CallOpts)
}

// MINCALLBACKGASLIMIT is a free data retrieval call binding the contract method 0xb3d03582.
//
// Solidity: function MIN_CALLBACK_GAS_LIMIT() view returns(uint256)
func (_OmniVRF *OmniVRFCallerSession) MINCALLBACKGASLIMIT() (*big.Int, error) {
	return _OmniVRF.Contract.MINCALLBACKGASLIMIT(&_OmniVRF.CallOpts)
}

// CalculateTaskFee is a free data retrieval call binding the contract method 0x09a9d71c.
//
// Solidity: function calculateTaskFee((address,uint32) operatorSet, bytes payload) view returns(uint96 fee)
func (_OmniVRF *OmniVRFCaller) CalculateTaskFee(opts *bind.CallOpts, operatorSet OperatorSet, payload []byte) (*big.Int, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "calculateTaskFee", operatorSet, payload)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateTaskFee is a free data retrieval call binding the contract method 0x09a9d71c.
//
// Solidity: function calculateTaskFee((address,uint32) operatorSet, bytes payload) view returns(uint96 fee)
func (_OmniVRF *OmniVRFSession) CalculateTaskFee(operatorSet OperatorSet, payload []byte) (*big.Int, error) {
	return _OmniVRF.Contract.CalculateTaskFee(&_OmniVRF.CallOpts, operatorSet, payload)
}

// CalculateTaskFee is a free data retrieval call binding the contract method 0x09a9d71c.
//
// Solidity: function calculateTaskFee((address,uint32) operatorSet, bytes payload) view returns(uint96 fee)
func (_OmniVRF *OmniVRFCallerSession) CalculateTaskFee(operatorSet OperatorSet, payload []byte) (*big.Int, error) {
	return _OmniVRF.Contract.CalculateTaskFee(&_OmniVRF.CallOpts, operatorSet, payload)
}

// CallbackGasDeposits is a free data retrieval call binding the contract method 0xe5860655.
//
// Solidity: function callbackGasDeposits(address ) view returns(uint256)
func (_OmniVRF *OmniVRFCaller) CallbackGasDeposits(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "callbackGasDeposits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CallbackGasDeposits is a free data retrieval call binding the contract method 0xe5860655.
//
// Solidity: function callbackGasDeposits(address ) view returns(uint256)
func (_OmniVRF *OmniVRFSession) CallbackGasDeposits(arg0 common.Address) (*big.Int, error) {
	return _OmniVRF.Contract.CallbackGasDeposits(&_OmniVRF.CallOpts, arg0)
}

// CallbackGasDeposits is a free data retrieval call binding the contract method 0xe5860655.
//
// Solidity: function callbackGasDeposits(address ) view returns(uint256)
func (_OmniVRF *OmniVRFCallerSession) CallbackGasDeposits(arg0 common.Address) (*big.Int, error) {
	return _OmniVRF.Contract.CallbackGasDeposits(&_OmniVRF.CallOpts, arg0)
}

// DecodeTaskData is a free data retrieval call binding the contract method 0x27b66026.
//
// Solidity: function decodeTaskData(bytes payload) pure returns((bytes32,uint256) taskData)
func (_OmniVRF *OmniVRFCaller) DecodeTaskData(opts *bind.CallOpts, payload []byte) (OmniVRFVRFTaskData, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "decodeTaskData", payload)

	if err != nil {
		return *new(OmniVRFVRFTaskData), err
	}

	out0 := *abi.ConvertType(out[0], new(OmniVRFVRFTaskData)).(*OmniVRFVRFTaskData)

	return out0, err

}

// DecodeTaskData is a free data retrieval call binding the contract method 0x27b66026.
//
// Solidity: function decodeTaskData(bytes payload) pure returns((bytes32,uint256) taskData)
func (_OmniVRF *OmniVRFSession) DecodeTaskData(payload []byte) (OmniVRFVRFTaskData, error) {
	return _OmniVRF.Contract.DecodeTaskData(&_OmniVRF.CallOpts, payload)
}

// DecodeTaskData is a free data retrieval call binding the contract method 0x27b66026.
//
// Solidity: function decodeTaskData(bytes payload) pure returns((bytes32,uint256) taskData)
func (_OmniVRF *OmniVRFCallerSession) DecodeTaskData(payload []byte) (OmniVRFVRFTaskData, error) {
	return _OmniVRF.Contract.DecodeTaskData(&_OmniVRF.CallOpts, payload)
}

// GetRandomness is a free data retrieval call binding the contract method 0x05b19402.
//
// Solidity: function getRandomness(bytes32 taskHash) view returns(bool fulfilled, uint256 randomness)
func (_OmniVRF *OmniVRFCaller) GetRandomness(opts *bind.CallOpts, taskHash [32]byte) (struct {
	Fulfilled  bool
	Randomness *big.Int
}, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "getRandomness", taskHash)

	outstruct := new(struct {
		Fulfilled  bool
		Randomness *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Randomness = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRandomness is a free data retrieval call binding the contract method 0x05b19402.
//
// Solidity: function getRandomness(bytes32 taskHash) view returns(bool fulfilled, uint256 randomness)
func (_OmniVRF *OmniVRFSession) GetRandomness(taskHash [32]byte) (struct {
	Fulfilled  bool
	Randomness *big.Int
}, error) {
	return _OmniVRF.Contract.GetRandomness(&_OmniVRF.CallOpts, taskHash)
}

// GetRandomness is a free data retrieval call binding the contract method 0x05b19402.
//
// Solidity: function getRandomness(bytes32 taskHash) view returns(bool fulfilled, uint256 randomness)
func (_OmniVRF *OmniVRFCallerSession) GetRandomness(taskHash [32]byte) (struct {
	Fulfilled  bool
	Randomness *big.Int
}, error) {
	return _OmniVRF.Contract.GetRandomness(&_OmniVRF.CallOpts, taskHash)
}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address requester, address callbackContract, uint256 callbackGasLimit, uint256 seed, uint256 blockNumber, bool fulfilled, uint256 result)
func (_OmniVRF *OmniVRFCaller) Requests(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Requester        common.Address
	CallbackContract common.Address
	CallbackGasLimit *big.Int
	Seed             *big.Int
	BlockNumber      *big.Int
	Fulfilled        bool
	Result           *big.Int
}, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "requests", arg0)

	outstruct := new(struct {
		Requester        common.Address
		CallbackContract common.Address
		CallbackGasLimit *big.Int
		Seed             *big.Int
		BlockNumber      *big.Int
		Fulfilled        bool
		Result           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Requester = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CallbackContract = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Seed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.BlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.Result = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address requester, address callbackContract, uint256 callbackGasLimit, uint256 seed, uint256 blockNumber, bool fulfilled, uint256 result)
func (_OmniVRF *OmniVRFSession) Requests(arg0 [32]byte) (struct {
	Requester        common.Address
	CallbackContract common.Address
	CallbackGasLimit *big.Int
	Seed             *big.Int
	BlockNumber      *big.Int
	Fulfilled        bool
	Result           *big.Int
}, error) {
	return _OmniVRF.Contract.Requests(&_OmniVRF.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address requester, address callbackContract, uint256 callbackGasLimit, uint256 seed, uint256 blockNumber, bool fulfilled, uint256 result)
func (_OmniVRF *OmniVRFCallerSession) Requests(arg0 [32]byte) (struct {
	Requester        common.Address
	CallbackContract common.Address
	CallbackGasLimit *big.Int
	Seed             *big.Int
	BlockNumber      *big.Int
	Fulfilled        bool
	Result           *big.Int
}, error) {
	return _OmniVRF.Contract.Requests(&_OmniVRF.CallOpts, arg0)
}

// TaskMailbox is a free data retrieval call binding the contract method 0xf42a9e13.
//
// Solidity: function taskMailbox() view returns(address)
func (_OmniVRF *OmniVRFCaller) TaskMailbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "taskMailbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TaskMailbox is a free data retrieval call binding the contract method 0xf42a9e13.
//
// Solidity: function taskMailbox() view returns(address)
func (_OmniVRF *OmniVRFSession) TaskMailbox() (common.Address, error) {
	return _OmniVRF.Contract.TaskMailbox(&_OmniVRF.CallOpts)
}

// TaskMailbox is a free data retrieval call binding the contract method 0xf42a9e13.
//
// Solidity: function taskMailbox() view returns(address)
func (_OmniVRF *OmniVRFCallerSession) TaskMailbox() (common.Address, error) {
	return _OmniVRF.Contract.TaskMailbox(&_OmniVRF.CallOpts)
}

// ValidatePreTaskCreation is a free data retrieval call binding the contract method 0xa7fbe930.
//
// Solidity: function validatePreTaskCreation(address caller, (address,uint96,(address,uint32),bytes) taskParams) view returns()
func (_OmniVRF *OmniVRFCaller) ValidatePreTaskCreation(opts *bind.CallOpts, caller common.Address, taskParams ITaskMailboxTypesTaskParams) error {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "validatePreTaskCreation", caller, taskParams)

	if err != nil {
		return err
	}

	return err

}

// ValidatePreTaskCreation is a free data retrieval call binding the contract method 0xa7fbe930.
//
// Solidity: function validatePreTaskCreation(address caller, (address,uint96,(address,uint32),bytes) taskParams) view returns()
func (_OmniVRF *OmniVRFSession) ValidatePreTaskCreation(caller common.Address, taskParams ITaskMailboxTypesTaskParams) error {
	return _OmniVRF.Contract.ValidatePreTaskCreation(&_OmniVRF.CallOpts, caller, taskParams)
}

// ValidatePreTaskCreation is a free data retrieval call binding the contract method 0xa7fbe930.
//
// Solidity: function validatePreTaskCreation(address caller, (address,uint96,(address,uint32),bytes) taskParams) view returns()
func (_OmniVRF *OmniVRFCallerSession) ValidatePreTaskCreation(caller common.Address, taskParams ITaskMailboxTypesTaskParams) error {
	return _OmniVRF.Contract.ValidatePreTaskCreation(&_OmniVRF.CallOpts, caller, taskParams)
}

// ValidatePreTaskResultSubmission is a free data retrieval call binding the contract method 0xba33565d.
//
// Solidity: function validatePreTaskResultSubmission(address caller, bytes32 taskHash, bytes cert, bytes result) view returns()
func (_OmniVRF *OmniVRFCaller) ValidatePreTaskResultSubmission(opts *bind.CallOpts, caller common.Address, taskHash [32]byte, cert []byte, result []byte) error {
	var out []interface{}
	err := _OmniVRF.contract.Call(opts, &out, "validatePreTaskResultSubmission", caller, taskHash, cert, result)

	if err != nil {
		return err
	}

	return err

}

// ValidatePreTaskResultSubmission is a free data retrieval call binding the contract method 0xba33565d.
//
// Solidity: function validatePreTaskResultSubmission(address caller, bytes32 taskHash, bytes cert, bytes result) view returns()
func (_OmniVRF *OmniVRFSession) ValidatePreTaskResultSubmission(caller common.Address, taskHash [32]byte, cert []byte, result []byte) error {
	return _OmniVRF.Contract.ValidatePreTaskResultSubmission(&_OmniVRF.CallOpts, caller, taskHash, cert, result)
}

// ValidatePreTaskResultSubmission is a free data retrieval call binding the contract method 0xba33565d.
//
// Solidity: function validatePreTaskResultSubmission(address caller, bytes32 taskHash, bytes cert, bytes result) view returns()
func (_OmniVRF *OmniVRFCallerSession) ValidatePreTaskResultSubmission(caller common.Address, taskHash [32]byte, cert []byte, result []byte) error {
	return _OmniVRF.Contract.ValidatePreTaskResultSubmission(&_OmniVRF.CallOpts, caller, taskHash, cert, result)
}

// HandlePostTaskCreation is a paid mutator transaction binding the contract method 0x09c5c450.
//
// Solidity: function handlePostTaskCreation(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFTransactor) HandlePostTaskCreation(opts *bind.TransactOpts, taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.contract.Transact(opts, "handlePostTaskCreation", taskHash)
}

// HandlePostTaskCreation is a paid mutator transaction binding the contract method 0x09c5c450.
//
// Solidity: function handlePostTaskCreation(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFSession) HandlePostTaskCreation(taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.Contract.HandlePostTaskCreation(&_OmniVRF.TransactOpts, taskHash)
}

// HandlePostTaskCreation is a paid mutator transaction binding the contract method 0x09c5c450.
//
// Solidity: function handlePostTaskCreation(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFTransactorSession) HandlePostTaskCreation(taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.Contract.HandlePostTaskCreation(&_OmniVRF.TransactOpts, taskHash)
}

// HandlePostTaskResultSubmission is a paid mutator transaction binding the contract method 0xdb6ecf67.
//
// Solidity: function handlePostTaskResultSubmission(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFTransactor) HandlePostTaskResultSubmission(opts *bind.TransactOpts, taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.contract.Transact(opts, "handlePostTaskResultSubmission", taskHash)
}

// HandlePostTaskResultSubmission is a paid mutator transaction binding the contract method 0xdb6ecf67.
//
// Solidity: function handlePostTaskResultSubmission(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFSession) HandlePostTaskResultSubmission(taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.Contract.HandlePostTaskResultSubmission(&_OmniVRF.TransactOpts, taskHash)
}

// HandlePostTaskResultSubmission is a paid mutator transaction binding the contract method 0xdb6ecf67.
//
// Solidity: function handlePostTaskResultSubmission(bytes32 taskHash) returns()
func (_OmniVRF *OmniVRFTransactorSession) HandlePostTaskResultSubmission(taskHash [32]byte) (*types.Transaction, error) {
	return _OmniVRF.Contract.HandlePostTaskResultSubmission(&_OmniVRF.TransactOpts, taskHash)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xe34c3d3a.
//
// Solidity: function requestRandomness(address callbackContract, uint256 callbackGasLimit) payable returns(bytes32 taskHash)
func (_OmniVRF *OmniVRFTransactor) RequestRandomness(opts *bind.TransactOpts, callbackContract common.Address, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _OmniVRF.contract.Transact(opts, "requestRandomness", callbackContract, callbackGasLimit)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xe34c3d3a.
//
// Solidity: function requestRandomness(address callbackContract, uint256 callbackGasLimit) payable returns(bytes32 taskHash)
func (_OmniVRF *OmniVRFSession) RequestRandomness(callbackContract common.Address, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _OmniVRF.Contract.RequestRandomness(&_OmniVRF.TransactOpts, callbackContract, callbackGasLimit)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xe34c3d3a.
//
// Solidity: function requestRandomness(address callbackContract, uint256 callbackGasLimit) payable returns(bytes32 taskHash)
func (_OmniVRF *OmniVRFTransactorSession) RequestRandomness(callbackContract common.Address, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _OmniVRF.Contract.RequestRandomness(&_OmniVRF.TransactOpts, callbackContract, callbackGasLimit)
}

// WithdrawCallbackGas is a paid mutator transaction binding the contract method 0x5ce2eec4.
//
// Solidity: function withdrawCallbackGas(uint256 amount) returns()
func (_OmniVRF *OmniVRFTransactor) WithdrawCallbackGas(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _OmniVRF.contract.Transact(opts, "withdrawCallbackGas", amount)
}

// WithdrawCallbackGas is a paid mutator transaction binding the contract method 0x5ce2eec4.
//
// Solidity: function withdrawCallbackGas(uint256 amount) returns()
func (_OmniVRF *OmniVRFSession) WithdrawCallbackGas(amount *big.Int) (*types.Transaction, error) {
	return _OmniVRF.Contract.WithdrawCallbackGas(&_OmniVRF.TransactOpts, amount)
}

// WithdrawCallbackGas is a paid mutator transaction binding the contract method 0x5ce2eec4.
//
// Solidity: function withdrawCallbackGas(uint256 amount) returns()
func (_OmniVRF *OmniVRFTransactorSession) WithdrawCallbackGas(amount *big.Int) (*types.Transaction, error) {
	return _OmniVRF.Contract.WithdrawCallbackGas(&_OmniVRF.TransactOpts, amount)
}

// OmniVRFCallbackFailedIterator is returned from FilterCallbackFailed and is used to iterate over the raw logs and unpacked data for CallbackFailed events raised by the OmniVRF contract.
type OmniVRFCallbackFailedIterator struct {
	Event *OmniVRFCallbackFailed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OmniVRFCallbackFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OmniVRFCallbackFailed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OmniVRFCallbackFailed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OmniVRFCallbackFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OmniVRFCallbackFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OmniVRFCallbackFailed represents a CallbackFailed event raised by the OmniVRF contract.
type OmniVRFCallbackFailed struct {
	TaskHash         [32]byte
	CallbackContract common.Address
	Reason           string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCallbackFailed is a free log retrieval operation binding the contract event 0x58a5347ef2b6c6fcc342cf2326c32d04d662829bb952a924df81526f75ec629b.
//
// Solidity: event CallbackFailed(bytes32 indexed taskHash, address indexed callbackContract, string reason)
func (_OmniVRF *OmniVRFFilterer) FilterCallbackFailed(opts *bind.FilterOpts, taskHash [][32]byte, callbackContract []common.Address) (*OmniVRFCallbackFailedIterator, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}
	var callbackContractRule []interface{}
	for _, callbackContractItem := range callbackContract {
		callbackContractRule = append(callbackContractRule, callbackContractItem)
	}

	logs, sub, err := _OmniVRF.contract.FilterLogs(opts, "CallbackFailed", taskHashRule, callbackContractRule)
	if err != nil {
		return nil, err
	}
	return &OmniVRFCallbackFailedIterator{contract: _OmniVRF.contract, event: "CallbackFailed", logs: logs, sub: sub}, nil
}

// WatchCallbackFailed is a free log subscription operation binding the contract event 0x58a5347ef2b6c6fcc342cf2326c32d04d662829bb952a924df81526f75ec629b.
//
// Solidity: event CallbackFailed(bytes32 indexed taskHash, address indexed callbackContract, string reason)
func (_OmniVRF *OmniVRFFilterer) WatchCallbackFailed(opts *bind.WatchOpts, sink chan<- *OmniVRFCallbackFailed, taskHash [][32]byte, callbackContract []common.Address) (event.Subscription, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}
	var callbackContractRule []interface{}
	for _, callbackContractItem := range callbackContract {
		callbackContractRule = append(callbackContractRule, callbackContractItem)
	}

	logs, sub, err := _OmniVRF.contract.WatchLogs(opts, "CallbackFailed", taskHashRule, callbackContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OmniVRFCallbackFailed)
				if err := _OmniVRF.contract.UnpackLog(event, "CallbackFailed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCallbackFailed is a log parse operation binding the contract event 0x58a5347ef2b6c6fcc342cf2326c32d04d662829bb952a924df81526f75ec629b.
//
// Solidity: event CallbackFailed(bytes32 indexed taskHash, address indexed callbackContract, string reason)
func (_OmniVRF *OmniVRFFilterer) ParseCallbackFailed(log types.Log) (*OmniVRFCallbackFailed, error) {
	event := new(OmniVRFCallbackFailed)
	if err := _OmniVRF.contract.UnpackLog(event, "CallbackFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OmniVRFRandomnessFulfilledIterator is returned from FilterRandomnessFulfilled and is used to iterate over the raw logs and unpacked data for RandomnessFulfilled events raised by the OmniVRF contract.
type OmniVRFRandomnessFulfilledIterator struct {
	Event *OmniVRFRandomnessFulfilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OmniVRFRandomnessFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OmniVRFRandomnessFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OmniVRFRandomnessFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OmniVRFRandomnessFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OmniVRFRandomnessFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OmniVRFRandomnessFulfilled represents a RandomnessFulfilled event raised by the OmniVRF contract.
type OmniVRFRandomnessFulfilled struct {
	TaskHash   [32]byte
	Randomness *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRandomnessFulfilled is a free log retrieval operation binding the contract event 0x9b0aa3f92f46e24caa76b000bdf0dd495b9b390c320cf6585ae10a12b7d09edb.
//
// Solidity: event RandomnessFulfilled(bytes32 indexed taskHash, uint256 randomness)
func (_OmniVRF *OmniVRFFilterer) FilterRandomnessFulfilled(opts *bind.FilterOpts, taskHash [][32]byte) (*OmniVRFRandomnessFulfilledIterator, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}

	logs, sub, err := _OmniVRF.contract.FilterLogs(opts, "RandomnessFulfilled", taskHashRule)
	if err != nil {
		return nil, err
	}
	return &OmniVRFRandomnessFulfilledIterator{contract: _OmniVRF.contract, event: "RandomnessFulfilled", logs: logs, sub: sub}, nil
}

// WatchRandomnessFulfilled is a free log subscription operation binding the contract event 0x9b0aa3f92f46e24caa76b000bdf0dd495b9b390c320cf6585ae10a12b7d09edb.
//
// Solidity: event RandomnessFulfilled(bytes32 indexed taskHash, uint256 randomness)
func (_OmniVRF *OmniVRFFilterer) WatchRandomnessFulfilled(opts *bind.WatchOpts, sink chan<- *OmniVRFRandomnessFulfilled, taskHash [][32]byte) (event.Subscription, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}

	logs, sub, err := _OmniVRF.contract.WatchLogs(opts, "RandomnessFulfilled", taskHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OmniVRFRandomnessFulfilled)
				if err := _OmniVRF.contract.UnpackLog(event, "RandomnessFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRandomnessFulfilled is a log parse operation binding the contract event 0x9b0aa3f92f46e24caa76b000bdf0dd495b9b390c320cf6585ae10a12b7d09edb.
//
// Solidity: event RandomnessFulfilled(bytes32 indexed taskHash, uint256 randomness)
func (_OmniVRF *OmniVRFFilterer) ParseRandomnessFulfilled(log types.Log) (*OmniVRFRandomnessFulfilled, error) {
	event := new(OmniVRFRandomnessFulfilled)
	if err := _OmniVRF.contract.UnpackLog(event, "RandomnessFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OmniVRFRandomnessRequestedIterator is returned from FilterRandomnessRequested and is used to iterate over the raw logs and unpacked data for RandomnessRequested events raised by the OmniVRF contract.
type OmniVRFRandomnessRequestedIterator struct {
	Event *OmniVRFRandomnessRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OmniVRFRandomnessRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OmniVRFRandomnessRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OmniVRFRandomnessRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OmniVRFRandomnessRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OmniVRFRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OmniVRFRandomnessRequested represents a RandomnessRequested event raised by the OmniVRF contract.
type OmniVRFRandomnessRequested struct {
	TaskHash         [32]byte
	Requester        common.Address
	CallbackContract common.Address
	Seed             *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRandomnessRequested is a free log retrieval operation binding the contract event 0xee5eec2b3210f6ea36e12fd442b35d8918a76da7ff35c15dc4b8d69480752aa0.
//
// Solidity: event RandomnessRequested(bytes32 indexed taskHash, address indexed requester, address callbackContract, uint256 seed)
func (_OmniVRF *OmniVRFFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, taskHash [][32]byte, requester []common.Address) (*OmniVRFRandomnessRequestedIterator, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OmniVRF.contract.FilterLogs(opts, "RandomnessRequested", taskHashRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &OmniVRFRandomnessRequestedIterator{contract: _OmniVRF.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequested is a free log subscription operation binding the contract event 0xee5eec2b3210f6ea36e12fd442b35d8918a76da7ff35c15dc4b8d69480752aa0.
//
// Solidity: event RandomnessRequested(bytes32 indexed taskHash, address indexed requester, address callbackContract, uint256 seed)
func (_OmniVRF *OmniVRFFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *OmniVRFRandomnessRequested, taskHash [][32]byte, requester []common.Address) (event.Subscription, error) {

	var taskHashRule []interface{}
	for _, taskHashItem := range taskHash {
		taskHashRule = append(taskHashRule, taskHashItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OmniVRF.contract.WatchLogs(opts, "RandomnessRequested", taskHashRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OmniVRFRandomnessRequested)
				if err := _OmniVRF.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRandomnessRequested is a log parse operation binding the contract event 0xee5eec2b3210f6ea36e12fd442b35d8918a76da7ff35c15dc4b8d69480752aa0.
//
// Solidity: event RandomnessRequested(bytes32 indexed taskHash, address indexed requester, address callbackContract, uint256 seed)
func (_OmniVRF *OmniVRFFilterer) ParseRandomnessRequested(log types.Log) (*OmniVRFRandomnessRequested, error) {
	event := new(OmniVRFRandomnessRequested)
	if err := _OmniVRF.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
