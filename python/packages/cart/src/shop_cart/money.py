from decimal import ROUND_HALF_UP, Decimal

CENT = Decimal("0.01")


def to_money(value) -> Decimal:
    """Turn a number or string into a Decimal rounded to whole cents."""
    return Decimal(str(value)).quantize(CENT, rounding=ROUND_HALF_UP)


def format_money(amount: Decimal) -> str:
    return f"${to_money(amount):.2f}"
