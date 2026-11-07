"""A behave runner for feature files that live apart from their step wiring.

The feature files are a language-neutral contract in specs/<team>/ at the
repo root. The behave wiring (environment.py and steps/) lives in the Python
tree. Plain behave cannot run that layout: it takes its base directory from
the first path, and that directory must hold the steps/ or environment.py
AND at least one .feature file.

This runner keeps the same idea but drops the second rule. Give it the wiring
directory first, then the spec paths:

    [tool.behave]
    runner = "shop_bdd_common.behave_runner:SpecsRunner"
    paths = ["features", "../../../specs/cart"]

The first path is used only to find environment.py and steps/; the remaining
paths are the feature files to run. If the first path is not a wiring
directory (for example `behave ../../../specs/cart/cart.feature` on the
command line), the wiring is taken from ./features.
"""

import os

from behave.exception import ConfigError
from behave.runner import Runner


class SpecsRunner(Runner):
    def _is_wiring(self, path: str) -> bool:
        return os.path.isdir(os.path.join(path, self.config.steps_dir)) or os.path.isfile(
            os.path.join(path, self.config.environment_file)
        )

    def setup_paths(self):
        paths = [getattr(p, "filename", p) for p in self.config.paths]
        if paths and self._is_wiring(paths[0]):
            wiring = paths.pop(0)
        else:
            wiring = self.DEFAULT_DIRECTORY
        wiring = os.path.abspath(wiring)
        if not self._is_wiring(wiring):
            raise ConfigError(f"No {self.config.steps_dir} directory or {self.config.environment_file} in {wiring!r}")
        if not paths:
            raise ConfigError("No feature paths: list the specs folder after the wiring folder")

        self.config.paths = paths
        self.config.base_dir = self.base_dir = wiring
        self.path_manager.add(wiring)
        if wiring != os.getcwd():
            self.path_manager.add(os.getcwd())
