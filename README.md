
# Manage Trend Micro Cloud One Workload Security/Deep Security Lists

**Add ability to include antivirus exclusion lists one into another**

## Usage
Open Coud One Workload Security (or Deep Security) console. Go to Policies section -> Common Objects -> Lists.
TMList support following lists:
1. Directory Lists
2. File Expension Lists
3. File Lists

To create list that combines ohter lists click New button, provide name and go to description section. Put into description section following line
```
Include: <list name>
Include: <list name>
...
```
Other lines can be added to description - they will be ignored by TMList.

After TMList run, this list will be populated with contents of specified lists.

**Warning:** Contents of the list with includes will be deleted!

**Note:** Cycle includes are not alowed 


## Options

TMList provides following ways to provide options:
1. Configuration file config.yaml. Application seeks for this file in its current folder or in folder of its executable
2. Environment variables
3. Command line parameters

Following options are available:

| Type | YAML Option<br/>Command line<br/>Env Variable | Description | Default |
| ---- | --------------------------------------------- | ----------- | ------- |
|String|address<br/>--address<br/>TMLIST_ADDRESS|Workload Security entrypoint URL or Deep Security Manager URL|none|
|String|api_key<br/>--api_key<br/>TMLIST_API_KEY|Cloud One or Deep Security API Key|none|
|Boolean|dir<br/>--dir<br/>TMLIST_DIR|Process directory lists|false|
|Boolean|ext<br/>--ext<br/>TMLIST_EXT|Process file extension lists|false|
|Boolean|file<br/>--file<br/>TMLIST_FILE|Process file lists|false|
|Boolean|dry<br/>--dry<br/>TMLIST_DRY|Dry run - do not modify any lists|false|

If none of dir, ext or file are provided they all supposed to be true and TMList processes all lists by default.
