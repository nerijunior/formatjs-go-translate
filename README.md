# formatjs-go-translate

A CLI tool that automatically translates missing entries in locale JSON files using the Google Cloud Translate API.

## How it works

The tool compares a source locale file (e.g. `defaultMessages.json`) against a target locale file (e.g. `pt-BR.json`). Any key that is missing from the target file, or whose value is prefixed with `[TODO]`, gets translated via Google Cloud Translate and written back to the target file.

Translations are sent in batches of 100 strings to stay within API limits.

## Requirements

- Go 1.18+
- A Google Cloud project with the **Cloud Translation API** enabled
- A service account credentials file (`google-creds.json`) in the working directory

## Installation

```bash
git clone https://github.com/nerijunior/formatjs-go-translate
cd formatjs-go-translate
go build -o formatjs-go-translate .
```

## Usage

```bash
./formatjs-go-translate [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--source-file` | `locales/defaultMessages.json` | Path to the source locale file |
| `--locale-file` | `locales/pt-BR.json` | Path to the target locale file to update |
| `--diff` | `false` | Print a comparison table without translating |

### Examples

Translate missing keys using the default paths:

```bash
./formatjs-go-translate
```

Translate using custom paths:

```bash
./formatjs-go-translate --source-file locales/en.json --locale-file locales/es.json
```

Preview differences without translating:

```bash
./formatjs-go-translate --diff
```

## Locale file format

Both source and target files must be JSON objects where each key maps to an object with a `string` field (and an optional `context` field):

```json
{
  "welcome.title": {
    "string": "Welcome to our store",
    "context": "Homepage hero heading"
  },
  "cart.empty": {
    "string": "Your cart is empty"
  }
}
```

Entries in the target file whose `string` value starts with `[TODO]` are treated as untranslated and will be retranslated.

## Google Cloud credentials

Place a service account key file named `google-creds.json` in the directory where you run the tool. The service account needs the **Cloud Translation API User** role.

To generate credentials:

1. Go to **IAM & Admin → Service Accounts** in the Google Cloud Console
2. Create a service account and grant it the `Cloud Translation API User` role
3. Create a JSON key and save it as `google-creds.json`
