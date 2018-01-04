## consistent_receiver_names

Checks that method receivers of a type are named consistently.

```go
func (i *Item) Execute() {} // OK
func (t *Item) Cancel() {}  // not OK
```

## exported_ident_doc

Checks that all exported identifiers are documented (including field names, variable,
constant, function, struct and interface declarations). Optionally checks that each
comment starts with the name of the item it describes.

Available options:
* `has_ident_prefix: bool` - ensure that every doc comment begins with the name of the
item it describes.

```go
// Runner is documented.
type Runner interface {

    // Run is documented.
    Run()

    // This method is not documented properly.
    Stop() error
}
```

## left_quantifiers

Checks that when a numerical quantifier appears in a binary expression it is the left
operand.

```go
_ = 5 * time.Minute // OK
_ = time.Minute * 5 // not OK
```

## line_length

Checks that the code lines are within specific length limits.

Available options:
* `max_length: int` - the maximum number of characters permitted on a single line.
* `tab_width: int` - the number of characters equivalent to a single tab.

## local_return

Checks that exported functions return exported (and internal) types only.

```go
func (i *Item) Do() result {} // not OK, `result` should be exported
```

## multi_word_ident_name

Checks the correctness of type names. Correct type names adhere to the following rules:
* PascalCase for exported types.
* camelCase for non-exported types.

```go
type processTracker struct{}  // OK
type ProcessTracker struct{}  // OK
type process_tracker struct{} // not OK
```

## return_error_last

Checks that `error` is the last value returned by a function.

```go
func Create() (int, error, bool) {} // not OK, `error` should be last
```

## test_package

Checks that tests are placed in `*_test` packages only.

```go
package feature // not OK, should be `feature_test`

import "testing"

func TestFeature(t *testing.T) {}
```
