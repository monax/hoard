# Regression test

The regression test in this directory is designed to test the byte-wise stability of Hoard over successive releases.

A regression in this context means that a newer version of Hoard cannot faithfully retrieve data stored by an older version of Hoard.

Sometimes a breaking change may have been intentionally introduced, in this case this test can be used to check that a Hoard upgrade and migration script combined produce no regression.

## Test cycle
The test is implemented as a simple go program that takes an input directory containing some plaintexts, a filesystem store, and some saved grants and cycles the plaintexts through Hoard by:

- Saving them to Hoard and saving the returned grant back to a 'grants' directory
- Trying to retrieve them using a previously saved grant (or using the new grant if one does not exist)
- Saving them back into the plaintexts directory

The program uses a non-random nonce for LINK refs and so the results are always deterministic. See [main.go](./main.go) for  more details.

## Breaking changes
There are three possible levels of byte-wise compatibility that may be observed between subsequent runs of this regression test:

1. Plaintexts are retrieved and saved to disk with the same contents
2. Hoard binary blobs are unchanged within the 'store' directory
3. Hoard grant files are unchanged within the 'grants' directory (never true for SymmetricGrants that use an AES nonce)

3 implies 2 implies 1, but the converse is not necessarily true.

The key criteria for a non-breaking change to Hoard is that plaintexts are preserved (1) across this cycle so that data can be faithfully retrieved. In some circumstances we will be interested in maintaining store-level (2) compatibility so that the same data is stored to the same location sharing storage. Finally grant-invariance (3) is useful if we wish to deduplicate entries at the grant-level (i.e. use the fact that two grants may be equal IFF their plaintexts are equal - note: this would only apply using the same salt/link nonce)

## Comparing runs
The comparison between runs of the regression test are intended to be made using git tooling and for the results of a subsequent run to be committed to the repository (before a version is released) in order to be used for comparison across the next significant version delta.

## Upgrades
In the case where a change to Hoard is intentionally breaking then this test should be extended to support a cycle of:

- Run upgrade script against store/grants
- Run regression cycle

And so check that the migration-upgrade combination has no regressions.
