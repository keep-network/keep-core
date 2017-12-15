package membership

// WatiForGroup waits for a group to be available for joining, then attempts to
// join it.
func WaitForGroup() {}

// WaitForGroupCompletion waits for this node's current group to be complete,
// then starts group initialization.
func WaitForGroupCompletion() {}

// InitializeMembership performs membership initialization actions for the
// current group. This includes generating a private key share and listening for
// other nodes' key proofs, as well as preparing to synthesize the public key
// and publish it to the chain.
func InitializeMembership() {}

// ActivateMembership performs membership activation actions; in particular, it
// waits for the activation delay to elapse and then puts the node into an
// state where it is actively listening for new beacon requests.
func ActivateMembership() {}
