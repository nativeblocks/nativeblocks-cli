# Nativeblocks CLI

### Installation

```bash
curl -fsSL https://nativeblocks.io/download/cli/installer.sh | bash
```

You can always find all command by help command

```bash
nativeblocks help
```

### Region

#### Set a region

```bash
nativeblocks region set "https://api.example.com"
```

#### Get the region

```bash
nativeblocks region get
```

### Auth

#### Auth with username and password

- -u, --username, username
- -p, --password, password

```bash
nativeblocks auth login --username "foo@example.com" --password "foobar1234"
nativeblocks auth login -u "foo@example.com" -p "foobar1234"
```

#### Auth with access token

- -a, --accessToken, access token

```bash
nativeblocks auth token --accessToken "123.123.123"
nativeblocks auth token -a "123.123.123"
```

### Organization

#### Set organization

```bash
nativeblocks organization set
```

#### Get organization

```bash
nativeblocks organization get
```

### Project

#### Set project

```bash
nativeblocks project set
```

#### Get project

```bash
nativeblocks project get
```

#### Generate project schema

Generates project base schema with found blocks and actions, you need to upload them on a public url to use for frame
and code-gen commands

- -p, --path, Project working path
- -e, --edition, Edition type (cloud or community)

```bash
nativeblocks project gen-schema -e cloud -p /Users/sample/projects/awesome_project
```

### Integration

#### Integration list

- -p, --platform, Platform of integration, ANDROID, IOS, REACT
- -k, --kind, Kind of integration, BLOCK, ACTION, LOGGER or ALL

```bash
nativeblocks integration list -p "REACT" -k "ALL"
```

#### Integration sync

Note: To sync an integration please make sure you pass the .nativeblocks directory

- -p, --path, Integration working path

```bash
nativeblocks integration sync -p /Users/sample/projects/awesome_project/integrations/button/.nativeblocks
```

#### Integration detail

- -p, --path, Integration working path
- -i, --integrationId, Integration working path

```bash
nativeblocks integration -i "2222-2222-2222-2222" -p /Users/sample/projects/awesome_project/integrations/button/.nativeblocks
```

### Frame

#### Frame generate

- -p, --path, Frame working path

```bash
nativeblocks frame gen -p "/Users/sample/projects/awesome_project/frame/login"
```

#### Frame push

- -p, --path, Frame working path

```bash
nativeblocks frame push -p "/Users/sample/projects/awesome_project/frame/login"
```

#### Frame pull

- -p, --path, Frame working path

```bash
nativeblocks frame pull -p "/Users/sample/projects/awesome_project/frame/login"
```

### Frame

#### Codegen typescript

- -a, --actionsSchemaUrl string Actions schema url
- -b, --blocksSchemaUrl string Blocks schema url
- -p, --path string Output path

```bash
nativeblocks code-gen typescript -p "/Users/sample/projects/src/integrations" -b https://publich-address.com/registered-blocks.json -a https://publich-address.com/registered-actions.json

```

#### Codegen php

- -a, --actionsSchemaUrl string Actions schema url
- -b, --blocksSchemaUrl string Blocks schema url
- -p, --path string Output path

```bash
nativeblocks code-gen php -p "/Users/sample/projects/src/Integrations" -b https://publich-address.com/registered-blocks.json -a https://publich-address.com/registered-actions.json

```

#### Codegen go

- -a, --actionsSchemaUrl string Actions schema url
- -b, --blocksSchemaUrl string Blocks schema url
- -p, --path string Output path

```bash
nativeblocks code-gen go -p "/Users/sample/projects/integration" -b https://publich-address.com/registered-blocks.json -a https://publich-address.com/registered-actions.json

```
