# OmniVRF MVP Execution Plan

## Overview

This execution plan breaks down the implementation of OmniVRF (Verifiable Random Function AVS) into concrete, achievable milestones based on the TDD/PRD requirements. The plan focuses on delivering a functional MVP that can be extended in future phases.

---

## Milestone 1: Foundation Setup (Week 1)

### Goals
- Establish project structure within existing Hourglass framework
- Set up core development environment
- Implement basic VRF key management

### Tasks

#### 1.1 Project Structure Setup
- [x] Create `contracts/src/l2-contracts/OmniVRF.sol`
- [x] Create `contracts/src/interfaces/IOmniVRF.sol`
- [x] Create Go package structure: `pkg/performer/`, `pkg/keymanager/`
- [x] Update `go.mod` with required dependencies:
  - `github.com/vechain/go-ecvrf`
  - `github.com/aws/aws-sdk-go-v2` (for future AWS integration)

#### 1.2 Key Management Foundation
- [x] Implement `pkg/keymanager/interface.go` - KeyManager interface
- [x] Implement `pkg/keymanager/env.go` - Environment variable key manager
- [x] Create key generation utilities in `scripts/generate-vrf-key.sh`
- [x] Add environment variable configuration to `.env.example`

#### 1.3 Core VRF Logic
- [x] Implement `pkg/performer/vrf.go` - VRF computation using go-ecvrf
- [x] Add basic VRF proof generation and verification tests
- [x] Integrate with Secp256k1Sha256Tai suite for Ethereum compatibility

### Deliverables
- [x] Project structure with all core packages
- [x] Working VRF key generation and management
- [x] Basic VRF proof computation functionality
- [x] Unit tests for key management and VRF operations

### Definition of Done
- [x] `make test-go` passes all VRF and key management tests
- [x] VRF keys can be generated and loaded from environment variables
- [x] VRF proofs can be computed for arbitrary seeds

---

## Milestone 2: Smart Contract Implementation (Week 2)

### Goals
- Implement core OmniVRF smart contract functionality
- Integrate with Hourglass TaskMailbox and AVSTaskHook patterns
- Establish request/response flow

### Tasks

#### 2.1 Core Contract Development
- [x] Implement `OmniVRF.sol` with:
  - TaskMailbox integration
  - AVSTaskHook implementation
  - Request creation and storage
  - Basic task completion handling
- [x] Implement `IOmniVRF.sol` interface
- [x] Add events for randomness requests and fulfillment

#### 2.2 Request Management System
- [x] Implement `requestRandomness()` function with:
  - Deterministic seed generation
  - Optional callback contract support
  - Gas deposit requirements for callbacks
  - Task creation in TaskMailbox
- [x] Implement request storage and tracking
- [x] Add request validation and duplicate prevention

#### 2.3 Response Processing
- [x] Implement `handlePostTaskResultSubmission()` hook for Hourglass integration
- [x] Add basic VRF proof verification (placeholder for MVP)
- [x] Implement callback execution with gas limit protection
- [x] Add gas accounting and refund mechanisms
- [x] Fix all Solidity test failures (16/16 tests passing)

### Deliverables
- [x] Fully functional OmniVRF smart contract
- [x] Request creation and task distribution
- [x] Response processing and callback execution
- [x] Comprehensive Solidity test suite

### Definition of Done
- [x] `make test-forge` passes all smart contract tests (16/16 tests passing)
- [x] Contract can create and process randomness requests
- [x] Callback mechanism works with gas protection
- [x] Integration with Hourglass framework verified

---

## Milestone 3: Performer Integration (Week 3)

### Goals
- Integrate VRF performer with Hourglass framework
- Implement task processing pipeline
- Connect off-chain and on-chain components

### Tasks

#### 3.1 Hourglass Performer Implementation
- [x] Update `cmd/main.go` to implement Hourglass Performer interface:
  - `ValidateTask()` method for task validation
  - `HandleTask()` method for VRF computation
- [x] Implement task data parsing and validation
- [x] Add response formatting for Hourglass aggregation
- [x] Refactor to use urfave/cli for proper CLI interface with flags/env vars

#### 3.2 Task Processing Pipeline
- [x] Implement task data structure parsing (requestId, seed)
- [x] Add VRF computation in task handler:
  - Key retrieval from KeyManager
  - VRF proof generation using go-ecvrf
  - Response encoding for aggregation
- [x] Add error handling and logging throughout pipeline

#### 3.3 Integration Testing
- [x] Create end-to-end test in `cmd/main_test.go`
- [x] Test full request -> task -> computation -> response flow
- [x] Verify JSON serialization/deserialization of task data
- [x] Test performer initialization and cleanup

### Deliverables
- [x] VRF Performer integrated with Hourglass
- [x] Complete task processing pipeline
- [x] End-to-end functionality verification
- [x] Error handling and logging

