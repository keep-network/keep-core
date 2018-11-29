pragma solidity ^0.5.00;

contract StakeDelegation {

    event Revoke(
        address indexed _delegator,
        address indexed _operator
    );

    // map address to a stake
    mapping(address => uint256) stakeOf;
    
    // maps (delegators) address to operators address
    mapping(address => address) operatorOf;

    // maps (operators) address to delegators address
    mapping(address => address) delegatorOf;

    modifier notNull(address _address)
    {
        require(_address != address(0), "Provided address is a zero (0) address");
        _;
    }

    // Check if we don’t have a delegator address mapped for provided address.
    modifier isNotOperator(address _address)
    {
        //notNull(_address)
        require(_address != address(0), "Provided address is a zero (0) address");
        
        require(delegatorOf[_address] == address(0), "Provided address is an operator");
        _;
    }

    // Check if we don’t have an operator address mapped for provided address. 
    modifier isNotDelegator(address _address)
    {
        //notNull(_address)
        require(_address != address(0), "Provided address is a zero (0) address");
        
        require(operatorOf[_address] == address(0), "Provided address is a delegator");
        _;
    }

    /* Check if we have a delegator address mapped for provided address and
     * if that address is not delegator.
     */
    modifier isOperator(address _address)
    {
        //notNull(_address)
        require(_address != address(0), "Provided address is a zero (0) address");
        
        //isNotDelegator(_address)
        require(operatorOf[_address] == address(0), "Provided address is a delegator");
    
        require(delegatorOf[_address] != address(0), "Provided address is not an operator");
        _;
    }

    /* Check if we have an operator address mapped for provided address and 
     * if that address is not an operator. 
     */
    modifier isDelegator(address _address)
    {
        //notNull(_address)
        require(_address != address(0), "Provided address is a zero (0) address");

        //isNotOperator(_address)
        require(delegatorOf[_address] == address(0), "Provided address is an operator");

        require(operatorOf[_address] != address(0), "Provided address is not a delegator");
        _;
    }


    // Address can be a delegator if it’s not null, and if it’s not an operator
    modifier canBeDelegator(address _address)
    {
        //isNotDelegator(_address)
        require(operatorOf[_address] == address(0), "Provided address is a delegator");

        //isNotOperator(_address)
        require(delegatorOf[_address] == address(0), "Provided address is an operator");

        require(stakeOf[_address] > 0, "Provided address can’t be a delegator");
        _;
    }

    /* Address can be an operator if it’s not null and if it’s not a delegator.
     * Also it needs to have 0 (or less) stake.
     */
    modifier canBeOperator(address _address)
    {
        //isNotDelegator(_address)
        require(operatorOf[_address] == address(0), "Provided address is a delegator");

        //isNotOperator(_address)
        require(delegatorOf[_address] == address(0), "Provided address is an operator");

        require(stakeOf[_address] <= 0, "Provided address can’t be an operator");
        _;
    }

    // Return stake of address only if it’s delegators address.
    function getStakeOfDelegator(address _address)
        public 
        view
        isDelegator(_address)
        returns(uint256)
    {
        return(stakeOf[_address]);
    }

    // Set stake to delegate by delegators address.
    function setStakeToDelegate(address _address, uint256 stake)
        public 
        payable // ???
        canBeDelegator(msg.sender)
        canBeOperator(_address)
        returns(bool)
    {
        if (stakeOf[msg.sender] <= 0) {
            stakeOf[msg.sender] = stake;
            //msg.sender.stake -= stake;        // ???
            return true;
        } else {
            return false;
        }
    }

    /* Returns delegators address of provided address after verifying that the
     * provided address is an operator address.
     */
    function getDelegatorOf(address _address)
        public
        view
        isOperator(_address)
        returns(address)
    {
        return(delegatorOf[_address]);
    }

    /* Returns operators address of provided address after verifying that the
     * provided address is a delegators address.
     */
    function getOperatorOf(address _address)
        public
        view
        isDelegator(_address)
        returns(address)
    {
        return(operatorOf[_address]);
    }

    function setDelegatorOf(address _address)
        public
        canBeOperator(msg.sender)
        canBeDelegator(_address)
    {
        delegatorOf[msg.sender] = _address;
    }


    function setOperatorOf(address _address)
        public
        canBeDelegator(msg.sender)
        canBeOperator(_address)
    {
        operatorOf[msg.sender] = _address;
    }

    /* Check if we have a bound between current address and respective operator
     * or delegator
     */
    function isBoundedTo(address _address)
        public
        view
        notNull(_address)
        returns(bool)
    {
        // if msg.sender is a delegator then the _address must be its operator
        if (operatorOf[msg.sender] == _address) {
            return(true);
        }
        // if msg.sender is an operator then the _address must be its delegator
        if (delegatorOf[msg.sender] == _address) { // msg.sender is operator
            return(true);
        }
        return(false);
    }

    function revokeDelgationByDelegator()
        public
        isOperator(msg.sender)
        isDelegator(delegatorOf[msg.sender])
    {
        address delegator = delegatorOf[msg.sender];

        emit Revoke(delegator, operatorOf[delegator]); // ???

        delete delegatorOf[msg.sender];
        delete operatorOf[delegator];
    }


    function revokeDelgationByOperator()
        public
        isDelegator(msg.sender)
        isOperator(delegatorOf[msg.sender])
    {
        address operator = operatorOf[msg.sender];

        emit Revoke(delegatorOf[operator], operator); // ???

        delete operatorOf[msg.sender];
        delete delegatorOf[operator];
    }
}
