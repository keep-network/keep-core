from mesa import Model
from mesa.time import SimultaneousActivation
import agent

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, token_distribution):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, token_distribution[i])
            self.schedule.add(node)

        #generate relay requests
        
    def step(self):
        '''Advance the model by one step'''
        self.schedule.step()