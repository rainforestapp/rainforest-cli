package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainforestapp/rainforest-cli/rainforest"
)

// apiClient defines the subset of rainforest.Client methods we use.
type apiClient interface {
	// Resources
	GetSites() ([]rainforest.Site, error)
	GetEnvironments() ([]rainforest.Environment, error)
	GetFolders() ([]rainforest.Folder, error)
	GetPlatforms() ([]rainforest.Platform, error)
	GetFeatures() ([]rainforest.Feature, error)
	GetRunGroups() ([]rainforest.RunGroup, error)

	// Tests
	GetTests(params *rainforest.RFTestFilters) ([]rainforest.RFTest, error)
	GetTest(testID int) (*rainforest.RFTest, error)
	GetTestIDs() ([]rainforest.TestIDPair, error)
	CreateTest(test *rainforest.RFTest) error
	UpdateTest(test *rainforest.RFTest, branchID int) error
	DeleteTest(testID int) error

	// Runs
	CreateRun(params rainforest.RunParams) (*rainforest.RunStatus, error)
	CheckRunStatus(runID int) (*rainforest.RunStatus, error)

	// Branches
	GetBranches(params ...string) ([]rainforest.Branch, error)
	CreateBranch(branch *rainforest.Branch) error
	MergeBranch(branchID int) error
	DeleteBranch(branchID int) error
}

type handlers struct {
	client apiClient
}

// jsonResult serializes v to JSON and returns it as a text tool result.
func jsonResult(v interface{}) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to serialize response: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// errResult returns an MCP error result.
func errResult(format string, args ...interface{}) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(fmt.Sprintf(format, args...)), nil
}

// --- Resource listing handlers ---

func (h *handlers) listSites(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sites, err := h.client.GetSites()
	if err != nil {
		return errResult("failed to list sites: %v", err)
	}
	return jsonResult(sites)
}

func (h *handlers) listEnvironments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envs, err := h.client.GetEnvironments()
	if err != nil {
		return errResult("failed to list environments: %v", err)
	}
	return jsonResult(envs)
}

func (h *handlers) listFolders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	folders, err := h.client.GetFolders()
	if err != nil {
		return errResult("failed to list folders: %v", err)
	}
	return jsonResult(folders)
}

func (h *handlers) listPlatforms(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	platforms, err := h.client.GetPlatforms()
	if err != nil {
		return errResult("failed to list platforms: %v", err)
	}
	return jsonResult(platforms)
}

func (h *handlers) listFeatures(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	features, err := h.client.GetFeatures()
	if err != nil {
		return errResult("failed to list features: %v", err)
	}
	return jsonResult(features)
}

func (h *handlers) listRunGroups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groups, err := h.client.GetRunGroups()
	if err != nil {
		return errResult("failed to list run groups: %v", err)
	}
	return jsonResult(groups)
}

// --- Test management handlers ---

func (h *handlers) listTests(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filters := &rainforest.RFTestFilters{}

	if tags, ok := req.GetArguments()["tags"].([]interface{}); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				filters.Tags = append(filters.Tags, s)
			}
		}
	}
	if v, ok := req.GetArguments()["site_id"].(float64); ok {
		filters.SiteID = int(v)
	}
	if v, ok := req.GetArguments()["folder_id"].(float64); ok {
		filters.SmartFolderID = int(v)
	}
	if v, ok := req.GetArguments()["feature_id"].(float64); ok {
		filters.FeatureID = int(v)
	}
	if v, ok := req.GetArguments()["run_group_id"].(float64); ok {
		filters.RunGroupID = int(v)
	}

	tests, err := h.client.GetTests(filters)
	if err != nil {
		return errResult("failed to list tests: %v", err)
	}

	// Return a summary view without full step details
	type testSummary struct {
		ID        int      `json:"id"`
		RFMLID    string   `json:"rfml_id"`
		Title     string   `json:"title"`
		State     string   `json:"state"`
		Tags      []string `json:"tags"`
		Type      string   `json:"type"`
		SiteID    int      `json:"site_id,omitempty"`
		FeatureID int      `json:"feature_id,omitempty"`
	}
	summaries := make([]testSummary, len(tests))
	for i, t := range tests {
		summaries[i] = testSummary{
			ID:        t.TestID,
			RFMLID:    t.RFMLID,
			Title:     t.Title,
			State:     t.State,
			Tags:      t.Tags,
			Type:      t.Type,
			SiteID:    t.SiteID,
			FeatureID: int(t.FeatureID),
		}
	}
	return jsonResult(summaries)
}

