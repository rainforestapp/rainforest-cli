package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer, client apiClient) {
	h := &handlers{client: client}

	// Resource listing tools
	s.AddTool(listSitesTool(), h.listSites)
	s.AddTool(listEnvironmentsTool(), h.listEnvironments)
	s.AddTool(listFoldersTool(), h.listFolders)
	s.AddTool(listPlatformsTool(), h.listPlatforms)
	s.AddTool(listFeaturesTool(), h.listFeatures)
	s.AddTool(listRunGroupsTool(), h.listRunGroups)

	// Test management tools
	s.AddTool(listTestsTool(), h.listTests)
	s.AddTool(getTestTool(), h.getTest)
	s.AddTool(createTestTool(), h.createTest)
	s.AddTool(deleteTestTool(), h.deleteTest)

	// Step-level tools
	s.AddTool(addTestStepTool(), h.addTestStep)
	s.AddTool(updateTestStepTool(), h.updateTestStep)
	s.AddTool(deleteTestStepTool(), h.deleteTestStep)
	s.AddTool(moveTestStepTool(), h.moveTestStep)

	// Run management tools
	s.AddTool(startRunTool(), h.startRun)
	s.AddTool(getRunStatusTool(), h.getRunStatus)
	s.AddTool(rerunFailedTool(), h.rerunFailed)

	// Branch management tools
	s.AddTool(listBranchesTool(), h.listBranches)
	s.AddTool(createBranchTool(), h.createBranch)
	s.AddTool(mergeBranchTool(), h.mergeBranch)
	s.AddTool(deleteBranchTool(), h.deleteBranch)
}

// --- Resource listing tool definitions ---

