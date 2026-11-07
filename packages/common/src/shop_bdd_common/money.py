from decimal import ROUND_HALF_UP, Decimal

CENT = Decimal("0.01")


def to_money(value) -> Decimal:
    """Turn a number or string into a Decimal rounded to whole cents."""
    return Decimal(str(value)).quantize(CENT, rounding=ROUND_HALF_UP)
