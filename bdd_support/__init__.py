"""Code shared by the behave and pytest-bdd step definitions.

Layers, from most readable to most technical:

    features/*.feature    WHAT the system should do (plain language, no code)
    step definitions      GLUE: match a sentence, pull out values, call the driver
    bdd_support.*_driver  HOW to do it: business actions in Python,
                          one driver per team (cart, payments, shipping, checkout)
    src/                  the application itself

Because both runners use this package, the "how" is written once.
"""
