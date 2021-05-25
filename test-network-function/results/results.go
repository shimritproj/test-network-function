package results

import (
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/test-network-function/test-network-function-claim/pkg/claim"
	"github.com/test-network-function/test-network-function/pkg/junit"
)

var results = map[claim.Identifier][]claim.Result{}

// RecordResult is a hook provided to save aspects of the ginkgo.GinkgoTestDescription for a given claim.Identifier.
// Multiple results for a given identifier are aggregated as an array under the same key.
func RecordResult(identifier claim.Identifier) {
	testContext := ginkgo.CurrentGinkgoTestDescription()
	results[identifier] = append(results[identifier], claim.Result{
		Duration:      int(testContext.Duration.Nanoseconds()),
		Filename:      testContext.FileName,
		IsMeasurement: testContext.IsMeasurement,
		LineNumber:    testContext.LineNumber,
		TestText:      testContext.FullTestText,
	})
}

// GetReconciledResults is a function added to aggregate a Claim's results.  Due to the limitations of
// test-network-function-claim's Go Client, results are generalized to map[string]interface{}.  This method is needed
// to take the results gleaned from JUnit output, and to combine them with the contexts built up by subsequent clals to
// RecordResult.  The combination of the two forms a Claim's results.
func GetReconciledResults(testResults map[string]junit.TestResult) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for key, vals := range results {
		// JSON cannot handle complex key types, so this flattens the complex key into a string format.
		strKey := fmt.Sprintf("{\"url\":\"%s\",\"version\":\"%s\"}", key.Url, key.Version)
		// initializes the result map, if necessary
		if _, ok := resultMap[strKey]; !ok {
			resultMap[strKey] = make([]claim.Result, 0)
		}
		// a codec which correlates claim.Result, JUnit results (testResults), and builds up the map
		// of claim's results.
		for _, val := range vals {
			val.Passed = testResults[val.TestText].Passed
			testFailReason := testResults[val.TestText].FailureReason
			val.FailureReason = testFailReason
			resultMap[strKey] = append(resultMap[strKey].([]claim.Result), val)
		}
	}
	return resultMap
}
