import logging
from math import log


MIN_TRUSTWORTHINESS = 0.8
# this is noted as "k" in the paper
TRUSTWORTHINESS_OFFSET_PARAM = 60
# this is noted as "e" in the paper
TRANSITIVE_TRUSTWORTHINESS_POWER = 0.5


def trustworthiness_from_score(score: float) -> float:
    """this function is noted as "f" in the paper"""
    return 1 - (1 - MIN_TRUSTWORTHINESS) * (
        1 - log(1 + (TRUSTWORTHINESS_OFFSET_PARAM - 1)*score)
        / log(TRUSTWORTHINESS_OFFSET_PARAM)
    )


def score_from_trustworthiness(trustworthiness: float) -> float:
    """this function is noted as $\\tilde{f}^{-1}$ in the paper"""

    t = trustworthiness
    k = TRUSTWORTHINESS_OFFSET_PARAM

    if t < MIN_TRUSTWORTHINESS:
        return 0

    return (
        ( 1 - k**(1 - (1-t)/0.2) )
        / (1 - k)
    )

def compute_aggregrated_trustworthiness(
    package,
    get_direct_dependencies,
    compute_intrinsic_score,
):

    s = compute_intrinsic_score(package)
    t = trustworthiness_from_score(s)

    logging.debug(f"intrinsic trustworthiness for {package}: {t} (score: {s})")

    t_hat = compute_transitive_trustworthiness(
        package,
        get_direct_dependencies,
        compute_intrinsic_score,
    )

    logging.debug(f"transitive trustworthiness for {package}: {t_hat}")

    return t * t_hat


def compute_transitive_trustworthiness(
    package,
    get_direct_dependencies,
    compute_intrinsic_score,
):
    e = TRANSITIVE_TRUSTWORTHINESS_POWER

    result = 1

    for dep in get_direct_dependencies(package):
        t_prime = compute_aggregrated_trustworthiness(
            dep,
            get_direct_dependencies,
            compute_intrinsic_score,
        )

        logging.debug(f"aggregated trustworthiness for {dep} (dependency of {package}): {t_prime}")

        t_hat = compute_transitive_trustworthiness(
            dep,
            get_direct_dependencies,
            compute_intrinsic_score,
        )

        logging.debug(f"transitive trustworthiness for {dep} (dependency of {package}): {t_hat}")

        result *= t_prime * t_hat**e

    return result

