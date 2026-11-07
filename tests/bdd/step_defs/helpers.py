from pytest_bdd import parsers

from bdd_support.types import TYPES


def p(pattern: str):
    """Shorthand: a parse-style step pattern that knows our Money/Qty types."""
    return parsers.parse(pattern, extra_types=TYPES)


def as_dicts(datatable):
    """pytest-bdd gives a table as a list of lists; the first row is the header."""
    header, *rows = datatable
    return [dict(zip(header, row)) for row in rows]
