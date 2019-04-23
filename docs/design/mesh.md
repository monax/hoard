# Service Mesh

Currently, Hoard can only be configured with one back-end store. This is particularly problematic when we want to distinguish public versus private stores in which we want to ensure that sensitive data is never stored on an externally accessible system. For example, we may want to start referencing shared documents on IPFS and limit exposure to a bucket holding our private data. This can be achieved currently by running multiple independent processes, but it is not very scalable - as such we may want to consider adopting a service mesh based model.

## Data Availability / Redundancy

In our single store model, data availability is contingent on two factors, the resiliency of our back-end and the hoard service itself. Orthogonally scaling the hoard daemon mitigates the latter issue but store availability may not be easily rectified. With a cloud provider, we could ask that the data is replicated to some backup, but this is not always an option.

## Routing

Without adding too much complexity, we want to extend the distribution of data.

### Multiple Stores

The obvious choice would be to extend Hoard with the ability to use multiple back-end stores - our configuration would simply list a number of named back-end stores which are explicitly targeted in the address field. Much like in the independent process model, we would need to be careful about which store is accessible in certain situations.

### Forwarding

Another option would be to forward based on the given address, so each instance would still effectively guard one configured store. If we prepend one or more host names to the address, then it would be a simple matter of GRPC forwarding amongst daemons. This is highly contingent on a better access control scheme however, otherwise it would be particularly easy for an adversary to read from an undesired location. If we were to define formal verbiage in the specifications for both stores and references we could simply prevent forwarding to services with higher privileges. That way, if a service is exposed to the internet on the same network as a private store, we could simply terminate access.