# Nativeblocks CLI

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
nativeblocks auth --username "foo@example.com" --password "foobar1234"
nativeblocks auth -u "foo@example.com" -p "foobar1234"
```

### Organization

#### Set organization id

```bash
nativeblocks organization set
```

#### Get organization id

```bash
nativeblocks region get
```

### Integration

#### Integration list

- -p, --platform, Platform of integration, ANDROID, IOS, REACT
- -k, --kind, Kind of integration, BLOCK, ACTION, LOGGER or ALL

```bash
nativeblocks integration list -p "REACT" -k "ALL"
```

#### Integration sync

- -p, --path, Integration working path

```bash
nativeblocks integration sync
```

#### Integration detail

- -p, --path, Integration working path
- -i, --integrationId, Integration working path

```bash
nativeblocks integration -i "2222-2222-2222-2222" -p "/Users/sample/projects/awesome_project/integrations/button"
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
