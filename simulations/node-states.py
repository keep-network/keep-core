import random
import simpy
import datetime
import numpy as np

SIM_CYCLES = 10
TOTAL_TOKEN_AMT = 100
# levers to pull: 
# stake amount and effect on ownership on group; given set group size, 
# failure rates of a single node (explodes!) when it comes bck its subject to connect (use lit)
# what if somebody is attacking the network, cost to bring down system, how much money do they need to get signingrights
class Node:
    #Node states based on Antonio's diagramming on Feb 15 2019
    #Assume staking mechanism is complete
    def __init__(self, env, identity, start_time):
        self.env = env
        self.id = identity
        self.starttime = start_time
        self.current_state = "not connected"
        self.process = env.process(self.Connect_Node(env))
        self.relay_request_time = 0
        self.relay_entry_watch_time = 0
        self.ingroup = False
        self.inrelay = False
        self.number_of_entries_generated = 0
        self.number_of_groups_joined = 0
        self.STAKING_AMT = np.random.lognormal(3,1,) #find total tokens from contract
        self.cycle_count = 0


    #Connecting to Ethereum
    def Connect_Node(self, env):
        while True:
            ethereum_conection_time = np.random.lognormal(3,1,)
            if ethereum_conection_time>=30:
                print (str(self.id) + " ethereum connection Failure" + " cycle="+str(self.cycle_count))
                self.current_state = " not connected"
                self.Connect_Node(env)
            else:
                print (str(self.id) + " ethereum connection success" + " cycle="+str(self.cycle_count))
                env.process(self.Forking_MainLoop(env))               
                self.current_state = " connected"
                yield env.exit()
    
    def Forking_MainLoop(self,env):
        print(str(self.id) + " Forking Main Loop" + " cycle="+str(self.cycle_count))
        env.process(self.Watching_RelayRequest(env))
        env.process(self.Watching_RelayEntry(env))  
        yield env.exit()
    
    # wait for relay request
    def Watching_RelayRequest(self, env):
        print(str(self.id)+" Watching Relay Request" + " cycle="+str(self.cycle_count))
        self.relay_request_time = np.random.normal(3,1,)
        env.process(self.Group_Selection(env))
        yield env.exit()
    
    #watching for relay entry
    def Watching_RelayEntry(self, env):
        print(str(self.id)+" Watching Relay Entry" + " cycle="+str(self.cycle_count))
        self.relay_entry_watch_time = np.random.normal(3,1,)
        yield env.exit()
    
    # Group Selection
    def Group_Selection(self, env):
        while True:
            if np.random.randint(10)<5:
                env.process(self.Group_Formation(env))
                yield env.exit()
                
            else:
                print(str(self.id)+" group formation failure" + " cycle="+str(self.cycle_count))
        
    # check if this node is a member of a signing group (assuming this is another process)
    def Group_Member_Check(self, env):
            if self.ingroup == True:
                env.process(self.Entry_Generation(env))
                print (str(self.id)+" in a group" + " cycle="+str(self.cycle_count))
                yield env.exit()
            else:
                print(str(self.id)+" not a group member" + " cycle="+str(self.cycle_count))
                env.process(self.Watching_RelayEntry(env))
                yield env.exit()
        
    # Generate Entry
    def Entry_Generation(self,env):
        print(str(self.id)+" generated entry" + " cycle="+str(self.cycle_count))
        self.number_of_entries_generated += 1
        self.ingroup = False
        self.cycle_count +=1
        if self. cycle_count > 10:
            yield env.exit()
        else:
            env.process(self.Forking_MainLoop(env))
        
    # Form Group
    def Group_Formation(self,env):
        print(str(self.id)+" formed group" + " cycle="+str(self.cycle_count))
        self.ingroup = True
        self.number_of_groups_joined +=1
        env.process(self.Group_Member_Check(env)) #doing it here instead of waiting for relay entry
        env.exit()

def node_failure_generator():
    failure = np.random.lognormal(1,0)

# Setup and start the simulation
print('Node States')

# Create an environment and start the setup process
env = simpy.Environment()
print("creating nodes")
nodes = [Node(env, 'Node %d' % i, datetime.datetime.now())
            for i in range(3)] #number of nodes
env.run()
print("xxxxxxxxxxxxxxxxxxxx")
print(" final node states ")
for n in nodes:
    print(str(n.id) + ", # of Entries = " 
    + str(n.number_of_entries_generated) 
    + ", # Groups Joined = " 
    + str(n.number_of_groups_joined)
    + ", Total relay request time = "
    + str(n.relay_request_time)
    + ", Total relay watch time = "
    + str(n.relay_entry_watch_time))
