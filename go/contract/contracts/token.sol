pragma solidity ^0.4.18;

// From: https://gist.github.com/karalabe/08f4b780e01c8452d989
// From: https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts

contract tokenRecipient {
	function receiveApproval (address _from, uint256 _value, address _token, bytes _extraData) public;
}

contract Token { 
    /* Public variables of the token */
    string public name;
    string public symbol;
    uint8 public decimals;

    /* This creates an array with all balances */
    mapping (address => uint256) public balanceOf;
    mapping (address => mapping (address => uint)) public allowance;
    mapping (address => mapping (address => uint)) public spentAllowance;

    /* This generates a public event on the blockchain that will notify clients */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /* Initializes contract with initial supply tokens to the creator of the contract */
    function Token(uint256 initialSupply, string tokenName, uint8 decimalUnits, string tokenSymbol) public {
        balanceOf[msg.sender] = initialSupply;              // Give the creator all initial tokens                    
        name = tokenName;                                   // Set the name for display purposes     
        symbol = tokenSymbol;                               // Set the symbol for display purposes    
        decimals = decimalUnits;                            // Amount of decimals for display purposes        
    }

    /* Send coins */
    function transfer(address _to, uint256 _value) public {
        if (balanceOf[msg.sender] < _value) revert();           // Check if the sender has enough   
        if (balanceOf[_to] + _value < balanceOf[_to]) revert(); // Check for overflows
        balanceOf[msg.sender] -= _value;                     // Subtract from the sender
        balanceOf[_to] += _value;                            // Add the same to the recipient            
        Transfer(msg.sender, _to, _value);                   // Notify anyone listening that this transfer took place
    }

    /* Allow another contract to spend some tokens in your behalf */

    function approveAndCall(address _spender, uint256 _value, bytes _extraData) public returns (bool success) {
        allowance[msg.sender][_spender] = _value;     
        tokenRecipient spender = tokenRecipient(_spender);
        spender.receiveApproval(msg.sender, _value, this, _extraData);  
		success = true;
    }

    /* A contract attempts to get the coins */

    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        if (balanceOf[_from] < _value) revert();                 // Check if the sender has enough   
        if (balanceOf[_to] + _value < balanceOf[_to]) revert();  // Check for overflows
        if (spentAllowance[_from][msg.sender] + _value > allowance[_from][msg.sender]) revert();   // Check allowance
        balanceOf[_from] -= _value;                          // Subtract from the sender
        balanceOf[_to] += _value;                            // Add the same to the recipient            
        spentAllowance[_from][msg.sender] += _value;
        Transfer(msg.sender, _to, _value); 
		success = true;
    } 

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () public {
        revert();     // Prevents accidental sending of ether
    }        
} 
