from mesa import Agent, Model
from mesa.time import RandomActivation

class Node(Agent):
    """ Node: One hardware device used to stake tokens on the network. 
    Each node will create virtual stakers proportional to
    the number of tokens owned by the node 
    
    Attributes:
    unique_id: unique int
    token_amount: int value of tokens staked by node
    node_status: status of node can be - not connected, connected
    """

    def __init__(self, unique_id, model, token_amt):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.node_status = "not connected" #change later to event - currently used for node failure process
    
    def step(self):
        pass

    #Connecting to a chain
    def connect_to_chain(self):
        pass
    
    def forking_mainLoop(self, env):
        pass
    

class Virtial_Stakers(Agent):
    """ Virtual Stakers aka Tickets """
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)

    def step():
        pass

    # wait for relay request
    def Watching_RelayRequest(self, env):
        pass

    #join group
    def join_group(self, env, group_object):
        pass

class Group(Agent):
    """ A Group """
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)

    def step():
        pass
    
    def sign():
        pass

    def expire():
        pass

class Beacon_Model(Model):
    """The model"""
    def __init__(self, N):
        pass