// testResponse is the structured response for get_test.
type testResponse struct {
	ID        int            `json:"id"`
	RFMLID    string         `json:"rfml_id"`
	Title     string         `json:"title"`
	State     string         `json:"state"`
	StartURI  string         `json:"start_uri"`
	SiteID    int            `json:"site_id,omitempty"`
	Tags      []string       `json:"tags"`
	Platforms []string       `json:"platforms"`
	FeatureID int            `json:"feature_id,omitempty"`
	Type      string         `json:"type"`
	Steps     []stepResponse `json:"steps"`
}

type stepResponse struct {
	Index    int    `json:"index"`
	Type     string `json:"type"` // "step" or "embedded_test"
	Action   string `json:"action,omitempty"`
	Response string `json:"response,omitempty"`
	Redirect bool   `json:"redirect"`
	RFMLID   string `json:"rfml_id,omitempty"` // for embedded tests
}

func (h *handlers) getTest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}

	test, err := h.client.GetTest(testID)
	if err != nil {
		return errResult("failed to get test %d: %v", testID, err)
	}

	// Unmarshal elements into steps
	testIDPairs, err := h.client.GetTestIDs()
	if err != nil {
		return errResult("failed to get test ID mappings: %v", err)
	}
	coll := rainforest.NewTestIDCollection(testIDPairs)
	if err := test.PrepareToWriteAsRFML(*coll, false); err != nil {
		return errResult("failed to parse test steps: %v", err)
	}

	resp := testResponse{
		ID:        test.TestID,
		RFMLID:    test.RFMLID,
		Title:     test.Title,
		State:     test.State,
		StartURI:  test.StartURI,
		SiteID:    test.SiteID,
		Tags:      test.Tags,
		Platforms: test.Platforms,
		FeatureID: int(test.FeatureID),
		Type:      test.Type,
	}
	if resp.Tags == nil {
		resp.Tags = []string{}
	}
	if resp.Platforms == nil {
		resp.Platforms = []string{}
	}

	for i, step := range test.Steps {
		switch s := step.(type) {
		case rainforest.RFTestStep:
			resp.Steps = append(resp.Steps, stepResponse{
				Index:    i,
				Type:     "step",
				Action:   s.Action,
				Response: s.Response,
				Redirect: s.Redirect,
			})
		case rainforest.RFEmbeddedTest:
			resp.Steps = append(resp.Steps, stepResponse{
				Index:    i,
				Type:     "embedded_test",
				RFMLID:   s.RFMLID,
				Redirect: s.Redirect,
			})
		}
	}
	if resp.Steps == nil {
		resp.Steps = []stepResponse{}
	}

	return jsonResult(resp)
}

