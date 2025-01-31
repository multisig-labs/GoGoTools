// SPDX-License-Identifier: Ecosystem

pragma solidity 0.8.25;

contract Viewer {
  enum ValidatorStatus {
    Unknown,
    PendingAdded,
    Active,
    PendingRemoved,
    Completed,
    Invalidated
  }

  struct Validator {
    ValidatorStatus status;
    bytes nodeID;
    uint64 startingWeight;
    uint64 messageNonce;
    uint64 weight;
    uint64 startedAt;
    uint64 endedAt;
  }

  struct ValidatorChurnPeriod {
    uint256 startedAt;
    uint64 initialWeight;
    uint64 totalWeight;
    uint64 churnAmount;
  }
  struct PoSValidatorInfo {
    address owner;
    uint16 delegationFeeBips;
    uint64 minStakeDuration;
    uint64 uptimeSeconds;
  }
  struct Delegator {
    DelegatorStatus status;
    address owner;
    bytes32 validationID;
    uint64 weight;
    uint64 startedAt;
    uint64 startingNonce;
    uint64 endingNonce;
  }
  enum DelegatorStatus {
    Unknown,
    PendingAdded,
    Active,
    PendingRemoved
  }

  struct ValidatorManagerStorage {
    /// @notice The subnetID associated with this validator manager.
    bytes32 _subnetID;
    /// @notice The number of seconds after which to reset the churn tracker.
    uint64 _churnPeriodSeconds;
    /// @notice The maximum churn rate allowed per churn period.
    uint8 _maximumChurnPercentage;
    /// @notice The churn tracker used to track the amount of stake added or removed in the churn period.
    ValidatorChurnPeriod _churnTracker;
    /// @notice Maps the validationID to the registration message such that the message can be re-sent if needed.
    mapping(bytes32 => bytes) _pendingRegisterValidationMessages;
    /// @notice Maps the validationID to the validator information.
    mapping(bytes32 => Validator) _validationPeriods;
    /// @notice Maps the nodeID to the validationID for validation periods that have not ended.
    mapping(bytes => bytes32) _registeredValidators;
    /// @notice Boolean that indicates if the initial validator set has been set.
    bool _initializedValidatorSet;
  }
  // solhint-enable private-vars-leading-underscore

  // keccak256(abi.encode(uint256(keccak256("avalanche-icm.storage.ValidatorManager")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 public constant VALIDATOR_MANAGER_STORAGE_LOCATION =
    0xe92546d698950ddd38910d2e15ed1d923cd0a7b3dde9e2a6a3f380565559cb00;

  uint8 public constant MAXIMUM_CHURN_PERCENTAGE_LIMIT = 20;
  uint64 public constant MAXIMUM_REGISTRATION_EXPIRY_LENGTH = 2 days;
  uint32 public constant ADDRESS_LENGTH = 20; // This is only used as a packed uint32
  uint8 public constant BLS_PUBLIC_KEY_LENGTH = 48;
  bytes32 public constant P_CHAIN_BLOCKCHAIN_ID = bytes32(0);

  struct PoSValidatorManagerStorage {
    /// @notice The minimum amount of stake required to be a validator.
    uint256 _minimumStakeAmount;
    /// @notice The maximum amount of stake allowed to be a validator.
    uint256 _maximumStakeAmount;
    /// @notice The minimum amount of time in seconds a validator must be staked for. Must be at least {_churnPeriodSeconds}.
    uint64 _minimumStakeDuration;
    /// @notice The minimum delegation fee percentage, in basis points, required to delegate to a validator.
    uint16 _minimumDelegationFeeBips;
    /**
     * @notice A multiplier applied to validator's initial stake amount to determine
     * the maximum amount of stake a validator can have with delegations.
     * Note: Setting this value to 1 would disable delegations to validators, since
     * the maximum stake would be equal to the initial stake.
     */
    uint64 _maximumStakeMultiplier;
    /// @notice The factor used to convert between weight and value.
    uint256 _weightToValueFactor;
    /// @notice The reward calculator for this validator manager.
    address _rewardCalculator;
    /// @notice The ID of the blockchain that submits uptime proofs. This must be a blockchain validated by the subnetID that this contract manages.
    bytes32 _uptimeBlockchainID;
    /// @notice Maps the validation ID to its requirements.
    mapping(bytes32 validationID => PoSValidatorInfo) _posValidatorInfo;
    /// @notice Maps the delegation ID to the delegator information.
    mapping(bytes32 delegationID => Delegator) _delegatorStakes;
    /// @notice Maps the delegation ID to its pending staking rewards.
    mapping(bytes32 delegationID => uint256) _redeemableDelegatorRewards;
    mapping(bytes32 delegationID => address) _delegatorRewardRecipients;
    /// @notice Maps the validation ID to its pending staking rewards.
    mapping(bytes32 validationID => uint256) _redeemableValidatorRewards;
    mapping(bytes32 validationID => address) _rewardRecipients;
  }

  bytes32 public constant POS_VALIDATOR_MANAGER_STORAGE_LOCATION =
    0x4317713f7ecbdddd4bc99e95d903adedaa883b2e7c2551610bd13e2c7e473d00;

  uint8 public constant MAXIMUM_STAKE_MULTIPLIER_LIMIT = 10;

  uint16 public constant MAXIMUM_DELEGATION_FEE_BIPS = 10000;

  uint16 public constant BIPS_CONVERSION_FACTOR = 10000;

  function _getValidatorManagerStorage()
    internal
    pure
    returns (ValidatorManagerStorage storage $)
  {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      $.slot := VALIDATOR_MANAGER_STORAGE_LOCATION
    }
  }

  function _getPoSValidatorManagerStorage()
    internal
    pure
    returns (PoSValidatorManagerStorage storage $)
  {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      $.slot := POS_VALIDATOR_MANAGER_STORAGE_LOCATION
    }
  }

  struct ContractSettings {
    bytes32 _subnetID;
    uint64 _churnPeriodSeconds;
    uint8 _maximumChurnPercentage;
    uint256 _minimumStakeAmount;
    uint256 _maximumStakeAmount;
    uint64 _minimumStakeDuration;
    uint16 _minimumDelegationFeeBips;
    uint64 _maximumStakeMultiplier;
    uint256 _weightToValueFactor;
    address _rewardCalculator;
    bytes32 _uptimeBlockchainID;
  }
  function getSettings() public view returns (ContractSettings memory) {
    ValidatorManagerStorage storage $ = _getValidatorManagerStorage();
    PoSValidatorManagerStorage storage $$ = _getPoSValidatorManagerStorage();
    return
      ContractSettings({
        _subnetID: $._subnetID,
        _churnPeriodSeconds: $._churnPeriodSeconds,
        _maximumChurnPercentage: $._maximumChurnPercentage,
        _minimumStakeAmount: $$._minimumStakeAmount,
        _maximumStakeAmount: $$._maximumStakeAmount,
        _minimumStakeDuration: $$._minimumStakeDuration,
        _minimumDelegationFeeBips: $$._minimumDelegationFeeBips,
        _maximumStakeMultiplier: $$._maximumStakeMultiplier,
        _weightToValueFactor: $$._weightToValueFactor,
        _rewardCalculator: $$._rewardCalculator,
        _uptimeBlockchainID: $$._uptimeBlockchainID
      });
  }
}
