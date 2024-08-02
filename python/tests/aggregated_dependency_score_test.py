import pytest

from aggregdepscore import *

@pytest.mark.parametrize("score", range(11))
def test_trustworthiness_inversion(score):
    x = score

    y = score_from_trustworthiness(trustworthiness_from_score(x))

    # note that we have to account for floating point arithmetic errors
    assert abs(x - y) < 0.0001, (x, y)