func (h *handlers) createTest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := requireString(req, "title")
	if err != nil {
		return errResult("%v", err)
	}

	test := &rainforest.RFTest{
		Title:    title,
		StartURI: getStringOr(req, "start_uri", "/"),
		State:    getStringOr(req, "state", "enabled"),
		Type:     getStringOr(req, "type", "test"),
		Source:   "rainforest-cli",
	}

	if tags, ok := req.GetArguments()["tags"].([]interface{}); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				test.Tags = append(test.Tags, s)
			}
		}
	}
	if test.Tags == nil {
		test.Tags = []string{}
	}

	if platforms, ok := req.GetArguments()["platforms"].([]interface{}); ok {
		for _, p := range platforms {
			if s, ok := p.(string); ok {
				test.Platforms = append(test.Platforms, s)
			}
		}
	}

	if v, ok := req.GetArguments()["site_id"].(float64); ok {
		test.SiteID = int(v)
	}
	if v, ok := req.GetArguments()["feature_id"].(float64); ok {
		test.FeatureID = rainforest.FeatureIDInt(int(v))
	}

	// Parse steps
	if stepsRaw, ok := req.GetArguments()["steps"].([]interface{}); ok {
		for i, stepRaw := range stepsRaw {
			stepMap, ok := stepRaw.(map[string]interface{})
			if !ok {
				return errResult("step %d: expected an object with 'action' and 'response' fields", i)
			}
			action, _ := stepMap["action"].(string)
			response, _ := stepMap["response"].(string)
			if action == "" || response == "" {
				return errResult("step %d: both 'action' and 'response' are required", i)
			}
			redirect := true
			if r, ok := stepMap["redirect"].(bool); ok {
				redirect = r
			}
			test.Steps = append(test.Steps, rainforest.RFTestStep{
				Action:   action,
				Response: response,
				Redirect: redirect,
			})
		}
	}

	// Generate a unique RFML ID
	test.RFMLID = fmt.Sprintf("mcp_%d", uniqueID())

	// Get test ID collection for marshalling
	testIDPairs, err := h.client.GetTestIDs()
	if err != nil {
		return errResult("failed to get test ID mappings: %v", err)
	}
	coll := rainforest.NewTestIDCollection(testIDPairs)

	// Prepare and create the test (initially without steps to get an ID)
	emptyTest := &rainforest.RFTest{
		RFMLID: test.RFMLID,
		Title:  test.Title,
		Type:   test.Type,
		Source: "rainforest-cli",
		Tags:   test.Tags,
	}
	if err := emptyTest.PrepareToUploadFromRFML(*coll); err != nil {
		return errResult("failed to prepare test: %v", err)
	}
	if err := h.client.CreateTest(emptyTest); err != nil {
		return errResult("failed to create test: %v", err)
	}

	// Refresh ID collection to get the new test's ID
	testIDPairs, err = h.client.GetTestIDs()
	if err != nil {
		return errResult("failed to refresh test IDs: %v", err)
	}
	coll = rainforest.NewTestIDCollection(testIDPairs)

	testID, err := coll.GetTestID(test.RFMLID)
	if err != nil {
		return errResult("failed to find created test: %v", err)
	}
	test.TestID = testID

	// Now update with full details including steps
	if err := test.PrepareToUploadFromRFML(*coll); err != nil {
		return errResult("failed to prepare test for update: %v", err)
	}
	if err := h.client.UpdateTest(test, rainforest.NO_BRANCH); err != nil {
		return errResult("failed to update test with steps: %v", err)
	}

	return jsonResult(map[string]interface{}{
		"id":      testID,
		"rfml_id": test.RFMLID,
		"title":   test.Title,
		"message": "Test created successfully",
	})
}

func (h *handlers) deleteTest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}

	if err := h.client.DeleteTest(testID); err != nil {
		return errResult("failed to delete test %d: %v", testID, err)
	}

	return jsonResult(map[string]string{
		"message": fmt.Sprintf("Test %d deleted successfully", testID),
	})
}

// --- Step-level handlers ---

// fetchAndParseSteps is a helper for step operations.
// It fetches a test, parses its elements into steps, and returns everything needed for modification.
func (h *handlers) fetchAndParseSteps(testID int) (*rainforest.RFTest, *rainforest.TestIDCollection, error) {
	test, err := h.client.GetTest(testID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get test %d: %v", testID, err)
	}

	testIDPairs, err := h.client.GetTestIDs()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get test ID mappings: %v", err)
	}
	coll := rainforest.NewTestIDCollection(testIDPairs)

	if err := test.PrepareToWriteAsRFML(*coll, false); err != nil {
		return nil, nil, fmt.Errorf("failed to parse test steps: %v", err)
	}

	return test, coll, nil
}