func listSitesTool() mcp.Tool {
	return mcp.NewTool("list_sites",
		mcp.WithDescription("List all available sites in your Rainforest account. Returns site IDs, names, and categories."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

func listEnvironmentsTool() mcp.Tool {
	return mcp.NewTool("list_environments",
		mcp.WithDescription("List all available environments in your Rainforest account. Returns environment IDs and names."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

func listFoldersTool() mcp.Tool {
	return mcp.NewTool("list_folders",
		mcp.WithDescription("List all available smart folders in your Rainforest account. Returns folder IDs and titles."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

func listPlatformsTool() mcp.Tool {
	return mcp.NewTool("list_platforms",
		mcp.WithDescription("List all available platforms (browsers) you can run tests against. Returns platform names and descriptions."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

func listFeaturesTool() mcp.Tool {
	return mcp.NewTool("list_features",
		mcp.WithDescription("List all available features in your Rainforest account. Returns feature IDs and titles."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

func listRunGroupsTool() mcp.Tool {
	return mcp.NewTool("list_run_groups",
		mcp.WithDescription("List all available run groups in your Rainforest account. Returns run group IDs and titles."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
}

// --- Test management tool definitions ---

func listTestsTool() mcp.Tool {
	return mcp.NewTool("list_tests",
		mcp.WithDescription("List tests in your Rainforest account, optionally filtered by tags, site, folder, feature, or run group."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithArray("tags",
			mcp.Description("Filter tests by tags"),
			mcp.WithStringItems(),
		),
		mcp.WithNumber("site_id",
			mcp.Description("Filter tests by site ID"),
		),
		mcp.WithNumber("folder_id",
			mcp.Description("Filter tests by folder ID"),
		),
		mcp.WithNumber("feature_id",
			mcp.Description("Filter tests by feature ID"),
		),
		mcp.WithNumber("run_group_id",
			mcp.Description("Filter tests by run group ID"),
		),
	)
}

func getTestTool() mcp.Tool {
	return mcp.NewTool("get_test",
		mcp.WithDescription("Get a single test by ID, including all its steps. Each step has an action (instruction) and a response (verification question)."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID to retrieve"),
		),
	)
}

func createTestTool() mcp.Tool {
	return mcp.NewTool("create_test",
		mcp.WithDescription("Create a new test in Rainforest with the given metadata and steps. Steps are action/response pairs where the action tells the tester what to do and the response asks a verification question (must end with '?')."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the test"),
		),
		mcp.WithString("start_uri",
			mcp.Description("Starting URI for the test (defaults to '/')"),
		),
		mcp.WithArray("tags",
			mcp.Description("Tags to apply to the test"),
			mcp.WithStringItems(),
		),
		mcp.WithArray("platforms",
			mcp.Description("Platforms to run the test on (e.g. 'chrome', 'firefox')"),
			mcp.WithStringItems(),
		),
		mcp.WithNumber("site_id",
			mcp.Description("Site ID for the test"),
		),
		mcp.WithNumber("feature_id",
			mcp.Description("Feature ID for the test"),
		),
		mcp.WithString("type",
			mcp.Description("Test type: 'test' or 'snippet' (defaults to 'test')"),
			mcp.Enum("test", "snippet"),
		),
		mcp.WithString("state",
			mcp.Description("Test state: 'enabled', 'disabled', or 'draft' (defaults to 'enabled')"),
			mcp.Enum("enabled", "disabled", "draft"),
		),
		// Steps as an array of objects
		mcp.WithArray("steps",
			mcp.Description("Test steps. Each step has an 'action' (what the tester does) and a 'response' (verification question ending with '?')"),
		),
	)
}

func deleteTestTool() mcp.Tool {
	return mcp.NewTool("delete_test",
		mcp.WithDescription("Delete a test from Rainforest by its ID. This action is irreversible."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID to delete"),
		),
	)
}

// --- Step-level tool definitions ---

func addTestStepTool() mcp.Tool {
	return mcp.NewTool("add_test_step",
		mcp.WithDescription("Add a new step to an existing test. A step consists of an action (instruction for the tester) and a response (verification question ending with '?'). By default the step is appended to the end."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID to add the step to"),
		),
		mcp.WithString("action",
			mcp.Required(),
			mcp.Description("The action/instruction for the tester"),
		),
		mcp.WithString("response",
			mcp.Required(),
			mcp.Description("The verification question (should end with '?')"),
		),
		mcp.WithNumber("position",
			mcp.Description("0-based position to insert the step at. Defaults to the end of the step list."),
		),
		mcp.WithBoolean("redirect",
			mcp.Description("Whether this step causes a page redirect (defaults to true)"),
		),
	)
}

func updateTestStepTool() mcp.Tool {
	return mcp.NewTool("update_test_step",
		mcp.WithDescription("Update an existing step in a test. You can update the action, response, and/or redirect flag. Only the fields you provide will be changed."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID containing the step"),
		),
		mcp.WithNumber("step_index",
			mcp.Required(),
			mcp.Description("0-based index of the step to update"),
		),
		mcp.WithString("action",
			mcp.Description("New action/instruction text"),
		),
		mcp.WithString("response",
			mcp.Description("New verification question text"),
		),
		mcp.WithBoolean("redirect",
			mcp.Description("Whether this step causes a page redirect"),
		),
	)
}

func deleteTestStepTool() mcp.Tool {
	return mcp.NewTool("delete_test_step",
		mcp.WithDescription("Delete a step from a test by its index."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID containing the step"),
		),
		mcp.WithNumber("step_index",
			mcp.Required(),
			mcp.Description("0-based index of the step to delete"),
		),
	)
}

func moveTestStepTool() mcp.Tool {
	return mcp.NewTool("move_test_step",
		mcp.WithDescription("Move a step from one position to another within a test."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithNumber("test_id",
			mcp.Required(),
			mcp.Description("The test ID containing the step"),
		),
		mcp.WithNumber("from_index",
			mcp.Required(),
			mcp.Description("0-based index of the step to move"),
		),
		mcp.WithNumber("to_index",
			mcp.Required(),
			mcp.Description("0-based index to move the step to"),
		),
	)
}

// --- Run management tool definitions ---

func startRunTool() mcp.Tool {
	return mcp.NewTool("start_run",
		mcp.WithDescription("Start a new test run on Rainforest. You must specify which tests to run using test_ids, tags, folder_id, feature_id, run_group_id, or pass 'all' as test_ids to run everything."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithArray("test_ids",
			mcp.Description("List of test IDs to run, or use the string 'all' to run all tests"),
		),
		mcp.WithArray("tags",
			mcp.Description("Filter tests by tags"),
			mcp.WithStringItems(),
		),
		mcp.WithNumber("site_id",
			mcp.Description("Filter tests by site ID"),
		),
		mcp.WithNumber("folder_id",
			mcp.Description("Filter tests by folder ID"),
		),
		mcp.WithNumber("feature_id",
			mcp.Description("Filter tests by feature ID"),
		),
		mcp.WithNumber("run_group_id",
			mcp.Description("Start a run using a run group"),
		),
		mcp.WithArray("platforms",
			mcp.Description("Platforms to run against (overrides test-level settings)"),
			mcp.WithStringItems(),
		),
		mcp.WithNumber("environment_id",
			mcp.Description("Environment ID to use for the run"),
		),
		mcp.WithString("execution_method",
			mcp.Description("Execution method for the run"),
			mcp.Enum("crowd", "automation", "automation_and_crowd", "on_premise"),
		),
		mcp.WithString("conflict",
			mcp.Description("How to handle conflicting runs"),
			mcp.Enum("cancel", "cancel-all"),
		),
		mcp.WithString("description",
			mcp.Description("Description for the run"),
		),
		mcp.WithString("release",
			mcp.Description("Release ID to associate with this run (e.g. commit SHA, build ID)"),
		),
		mcp.WithString("branch",
			mcp.Description("Branch name to run tests on"),
		),
		mcp.WithNumber("automation_max_retries",
			mcp.Description("Max retries for automation within the same run before reporting failure"),
		),
	)
}

func getRunStatusTool() mcp.Tool {
	return mcp.NewTool("get_run_status",
		mcp.WithDescription("Get the current status of a test run, including progress, pass/fail counts, and result."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithNumber("run_id",
			mcp.Required(),
			mcp.Description("The run ID to check"),
		),
	)
}

func rerunFailedTool() mcp.Tool {
	return mcp.NewTool("rerun_failed",
		mcp.WithDescription("Rerun only the failed tests from a previous run. The environment, execution method, and other settings are copied from the original run."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("run_id",
			mcp.Required(),
			mcp.Description("The run ID to rerun failed tests from"),
		),
		mcp.WithString("conflict",
			mcp.Description("How to handle conflicting runs"),
			mcp.Enum("cancel", "cancel-all"),
		),
	)
}

// --- Branch management tool definitions ---

func listBranchesTool() mcp.Tool {
	return mcp.NewTool("list_branches",
		mcp.WithDescription("List test branches, optionally filtered by name."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("name",
			mcp.Description("Filter branches by name"),
		),
	)
}

func createBranchTool() mcp.Tool {
	return mcp.NewTool("create_branch",
		mcp.WithDescription("Create a new test branch."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the new branch"),
		),
	)
}

func mergeBranchTool() mcp.Tool {
	return mcp.NewTool("merge_branch",
		mcp.WithDescription("Merge a branch into the main branch."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("branch_id",
			mcp.Required(),
			mcp.Description("ID of the branch to merge"),
		),
	)
}

func deleteBranchTool() mcp.Tool {
	return mcp.NewTool("delete_branch",
		mcp.WithDescription("Delete an existing branch."),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithNumber("branch_id",
			mcp.Required(),
			mcp.Description("ID of the branch to delete"),
		),
	)
}
