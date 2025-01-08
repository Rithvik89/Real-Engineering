C: Consistency
A: Availability
P: Partition Tolerence

Some time back in early 2000's (Eric Brewer in 2000) people started to see that all the above three principle will not be working for a distributed system at the same time.

What does these actually mean:

consistency:

* A system is said to be in consistent state , if the state of data is same (consistent) on all the nodes in the system. meaning all the clients should receive the same data if they are coming
with a same request at a given point of time.

Availability:

* At any given point of time, irrespective of the health of the nodes in the cluster, it should give a response back to the client.

Partition tolerance:

* Partition tolerance is the ability of a system to keep working even when there are communication breakdowns between nodes in a distributed system


Only 2 of the above 3 principles can be acheived at a time.

Lets talk about few :

CP:

Generally speaking, a CP system means that when 'network partition' happened, the system will try to ensure consistency (or 'linearizability' from some blogs) instead of making the system 'available'.

Network partition essentially means that the cluster is partitioned to two or more parts and they can't communicate with each other due to network failures. During the time of this network issue, if a write request was made to one side of the cluster (a region server), the system will not accept the write until the issue is resolved. Thus this essentially make the system not really 'available' when network partition happened. In HBase, network partition can result in region in transition and those regions affected will not be able to accept read/write. This is the reason why some post may consider HBase as a CP system. For a system favoring availability, read will not be blocked when system is in a partitioned state.

AP: 

Cassandra achieves Partition Tolerance (P) and Availability (A) in the CAP theorem through:

Multi-master architecture: Any node can accept writes, ensuring availability even during network partitions
Eventual consistency model:
Writes are accepted locally and asynchronously propagated
Conflicts resolved using timestamps and last-write-wins strategy
Sacrifices immediate consistency (C) for availability
Gossip protocol:
Nodes periodically exchange state information
Maintains cluster state without central coordination
Continues operating during network partitions