// saveSteps marshals steps back to elements and pushes the update.
func (h *handlers) saveSteps(test *rainforest.RFTest, coll *rainforest.TestIDCollection) error {
	test.Source = "rainforest-cli"
	if test.Tags == nil {
		test.Tags = []string{}
	}
	if err := test.PrepareToUploadFromRFML(*coll); err != nil {
		return fmt.Errorf("failed to prepare test: %v", err)
	}
	if err := h.client.UpdateTest(test, rainforest.NO_BRANCH); err != nil {
		return fmt.Errorf("failed to update test: %v", err)
	}
	return nil
}

func (h *handlers) addTestStep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}
	action, err := requireString(req, "action")
	if err != nil {
		return errResult("%v", err)
	}
	response, err := requireString(req, "response")
	if err != nil {
		return errResult("%v", err)
	}

	redirect := true
	if r, ok := req.GetArguments()["redirect"].(bool); ok {
		redirect = r
	}

	test, coll, err := h.fetchAndParseSteps(testID)
	if err != nil {
		return errResult("%v", err)
	}

	newStep := rainforest.RFTestStep{
		Action:   action,
		Response: response,
		Redirect: redirect,
	}

	// Insert at position or append
	if posRaw, ok := req.GetArguments()["position"].(float64); ok {
		pos := int(posRaw)
		if pos < 0 || pos > len(test.Steps) {
			return errResult("position %d is out of range (0-%d)", pos, len(test.Steps))
		}
		// Insert at position
		test.Steps = append(test.Steps, nil)
		copy(test.Steps[pos+1:], test.Steps[pos:])
		test.Steps[pos] = newStep
	} else {
		test.Steps = append(test.Steps, newStep)
	}

	if err := h.saveSteps(test, coll); err != nil {
		return errResult("%v", err)
	}

	return jsonResult(map[string]interface{}{
		"message":     "Step added successfully",
		"test_id":     testID,
		"total_steps": len(test.Steps),
	})
}

func (h *handlers) updateTestStep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}
	stepIndex, err := requireInt(req, "step_index")
	if err != nil {
		return errResult("%v", err)
	}

	test, coll, err := h.fetchAndParseSteps(testID)
	if err != nil {
		return errResult("%v", err)
	}

	if stepIndex < 0 || stepIndex >= len(test.Steps) {
		return errResult("step_index %d is out of range (0-%d)", stepIndex, len(test.Steps)-1)
	}

	step, ok := test.Steps[stepIndex].(rainforest.RFTestStep)
	if !ok {
		return errResult("step at index %d is an embedded test reference, not an editable step", stepIndex)
	}

	if action, ok := req.GetArguments()["action"].(string); ok {
		step.Action = action
	}
	if response, ok := req.GetArguments()["response"].(string); ok {
		step.Response = response
	}
	if redirect, ok := req.GetArguments()["redirect"].(bool); ok {
		step.Redirect = redirect
	}

	test.Steps[stepIndex] = step

	if err := h.saveSteps(test, coll); err != nil {
		return errResult("%v", err)
	}

	return jsonResult(map[string]interface{}{
		"message":  "Step updated successfully",
		"test_id":  testID,
		"step":     stepIndex,
		"action":   step.Action,
		"response": step.Response,
	})
}

