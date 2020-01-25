### list all nodes and rungroups
dsh -q

### run <command> on all nodes in the <rungroup>
dsh -g <rungroup> <command>

### same as before but for a list of hosts instead of a run group
dsh -w <hostname>[,<hostname>,<hostname>,...] <command>

### run command on all nodes in a rungroup but only do it <#> nodes at a time (keeps services from restarting all at the same time
dsh -f <#> -g <rungroup> <command>

### run command as root against a group
dsh -g <rungroup> -l root <command>

### run against the role-rtb-log
dsh -f 3 -g role-rtb-log sudo puppet agent -t --noop

### run against a specific cluster
CLUSTER=/etc/clusterit.diplomat.apro-us-west-2 dsh -g role-cluster-controller sudo puppet agent -t