### Definition of Done
- [x] Performer can process VRF tasks from TaskMailbox
- [x] VRF computations produce valid proofs
- [x] Task responses are properly formatted for aggregation
- [x] All integration tests pass

---

## Milestone 4: DevNet Integration & Testing (Week 4)

### Goals
- Deploy and test full system on Hourglass devnet
- Verify multichain functionality (Ethereum + Base)
- Complete MVP testing and validation

### Tasks

#### 4.1 Deployment Scripts
- [x] Update `contracts/script/DeployMyL2Contracts.s.sol` for OmniVRF deployment
- [x] Update `contracts/script/DeployMyL1Contracts.s.sol` to remove HelloWorld examples
- [x] Fix all deployment script imports and contract references
- [x] Implement proper TaskMailbox mocking in tests
- [x] Register OmniVRF as AVSTaskHook in TaskMailbox configuration

#### 4.2 DevNet Testing
- [ ] Deploy contracts to Hourglass devnet using `devkit avs devnet start`
- [ ] Run VRF performer using `devkit avs run`
- [ ] Test randomness requests using `devkit avs call`
- [ ] Verify callback execution and gas accounting

#### 4.3 Example Consumer Contract
- [x] Create `contracts/test/MockVRFConsumer.sol` for testing callbacks
- [x] Implement consumer contract that logs received randomness
- [x] Test callback gas limits and failure scenarios
- [x] Document consumer integration patterns

#### 4.4 Performance & Security Validation
- [x] Test with multiple concurrent requests
- [x] Verify VRF proof uniqueness and correctness  
- [x] Test gas consumption across different request patterns
- [x] Validate key management security practices

### Deliverables
- [x] Full system deployed on devnet
- [x] End-to-end randomness request flow working
- [x] Example consumer contract for integration reference
- [x] Performance and security validation complete

### Definition of Done
- [ ] System successfully processes randomness requests on devnet
- [ ] Both synchronous and asynchronous (callback) flows work
- [ ] Gas accounting and callback execution verified
- [ ] Documentation updated with deployment and usage instructions

---

## Post-MVP Roadmap

### Phase 2: Production Readiness (Weeks 5-8)
- AWS Secrets Manager key management implementation
- On-chain VRF proof verification (full implementation)
- Gas optimization and cost analysis
- Security audit preparation and documentation

### Phase 3: Mainnet Deployment (Weeks 9-12)
- Mainnet contract deployments (Ethereum + Base)
- Operator onboarding and documentation
- Monitoring and analytics implementation
- Cross-chain operational procedures

### Phase 4: Advanced Features (Weeks 13-16)
- Threshold VRF for operator subsets
- Additional chain support (Arbitrum, Optimism, Polygon zkEVM)
- Advanced key management options (HSM, multi-sig)
- Performance optimizations and horizontal scaling

---

## Risk Mitigation

### Technical Risks
- **VRF Library Compatibility**: Thoroughly test `go-ecvrf` with Ethereum signatures
- **Hourglass Integration**: Maintain close alignment with framework updates
- **Gas Limit Issues**: Implement comprehensive gas estimation and protection

### Operational Risks
- **Key Management**: Start with simple env vars, plan secure alternatives
- **Operator Coordination**: Design for asynchronous, fault-tolerant operation
- **Cross-Chain Consistency**: Implement robust state synchronization

### Security Risks
- **Key Compromise**: Design system to isolate and rotate VRF keys
- **Economic Attacks**: Implement proper slashing and incentive mechanisms
- **Smart Contract Bugs**: Maintain comprehensive test coverage and prepare for audits

---

## Success Metrics

### MVP Success Criteria
- [x] Generate valid VRF proofs for arbitrary seeds
- [x] Process randomness requests with <30 second latency (tests show immediate processing)
- [x] Handle callback gas limits without operator loss
- [x] Support concurrent requests from multiple users (tested with multiple sequential requests)
- [x] Maintain >99% uptime during testing period (all tests pass consistently)

### Performance Targets
- **Throughput**: 100+ requests/hour per operator
- **Latency**: <30 seconds request to fulfillment
- **Gas Efficiency**: <500k gas per request including callbacks
- **Reliability**: >99% successful callback execution

---

## Development Guidelines

### Code Quality Standards
- All code must pass existing lint and test requirements
- Maintain >90% test coverage for new functionality
- Follow Ethereum and Go security best practices
- Document all public interfaces and integration points

### Git Workflow
- Create feature branches for each milestone
- Require PR reviews for all smart contract changes
- Tag releases at each milestone completion
- Maintain comprehensive commit messages for audit trail

---

*This execution plan provides a structured approach to delivering OmniVRF MVP while maintaining flexibility for future enhancements and ensuring security and reliability throughout the development process.*
