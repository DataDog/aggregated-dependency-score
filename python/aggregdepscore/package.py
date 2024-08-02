from dataclasses import dataclass
import re

_REGEXPS = {
    "npm": [
        re.compile(r"^https://www.npmjs.com/package/(?P<name>[^/]+)/v/(?P<version>[^/]+)"),
    ],
    "pypi": [
        re.compile(r"^https://pypi.org/project/(?P<name>[^/]+)/(?P<version>[^/]+)"),
    ],
}

_DEPS_DOT_DEV_PACKAGE_REGEXP = re.compile(r"https://deps.dev/(?P<ecosystem>[^/]+)/(?P<name>[^/]+)/(?P<version>[^/]+)")


@dataclass
class Package:
    ecosystem: str
    name: str
    version: str

    @classmethod
    def from_any(cls, p):
        if isinstance(p, Package):
            return p

        assert isinstance(p, str), type(p)

        if match := _DEPS_DOT_DEV_PACKAGE_REGEXP.match(p):
            return Package(match.group("ecosystem"), match.group("name"), match.group("version"))

        for ecosystem, regexps in _REGEXPS.items():
            for regexp in regexps:
                if match := regexp.match(p):
                    return Package(ecosystem, match.group("name"), match.group("version"))

        raise ValueError(f"could not parse package")

    def __str__(self):
        """TODO follow the PURL spec more closely"""
        return f"pkg:{self.ecosystem}/{self.name}@{self.version}"

