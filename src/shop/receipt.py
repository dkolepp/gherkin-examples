from shop.cart import Cart


def render_receipt(cart: Cart) -> str:
    rows = [
        f"{product.name:<14}{qty:>3} x {product.price:>6.2f} = {line:>7.2f}"
        for product, qty, line in cart.lines()
    ]
    rows.append(f"{'Subtotal':<30}{cart.subtotal:>7.2f}")
    if cart.discount_amount:
        label = f"Discount ({cart.discount.code})"
        rows.append(f"{label:<30}{-cart.discount_amount:>7.2f}")
    rows.append(f"{'TOTAL':<30}{cart.total:>7.2f}")
    return "\n".join(rows)
