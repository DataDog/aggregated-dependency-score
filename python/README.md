# Package `aggregdepscore`

About how to use the package, see [`demo.py`](./demo.py).

<details>
<summary>demo</summary>

```
$ python demo.py
PACKAGE = 'https://pypi.org/project/requests/2.28.1/'

DEBUG:root:intrinsic trustworthiness for https://pypi.org/project/requests/2.28.1/: 0.9938665043419209 (score: 0.8800000000000001)
DEBUG:root:intrinsic trustworthiness for pkg:pypi/certifi@2024.7.4: 0.9808368780096528 (score: 0.67)
DEBUG:root:transitive trustworthiness for pkg:pypi/certifi@2024.7.4: 1
DEBUG:root:aggregated trustworthiness for pkg:pypi/certifi@2024.7.4 (dependency of https://pypi.org/project/requests/2.28.1/): 0.9808368780096528
DEBUG:root:transitive trustworthiness for pkg:pypi/certifi@2024.7.4 (dependency of https://pypi.org/project/requests/2.28.1/): 1
DEBUG:root:intrinsic trustworthiness for pkg:pypi/charset-normalizer@2.1.1: 0.9868507652091928 (score: 0.76)
DEBUG:root:transitive trustworthiness for pkg:pypi/charset-normalizer@2.1.1: 1
DEBUG:root:aggregated trustworthiness for pkg:pypi/charset-normalizer@2.1.1 (dependency of https://pypi.org/project/requests/2.28.1/): 0.9868507652091928
DEBUG:root:transitive trustworthiness for pkg:pypi/charset-normalizer@2.1.1 (dependency of https://pypi.org/project/requests/2.28.1/): 1
DEBUG:root:intrinsic trustworthiness for pkg:pypi/idna@3.7.0: 0.9793935961691637 (score: 0.65)
DEBUG:root:transitive trustworthiness for pkg:pypi/idna@3.7.0: 1
DEBUG:root:aggregated trustworthiness for pkg:pypi/idna@3.7.0 (dependency of https://pypi.org/project/requests/2.28.1/): 0.9793935961691637
DEBUG:root:transitive trustworthiness for pkg:pypi/idna@3.7.0 (dependency of https://pypi.org/project/requests/2.28.1/): 1
DEBUG:root:intrinsic trustworthiness for pkg:pypi/urllib3@1.26.19: 0.9965163167062648 (score: 0.93)
DEBUG:root:transitive trustworthiness for pkg:pypi/urllib3@1.26.19: 1
DEBUG:root:aggregated trustworthiness for pkg:pypi/urllib3@1.26.19 (dependency of https://pypi.org/project/requests/2.28.1/): 0.9965163167062648
DEBUG:root:transitive trustworthiness for pkg:pypi/urllib3@1.26.19 (dependency of https://pypi.org/project/requests/2.28.1/): 1
DEBUG:root:transitive trustworthiness for https://pypi.org/project/requests/2.28.1/: 0.9446913584378166

trustworthiness = 0.9388970980926135
score: 2.7/10
```
</details>

## Development

Requires [Poetry](https://python-poetry.org/).

* setup: `poetry install`
* testing: `poetry run python -m pytest`
