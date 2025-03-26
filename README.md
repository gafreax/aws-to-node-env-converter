# Aws to node environment converter (atnec)

A simple CLI utility to transform environment files by removing comments and converting the format.

## Features

- Remove comments from environment files
- Convert `KEY=VALUE` to `KEY='VALUE'` format
- Simple and fast CLI tool

## Installation

```bash
git clone https://github.com/yourusername/dotenv-transformer.git
cd dotenv-transformer
make install
```

## Usage

```bash
# Basic usage
atnec input.env output.env

# Example in package.json
"scripts": {
  "dev:stage": "atnec aws-file-stage.env .env-stage && node server.js"
}
```

## Development

- Build: `make build`
- Install: `make install`
- Clean: `make clean`