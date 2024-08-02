from aggregdepscore.package import *

def test_package_from_any():
    assert Package.from_any("https://deps.dev/npm/express/4.18.3") == Package("npm", "express", "4.18.3")
    assert Package.from_any("https://pypi.org/project/requests/2.32.3/") == Package("pypi", "requests", "2.32.3")
    assert Package.from_any("https://www.npmjs.com/package/react/v/18.3.1") == Package("npm", "react", "18.3.1")
