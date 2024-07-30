import logging
import urllib.parse

# TODO add caching
# and possibly rate limiting
import requests

from packages import Package

_ECOSYSTEMS = {
    "npm",
    "pypi",
    "go",
    "maven",
}

def get_direct_dependencies(package):
    p = Package.from_any(package)

    if p.ecosystem not in _ECOSYSTEMS:
        raise ValueError(f"ecosystem not supported with Deps.dev: {p.ecosystem}")

    r = requests.get("/".join([
        "https://api.deps.dev/v3",
        "systems",
        p.ecosystem,
        "packages",
        p.name,
        "versions",
        p.version+":dependencies",
    ]))
    r.raise_for_status()

    data = r.json()

    result = []

    for dep in data["nodes"]:
        if dep["relation"] == "DIRECT":
            try:
                result.append(Package(
                    p.ecosystem,
                    dep["versionKey"]["name"],
                    dep["versionKey"]["version"],
                ))
            except Exception as e:
                e.add_note(f"{dep = }")
                raise e from None

    return result


def compute_intrinsic_score(package):
    p = Package.from_any(package)

    if p.ecosystem not in _ECOSYSTEMS:
        raise ValueError(f"ecosystem not supported with Deps.dev: {p.ecosystem}")

    project = _get_package_project(p)

    if not project:
        raise ValueError(f"no project found for {p}")

    r = requests.get("/".join([
        "https://api.deps.dev/v3",
        "projects",
        # note the safe="" so that slashes are encoded as well
        urllib.parse.quote(project, safe=""),
    ]))
    r.raise_for_status()

    data = r.json()

    scorecard = data.get("scorecard")

    if not scorecard:
        raise ValueError(f"no OSSF scorecard available in deps.dev ({project = })")

    score = float(scorecard["overallScore"])

    assert 0 <= score <= 10, score

    return score / 10


def _get_package_project(package: Package):
    r = requests.get("/".join([
        "https://api.deps.dev/v3",
        "systems",
        package.ecosystem,
        "packages",
        package.name,
        "versions",
        package.version,
    ]))

    r.raise_for_status()

    data = r.json()

    projects = data.get("relatedProjects", [])

    if not projects:
        return None

    if len(projects) > 1:
        source_repos = [
            p for p in projects
            if p.get("relationType") == "SOURCE_REPO"
        ]

        if len(source_repos) == 1:
            return source_repos[0]["projectKey"]["id"]

        logging.warning(f"multiple related projects found for {package} ({ projects = }); picking the first one")

    return projects[0]["projectKey"]["id"]
