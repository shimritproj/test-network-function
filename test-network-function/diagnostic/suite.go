// Copyright (C) 2020 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package diagnostic

import (
	"encoding/json"
	"path"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	expect "github.com/ryandgoulding/goexpect"
)

const (
	// testSuiteSpec contains the name of the Ginkgo test specification.
	testSuiteSpec = "diagnostic"
	// defaultTimeoutSeconds contains the default timeout in secons.
	defaultTimeoutSeconds = 20
)

var (
	// defaultTestTimeout is the timeout for the test.
	defaultTestTimeout = time.Duration(defaultTimeoutSeconds) * time.Second

	// nodeSummary stores the raw JSON output of `oc get nodes -o json`
	nodeSummary = make(map[string]interface{})

	// nodesTestPath is the file location of the nodes.json test case relative to the project root.
	nodesTestPath = path.Join("pkg", "tnf", "handlers", "node", "nodes.json")

	// pathRelativeToRoot is used to calculate relative filepaths for the `test-network-function` executable entrypoint.
	pathRelativeToRoot = path.Join("..")

	// relativeNodesTestPath is the relative path to the nodes.json test case.
	relativeNodesTestPath = path.Join(pathRelativeToRoot, nodesTestPath)

	// relativeSchemaPath is the relative path to the generic-test.schema.json JSON schema.
	relativeSchemaPath = path.Join(pathRelativeToRoot, schemaPath)

	// schemaPath is the path to the generic-test.schema.json JSON schema relative to the project root.
	schemaPath = path.Join("schemas", "generic-test.schema.json")
)

// createShell sets up a local shell expect.Expecter, checking errors along the way.
func createShell() *interactive.Context {
	context, err := interactive.SpawnShell(interactive.CreateGoExpectSpawner(), defaultTestTimeout, expect.Verbose(true))
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(context).ToNot(gomega.BeNil())
	return context
}

var _ = ginkgo.Describe(testSuiteSpec, func() {
	ginkgo.When("a cluster is set up and installed with OpenShift", func() {
		ginkgo.It("should report all available nodeSummary", func() {
			context := createShell()

			test, handlers, jsonParseResult, err := generic.NewGenericFromJSONFile(relativeNodesTestPath, relativeSchemaPath)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(jsonParseResult).ToNot(gomega.BeNil())
			gomega.Expect(jsonParseResult.Valid()).To(gomega.BeTrue())
			gomega.Expect(handlers).ToNot(gomega.BeNil())
			gomega.Expect(test).ToNot(gomega.BeNil())

			tester, err := tnf.NewTest(context.GetExpecter(), *test, handlers, context.GetErrorChannel())
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(tester).ToNot(gomega.BeNil())

			result, err := tester.Run()
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result).To(gomega.Equal(tnf.SUCCESS))

			genericTest := (*test).(*generic.Generic)
			gomega.Expect(genericTest).ToNot(gomega.BeNil())
			matches := genericTest.Matches
			gomega.Expect(len(matches)).To(gomega.Equal(1))
			match := genericTest.GetMatches()[0]
			err = json.Unmarshal([]byte(match.Match), &nodeSummary)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})

// GetNodeSummary returns the result of running `oc get nodes -o json`.
func GetNodeSummary() map[string]interface{} {
	return nodeSummary
}
