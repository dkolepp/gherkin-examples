import json
from importlib import resources
from pathlib import Path


def load_json(package: str, file_name: str):
    """Load JSON test data.

    A file at that path in the current folder wins, so a journey can bring its
    own data. Otherwise the file comes from the package's bundled data/ folder,
    like an Ansible role's defaults.
    """
    local = Path(file_name)
    if local.is_file():
        return json.loads(local.read_text())
    return json.loads(resources.files(package).joinpath("data", file_name).read_text())
