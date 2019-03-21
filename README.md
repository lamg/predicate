# Predicate

Predicate is a simple library for parsing, evaluating predicates (boolean functions) and getting a textual representation of them. The syntax is based on https://www.cs.utexas.edu/users/EWD/transcriptions/EWD13xx/EWD1300.html, which I have formalized in the following grammar:

```ebnf
predicate = term ('≡'|'≢') term {('≡'|'≢') term}| term.
term = implication | consequence | junction.
implication = junction '⇒' junction {'⇒' junction}.
consequence = junction '⇐' junction {'⇐' junction}.
junction = disjunction | conjunction | factor.
disjunction = factor '∨' factor {'∨' factor}.
conjunction = factor '∧' factor {'∧' factor}.
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
```

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
