pragma solidity ^0.5.4;

library SetUtils {
    struct UintSet {
        // members[positions[item] - 1] = item
        uint256[] members;
        mapping(uint256 => uint256) positions;
    }

    function contains(UintSet storage self, uint256 item)
        internal view returns (bool) {
        return (self.positions[item] != 0);
    }

    // Insert item to set
    // If already present, do nothing
    function insert(UintSet storage self, uint256 item) internal {
        if (!contains(self, item)) {
            self.members.push(item);
            self.positions[item] = self.members.length;
        }
    }

    // Remove item from set
    // If not already present, do nothing
    function remove(UintSet storage self, uint256 item) internal {
        uint256 positionPlusOne = self.positions[item];
        if (positionPlusOne != 0) {
            uint256 memberCount = self.members.length;
            if (positionPlusOne != memberCount) {
                // Not the last member,
                // so we need to move the last member into the emptied position.
                uint256 lastMember = self.members[memberCount - 1];
                self.members[positionPlusOne - 1] = lastMember;
                self.positions[lastMember] = positionPlusOne;
            }
            self.members.length--;
            self.positions[item] = 0;
        }
    }

    function enumerate(UintSet storage self)
        internal view returns (uint256[] memory) {
        return self.members;
    }
}
