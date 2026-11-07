# Changelog

## 1.1.0
- New `RunSpecs`: runs feature files that live in another folder (the shared
  `specs/` at the repo root) with godog, strict, each scenario as a Go
  subtest. The Go counterpart of the Python `SpecsRunner`.

## 1.0.0
- First release: `Install`/`OutcomeFrom` and the shared "I should be told"
  step, `Money`/`Qty` parsing, `Outcome`, `Rows` for data tables, the
  bundled-data loader and the step deprecation helper.
