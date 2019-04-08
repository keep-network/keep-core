from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, token_distribution):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = False
        self.active_groups = []

        #create ticket distribution
        tickets = np.ones(nodes)*10 #this can be any distribution we decide

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, tickets[i])
            self.schedule.add(node)

        #bootstrap active groups
        
    def step(self):
        '''Advance the model by one step'''
        #generate relay requests
        self.relay_request = bool(np.random.randint(0,1))
        print(self.active_groups)

        #advance the agents
        self.schedule.step()


