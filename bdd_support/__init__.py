"""Code shared by the behave and pytest-bdd step definitions.

Layers, from most readable to most technical:

    features/*.feature    WHAT the system should do (plain language, no code)
    step definitions      GLUE: match a sentence, pull out values, call the driver
    bdd_support.driver    HOW to do it: business actions in Python
    src/shop              the application itself

Because both runners use this package, the "how" is written once.
"""
