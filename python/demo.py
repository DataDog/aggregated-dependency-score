import logging

from aggregdepscore import *
from aggregdepscore import deps_dot_dev

logging.basicConfig(level=logging.DEBUG)

logging.getLogger('urllib3.connectionpool').setLevel(logging.INFO)
logging.getLogger('asyncio').setLevel(logging.INFO)

PACKAGE = "https://pypi.org/project/requests/2.28.1/"

print(f"{PACKAGE = }")
print()

trustworthiness = compute_aggregrated_trustworthiness(
    PACKAGE,
    get_direct_dependencies=deps_dot_dev.get_direct_dependencies,
    compute_intrinsic_score=deps_dot_dev.compute_intrinsic_score
)

score = score_from_trustworthiness(trustworthiness)

print()
print(f"{trustworthiness = }")
print(f"score: {score*10:.1f}/10")
