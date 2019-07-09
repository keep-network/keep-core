from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np
from mesa.datacollection import DataCollector
import logging as log
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, ticket_distribution, active_group_threshold, 
    group_size, max_malicious_threshold, group_expiry, 
    node_failure_percent, node_death_percent,
    signature_delay, min_nodes, node_connection_delay, node_mainloop_connection_delay, 
    log_filename, run_number, misbehaving_nodes, dkg_block_delay, compromised_threshold,
    failed_signature_threshold):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = False
        self.active_groups = []
        self.num_active_groups = 0
        self.active_nodes = []
        self.num_active_nodes = 0
        self.inactive_nodes = []
        self.active_group_threshold = active_group_threshold # number of groups that will always be maintained in an active state
        self.max_malicious_threshold = max_malicious_threshold # threshold above which a signature is deemed to be compromised, typically 51%
        self.group_size = group_size
        self.ticket_distribution = ticket_distribution
        self.newest_id = 0
        self.newest_group_id = 0
        self.newest_signature_id = 0
        self.group_expiry = group_expiry
        self.bootstrap_complete = False # indicates when the initial active group list bootstrap is complete
        self.group_formation_threshold = min_nodes # min nodes required to form a group
        self.timer = 0
        self.unsuccessful_signature_events = []
        self.signature_delay = signature_delay
        self.dkg_block_delay = dkg_block_delay
        self.compromised_threshold = compromised_threshold
        self.median_malicious_group_percents = 0
        self.median_dominated_signatures_percents = 0
        self.perc_dominated_signatures = 0
        self.perc_compromised_groups = 0
        self.total_signatures = 0
        self.failed_signature_threshold = failed_signature_threshold
        self.perc_failed_signatures = 0
        self.datacollector = DataCollector(
            model_reporters = {"# of Active Groups":"num_active_groups",
             "# of Active Nodes":"num_active_nodes",
             "# of Signatures":"total_signatures",
             "Median Malicious Group %": "median_malicious_group_percents",
             "% Compromised Groups": "perc_compromised_groups",
             "Median Dominator %":"median_dominated_signatures_percents",
             "% Dominated signatures":"perc_dominated_signatures",
             "Failed Singature %" : "perc_failed_signatures" },
            agent_reporters={"Type_ID": lambda x : x.node_id if x.type == "node" else ( x.group_id if x.type == "group" else x.signature_id) , 
            "Type" : "type",
            "Node Status (Connection_Mainloop_Stake)": lambda x : str(x.connection_status + x.mainloop_status + x.stake_status) if x.type == "node" else None,
            "Status": lambda x: x.status if x.type == "group" or x.type == "signature" else None, 
            "Malicious": lambda x : x.malicious if x.type == "node" else None,
            "DKG Block Delay" : lambda x : x.dkg_block_delay if x.type =="group" else None,
            "Ownership Distribution" : lambda x : x.ownership_distr if x.type =="group" or x.type == "signature" else None,
            "Malicious %" : lambda x : x.malicious_percent if x.type == "group" else None,
            "Offline %" : lambda x : x.offline_percent if x.type == "group" or x.type == "signature" else None,
            "Dominator %": lambda x : x.dominator_percent if x.type == "signature" else None})


        #create log file
        log.basicConfig(filename=log_filename + str(run_number), filemode='w', format='%(name)s - %(levelname)s - %(message)s')

        print("creating nodes")
        #create nodes
        for i in range(nodes):
            node = agent.Node(i, i, self, 
            self.ticket_distribution[i], 
            node_failure_percent, 
            node_death_percent, 
            node_connection_delay, node_mainloop_connection_delay, misbehaving_nodes)
            self.newest_id = i
            self.schedule.add(node)
        self.newest_id +=1


    def step(self):
        '''Advance the model by one step'''
 
        log.debug("Number of nodes in the forked state = " + str(len(self.active_nodes)))

        #bootstrap active groups as nodes become available. Can only happen once enough nodes are online
        temp_bootstrap_groups = []
        if self.bootstrap_complete == False:
            log.debug("bootstrapping active groups")
            if len(self.active_nodes)>=self.group_formation_threshold:
                for i in range(self.active_group_threshold):
                    temp_bootstrap_groups.append(self.group_registration())
                self.bootstrap_complete = True
            self.active_groups = temp_bootstrap_groups
        
        #generate relay requests
        self.relay_request = np.random.choice([True,False]) # make this variable so it can be what-if'd
        log.debug("relay request recieved? = "+ str(self.relay_request))

        if self.relay_request:
            try:
                log.debug('     selecting group at random')
                try:
                    # pick an active group from the active group list and create a signature object
                    signature = agent.Signature(self.newest_id, self.newest_signature_id, self, self.active_groups[np.random.randint(len(self.active_groups))]) 
                except Exception as e: log.warning(e)
                self.schedule.add(signature)
            except:
                log.debug('     no active groups available')

            log.debug('     registering new group')
            self.group_registration()
        else:
            log.debug("     No relay request")
        self.timer += 1

        #calculate model measurements
        self.calculate_compromised_groups()
        self.calculate_dominated_signatures()
        
        #advance the agents
        self.schedule.step()
        self.num_active_nodes = len(self.active_nodes)
        self.num_active_groups = len(self.active_groups)
        self.datacollector.collect(self)

    def group_registration(self):
        ticket_list = []
        group_members = []

        if len(self.active_nodes)<self.group_formation_threshold: 
            log.debug("             Not enough nodes to register a group")

        else:
            # make each active node generate tickets and save them to a list
            max_tickets = int(max(self.ticket_distribution))
            for node in self.active_nodes:
                adjusted_ticket_list = []
                node.generate_tickets()
                adjusted_ticket_list = np.concatenate([node.ticket_list,np.ones(int(max_tickets)-len(node.ticket_list))])  #adds 2's the ends of the list so that the 2D array can have equal length rows
                ticket_list.append(adjusted_ticket_list)

            #iteratively add group members by lowest value
            while len(group_members) < self.group_size:

                min_index = np.where(ticket_list == np.min(ticket_list)) # find the index of the minimum value in the array
                for i,index in enumerate(min_index[0]): #if there are repeated values, iterate through and add the indexes to the group
                    group_members.append(self.active_nodes[index])
                    ticket_list[index][min_index[1][i]] = 2 # Set the value of the ticket to a high value so it doesn't get counted again
            
            #create a group agent which can track expiry, sign, etc
            group_object = agent.Group(self.newest_id, self.newest_group_id, self, group_members, self.group_expiry)


            #add group to schedule
            self.schedule.add(group_object)

            #add group to active group list
            self.active_groups.append(group_object)
            
            return group_object

    def refresh_active_group_list(self):
        temp_list = []

        for group in self.schedule.agents:
            if group.type == "group":
                if group.status == "active":
                    temp_list.append(group)
        self.active_groups = temp_list

    def refresh_connected_nodes_list(self):
        log.debug("refreshing active nodes list")
        temp_active_node_list = []
        temp_inactive_node_list = []
        temp_active_nodes = [] #take out later - id of active nodes
        for agent in self.schedule.agents:
            if agent.type == "node":
                if agent.mainloop_status == "forked": 
                    temp_active_node_list.append(agent) #adds the node to the active list only if it is in the forked state
                    temp_active_nodes.append(agent.node_id)
                
                else:
                    temp_inactive_node_list.append(agent.node_id)
        self.active_nodes = temp_active_node_list
        self.inactive_nodes = temp_inactive_node_list

    def calculate_compromised_groups(self):
    #Calculate compromised groups
        malicious_array = []
        compromised_count = 0
        total_groups = 0
        for group in self.schedule.agents:
            if group.type == "group":
                total_groups +=1
                malicious_array.append(group.malicious_percent) #creates an array of malicious percents for each group
                if group.compromised_percent >= self.compromised_threshold:
                    compromised_count +=1

        self.median_malicious_group_percents = np.median(malicious_array)
        self.perc_compromised_groups = compromised_count/(total_groups+0.000000000000000001)

    def calculate_dominated_signatures(self):
        dominator_array = []
        dominator_count = 0
        total_signatures = 0
        failed_signatures = 0
        for signature in self.schedule.agents:
            if signature.type == "signature":
                total_signatures +=1
                dominator_array.append(signature.dominator_percent)
                dominator_count += (signature.dominator_id>=0)
                failed_signatures += (signature.offline_percent>=self.failed_signature_threshold)

        self.perc_failed_signatures = failed_signatures/(total_signatures+0.00000000000000001)
        self.median_dominated_signatures_percents = np.median(dominator_array)
        self.perc_dominated_signatures = dominator_count/(total_signatures+0.00000000000000001)
        self.total_signatures = total_signatures


def create_cdf(nodes,ticket_distr):
# Create CDF's - used to determine max ownership ticket index
    cdf = np.zeros(nodes)
    for node,ticketmax in enumerate(ticket_distr):
        
        cdf[node]=sum(ticket_distr[0:node+1])
    return cdf


            
            





    










        

        







