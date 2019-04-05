from mesa import Agent, Model
from mesa.time import RandomActivation
import numpy as np

class Node(Agent):
    """ Node: One hardware device used to stake tokens on the network. 
    Each node will create virtual stakers proportional to
    the number of tokens owned by the node 
    
    Attributes:
    unique_id: unique int
    token_amount: int value of tokens staked by node
    node_status: status of node can be - not connected, connected
    """
    def __init__(self, unique_id, model, tickets, relay_request):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.num_tickets = 
        self.ticket_list = []
        self.connection_status = "not connected" #change later to event - currently used for node failure process
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"
        self.connection_delay = np.random.randint(0,100) #uniform randomly assigned connection delay step value
        self.mainloop_fork_delay = np.random.randint(0,100) #uniform randomly assigned connection delay step value
        self.relay = relay_request #relay request object used to trigger the generation of tickets
    
    def step(self):
        #connect to chain
        if self.connection_delay>0:
            self.connection_delay -=1
        elif self.connection_status == "not connected":
            self.connection_status = "connected"
        else:
        #once connected fork the main loop
            if self.mainloop_fork_delay>0:
                self.mainloop_fork_delay -=1
            elif self.mainloop_fork_delay == "not forked":
                self.mainloop_status = "forked"
            elif self.stake_status == "not staked":
                self.generate_tickets()

        if self.stake_status = "staked" and self.relay.triggered



    def advance(self):
        pass

    def generate_tickets(self):
        #add code here to create virtual staker distribution
        self.stake_status = "staked"
        self.ticket_list.append(np.random.random_sample(self.num_tickets))
        
    def node_disconnect(self):
        # disconnect the node from the network; causes it to go through the entire reconnection sequence in the next step
        self.connection_status = "not connected"
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"


class Group(Agent):
    """ A Group """
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)
        self.active_members = []
        self.last_signature = "none"

    def step(self):
        """ At each step check how many members are still active """
        pass

    def register_members(self, virtual_staker_list):
        """detect which members with winning tickets are active and add them to the list"""
        pass
        
    def sign(self, sign_threshold):
        """Check if enough members are available to perform a signature using the sign_threshold, and then perform a signature"""
        if len(self.active_members)<sign_threshold:
            self.last_signature = "failed"
        else:
            self.last_signature = "success"

    def expire(self):
        pass

class Relay_Trigger(Agent):
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)
        self.trigger = False

    def step():
        self.trigger = bool(np.random.randint(0,1)

