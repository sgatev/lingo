# lingo

Lingo helps you define and enforce project-specific Go lingo.

## Installation

To install the latest version of lingo, execute:

```sh
go get -u github.com/s2gatev/lingo
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

## Copyright

Copyright (c) 2017 Stanislav Gatev. See [LICENSE](LICENSE) for
further details.
