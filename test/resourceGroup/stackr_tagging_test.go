package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/spf13/pflag" // godog v0.11.0 and later

	"github.com/open-policy-agent/opa/rego"
)

var terraformOptions *terraform.Options
var tm *testing.T = new(testing.T)
var jsonPlan string
var regoDir string

// godog.TestSuite
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func deferTearDown(t *testing.T, terraformOptions *terraform.Options) {
	log.Println("IN TEAR DOWN")

	// write the json file to file so we can use it later on
	err := os.WriteFile("plan.json", []byte(jsonPlan), 0644)
	if err != nil {
		panic(err)
	}

	// we only run against the plan so far
	//defer terraform.Destroy(t, terraformOptions)
}

func setup(t *testing.T) *terraform.Options {
	log.Println("IN SETUP")

	resourcegroup_name := RandomName("stg", "rg")
	resourcegroup_location := "westeurope"

	providerLocation := "../../providers.tf"
	testLocation := "./"
	CopyFile(providerLocation, testLocation)

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../",
		PlanFilePath: "plan.out",

		Vars: map[string]interface{}{
			"stackrName":     resourcegroup_name,
			"stackrLocation": resourcegroup_location,
		},
	})

	log.Printf("terraformOptions: %s", tfOptions.TerraformDir)

	return tfOptions
}

// godog.TestSuite
func init() {
	log.Println("IN INIT")

	regoDir = os.Getenv("REGO_DIR")
	log.Println("REGO_DIR = " + regoDir)

	//godog.BindFlags("godog.", pflag.CommandLine, &opts) // godog v0.10.0 and earlier
	godog.BindCommandLineFlags("godog.", &opts) // godog v0.11.0 and later
}

func TestMain(m *testing.M) {
	log.Println("IN MAIN")

	pflag.Parse()
	opts.Paths = pflag.Args()

	log.Println("BEFORE SETUP")
	terraformOptions = setup(tm)

	log.Println("BEFORE godog.run")

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	log.Println("BEFORE DEFER_TEARDOWN")
	deferTearDown(tm, terraformOptions)

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	log.Println("IN INITIALIZE_SCENARIO")
	ctx.Step(`^A resource group is planned via Terraform$`, aResourceGroupIsPlanned)
	ctx.Step(`^I expect to have at least the following tags present$`, iExpectToHaveAtLeastTheFollowingTagsPresent)
	ctx.Step(`^I expect the location of the resource group to be one of the following$`, iExpectTheLocationOfTheResourceGroupToBeOneOfTheFollowing)
}

func aResourceGroupIsPlanned() error {
	log.Println("IN GIVEN")
	jsonPlan = terraform.InitAndPlanAndShow(tm, terraformOptions)
	//terraform.InitAndApply(tm, terraformOptions)
	return nil
}

func iExpectToHaveAtLeastTheFollowingTagsPresent() error {
	log.Println("IN THEN")

	ctx := context.Background()

	jd := json.NewDecoder(bytes.NewBufferString(jsonPlan))
	jd.UseNumber()

	var input interface{}
	if err := jd.Decode(&input); err != nil {
		return err
	}

	// Create query that returns a single boolean value.
	rego := rego.New(
		rego.Query("data.stackr.allow = true"),
		// FIXME this should be a proper variable, not a relative path
		rego.Load([]string{regoDir + "/rules/required_tags.rego"}, nil),
		rego.Input(input))

	// Run evaluation.
	rs, err := rego.Eval(ctx)
	log.Println("Allowed: " + strconv.FormatBool(rs.Allowed()))

	// Check if we should fail the scenario
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}

	return nil
}

func iExpectTheLocationOfTheResourceGroupToBeOneOfTheFollowing() error {
	log.Println("IN AND")

	ctx := context.Background()

	jd := json.NewDecoder(bytes.NewBufferString(jsonPlan))
	jd.UseNumber()

	var input interface{}
	if err := jd.Decode(&input); err != nil {
		return err
	}

	// Create query that returns a single boolean value.
	rego := rego.New(
		rego.Query("data.stackr.allow = true"),
		rego.Load([]string{regoDir + "/rules/required_locations.rego",
			regoDir + "/data/allowed_locations.json"}, nil),
		rego.Input(input))

	// Run evaluation.
	rs, err := rego.Eval(ctx)
	log.Println("Allowed: " + strconv.FormatBool(rs.Allowed()))

	// Check if we should fail the scenario
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}

	return nil
}
