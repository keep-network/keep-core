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
    def __init__(self, unique_id, model, tickets):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "node"
        self.num_tickets = int(tickets)
        self.ticket_list = []
        self.connection_status = "not connected" #change later to event - currently used for node failure process
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"
        self.connection_delay = np.random.randint(0,100) #uniform randomly assigned connection delay step value
        self.mainloop_fork_delay = np.random.randint(0,100) #uniform randomly assigned connection delay step value
        self.timer = self.model.timer
    
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

        #simulate node failure
        

    def advance(self):
        if self.model.relay_request:
            self.generate_tickets()

    def generate_tickets(self):
        #generates tickets using the uniform distribution
        self.stake_status = "staked"
        self.ticket_list.append(np.random.random_sample(self.num_tickets))
        
    def node_disconnect(self):
        # disconnect the node from the network; causes it to go through the entire reconnection sequence in the next step
        self.connection_status = "not connected"
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"


class Group(Agent):
    """ A Group """
    def __init__(self, unique_id, model, members, expiry, sign_threshold):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "group"
        self.members = members
        self.last_signature = "none"
        self.status = "Active"
        self.expiry = expiry # of steps before expiration
        self.timer = self.model.timer
        self.model.newest_group_id +=1

    def step(self):
        """ At each step check if the group as expired """
        self.expiry -=1
        #print('group ID '+ str(self.id) + ' expiry ' + str(self.expiry))
        if self.expiry == 0: 
            self.status = "Expired"

    def advance(self):
        pass

        

