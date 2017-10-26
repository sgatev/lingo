# lingo

[![Build Status](https://travis-ci.org/s2gatev/lingo.svg?branch=master)](https://travis-ci.org/s2gatev/lingo)

Lingo helps you define and enforce project-specific Go lingo.

## Description

Lingo is a tool for fully-automated static code analysis of Go code. It is designed
to be integrated in continuous delivery pipelines and act as a single source of truth
for a project's code style.

## Installation

To install the latest version of lingo, execute:

```sh
go get -u github.com/s2gatev/lingo
```

## Configuration

The static analysis checks performed by lingo are defined in a configuration file
that needs to be provided upon execution. The following example defines five
static analysis checks that are executed against all `.go` files except those
under `vendor/` and those with names ending in `_test`:

```yaml
matchers:
  -
    type: 'glob'
    config:
      pattern: '**/*.go'
  -
    type: 'not'
    config:
      type: 'glob'
      config:
        pattern: '**/vendor/**/*'
  -
    type: 'not'
    config:
      type: 'glob'
      config:
        pattern: '**/*_test.go'

checkers:
  local_return:
  multi_word_ident_name:
  exported_ident_doc:
  test_package:
  consistent_receiver_names:
```

## Checking

To check all files rooted at the current directory for lingo viloations execute:

```sh
lingo check ./... --config lingo.yml
```

## Contributing

1. Fork the project
2. Clone your fork (`git clone https://github.com/username/lingo && cd lingo`)
3. Create a feature branch (`git checkout -b new-feature`)
4. Make changes and add them (`git add .`)
5. Make sure tests are passing and coverage is good (`go test ./... -race -cover`)
6. Make sure code style is matching the lingo of the project (`lingo ./...`)
7. Commit your changes (`git commit -m 'Add some feature'`)
8. Push the branch (`git push origin new-feature`)
9. Create a new pull request

## Credits

Lingo draws huge inspiration from [RuboCop](https://github.com/bbatsov/rubocop) and
other Go tools such as [dep](https://github.com/golang/dep), [cobra](https://github.com/spf13/cobra)
and [golint](https://github.com/golang/lint).

## Copyright

Copyright (c) 2017 Stanislav Gatev. See [LICENSE](LICENSE) for
further details.
