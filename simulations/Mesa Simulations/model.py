from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, token_distribution):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = []

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, token_distribution[i])
            self.schedule.add(node)
        
    def step(self):
        '''Advance the model by one step'''
        #generate relay requests
        self.relay_request.append(bool(np.random.randint(0,1))
        self.schedule.step()
