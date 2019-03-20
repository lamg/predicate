# Predicate

Predicate is a simple library for parsing, evaluating predicates (boolean functions) and getting a textual representation of them.

## Install

```sh
git clone git@github.com:lamg/predicate.git
cd predicate/cmd/reduce && go install
```

## Example

```sh
reduce 'A ∧ true'
reduce ¬false
reduce 'A ≡ true ⇒ false'
```

outputs:

```
A
true
¬A
```
