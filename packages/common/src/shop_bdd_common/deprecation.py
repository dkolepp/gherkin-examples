"""Retire a step sentence without breaking the feature files that use it.

Changing a sentence is a breaking change for everyone who uses it. So for one
minor release, keep the old sentence working and warn:

    @then(old_pattern)
    @deprecated(old='...', new='...', since="1.1.0")
    def old_step(...):
        return new_step(...)

Then remove it in the next major release.
"""

import functools
import warnings


class StepDeprecationWarning(UserWarning):
    """Shown by default (unlike DeprecationWarning) in both behave and pytest."""


def deprecated(old: str, new: str, since: str):
    def decorate(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            warnings.warn(
                f'Step "{old}" is deprecated since {since}. Use "{new}" instead.',
                StepDeprecationWarning,
                stacklevel=2,
            )
            return func(*args, **kwargs)

        return wrapper

    return decorate
