# stackrModule
Represents the core stackr module building the Azure base structures

## Useful links
[godog](https://github.com/cucumber/godog)

[cucumber-html-reporter](https://github.com/gkushang/cucumber-html-reporter)

[godog - multiple formatters](https://github.com/cucumber/godog/issues/346)

[terratest](https://terratest.gruntwork.io)

[sphinx](https://www.sphinx-doc.org/en/master/tutorial/getting-started.html)

[Azure Blob Storage Python Library](https://learn.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-python?tabs=managed-identity%2Croles-azure-cli%2Csign-in-azure-cli)
OR
[Use Github to upload to Azure Static Website](https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blobs-static-site-github-actions?tabs=userlevel)

## Build Instructions for godog BDD tests

- Initialize module

    `go mod init github.com/tunix78/stackrModule`
- Pull all dependencies

    `go mod tidy`
- Run godog tests and output test results into json file
    - Needs to be run in the folder with the golang test file (e.g. test/resourceGroup)
    `go test -v --godog.random --godog.format=cucumber:results/cucumber.json`

- Install cucumber-html-reporter

    `npm install cucumber-html-reporter --save-dev`
- Upgrade npm (if required)

    `npm install -g npm@9.6.1`
- Create index.js (see cucumber-html-reporter link) into root test folder
- Run cucumber-html-report against the files specified in index.js to create the html report
    - The html file will be generated into the file specified in index.js
    - A browser window will automatically be opened usually

    `node index.js`

## Build instructions for sphinx RST documentation

- Install sphinx

    ```
    python -m venv .venv
    source .venv/bin/activate
    python -m pip install sphinx
    ```
- Build the documentation

    `sphinx-build -b html docs/source/ docs/build/html`

## When checking if storage was uploaded to Azure

- Do NOT use the blob storage level link

    `https://svensappstorage.blob.core.windows.net/$web/docs/stackrModule/resourceGroup/index.html`

- Use the primary endpoint of the storage account as the base link

    `https://svensappstorage.z6.web.core.windows.net/`

## Instructions to connect to MongoDB and store test results as well as decision logs

### Connect to MongoDB using VS Code

- Install extension "MongoDB for VS Code"
- Connect to MongoDB using this string: "mongodb+srv://svengauggel:<password>@stackrstore.0omkqsa.mongodb.net/test"