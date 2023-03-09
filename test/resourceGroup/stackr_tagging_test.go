package main

import (
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/spf13/pflag" // godog v0.11.0 and later
	"github.com/stretchr/testify/assert"
)

var terraformOptions *terraform.Options
var tm *testing.T = new(testing.T)

// godog.TestSuite
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func tearDown(t *testing.T, terraformOptions *terraform.Options) {
	log.Println("IN TEAR DOWN")
	terraform.Destroy(t, terraformOptions)
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

		Vars: map[string]interface{}{
			"stackrName":     resourcegroup_name,
			"stackrLocation": resourcegroup_location,
		},
	})

	terraform.InitAndApply(t, tfOptions)

	log.Printf("terraformOptions: %s", tfOptions.TerraformDir)

	return tfOptions
}

func aResourceGroupIsCreated() error {
	log.Println("IN GIVEN")
	return nil
}

func iCheckForTagsAgainstTehResourceGroup() error {
	log.Println("IN WHEN")
	return nil
}

func iExpectToHaveAtLeastTheFollowingTagsPresent() error {
	log.Println("IN THEN")
	assert.Equal(tm, "ABC", "DEF", "The two flags should be the same")
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	log.Println("IN INITIALIZE_SCENARIO")
	ctx.Step(`^A resource group is created$`, aResourceGroupIsCreated)
	ctx.Step(`^I check for tags against the resource group$`, iCheckForTagsAgainstTehResourceGroup)
	ctx.Step(`^I expect to have at least the following tags present$`, iExpectToHaveAtLeastTheFollowingTagsPresent)
}

// godog.TestSuite
func init() {
	log.Println("IN INIT")
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

	log.Println("BEFORE TEARDOWN")

	tearDown(tm, terraformOptions)

	os.Exit(status)
}
