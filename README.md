Introduction
============

This repository contains the Kubernetes E2E test framework set up in
such a way that it runs the Kubernetes storage tests.

Usage
=====

No cloud-specific code gets imported, so tests can only be run against
a cluster that has already been set up. Set the the usual
`KUBECONFIG=<config file>` and then run `go test -v ./test/e2e` or
`ginkgo ./test/e2e`.

Note that all `test/e2e/storage` tests are imported. To run only the
CSI tests, use `go test -v ./test/e2e -args -ginkgo.focus=CSI.Volumes`.

Adding `-provider=local` suppresses a message about the flag not being
set and treating the run as "conformance test", but has no other
effect in practice.

The
[upstream documentation](https://github.com/kubernetes/community/blob/master/contributors/devel/e2e-tests.md#local-clusters)
uses `hack/e2e.go` as wrapper around the test execution. This is not
necessary for the test suite defined in this repository.

Adding Tests
============

New tests can be written in their own packages under `test/e2e` and
then need to be added to the import list in `test/e2e_test.go`.

