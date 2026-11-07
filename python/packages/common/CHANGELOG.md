# Changelog

## 1.1.0
- New `behave_runner.SpecsRunner`: runs feature files that live in another
  folder (the shared `specs/` at the repo root) than their `environment.py`
  and `steps/`. Use it with `runner = "shop_bdd_common.behave_runner:SpecsRunner"`.

## 1.0.0
- First release: shared "I should be told" step, Money/Qty types, Outcome,
  bundled-data loader and step deprecation helper.