func (h *handlers) deleteTestStep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}
	stepIndex, err := requireInt(req, "step_index")
	if err != nil {
		return errResult("%v", err)
	}

	test, coll, err := h.fetchAndParseSteps(testID)
	if err != nil {
		return errResult("%v", err)
	}

	if stepIndex < 0 || stepIndex >= len(test.Steps) {
		return errResult("step_index %d is out of range (0-%d)", stepIndex, len(test.Steps)-1)
	}

	test.Steps = append(test.Steps[:stepIndex], test.Steps[stepIndex+1:]...)

	if err := h.saveSteps(test, coll); err != nil {
		return errResult("%v", err)
	}

	return jsonResult(map[string]interface{}{
		"message":     "Step deleted successfully",
		"test_id":     testID,
		"total_steps": len(test.Steps),
	})
}

func (h *handlers) moveTestStep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	testID, err := requireInt(req, "test_id")
	if err != nil {
		return errResult("%v", err)
	}
	fromIndex, err := requireInt(req, "from_index")
	if err != nil {
		return errResult("%v", err)
	}
	toIndex, err := requireInt(req, "to_index")
	if err != nil {
		return errResult("%v", err)
	}

	test, coll, err := h.fetchAndParseSteps(testID)
	if err != nil {
		return errResult("%v", err)
	}

	if fromIndex < 0 || fromIndex >= len(test.Steps) {
		return errResult("from_index %d is out of range (0-%d)", fromIndex, len(test.Steps)-1)
	}
	if toIndex < 0 || toIndex >= len(test.Steps) {
		return errResult("to_index %d is out of range (0-%d)", toIndex, len(test.Steps)-1)
	}

	// Remove the step from its current position
	step := test.Steps[fromIndex]
	test.Steps = append(test.Steps[:fromIndex], test.Steps[fromIndex+1:]...)
	// Insert at new position
	test.Steps = append(test.Steps[:toIndex], append([]interface{}{step}, test.Steps[toIndex:]...)...)

	if err := h.saveSteps(test, coll); err != nil {
		return errResult("%v", err)
	}

	return jsonResult(map[string]interface{}{
		"message":     "Step moved successfully",
		"test_id":     testID,
		"from":        fromIndex,
		"to":          toIndex,
		"total_steps": len(test.Steps),
	})
}

// --- Run management handlers ---

func (h *handlers) startRun(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := rainforest.RunParams{}

	// Parse test_ids — can be array of numbers or contain "all"
	if testIDsRaw, ok := req.GetArguments()["test_ids"].([]interface{}); ok && len(testIDsRaw) > 0 {
		if s, ok := testIDsRaw[0].(string); ok && s == "all" {
			params.Tests = "all"
		} else {
			ids := []int{}
			for _, v := range testIDsRaw {
				if n, ok := v.(float64); ok {
					ids = append(ids, int(n))
				}
			}
			params.Tests = ids
		}
	}

	if tags, ok := req.GetArguments()["tags"].([]interface{}); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				params.Tags = append(params.Tags, s)
			}
		}
	}

	if v, ok := req.GetArguments()["site_id"].(float64); ok {
		params.SiteID = int(v)
	}
	if v, ok := req.GetArguments()["folder_id"].(float64); ok {
		params.SmartFolderID = int(v)
	}
	if v, ok := req.GetArguments()["feature_id"].(float64); ok {
		params.FeatureID = int(v)
	}
	if v, ok := req.GetArguments()["run_group_id"].(float64); ok {
		params.RunGroupID = int(v)
	}
	if v, ok := req.GetArguments()["environment_id"].(float64); ok {
		params.EnvironmentID = int(v)
	}
	if v, ok := req.GetArguments()["automation_max_retries"].(float64); ok {
		params.AutomationMaxRetries = int(v)
	}

	if platforms, ok := req.GetArguments()["platforms"].([]interface{}); ok {
		for _, p := range platforms {
			if s, ok := p.(string); ok {
				params.Browsers = append(params.Browsers, s)
			}
		}
	}

	if v, ok := req.GetArguments()["execution_method"].(string); ok {
		params.ExecutionMethod = v
	}
	if v, ok := req.GetArguments()["conflict"].(string); ok {
		params.Conflict = v
	}
	if v, ok := req.GetArguments()["description"].(string); ok {
		params.Description = v
	}
	if v, ok := req.GetArguments()["release"].(string); ok {
		params.Release = v
	}

	// Handle branch by name → ID lookup
	if branchName, ok := req.GetArguments()["branch"].(string); ok && branchName != "" {
		branches, err := h.client.GetBranches(branchName)
		if err != nil {
			return errResult("failed to look up branch '%s': %v", branchName, err)
		}
		if len(branches) == 0 {
			return errResult("branch '%s' not found", branchName)
		}
		params.BranchID = branches[0].ID
	}

	status, err := h.client.CreateRun(params)
	if err != nil {
		return errResult("failed to start run: %v", err)
	}

	return jsonResult(status)
}

