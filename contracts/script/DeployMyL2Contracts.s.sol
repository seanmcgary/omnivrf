// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.27;

import {Script, console} from "forge-std/Script.sol";
import {stdJson} from "forge-std/StdJson.sol";

import {IBN254CertificateVerifier} from
    "@eigenlayer-contracts/src/contracts/interfaces/IBN254CertificateVerifier.sol";
import {IECDSACertificateVerifier} from "@eigenlayer-contracts/src/contracts/interfaces/IECDSACertificateVerifier.sol";
import {ITaskMailbox, ITaskMailboxTypes} from "@hourglass-monorepo/src/interfaces/core/ITaskMailbox.sol";
import {IKeyRegistrarTypes} from "@eigenlayer-contracts/src/contracts/interfaces/IKeyRegistrar.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IAVSTaskHook} from "@hourglass-monorepo/src/interfaces/avs/l2/IAVSTaskHook.sol";

import {OmniVRF} from "@project/l2-contracts/OmniVRF.sol";
import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";

contract DeployMyL2Contracts is Script {
    using stdJson for string;

    struct Context {
        address avs;
        uint256 avsPrivateKey;
        uint256 deployerPrivateKey;
        IBN254CertificateVerifier certificateVerifier;
        IECDSACertificateVerifier ecdsaCertificateVerifier;
        ITaskMailbox taskMailbox;
    }

    struct Output {
        string name;
        address contractAddress;
    }

    function run(string memory environment, string memory _context) public {
        // Read the context
        Context memory context = _readContext(environment, _context);

        vm.startBroadcast(context.deployerPrivateKey);
        console.log("Deployer address:", vm.addr(context.deployerPrivateKey));

        // Deploy OmniVRF
        OmniVRF omniVRF = new OmniVRF(address(context.taskMailbox));
        console.log("OmniVRF deployed to:", address(omniVRF));

        vm.stopBroadcast();

        vm.startBroadcast(context.avsPrivateKey);
        console.log("AVS address:", context.avs);

        // Set up operator set
        OperatorSet memory operatorSet = OperatorSet({
            avs: context.avs,
            id: 1
        });

        // Configure task mailbox with OmniVRF as the task hook
        ITaskMailboxTypes.ExecutorOperatorSetTaskConfig memory taskConfig = ITaskMailboxTypes
            .ExecutorOperatorSetTaskConfig({
            curveType: IKeyRegistrarTypes.CurveType.BN254, // Using BN254 curve for VRF
            taskHook: IAVSTaskHook(address(omniVRF)),
            feeToken: IERC20(address(0)), // No fee token for MVP
            feeCollector: address(0), // No fee collector for MVP
            taskSLA: 3600, // 1 hour SLA for MVP
            stakeProportionThreshold: 10_000, // 100% threshold
            taskMetadata: bytes("") // No metadata for MVP
        });
        
        console.log("Registering OmniVRF as task hook...");
        context.taskMailbox.setExecutorOperatorSetTaskConfig(operatorSet, taskConfig);
        console.log("OmniVRF registered as task hook for operator set");

        vm.stopBroadcast();

        // Write OmniVRF contracts to output file
        Output[] memory outputs = new Output[](1);
        outputs[0] = Output({name: "OmniVRF", contractAddress: address(omniVRF)});
        _writeOutputToJson(environment, outputs);
    }

    function _readContext(
        string memory environment,
        string memory _context
    ) internal view returns (Context memory) {
        // Parse the context
        Context memory context;
        context.avs = stdJson.readAddress(_context, ".context.avs.address");
        context.avsPrivateKey = uint256(stdJson.readBytes32(_context, ".context.avs.avs_private_key"));
        context.deployerPrivateKey = uint256(stdJson.readBytes32(_context, ".context.deployer_private_key"));
        context.certificateVerifier = IBN254CertificateVerifier(stdJson.readAddress(_context, ".context.eigenlayer.l2.bn254_certificate_verifier"));
        context.ecdsaCertificateVerifier = IECDSACertificateVerifier(stdJson.readAddress(_context, ".context.eigenlayer.l2.ecdsa_certificate_verifier"));
        context.taskMailbox = ITaskMailbox(_readHourglassConfigAddress(environment, "taskMailbox"));

        return context;
    }

    function _readHourglassConfigAddress(
        string memory environment,
        string memory key
    ) internal view returns (address) {
        // Load the Hourglass config file
        string memory hourglassConfigFile =
                            string.concat("script/", environment, "/output/deploy_hourglass_core_output.json");
        string memory hourglassConfig = vm.readFile(hourglassConfigFile);

        // Parse and return the address
        return stdJson.readAddress(hourglassConfig, string.concat(".addresses.", key));
    }

    function _readAVSL2ConfigAddress(string memory environment, string memory key) internal view returns (address) {
        // Load the AVS L2 config file
        string memory avsL2ConfigFile = string.concat("script/", environment, "/output/deploy_avs_l2_output.json");
        string memory avsL2Config = vm.readFile(avsL2ConfigFile);

        // Parse and return the address
        return stdJson.readAddress(avsL2Config, string.concat(".addresses.", key));
    }

    function _writeOutputToJson(
        string memory environment,
        Output[] memory outputs
    ) internal {
        uint256 length = outputs.length;

        if (length > 0) {
            // Add the addresses object
            string memory addresses = "addresses";

            for (uint256 i = 0; i < outputs.length - 1; i++) {
                vm.serializeAddress(addresses, outputs[i].name, outputs[i].contractAddress);
            }
            addresses = vm.serializeAddress(addresses, outputs[length - 1].name, outputs[length - 1].contractAddress);

            // Add the chainInfo object
            string memory chainInfo = "chainInfo";
            chainInfo = vm.serializeUint(chainInfo, "chainId", block.chainid);

            // Finalize the JSON
            string memory finalJson = "final";
            vm.serializeString(finalJson, "addresses", addresses);
            finalJson = vm.serializeString(finalJson, "chainInfo", chainInfo);

            // Write to output file
            string memory outputFile = string.concat("script/", environment, "/output/deploy_custom_contracts_l2_output.json");
            vm.writeJson(finalJson, outputFile);
        }
    }
}
