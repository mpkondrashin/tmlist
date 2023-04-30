 
# Manage Trend Micro Cloud One Workload Security/Deep Security Lists

**Include antivirus exclusion lists one into another**

## Usage

To use TMList one has to create exclusion lists with includes, provide API Key and REST API entry point URL, and run the program.

### Create exclusion lists with includes

Open Coud One Workload Security (or Deep Security) console. Go to Policies section -> Common Objects -> Lists.
TMList support following lists:

1. Directory Lists
2. File Expension Lists
3. File Lists

To create a list that combines other lists, click New button, provide a name, and go to the description section. Put into description section the following lines:
```
Include: <list name>
Include: <list name>
...
```
Any other lines can be added to the description; they will be ignored by TMList.
 
After TMList is run, this list will be populated with the contents of the specified lists.

**Warning:** The contents of the list with includes will be deleted! After TMList is run it will contain only included lists.

**Note:** Cycle includes are not alowed

### Get an API Key

Before generating API Key itself, custom role should be created to avoid using default Full Control role.

#### Create a custom restricted role

1. Open the Cloud One Workload Security/Deep Security console
2. Go to Administration -> User Management -> Roles
3. Click New button
4. Give the new role a name, e.g., "tmlists" and allow access only to the web services API:<br/><img width="640" src="images/role01.png"/><br/>
5. Switch to Computer Rights tab. Turn off all rights <br/><img width="640" src="images/role02.png"/><br/>
6. Switch to API Key Rights tab. Pick Custom and do not select any checkboxes <br/><img width="640" src="images/role03.png"/><br/>
7. Switch to Common Object Rights tab. Change the dropdown next to File Lists (All), File Extension Lists (All), and 7. Directory Lists (All) to Full. Custom and check only "Can Edit..." checkbox. <br/><img width="640" src="images/role04.png"/><br/>
8. Press Ok button at the bottom

#### Generate API Key

1. Go to Administration -> User Management -> API Keys
2. Click New button
3. Give the API Key a name, e.g., "TMList"
4. In the Role dropdown pick name of the role created on the previous step <br/><img width="640" src="images/key01.png"/><br/>
5. Press Next button
6. Save the API key to the configuration file (see below)

**Note:** Alternatevly Cloud One API Key can be used with a role that provides Full Control for Endpoint & Workload Security. This approach does not provide appropriate permissions granularity.
 
### Get Entry Point Address
Entry point is URL where program will connect to submit Web API calls

#### For Deep Security users
For Deep Security, address should be the address of one of Deep Security Managers with port 4119 (this is the  default value) and /api path:
```
https://<dsm address>:4119/api
```
 
#### For Cloud One Endpoint & Workload Security users

Open Cloud One console
Click your account at the top of the screen and pick Account Settings
Note your Region on the Account Settings Screen

API entry point will be
```
https://workload.<region>.cloudone.trendmicro.com/api
```
 
### Create a configuration file
In the same directory as  TMList executable create file named config.yaml and open it in your favorite text editor.
Fill it with the following content:
```yaml
address: <your API entry point URL>
api_key: <your apy key>
```

## Get TMList program

Download TMList program for your platform from [Releases](https://github.com/mpkondrashin/tmlist/releases/latest) page

## Run

TMList can be run without parameters. For Linux/macOS:
```commandline
./tmlist
```
For Windows:
```commandline
tmlist.exe
```
It will process all of the supported exclusion lists one by one.

## Options

TMList offers three ways to provide options:

1. Configuration file config.yaml. The application seeks for this file in the current folder or in folder of its executable
2. Environment variables
3. Command line parameters

The following options are available:

| Type | YAML Option<br/>Command line<br/>Env Variable | Description | Default |
| ---- | --------------------------------------------- | ----------- | ------- |
|String|address<br/>--address<br/>TMLIST_ADDRESS|Workload Security entrypoint URL or Deep Security Manager URL|none|
|String|api_key<br/>--api_key<br/>TMLIST_API_KEY|Cloud One or Deep Security API Key|none|
|Boolean|dir<br/>--dir<br/>TMLIST_DIR|Process directory lists|false|
|Boolean|ext<br/>--ext<br/>TMLIST_EXT|Process file extension lists|false|
|Boolean|file<br/>--file<br/>TMLIST_FILE|Process file lists|false|
|Boolean|dry<br/>--dry<br/>TMLIST_DRY|Dry run - do not modify any lists|false|

**Note:** If none of the --dir, --ext or --file options are provided, they all supposed to be true and TMList processes all lists.

**Note:** If the same parameter is provided more than one way, then the following precedence will take place:

1. Environment variables override the configuration file parameters
2. Command line options override both environment variables and configuration file parameters

## Return Codes

If TMList successfully finishes lists modification it returnes code 0. In case of the error, non zero return c  ode can be checked to diagnose a problem.

| Return Code | Description |
| ----------- | ----------- |
|0|Ok|
|2|Command line error|
|3|Other error|
|4|API error|
|5|Cycle Dependence|
|6|List Not Found|