func (h *handlers) getRunStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runID, err := requireInt(req, "run_id")
	if err != nil {
		return errResult("%v", err)
	}

	status, err := h.client.CheckRunStatus(runID)
	if err != nil {
		return errResult("failed to get run status for %d: %v", runID, err)
	}

	return jsonResult(status)
}

func (h *handlers) rerunFailed(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runID, err := requireInt(req, "run_id")
	if err != nil {
		return errResult("%v", err)
	}

	params := rainforest.RunParams{
		RunID: runID,
	}
	if v, ok := req.GetArguments()["conflict"].(string); ok {
		params.Conflict = v
	}

	status, err := h.client.CreateRun(params)
	if err != nil {
		return errResult("failed to rerun failed tests for run %d: %v", runID, err)
	}

	return jsonResult(status)
}

// --- Branch management handlers ---

func (h *handlers) listBranches(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args []string
	if name, ok := req.GetArguments()["name"].(string); ok && name != "" {
		args = append(args, name)
	}

	branches, err := h.client.GetBranches(args...)
	if err != nil {
		return errResult("failed to list branches: %v", err)
	}

	return jsonResult(branches)
}

func (h *handlers) createBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := requireString(req, "name")
	if err != nil {
		return errResult("%v", err)
	}

	branch := &rainforest.Branch{Name: name}
	if err := h.client.CreateBranch(branch); err != nil {
		return errResult("failed to create branch '%s': %v", name, err)
	}

	return jsonResult(map[string]string{
		"message": fmt.Sprintf("Branch '%s' created successfully", name),
	})
}

func (h *handlers) mergeBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branchID, err := requireInt(req, "branch_id")
	if err != nil {
		return errResult("%v", err)
	}

	if err := h.client.MergeBranch(branchID); err != nil {
		return errResult("failed to merge branch %d: %v", branchID, err)
	}

	return jsonResult(map[string]string{
		"message": fmt.Sprintf("Branch %d merged successfully", branchID),
	})
}

func (h *handlers) deleteBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branchID, err := requireInt(req, "branch_id")
	if err != nil {
		return errResult("%v", err)
	}

	if err := h.client.DeleteBranch(branchID); err != nil {
		return errResult("failed to delete branch %d: %v", branchID, err)
	}

	return jsonResult(map[string]string{
		"message": fmt.Sprintf("Branch %d deleted successfully", branchID),
	})
}

// --- Helper functions ---

func requireInt(req mcp.CallToolRequest, key string) (int, error) {
	v, err := req.RequireInt(key)
	if err != nil {
		return 0, fmt.Errorf("%s is required and must be a number", key)
	}
	return v, nil
}

func requireString(req mcp.CallToolRequest, key string) (string, error) {
	v, err := req.RequireString(key)
	if err != nil || v == "" {
		return "", fmt.Errorf("%s is required and must be a non-empty string", key)
	}
	return v, nil
}

func getStringOr(req mcp.CallToolRequest, key, defaultVal string) string {
	return req.GetString(key, defaultVal)
}

// uniqueID returns a unique identifier based on timestamp.
func uniqueID() int64 {
	return time.Now().UnixNano()
}
